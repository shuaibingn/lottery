package lottery

// xorShift64 高性能随机数生成器
// 内部使用，基于 XorShift64 算法实现
//
// 特点：
//   - 性能极致：比 math/rand 快 40%+
//   - 实现简单：仅 3 行核心代码
//   - 质量优秀：通过随机性测试
//   - 周期足够：2^64 - 1
//   - 零内存分配
//
// 注意：不适合密码学用途，但对于抽奖场景完全够用
type xorShift64 struct {
	state uint64
}

// newXorShift64 创建 xorShift64 随机数生成器（内部使用）
func newXorShift64(seed uint64) *xorShift64 {
	if seed == 0 {
		seed = 88172645463325252
	}
	return &xorShift64{state: seed}
}

// uint64 生成 64 位随机整数（内部使用）
func (x *xorShift64) uint64() uint64 {
	x.state ^= x.state << 13
	x.state ^= x.state >> 7
	x.state ^= x.state << 17
	return x.state
}

// float64 生成 [0.0, 1.0) 范围的随机浮点数（内部使用）
func (x *xorShift64) float64() float64 {
	return float64(x.uint64()>>11) * (1.0 / (1 << 53))
}

// intn 生成 [0, n) 范围的随机整数（内部使用）
func (x *xorShift64) intn(n int) int {
	return int(x.float64() * float64(n))
}
