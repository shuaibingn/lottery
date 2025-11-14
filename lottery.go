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
	// 构建 Alias Method 查找表
	table, err := buildAliasTables(data)
	if err != nil {
		return nil, err
	}

	// 生成随机种子
	seed, err := generateSecureSeed()
	if err != nil {
		return nil, errors.New("failed to generate random seed: " + err.Error())
	}

	return &AliasMethod{
		prob:  table.prob,
		alias: table.alias,
		keys:  table.keys,
		n:     table.n,
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
