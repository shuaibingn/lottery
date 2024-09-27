# lottery

可以根据传递的概率快速抽奖, 并且解决了概率为浮点数时的精度问题

## 安装

使用`go get`安装本模块

```shell
go get -u github.com/shuaibingn/lottery
```

## 使用

1. 首先需要使用匿名结构体`*lottery.DrawBase`来实现抽奖接口
2. 构建一个抽奖切片, 里面包含ID, 概率等
3. 使用`InitLotteries`函数初始化抽奖, 并且传递一个放大倍数, 保证概率和该数字的积为整数
4. 调用`Draw`方法抽奖

```go
package main

import (
	"fmt"

	"github.com/shuaibingn/lottery"
)

type Card struct {
	*lottery.DrawBase `json:",inline"`
	Name              string `json:"name"`
}

func main() {
	cards := []lottery.Lottery{
		&Card{&lottery.DrawBase{ID: "1", Probability: 0.1}, "SSR"},
		&Card{&lottery.DrawBase{ID: "2", Probability: 0.2}, "SR"},
		&Card{&lottery.DrawBase{ID: "3", Probability: 0.3}, "R"},
		&Card{&lottery.DrawBase{ID: "4", Probability: 0.4}, "N"},
	}

	lotteries, err := lottery.InitLotteries(cards, 10)
	if err != nil {
		panic(err)
	}

	id := lotteries.Draw()
	fmt.Println(id)
}
```
