package test

import (
	"strings"
	"testing"

	"github.com/shuaibingn/lottery"
)

// TestAliasMethodProbabilityValidation 测试 Alias Method 的概率验证
func TestAliasMethodProbabilityValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        []lottery.Item
		expectError bool
		errorMsg    string
	}{
		{
			name: "正确的概率和（1.0）",
			data: []lottery.Item{
				&Data{&lottery.DrawBase{ID: "1", Probability: 0.5}},
				&Data{&lottery.DrawBase{ID: "2", Probability: 0.5}},
			},
			expectError: false,
		},
		{
			name: "概率和过小（< 1.0）",
			data: []lottery.Item{
				&Data{&lottery.DrawBase{ID: "1", Probability: 0.3}},
				&Data{&lottery.DrawBase{ID: "2", Probability: 0.5}},
			},
			expectError: true,
			errorMsg:    "sum of probabilities must be approximately 1.0",
		},
		{
			name: "概率和过大（> 1.0）",
			data: []lottery.Item{
				&Data{&lottery.DrawBase{ID: "1", Probability: 0.7}},
				&Data{&lottery.DrawBase{ID: "2", Probability: 0.5}},
			},
			expectError: true,
			errorMsg:    "sum of probabilities must be approximately 1.0",
		},
		{
			name: "负概率",
			data: []lottery.Item{
				&Data{&lottery.DrawBase{ID: "1", Probability: -0.1}},
				&Data{&lottery.DrawBase{ID: "2", Probability: 1.1}},
			},
			expectError: true,
			errorMsg:    "probability cannot be negative",
		},
		{
			name: "接近 1.0 的概率（允许误差）",
			data: []lottery.Item{
				&Data{&lottery.DrawBase{ID: "1", Probability: 0.33333}},
				&Data{&lottery.DrawBase{ID: "2", Probability: 0.33333}},
				&Data{&lottery.DrawBase{ID: "3", Probability: 0.33334}},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am, err := lottery.New(tt.data)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("错误信息不匹配：期望包含 %q, 实际得到 %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了：%v", err)
				}
				if am == nil {
					t.Errorf("期望返回非 nil 实例")
				}
			}
		})
	}
}

// TestAliasMethodPoolProbabilityValidation 测试 Alias Method Pool 的概率验证
func TestAliasMethodPoolProbabilityValidation(t *testing.T) {
	tests := []struct {
		name        string
		data        []lottery.Item
		expectError bool
		errorMsg    string
	}{
		{
			name: "正确的概率和（1.0）",
			data: []lottery.Item{
				&Data{&lottery.DrawBase{ID: "1", Probability: 0.5}},
				&Data{&lottery.DrawBase{ID: "2", Probability: 0.5}},
			},
			expectError: false,
		},
		{
			name: "概率和不正确",
			data: []lottery.Item{
				&Data{&lottery.DrawBase{ID: "1", Probability: 0.3}},
				&Data{&lottery.DrawBase{ID: "2", Probability: 0.5}},
			},
			expectError: true,
			errorMsg:    "sum of probabilities must be approximately 1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amp, err := lottery.NewPool(tt.data)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("错误信息不匹配：期望包含 %q, 实际得到 %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了：%v", err)
				}
				if amp == nil {
					t.Errorf("期望返回非 nil 实例")
				}
			}
		})
	}
}

