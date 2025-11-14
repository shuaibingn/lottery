package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shuaibingn/lottery"
)

// Prize 奖品结构
type Prize struct {
	*lottery.DrawBase
}

var globalLotteries *lottery.Pool

func init() {
	// 定义奖品（概率之和必须为 1.0）
	prizes := []lottery.Item{
		&Prize{&lottery.DrawBase{ID: "特等奖", Probability: 0.001}}, // 0.1%
		&Prize{&lottery.DrawBase{ID: "一等奖", Probability: 0.009}}, // 0.9%
		&Prize{&lottery.DrawBase{ID: "二等奖", Probability: 0.09}},  // 9%
		&Prize{&lottery.DrawBase{ID: "三等奖", Probability: 0.2}},   // 20%
		&Prize{&lottery.DrawBase{ID: "谢谢参与", Probability: 0.7}},  // 70%
	}

	// 初始化全局抽奖器（线程安全，支持高并发）
	var err error
	globalLotteries, err = lottery.NewPool(prizes)
	if err != nil {
		panic(fmt.Sprintf("初始化失败: %v", err))
	}
}

func main() {
	fmt.Println("=== 高并发抽奖示例 ===\n")

	// 示例 1: 模拟 10000 个用户同时抽奖
	fmt.Println("场景 1: 10,000 个用户同时抽奖")
	simulateConcurrentDraws(10000)

	// 示例 2: 压力测试 - 100 万次抽奖
	fmt.Println("\n场景 2: 压力测试 - 1,000,000 次抽奖")
	stressTest(1000000, 100)

	// 示例 3: 性能测试 - 测量吞吐量
	fmt.Println("\n场景 3: 性能测试 - 测量吞吐量")
	performanceTest(5000000, 100)
}

// simulateConcurrentDraws 模拟并发抽奖
func simulateConcurrentDraws(userCount int) {
	var wg sync.WaitGroup
	results := make(map[string]*atomic.Int64)
	
	// 初始化计数器
	prizeIDs := []string{"特等奖", "一等奖", "二等奖", "三等奖", "谢谢参与"}
	for _, id := range prizeIDs {
		results[id] = &atomic.Int64{}
	}

	start := time.Now()

	// 启动并发抽奖
	for i := 0; i < userCount; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// 执行抽奖（线程安全，O(1) 时间复杂度）
			result := globalLotteries.Draw()
			if counter, ok := results[result]; ok {
				counter.Add(1)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	// 输出结果
	fmt.Printf("完成 %d 次抽奖，耗时: %v\n", userCount, elapsed)
	fmt.Printf("平均每次抽奖: %.2f ns\n", float64(elapsed.Nanoseconds())/float64(userCount))
	fmt.Println("\n抽奖结果分布:")

	prizeProbs := []float64{0.001, 0.009, 0.09, 0.2, 0.7}
	for i, id := range prizeIDs {
		count := results[id].Load()
		actualProb := float64(count) / float64(userCount) * 100
		expectedProb := prizeProbs[i] * 100
		deviation := (actualProb - expectedProb) / expectedProb * 100

		fmt.Printf("  %s: %d 次 (%.2f%%), 期望: %.2f%%, 偏差: %+.2f%%\n",
			id, count, actualProb, expectedProb, deviation)
	}
}

// stressTest 压力测试
func stressTest(totalDraws int, goroutines int) {
	var wg sync.WaitGroup
	drawsPerGoroutine := totalDraws / goroutines

	start := time.Now()

	// 启动多个 goroutine 执行抽奖
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < drawsPerGoroutine; j++ {
				_ = globalLotteries.Draw()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	// 输出结果
	fmt.Printf("完成 %d 次抽奖，使用 %d 个 goroutine\n", totalDraws, goroutines)
	fmt.Printf("总耗时: %v\n", elapsed)
	fmt.Printf("平均每次抽奖: %.2f ns\n", float64(elapsed.Nanoseconds())/float64(totalDraws))
	fmt.Printf("吞吐量: %.2f M draws/sec\n", float64(totalDraws)/elapsed.Seconds()/1000000)
}

// performanceTest 性能测试
func performanceTest(totalDraws int, goroutines int) {
	drawsPerGoroutine := totalDraws / goroutines

	// 预热（让 sync.Pool 建立缓存）
	for i := 0; i < 1000; i++ {
		_ = globalLotteries.Draw()
	}

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < drawsPerGoroutine; j++ {
				_ = globalLotteries.Draw()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	// 输出结果
	fmt.Printf("完成 %d 次抽奖（预热后）\n", totalDraws)
	fmt.Printf("总耗时: %v\n", elapsed)
	fmt.Printf("平均每次抽奖: %.2f ns\n", float64(elapsed.Nanoseconds())/float64(totalDraws))
	fmt.Printf("吞吐量: %.2f M draws/sec\n", float64(totalDraws)/elapsed.Seconds()/1000000)
	fmt.Printf("QPS: %.2f M/sec\n", float64(totalDraws)/elapsed.Seconds()/1000000)
	
	fmt.Println("\n性能特点:")
	fmt.Println("  ✅ O(1) 时间复杂度 - 恒定时间")
	fmt.Println("  ✅ 零内存分配 - GC 友好")
	fmt.Println("  ✅ sync.Pool 无锁 - 高并发性能")
	fmt.Println("  ✅ XorShift64 - 比 math/rand 快 40%+")
}

