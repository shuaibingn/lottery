package lottery

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
)

type LotteriesPool struct {
	lotteries  []Lottery
	mul        float64
	randPool   sync.Pool     // 随机数生成器对象池
	seedSource atomic.Uint64 // 原子计数器，确保每个池中对象的种子不同
}

func (lotteries *LotteriesPool) Draw() string {
	rng := lotteries.randPool.Get().(*rand.Rand)
	defer lotteries.randPool.Put(rng)

	randomNumber := int64(rng.Float64() * lotteries.mul)

	var cumulativeProbability int64 = 0
	for _, lottery := range lotteries.lotteries {
		cumulativeProbability += lottery.getProbabilityInt64()
		if randomNumber < cumulativeProbability {
			return lottery.getID()
		}
	}
	return ""
}

func NewLotteriesPool(data []Lottery) (*LotteriesPool, error) {
	mul, err := calculateMul(data)
	if err != nil {
		return nil, err
	}

	return InitLotteriesPool(data, mul)
}

func InitLotteriesPool(data []Lottery, mul float64) (*LotteriesPool, error) {
	if err := validateAndSetProbabilities(data, mul); err != nil {
		return nil, err
	}

	seed, err := generateSecureSeed()
	if err != nil {
		return nil, errors.New("failed to generate random seed: " + err.Error())
	}

	lotteries := &LotteriesPool{
		lotteries: data,
		mul:       mul,
	}
	lotteries.seedSource.Store(uint64(seed))

	lotteries.randPool = sync.Pool{
		New: func() interface{} {
			newSeed := int64(lotteries.seedSource.Add(1))
			return rand.New(rand.NewSource(newSeed))
		},
	}

	return lotteries, nil
}
