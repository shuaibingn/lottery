package lottery

// Item 抽奖奖品接口
// 所有奖品类型都需要实现此接口
type Item interface {
	// getProbability 获取奖品的抽中概率（0.0 - 1.0）
	getProbability() float64
	
	// getID 获取奖品的唯一标识符
	getID() string
}

// DrawBase 抽奖基础结构
// 用户可以嵌入此结构来快速实现 Item 接口
//
// 示例：
//   type Prize struct {
//       *lottery.DrawBase
//       // 其他自定义字段
//   }
type DrawBase struct {
	ID          string  `json:"id"`          // 奖品ID
	Probability float64 `json:"probability"` // 抽中概率（0.0 - 1.0，所有奖品概率之和应为 1.0）
}

// getProbability 获取浮点概率
func (b *DrawBase) getProbability() float64 {
	return b.Probability
}

// getID 获取奖项ID
func (b *DrawBase) getID() string {
	return b.ID
}

