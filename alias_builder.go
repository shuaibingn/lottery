package lottery

import "errors"

// aliasTable Alias Method 算法的查找表
//
// Alias Method 是一种将不均匀概率分布转换为均匀分布的算法
// 通过预处理构建查找表，使得抽奖时可以在 O(1) 时间完成
//
// 数据结构：
//   - prob: 概率表，存储每个桶的接受概率
//   - alias: 别名表，存储每个桶的备选项索引
//   - keys: ID 映射表，将索引映射到奖品ID
//   - n: 奖品数量
type aliasTable struct {
	prob  []float64 // 概率表
	alias []int     // 别名表
	keys  []string  // ID映射表
	n     int       // 选项数量
}

// buildAliasTables 构建 Alias Method 算法所需的查找表
//
// 这是 Alias Method 算法的核心实现（Vose's Alias Method）
// 将不均匀的概率分布转换为均匀分布 + 别名表，使得抽奖时可以 O(1) 完成
//
// 算法步骤：
//   1. 验证输入数据（非空、概率非负）
//   2. 验证概率和是否约等于 1.0（±0.0001 误差）
//   3. 缩放概率：将每个概率乘以 n，使平均概率为 1.0
//   4. 分离 small 和 large：
//      - small: 缩放后概率 < 1.0 的项
//      - large: 缩放后概率 >= 1.0 的项
//   5. 配对构建表：
//      - 从 small 和 large 各取一项配对
//      - 填充 small 项的概率和别名
//      - 更新 large 项的剩余概率
//      - 重新分类更新后的 large 项
//   6. 处理剩余项（由于浮点误差）
//
// 时间复杂度：O(n)
// 空间复杂度：O(n)
//
// 参数：
//   - data: 奖品列表，实现 Lottery 接口
//
// 返回：
//   - *aliasTable: 构建好的查找表，可用于 O(1) 抽奖
//   - error: 验证失败时返回错误
//
// 错误情况：
//   - data 为空
//   - 存在负概率
//   - 概率和不约等于 1.0
func buildAliasTables(data []Item) (*aliasTable, error) {
	// ===== 步骤1: 验证输入并初始化 =====
	if len(data) == 0 {
		return nil, errors.New("lotteries must be greater than 0")
	}

	n := len(data)
	prob := make([]float64, n)
	alias := make([]int, n)
	keys := make([]string, n)

	// ===== 步骤2: 提取概率并验证 =====
	sum := 0.0
	for i, d := range data {
		keys[i] = d.getID()
		p := d.getProbability()
		if p < 0 {
			return nil, errors.New("probability cannot be negative")
		}
		sum += p
	}

	// ===== 步骤3: 验证概率和 =====
	// 允许浮点误差 ±0.0001 (0.01%)
	if sum < 0.9999 || sum > 1.0001 {
		return nil, errors.New("sum of probabilities must be approximately 1.0")
	}

	// ===== 步骤4: 缩放概率 =====
	// 将每个概率乘以 n，使平均概率变为 1.0
	// 这样可以将问题转换为"如何将 n 个桶填充到 n 个位置"
	scaled := make([]float64, n)
	for i, d := range data {
		scaled[i] = d.getProbability() * float64(n) / sum
	}

	// ===== 步骤5: 分离 small 和 large =====
	// small: 缩放后概率 < 1.0 的项（需要填充）
	// large: 缩放后概率 >= 1.0 的项（有多余概率）
	small := make([]int, 0, n)
	large := make([]int, 0, n)

	for i, p := range scaled {
		if p < 1.0 {
			small = append(small, i)
		} else {
			large = append(large, i)
		}
	}

	// ===== 步骤6: 构建概率表和别名表 =====
	// Vose's Alias Method 的核心：配对 small 和 large
	for len(small) > 0 && len(large) > 0 {
		// 取出一个 small 和一个 large
		s := small[len(small)-1]
		small = small[:len(small)-1]

		l := large[len(large)-1]
		large = large[:len(large)-1]

		// 设置 small 项的概率和别名
		// prob[s]: 从桶 s 中选择 keys[s] 的概率
		// alias[s]: 如果不选 keys[s]，则选 keys[l]
		prob[s] = scaled[s]
		alias[s] = l

		// 更新 large 项的剩余概率
		// large 项用掉了 (1.0 - scaled[s]) 的概率来填充 small 项
		scaled[l] = scaled[l] + scaled[s] - 1.0

		// 根据剩余概率，重新分类 large 项
		if scaled[l] < 1.0 {
			small = append(small, l)
		} else {
			large = append(large, l)
		}
	}

	// ===== 步骤7: 处理剩余的项 =====
	// 由于浮点误差，可能还有剩余的 small 或 large
	// 这些项的概率应该非常接近 1.0，直接设置为 1.0
	for len(large) > 0 {
		l := large[len(large)-1]
		large = large[:len(large)-1]
		prob[l] = 1.0
	}

	for len(small) > 0 {
		s := small[len(small)-1]
		small = small[:len(small)-1]
		prob[s] = 1.0
	}

	return &aliasTable{
		prob:  prob,
		alias: alias,
		keys:  keys,
		n:     n,
	}, nil
}

