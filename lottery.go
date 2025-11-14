package lottery

import (
	"errors"
)

// Lottery 高性能抽奖器（无锁版本）
// 基于 Alias Method 算法实现，提供 O(1) 时间复杂度的恒定抽奖性能
//
// 特点：
//   - O(1) 时间复杂度（恒定抽奖时间）
//   - 使用 XorShift64 高性能随机数生成器（比 math/rand 快 40%+）
//   - 非线程安全，适合单线程或每个 goroutine 独立实例
//
// 如需线程安全版本，请使用 Pool
type Lottery struct {
	prob  []float64   // 概率表
	alias []int       // 别名表
	keys  []string    // ID映射表
	n     int         // 选项数量
	rand  *xorShift64 // 随机数生成器
}

// New 创建一个新的抽奖器实例
//
// 时间复杂度：初始化 O(n)，抽奖 O(1)
//
// 参数：
//   - data: 奖品列表，所有奖品的概率之和必须约等于 1.0
//
// 返回：
//   - *Lottery: 抽奖器实例
//   - error: 如果概率验证失败则返回错误
//
// 注意：返回的实例不是线程安全的
// 如需在多个 goroutine 间共享使用，请使用 NewPool() 创建线程安全版本
func New(data []Item) (*Lottery, error) {
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

	return &Lottery{
		prob:  table.prob,
		alias: table.alias,
		keys:  table.keys,
		n:     table.n,
		rand:  newXorShift64(uint64(seed)),
	}, nil
}

// Draw 执行一次抽奖
//
// 时间复杂度：O(1) - 无论多少奖项都是恒定时间
//
// 返回：
//   - string: 抽中的奖品 ID
//
// 注意：此方法不是线程安全的
func (l *Lottery) Draw() string {
	// 使用 XorShift64 生成随机数（比 math/rand 快 40%+）
	i := l.rand.intn(l.n)
	r := l.rand.float64()

	if r < l.prob[i] {
		return l.keys[i]
	}
	return l.keys[l.alias[i]]
}
