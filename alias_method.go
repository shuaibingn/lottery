package lottery

import (
	"errors"
)

// AliasMethod 别名方法实现的抽奖器（无锁版本）
// 特点：
//   - O(1) 时间复杂度（恒定抽奖时间）
//   - 使用 XorShift64 高性能随机数生成器（比 math/rand 快 40%+）
//   - 非线程安全，适合单线程或每个 goroutine 独立实例
//
// 如需线程安全版本，请使用 AliasMethodPool
type AliasMethod struct {
	prob  []float64   // 概率表
	alias []int       // 别名表
	keys  []string    // ID映射表
	n     int         // 选项数量
	rand  *XorShift64 // 随机数生成器
}

// NewAliasMethod 创建无锁版本的 Alias Method 抽奖器
// 时间复杂度：初始化 O(n)，抽奖 O(1)
// 注意：返回的实例不是线程安全的
func NewAliasMethod(data []Lottery) (*AliasMethod, error) {
	if len(data) == 0 {
		return nil, errors.New("lotteries must be greater than 0")
	}

	n := len(data)
	prob := make([]float64, n)
	alias := make([]int, n)
	keys := make([]string, n)

	// 提取概率并验证
	sum := 0.0
	for i, d := range data {
		keys[i] = d.getID()
		p := d.getProbability()
		if p < 0 {
			return nil, errors.New("probability cannot be negative")
		}
		sum += p
	}

	// 验证概率和是否接近 1.0（允许浮点误差）
	if sum < 0.9999 || sum > 1.0001 {
		return nil, errors.New("sum of probabilities must be approximately 1.0")
	}

	// 缩放概率（使期望值为1）
	scaled := make([]float64, n)
	for i, d := range data {
		scaled[i] = d.getProbability() * float64(n) / sum
	}

	// 分离 small 和 large
	small := make([]int, 0, n)
	large := make([]int, 0, n)

	for i, p := range scaled {
		if p < 1.0 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	// 构建概率表和别名表
	for len(small) > 0 && len(large) > 0 {
		s := small[len(small)-1]
		small = small[:len(small)-1]

		l := large[len(large)-1]
		large = large[:len(large)-1]

		prob[s] = scaled[s]
		alias[s] = l

		scaled[l] = scaled[l] + scaled[s] - 1.0

		if scaled[l] < 1.0 {
			small = append(small, l)
		} else {
			large = append(large, l)
		}
	}

	// 处理剩余的
	for len(large) > 0 {
		l := large[len(large)-1]
		large = large[:len(large)-1]
		prob[l] = 1.0
	}

	for len(small) > 0 {
		s := small[len(small)-1]
		small = small[:len(small)-1]
		prob[s] = 1.0
	}

	// 生成随机种子
	seed, err := generateSecureSeed()
	if err != nil {
		return nil, errors.New("failed to generate random seed: " + err.Error())
	}

	return &AliasMethod{
		prob:  prob,
		alias: alias,
		keys:  keys,
		n:     n,
		rand:  newXorShift64(uint64(seed)),
	}, nil
}

// Draw 执行一次抽奖
// 时间复杂度：O(1) - 无论多少奖项都是恒定时间
// 注意：此方法不是线程安全的
func (am *AliasMethod) Draw() string {
	// 使用 XorShift64 生成随机数（比 math/rand 快 40%+）
	i := am.rand.Intn(am.n)
	r := am.rand.Float64()

	if r < am.prob[i] {
		return am.keys[i]
	}
	return am.keys[am.alias[i]]
}
