package test

import (
	"fmt"
	"testing"
	"time"

	pkg "lottery"
)

type Data struct {
	*pkg.DrawBase
}

func TestDraw(t *testing.T) {
	lotteries, err := pkg.InitLotteries(10000,
		&Data{
			DrawBase: &pkg.DrawBase{
				ID:          "SSR",
				Probability: 0.1,
			},
		},
		&Data{
			DrawBase: &pkg.DrawBase{
				ID:          "SR",
				Probability: 0.2,
			},
		},
		&Data{
			DrawBase: &pkg.DrawBase{
				ID:          "R",
				Probability: 0.3,
			},
		},
		&Data{
			DrawBase: &pkg.DrawBase{
				ID:          "N",
				Probability: 0.4,
			},
		},
	)
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
