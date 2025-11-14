package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/shuaibingn/lottery"
)

// Data 测试用奖品结构
type Data struct {
	*lottery.DrawBase
}

// TestLottery 测试无锁版本抽奖器（单线程）
func TestLottery(t *testing.T) {
	data := []lottery.Item{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	lotteries, err := lottery.New(data)
	if err != nil {
		panic(err)
	}

	start := time.Now().UnixNano()
	result := make(map[string]int)
	for i := 0; i < 100000; i++ {
		id := lotteries.Draw()
		if _, ok := result[id]; ok {
			result[id]++
			continue
		}
		result[id] = 1
	}
	end := time.Now().UnixNano()
	fmt.Println("无锁版本（单线程）:", result, end-start, "ns")
}

// TestPool 测试线程安全抽奖器（单线程）
func TestPool(t *testing.T) {
	data := []lottery.Item{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	pool, err := lottery.NewPool(data)
	if err != nil {
		panic(err)
	}

	start := time.Now().UnixNano()
	result := make(map[string]int)
	for i := 0; i < 100000; i++ {
		id := pool.Draw()
		if _, ok := result[id]; ok {
			result[id]++
			continue
		}
		result[id] = 1
	}
	end := time.Now().UnixNano()
	fmt.Println("Pool 版本（单线程）:", result, end-start, "ns")
}

// TestPoolConcurrent 测试线程安全抽奖器的并发安全性
func TestPoolConcurrent(t *testing.T) {
	data := []lottery.Item{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	pool, err := lottery.NewPool(data)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	workers := 100
	drawsPerWorker := 1000

	start := time.Now()

	// 100 个 goroutine 并发抽奖
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < drawsPerWorker; j++ {
				_ = pool.Draw()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	totalDraws := workers * drawsPerWorker
	t.Logf("✓ Pool 并发测试完成: %d个goroutine, 总计%d次抽奖, 耗时%v", workers, totalDraws, elapsed)
	t.Logf("  平均每次抽奖: %v", elapsed/time.Duration(totalDraws))
}

// TestProbability 测试抽奖概率准确性
func TestProbability(t *testing.T) {
	data := []lottery.Item{
		&Data{&lottery.DrawBase{ID: "rare", Probability: 0.01}},   // 1%
		&Data{&lottery.DrawBase{ID: "common", Probability: 0.99}}, // 99%
	}

	lotteries, err := lottery.New(data)
	if err != nil {
		t.Fatal(err)
	}

	// 大量测试验证概率
	const testCount = 1000000
	result := make(map[string]int)

	for i := 0; i < testCount; i++ {
		id := lotteries.Draw()
		result[id]++
	}

	// 验证概率误差在 ±0.5% 以内
	rareRate := float64(result["rare"]) / testCount
	commonRate := float64(result["common"]) / testCount

	t.Logf("rare: %.4f%% (期望 1%%), common: %.4f%% (期望 99%%)", rareRate*100, commonRate*100)

	if rareRate < 0.005 || rareRate > 0.015 {
		t.Errorf("rare 概率异常: %.4f%%, 期望 1%%", rareRate*100)
	}
	if commonRate < 0.985 || commonRate > 0.995 {
		t.Errorf("common 概率异常: %.4f%%, 期望 99%%", commonRate*100)
	}
}
