package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/shuaibingn/lottery"
)

type Data struct {
	*lottery.DrawBase
}

// TestDraw 测试无锁版本（单线程）
func TestDraw(t *testing.T) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	// 使用新的 API，自动计算 mul
	lotteries, err := lottery.NewLotteries(data)
	if err != nil {
		panic(err)
	}

	start := time.Now().UnixNano()
	result := make(map[string]int)
	for i := 0; i < 100000; i++ {
		id, err := lotteries.Draw()
		if err != nil {
			t.Errorf("Draw failed: %v", err)
			continue
		}
		if _, ok := result[id]; ok {
			result[id]++
			continue
		}
		result[id] = 1
	}
	end := time.Now().UnixNano()
	fmt.Println("无锁版本（单线程）:", result, end-start, "ns")
}

// TestDrawPool 测试 sync.Pool 版本（单线程）
func TestDrawPool(t *testing.T) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	// 使用新的 API，自动计算 mul
	lotteries, err := lottery.NewLotteriesPool(data)
	if err != nil {
		panic(err)
	}

	start := time.Now().UnixNano()
	result := make(map[string]int)
	for i := 0; i < 100000; i++ {
		id, err := lotteries.Draw()
		if err != nil {
			t.Errorf("Draw failed: %v", err)
			continue
		}
		if _, ok := result[id]; ok {
			result[id]++
			continue
		}
		result[id] = 1
	}
	end := time.Now().UnixNano()
	fmt.Println("sync.Pool版本（单线程）:", result, end-start, "ns")
}

// TestDrawPoolConcurrent 测试 sync.Pool 版本的并发安全性
func TestDrawPoolConcurrent(t *testing.T) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	// 使用新的 API，自动计算 mul
	lotteries, err := lottery.NewLotteriesPool(data)
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
				_, err := lotteries.Draw()
				if err != nil {
					t.Errorf("Draw failed: %v", err)
				}
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	
	totalDraws := workers * drawsPerWorker
	t.Logf("✓ 并发测试完成: %d个goroutine, 总计%d次抽奖, 耗时%v", workers, totalDraws, elapsed)
	t.Logf("  平均每次抽奖: %v", elapsed/time.Duration(totalDraws))
}
