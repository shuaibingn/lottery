package lottery

import (
	"errors"
	"math/rand"
)

// Lotteries 无锁版本抽奖器
// 注意：此版本不是线程安全的，适用于：
//   - 单线程场景
//   - 每个 goroutine 独立创建实例
//   - 外部自行处理并发控制
//
// 如需线程安全的高性能版本，请使用 LotteriesPool
type Lotteries struct {
	lotteries []Lottery
	mul       float64
	localRand *rand.Rand
}

// Draw 执行一次抽奖
func (lotteries *Lotteries) Draw() string {
	randomNumber := int64(lotteries.localRand.Float64() * lotteries.mul)

	var cumulativeProbability int64 = 0
	for _, lottery := range lotteries.lotteries {
		cumulativeProbability += lottery.getProbabilityInt64()
		if randomNumber < cumulativeProbability {
			return lottery.getID()
		}
	}
	return ""
}

func NewLotteries(data []Lottery) (*Lotteries, error) {
	mul, err := calculateMul(data)
	if err != nil {
		return nil, err
	}

	return InitLotteries(data, mul)
}

func InitLotteries(data []Lottery, mul float64) (*Lotteries, error) {
	if err := validateAndSetProbabilities(data, mul); err != nil {
		return nil, err
	}

	seed, err := generateSecureSeed()
	if err != nil {
		return nil, errors.New("failed to generate random seed: " + err.Error())
	}

	return &Lotteries{
		lotteries: data,
		mul:       mul,
		localRand: rand.New(rand.NewSource(seed)),
	}, nil
}
