package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/shuaibingn/lottery"
)

type Data struct {
	*lottery.DrawBase
}

func TestDraw(t *testing.T) {
	data := []lottery.Lottery{
		&Data{&lottery.DrawBase{ID: "1", Probability: 0.1}},
		&Data{&lottery.DrawBase{ID: "2", Probability: 0.2}},
		&Data{&lottery.DrawBase{ID: "3", Probability: 0.3}},
		&Data{&lottery.DrawBase{ID: "4", Probability: 0.4}},
	}

	lotteries, err := lottery.InitLotteries(data, 10000)
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
	fmt.Println(result, end-start)
}
