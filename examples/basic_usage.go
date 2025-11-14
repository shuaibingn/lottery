package main

import (
	"fmt"
	"github.com/shuaibingn/lottery"
)

// Prize 奖品结构
type Prize struct {
	*lottery.DrawBase
}

func main() {
	fmt.Println("=== 基础用法示例 ===\n")

	// 定义奖品（概率之和必须为 1.0）
	prizes := []lottery.Lottery{
		&Prize{&lottery.DrawBase{ID: "特等奖", Probability: 0.001}}, // 0.1%
		&Prize{&lottery.DrawBase{ID: "一等奖", Probability: 0.009}}, // 0.9%
		&Prize{&lottery.DrawBase{ID: "二等奖", Probability: 0.09}},  // 9%
		&Prize{&lottery.DrawBase{ID: "三等奖", Probability: 0.2}},   // 20%
		&Prize{&lottery.DrawBase{ID: "谢谢参与", Probability: 0.7}},  // 70%
	}

	// 初始化抽奖器（线程安全版本，推荐）
	lotteries, err := lottery.NewAliasMethodPool(prizes)
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		return
	}

	// 单次抽奖
	fmt.Println("单次抽奖示例:")
	result := lotteries.Draw()
	fmt.Printf("恭喜您抽中: %s\n\n", result)

	// 多次抽奖测试
	fmt.Println("10次抽奖测试:")
	for i := 0; i < 10; i++ {
		result := lotteries.Draw()
		fmt.Printf("第%d次: %s\n", i+1, result)
	}

	// 概率验证（大量抽奖）
	fmt.Println("\n概率验证 (100,000次抽奖):")
	counts := make(map[string]int)
	testTimes := 100000

	for i := 0; i < testTimes; i++ {
		result := lotteries.Draw()
		counts[result]++
	}

	fmt.Printf("\n抽奖次数: %d\n", testTimes)
	fmt.Println("结果分布:")
	
	// 按照定义顺序输出
	prizeIDs := []string{"特等奖", "一等奖", "二等奖", "三等奖", "谢谢参与"}
	prizeProbs := []float64{0.001, 0.009, 0.09, 0.2, 0.7}
	
	for i, id := range prizeIDs {
		count := counts[id]
		actualProb := float64(count) / float64(testTimes) * 100
		expectedProb := prizeProbs[i] * 100
		deviation := (actualProb - expectedProb) / expectedProb * 100
		
		fmt.Printf("  %s: %d 次 (%.2f%%), 期望: %.2f%%, 偏差: %+.2f%%\n",
			id, count, actualProb, expectedProb, deviation)
	}
}

