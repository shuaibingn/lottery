package test

import (
	"testing"

	"github.com/shuaibingn/lottery"
)

// ========================================
// 单线程性能对比
// ========================================

// BenchmarkLockFree_SingleThread 无锁版本（单线程）
func BenchmarkLockFree_SingleThread(b *testing.B) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	lotteries, err := lottery.InitLotteries(data, 10000)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = lotteries.Draw()
	}
}

// BenchmarkPool_SingleThread sync.Pool 版本（单线程）
func BenchmarkPool_SingleThread(b *testing.B) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	lotteries, err := lottery.InitLotteriesPool(data, 10000)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = lotteries.Draw()
	}
}

// ========================================
// 并发性能对比（关键测试）
// ========================================

// BenchmarkPool_Parallel sync.Pool 版本（并发）
func BenchmarkPool_Parallel(b *testing.B) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	lotteries, err := lottery.InitLotteriesPool(data, 10000)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = lotteries.Draw()
		}
	})
}

// ========================================
// 大量奖项场景性能测试
// ========================================

// BenchmarkLockFree_LargeItems 无锁版本（100个奖项）
func BenchmarkLockFree_LargeItems(b *testing.B) {
	data := make([]lottery.Lottery, 100)
	for i := 0; i < 100; i++ {
		data[i] = &Data{&lottery.DrawBase{
			ID:          string(rune('0' + i%10)),
			Probability: 0.01,
		}}
	}

	lotteries, err := lottery.InitLotteries(data, 100000)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = lotteries.Draw()
	}
}

// BenchmarkPool_LargeItems sync.Pool 版本（100个奖项）
func BenchmarkPool_LargeItems(b *testing.B) {
	data := make([]lottery.Lottery, 100)
	for i := 0; i < 100; i++ {
		data[i] = &Data{&lottery.DrawBase{
			ID:          string(rune('0' + i%10)),
			Probability: 0.01,
		}}
	}

	lotteries, err := lottery.InitLotteriesPool(data, 100000)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = lotteries.Draw()
	}
}

// BenchmarkPool_LargeItems_Parallel sync.Pool 版本（100个奖项，并发）
func BenchmarkPool_LargeItems_Parallel(b *testing.B) {
	data := make([]lottery.Lottery, 100)
	for i := 0; i < 100; i++ {
		data[i] = &Data{&lottery.DrawBase{
			ID:          string(rune('0' + i%10)),
			Probability: 0.01,
		}}
	}

	lotteries, err := lottery.InitLotteriesPool(data, 100000)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = lotteries.Draw()
		}
	})
}

// ========================================
// 真实场景模拟
// ========================================

// BenchmarkPool_RealWorldScenario 真实场景：5个奖项，高并发
func BenchmarkPool_RealWorldScenario(b *testing.B) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "特等奖", Probability: 0.001}},
		&Data{&lottery.DrawBase{ID: "一等奖", Probability: 0.01}},
		&Data{&lottery.DrawBase{ID: "二等奖", Probability: 0.089}},
		&Data{&lottery.DrawBase{ID: "三等奖", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "谢谢参与", Probability: 0.7}},
	}

	lotteries, err := lottery.InitLotteriesPool(data, 100000)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = lotteries.Draw()
		}
	})
}

