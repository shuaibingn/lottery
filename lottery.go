package lottery

import (
	"errors"
	"hash/maphash"
	"math/rand"
)

type Lottery interface {
	getProbability() float64
	getProbabilityInt64() int64
	getID() string
	setInt64Probability(int64)
}

type Lotteries struct {
	lotteries []Lottery
	mul       float64
	localRand *rand.Rand
}

func (lotteries Lotteries) Draw() string {
	var (
		cumulativeProbability int64 = 0
		randomNumber                = int64(lotteries.localRand.Float64() * lotteries.mul)
	)

	for _, lottery := range lotteries.lotteries {
		cumulativeProbability += lottery.getProbabilityInt64()
		if randomNumber < cumulativeProbability {
			return lottery.getID()
		}
	}
	return ""
}

func InitLotteries(mul float64, data ...Lottery) (*Lotteries, error) {
	if len(data) == 0 {
		return nil, errors.New("lotteries must be greater than 0")
	}

	var (
		sumProbabilities int64 = 0
		maxProbability         = 1 * int64(mul)
	)

	for _, d := range data {
		probability := int64(d.getProbability() * mul)
		sumProbabilities += probability
		d.setInt64Probability(probability)
	}

	if sumProbabilities != maxProbability {
		return nil, errors.New("cumulative probability must be approximately 1 when scaled")
	}

	return &Lotteries{
		lotteries: data,
		mul:       mul,
		localRand: rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64()))),
	}, nil
}

type DrawBase struct {
	ID             string  `json:"id"`
	Probability    float64 `json:"probability"`
	intProbability int64
}

func (b *DrawBase) getProbability() float64 {
	return b.Probability
}

func (b *DrawBase) getProbabilityInt64() int64 {
	return b.intProbability
}

func (b *DrawBase) getID() string {
	return b.ID
}

func (b *DrawBase) setInt64Probability(probability int64) {
	b.intProbability = probability
}
