package lottery

import (
	"errors"
	"sync"
	"sync/atomic"
)

// AliasMethodPool 高性能线程安全的 Alias Method 抽奖器
// 特点：
//   - O(1) 时间复杂度（恒定抽奖时间）
//   - 使用 XorShift64 高性能随机数生成器（比 math/rand 快 40%+）
//   - 线程安全，使用 sync.Pool 实现无锁并发
//   - 零内存分配，GC 友好
type AliasMethodPool struct {
	prob       []float64 // 概率表
	alias      []int     // 别名表
	keys       []string  // ID映射表
	n          int       // 选项数量
	randPool   sync.Pool // XorShift64 随机数生成器对象池
	seedSource atomic.Uint64
}

// NewAliasMethodPool 创建线程安全的 Alias Method 抽奖器
// 时间复杂度：初始化 O(n)，抽奖 O(1)
func NewAliasMethodPool(data []Lottery) (*AliasMethodPool, error) {
	// 构建 Alias Method 查找表
	table, err := buildAliasTables(data)
	if err != nil {
		return nil, err
	}

	// 生成初始种子
	seed, err := generateSecureSeed()
	if err != nil {
		return nil, errors.New("failed to generate random seed: " + err.Error())
	}

	am := &AliasMethodPool{
		prob:  table.prob,
		alias: table.alias,
		keys:  table.keys,
		n:     table.n,
	}
	am.seedSource.Store(uint64(seed))

	// 创建 XorShift64 随机数生成器对象池
	// 每个 goroutine 获取独立的生成器实例，实现无锁并发
	am.randPool = sync.Pool{
		New: func() interface{} {
			newSeed := am.seedSource.Add(1)
			return newXorShift64(newSeed)
		},
	}

	return am, nil
}

// Draw 执行一次抽奖
// 时间复杂度：O(1) - 无论多少奖项都是恒定时间
// 线程安全，适合高并发场景
func (am *AliasMethodPool) Draw() string {
	// 从池中获取 XorShift64 随机数生成器（几乎无锁）
	rng := am.randPool.Get().(*XorShift64)
	defer am.randPool.Put(rng)

	// O(1) 抽奖（使用 XorShift64，比 math/rand 快 40%+）
	i := rng.Intn(am.n)
	r := rng.Float64()

	if r < am.prob[i] {
		return am.keys[i]
	}
	return am.keys[am.alias[i]]
}
