package lottery

import (
	"crypto/rand"
	"encoding/binary"
)

// generateSecureSeed 使用 crypto/rand 生成安全的随机种子
//
// 特点：
//   - 使用密码学安全的随机数生成器（crypto/rand）
//   - 即使在同一纳秒内多次调用，也能保证种子不同
//   - 适合高并发场景，避免种子冲突
//
// 原理：
//   - crypto/rand 从操作系统的熵池读取真随机数
//   - 每次调用都会生成不同的 64 位随机数
//
// 返回：
//   - int64: 生成的随机种子
//   - error: 如果读取随机数失败则返回错误
func generateSecureSeed() (int64, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(b[:])), nil
}

