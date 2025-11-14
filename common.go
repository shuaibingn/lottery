package lottery

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
)

// Lottery 抽奖接口
// 所有奖项必须实现此接口
type Lottery interface {
	getProbability() float64
	getProbabilityInt64() int64
	getID() string
	setInt64Probability(int64)
}

// DrawBase 抽奖基础结构体
// 包含奖项的基本属性，可以直接嵌入使用
type DrawBase struct {
	ID             string  `json:"id"`
	Probability    float64 `json:"probability"`
	intProbability int64
}

// getProbability 获取浮点概率
func (b *DrawBase) getProbability() float64 {
	return b.Probability
}

// getProbabilityInt64 获取整数概率（已放大）
func (b *DrawBase) getProbabilityInt64() int64 {
	return b.intProbability
}

// getID 获取奖项ID
func (b *DrawBase) getID() string {
	return b.ID
}

// setInt64Probability 设置整数概率
func (b *DrawBase) setInt64Probability(probability int64) {
	b.intProbability = probability
}

// generateSecureSeed 使用 crypto/rand 生成安全的随机种子
// 即使在同一纳秒内多次调用，也能保证种子不同
func generateSecureSeed() (int64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(b[:])), nil
}

// calculateMul 自动计算合适的放大倍数
// 根据概率的精度自动选择 100, 1000, 10000, 100000 或 1000000
// 确保：
//  1. 所有概率乘以 mul 后都能转换为非零整数
//  2. 所有概率之和乘以 mul 等于 mul（概率和为1）
func calculateMul(data []Lottery) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("lotteries must be greater than 0")
	}

	// 候选的 mul 值，从小到大尝试
	candidates := []float64{100, 1000, 10000, 100000, 1000000}

	for _, mul := range candidates {
		valid := true
		var sum int64 = 0

		// 检查这个 mul 是否适合所有概率
		for _, d := range data {
			prob := d.getProbability()
			intProb := int64(prob * mul)

			// 如果概率不为0，但乘以mul后变成0，说明精度不够
			if prob > 0 && intProb == 0 {
				valid = false
				break
			}

			sum += intProb
		}

		// 检查概率和是否等于 mul
		if valid && sum == int64(mul) {
			return mul, nil
		}
	}

	// 如果所有候选值都不合适，返回错误
	return 0, errors.New("cannot find suitable mul value, probabilities may be too precise or sum is not 1.0")
}

// validateAndSetProbabilities 验证并设置概率
// 这是 InitLotteries 和 InitLotteriesPool 共用的逻辑
// 返回：是否验证通过，如果不通过则返回错误
func validateAndSetProbabilities(data []Lottery, mul float64) error {
	if len(data) == 0 {
		return errors.New("lotteries must be greater than 0")
	}

	var (
		sumProbabilities int64 = 0
		maxProbability         = 1 * int64(mul)
	)

	// 计算并设置整数概率
	for _, d := range data {
		probability := int64(d.getProbability() * mul)
		sumProbabilities += probability
		d.setInt64Probability(probability)
	}

	// 验证概率和
	if sumProbabilities != maxProbability {
		return errors.New("cumulative probability must be approximately 1 when scaled")
	}

	return nil
}

