package test

import (
	"testing"

	"github.com/shuaibingn/lottery"
)

// BenchmarkAliasMethod_4Items 测试 4 个奖项的性能
func BenchmarkAliasMethod_4Items(b *testing.B) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	aliasMethod, err := lottery.NewAliasMethod(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = aliasMethod.Draw()
	}
}

// BenchmarkAliasMethod_100Items 测试 100 个奖项的性能
func BenchmarkAliasMethod_100Items(b *testing.B) {
	data := make([]lottery.Lottery, 100)
	prob := 1.0 / 100.0
	for i := 0; i < 100; i++ {
		data[i] = &Data{&lottery.DrawBase{ID: string(rune('A' + i%26)), Probability: prob}}
	}

	aliasMethod, err := lottery.NewAliasMethod(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = aliasMethod.Draw()
	}
}

// BenchmarkAliasMethodPool_4Items 测试 Pool 版本的性能（4 个奖项）
func BenchmarkAliasMethodPool_4Items(b *testing.B) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	aliasMethodPool, err := lottery.NewAliasMethodPool(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = aliasMethodPool.Draw()
	}
}

// BenchmarkAliasMethodPool_100Items 测试 Pool 版本的性能（100 个奖项）
func BenchmarkAliasMethodPool_100Items(b *testing.B) {
	data := make([]lottery.Lottery, 100)
	prob := 1.0 / 100.0
	for i := 0; i < 100; i++ {
		data[i] = &Data{&lottery.DrawBase{ID: string(rune('A' + i%26)), Probability: prob}}
	}

	aliasMethodPool, err := lottery.NewAliasMethodPool(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = aliasMethodPool.Draw()
	}
}

// BenchmarkAliasMethodPool_Parallel 测试并发性能
func BenchmarkAliasMethodPool_Parallel(b *testing.B) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	aliasMethodPool, err := lottery.NewAliasMethodPool(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = aliasMethodPool.Draw()
		}
	})
}

