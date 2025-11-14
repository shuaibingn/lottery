package lottery

import (
	"errors"
	mathrand "math/rand"
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
	localRand *mathrand.Rand
}

// Draw 执行一次抽奖
// 注意：此方法不是线程安全的
func (lotteries *Lotteries) Draw() (string, error) {
	randomNumber := int64(lotteries.localRand.Float64() * lotteries.mul)

	var cumulativeProbability int64 = 0
	for _, lottery := range lotteries.lotteries {
		cumulativeProbability += lottery.getProbabilityInt64()
		if randomNumber < cumulativeProbability {
			return lottery.getID(), nil
		}
	}
	return "", errors.New("no lottery item matched, check probability configuration")
}

// NewLotteries 初始化无锁版本抽奖器（自动计算 mul）
// mul 值会根据概率自动计算，无需手动指定
// 注意：返回的 Lotteries 不是线程安全的
// 如需线程安全的高性能版本，请使用 NewLotteriesPool
func NewLotteries(data []Lottery) (*Lotteries, error) {
	// 自动计算合适的 mul 值
	mul, err := calculateMul(data)
	if err != nil {
		return nil, err
	}

	return InitLotteries(data, mul)
}

// InitLotteries 初始化无锁版本抽奖器（手动指定 mul）
// 推荐使用 NewLotteries，它会自动计算合适的 mul 值
// 注意：返回的 Lotteries 不是线程安全的
// 如需线程安全的高性能版本，请使用 InitLotteriesPool
func InitLotteries(data []Lottery, mul float64) (*Lotteries, error) {
	// 验证并设置概率
	if err := validateAndSetProbabilities(data, mul); err != nil {
		return nil, err
	}

	// 生成随机种子
	seed, err := generateSecureSeed()
	if err != nil {
		return nil, errors.New("failed to generate random seed: " + err.Error())
	}

	return &Lotteries{
		lotteries: data,
		mul:       mul,
		localRand: mathrand.New(mathrand.NewSource(seed)),
	}, nil
}
