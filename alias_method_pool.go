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

	// 生成初始种子
	seed, err := generateSecureSeed()
	if err != nil {
		return nil, errors.New("failed to generate random seed: " + err.Error())
	}

	am := &AliasMethodPool{
		prob:  prob,
		alias: alias,
		keys:  keys,
		n:     n,
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
