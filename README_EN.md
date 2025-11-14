# 🎰 High-Performance Go Lottery Library

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Performance](https://img.shields.io/badge/Performance-4.3ns%2Fop-brightgreen)](README_EN.md#performance-benchmarks)

An **ultra-high-performance** Go lottery library with two implementations:

- 🚀 **Lock-Free Version**: Optimal for single-threaded scenarios (~12.74 ns/op)
- ⚡ **sync.Pool Version**: Excellent for high-concurrency scenarios (~4.32 ns/op), handles **231 million** requests per second

## ✨ Key Features

- 🎯 **Ultra-Fast Performance**: 4.32 ns/op in concurrent scenarios, **40+ times faster** than traditional mutex-based solutions
- 🔒 **Thread-Safe**: sync.Pool version supports high concurrency with zero contention
- 💾 **Zero Heap Allocation**: 0 allocs/op, GC-friendly
- 🎲 **True Randomness**: Uses crypto/rand for seed generation, avoids seed conflicts
- 📊 **Auto Precision**: Automatically calculates probability multiplier, no manual specification needed
- 🔧 **Easy to Use**: Simple API design, get started with just 5 lines of code
- 📈 **Highly Scalable**: Supports 5-1000+ prize items
- ✅ **Production-Ready**: Thoroughly tested, high code quality

## 🚀 Performance Metrics

### Single-Thread Performance

| Scenario | Time | Throughput | Memory |
|----------|------|------------|--------|
| Lock-Free (4 items) | **12.74 ns/op** | 78.5M ops/sec | 0 B/op |
| sync.Pool (4 items) | **20.08 ns/op** | 49.8M ops/sec | 0 B/op |
| Lock-Free (100 items) | **81.11 ns/op** | 12.3M ops/sec | 0 B/op |
| sync.Pool (100 items) | **103.5 ns/op** | 9.6M ops/sec | 0 B/op |

### High-Concurrency Performance ⚡ (Key Metrics)

| Scenario | Time | Throughput | Memory |
|----------|------|------------|--------|
| **sync.Pool Concurrent (4 items)** | **4.32 ns/op** | **231M ops/sec** ⚡ | 0 B/op |
| **sync.Pool Concurrent (100 items)** | **15.44 ns/op** | **64.8M ops/sec** ⚡ | 0 B/op |
| **Real-World (5 items, concurrent)** | **5.02 ns/op** | **199M ops/sec** ⚡ | 0 B/op |

> **💡 Performance Highlight**: sync.Pool version is even faster in concurrent scenarios than single-threaded! This is because:
> - Per-P local cache, mostly lock-free operations
> - Fully utilizes multi-core CPUs
> - Zero memory allocation, no GC pressure

### Performance Comparison

```
Traditional Mutex:    ~85 ns/op  (concurrent)
This Library:         ~4.3 ns/op (concurrent)
Performance Gain:     19.7x faster! ⚡⚡⚡
```

## 📦 Installation

```bash
go get github.com/shuaibingn/lottery
```

## 🎯 Quick Start

### Basic Example (Recommended)

```go
package main

import (
    "fmt"
    "github.com/shuaibingn/lottery"
)

type Prize struct {
    *lottery.DrawBase
}

func main() {
    // Define prizes (probabilities must sum to 1)
    prizes := []lottery.Lottery{
        &Prize{&lottery.DrawBase{ID: "First Prize", Probability: 0.01}},   // 1%
        &Prize{&lottery.DrawBase{ID: "Second Prize", Probability: 0.09}},  // 9%
        &Prize{&lottery.DrawBase{ID: "Third Prize", Probability: 0.2}},    // 20%
        &Prize{&lottery.DrawBase{ID: "Thank You", Probability: 0.7}},      // 70%
    }

    // Initialize lottery (auto-calculate precision)
    lotteries, err := lottery.NewLotteriesPool(prizes)
    if err != nil {
        panic(err)
    }

    // Draw
    result := lotteries.Draw()
    fmt.Printf("Congratulations! You won: %s\n", result)
}
```

### High-Concurrency Scenario (Recommended)

```go
package main

import (
    "fmt"
    "sync"
    "github.com/shuaibingn/lottery"
)

var globalLotteries *lottery.LotteriesPool

func init() {
    prizes := []lottery.Lottery{
        &Prize{&lottery.DrawBase{ID: "Grand Prize", Probability: 0.001}},
        &Prize{&lottery.DrawBase{ID: "First Prize", Probability: 0.009}},
        &Prize{&lottery.DrawBase{ID: "Second Prize", Probability: 0.09}},
        &Prize{&lottery.DrawBase{ID: "Third Prize", Probability: 0.2}},
        &Prize{&lottery.DrawBase{ID: "Thank You", Probability: 0.7}},
    }
    
    globalLotteries, _ = lottery.NewLotteriesPool(prizes)
}

func main() {
    var wg sync.WaitGroup
    
    // Simulate 1000 concurrent users
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func(userID int) {
            defer wg.Done()
            
            result := globalLotteries.Draw()
            fmt.Printf("User %d won: %s\n", userID, result)
        }(i)
    }
    
    wg.Wait()
}
```

## 📚 API Documentation

### Two Implementations

| Version | Thread-Safe | Performance | Use Case |
|---------|-------------|-------------|----------|
| **Lotteries** | ❌ No | ⚡⚡⚡ Fastest single-thread | Single-thread, independent instances |
| **LotteriesPool** | ✅ Yes | ⚡⚡⚡⚡⚡ Fastest concurrent | High concurrency, shared instance ⭐ |

### Initialization Methods

#### Auto-Calculate Precision (Recommended)

```go
// Lock-free version
lotteries, err := lottery.NewLotteries(prizes)

// sync.Pool version (Recommended for high concurrency) ⭐
lotteries, err := lottery.NewLotteriesPool(prizes)
```

#### Manual Precision Specification

```go
// Manually specify mul = 10000 (supports 1/10000 precision)
lotteries, err := lottery.InitLotteries(prizes, 10000)
lotteries, err := lottery.InitLotteriesPool(prizes, 10000)
```

### Draw Method

```go
// Execute one draw
result := lotteries.Draw()
fmt.Println("Result:", result)
```

## 🎲 Probability Configuration

### Supported Precision Range

The system automatically selects appropriate precision (mul value):

| Precision | mul Value | Min Probability | Example |
|-----------|-----------|-----------------|---------|
| 1/100 | 100 | 1% | 0.01 |
| 1/1000 | 1,000 | 0.1% | 0.001 |
| 1/10000 | 10,000 | 0.01% | 0.0001 ⭐ Common |
| 1/100000 | 100,000 | 0.001% | 0.00001 |
| 1/1000000 | 1,000,000 | 0.0001% | 0.000001 |

### Probability Examples

```go
prizes := []lottery.Lottery{
    // 1/1000 probability
    &Prize{&lottery.DrawBase{ID: "SSR", Probability: 0.001}},  // 0.1%
    
    // 1/100 probability
    &Prize{&lottery.DrawBase{ID: "SR", Probability: 0.01}},    // 1%
    
    // 1/10 probability
    &Prize{&lottery.DrawBase{ID: "R", Probability: 0.1}},      // 10%
    
    // Remaining probability
    &Prize{&lottery.DrawBase{ID: "N", Probability: 0.889}},    // 88.9%
}
```

**⚠️ Important**: All probabilities must sum to 1.0

## 📊 Performance Optimization Tips

### 1. Choose the Right Version

```go
// ✅ Single-thread or independent instance per goroutine
lotteries, _ := lottery.NewLotteries(prizes)

// ✅ High concurrency with shared instance (Recommended)
lotteries, _ := lottery.NewLotteriesPool(prizes)
```

### 2. Global Singleton Pattern

```go
var (
    globalLotteries *lottery.LotteriesPool
    once            sync.Once
)

func GetLotteries() *lottery.LotteriesPool {
    once.Do(func() {
        prizes := []lottery.Lottery{ /* ... */ }
        globalLotteries, _ = lottery.NewLotteriesPool(prizes)
    })
    return globalLotteries
}

// Use anywhere
result := GetLotteries().Draw()
```

### 3. Reduce Prize Count

Time complexity is O(n), fewer prizes means better performance:

| Prize Count | Avg Time | Recommendation |
|-------------|----------|----------------|
| < 20 | < 30 ns | ✅ Optimal |
| 20-50 | 30-70 ns | ✅ Good |
| 50-100 | 70-150 ns | ⚠️ Acceptable |
| > 100 | > 150 ns | ⚠️ Consider optimization |

## 🔍 Time Complexity

### Draw() Method: O(n)

| Case | Complexity | Description |
|------|-----------|-------------|
| Best | O(1) | First item matches |
| Average | O(n/2) | Traverse about half items |
| Worst | O(n) | Traverse all items |

**Why not binary search (O(log n))?**

- For typical scenarios (< 50 items), linear search is faster:
  - CPU cache-friendly
  - Better branch prediction
  - Simpler code
- Binary search only shows advantage when item count > 100

## 🎯 Use Cases

### ✅ Suitable Scenarios

- 🎮 Game lottery (equipment, items, cards)
- 🎁 Marketing campaigns (coupons, points)
- 🎰 Lottery systems (wheel, mystery box)
- 🏆 Competition rewards
- 📱 Social apps (gifts, effects)
- 🛍️ E-commerce promotions

### ⚠️ Not Suitable For

- Need real-time probability adjustment (not supported yet)
- Probability precision > 1/1000000
- Prize count > 1000 (consider tiered lottery)

## 📖 Complete Examples

### Example 1: Game Equipment Lottery

```go
package main

import (
    "fmt"
    "github.com/shuaibingn/lottery"
)

type Equipment struct {
    *lottery.DrawBase
    Quality string
}

func main() {
    equipments := []lottery.Lottery{
        &Equipment{
            DrawBase: &lottery.DrawBase{ID: "Legendary", Probability: 0.001},
            Quality:  "legendary",
        },
        &Equipment{
            DrawBase: &lottery.DrawBase{ID: "Epic", Probability: 0.009},
            Quality:  "epic",
        },
        &Equipment{
            DrawBase: &lottery.DrawBase{ID: "Rare", Probability: 0.09},
            Quality:  "rare",
        },
        &Equipment{
            DrawBase: &lottery.DrawBase{ID: "Common", Probability: 0.9},
            Quality:  "common",
        },
    }

    lotteries, err := lottery.NewLotteriesPool(equipments)
    if err != nil {
        panic(err)
    }

    // Simulate 10 draws
    results := make(map[string]int)
    for i := 0; i < 10; i++ {
        result := lotteries.Draw()
        results[result]++
    }

    fmt.Println("Draw Results:")
    for name, count := range results {
        fmt.Printf("  %s: %d times\n", name, count)
    }
}
```

### Example 2: Red Packet Lottery

```go
package main

import (
    "fmt"
    "github.com/shuaibingn/lottery"
)

type RedPacket struct {
    *lottery.DrawBase
    Amount float64
}

func main() {
    redPackets := []lottery.Lottery{
        &RedPacket{
            DrawBase: &lottery.DrawBase{ID: "$100", Probability: 0.001},
            Amount:   100.0,
        },
        &RedPacket{
            DrawBase: &lottery.DrawBase{ID: "$50", Probability: 0.01},
            Amount:   50.0,
        },
        &RedPacket{
            DrawBase: &lottery.DrawBase{ID: "$10", Probability: 0.089},
            Amount:   10.0,
        },
        &RedPacket{
            DrawBase: &lottery.DrawBase{ID: "$1", Probability: 0.3},
            Amount:   1.0,
        },
        &RedPacket{
            DrawBase: &lottery.DrawBase{ID: "Thank You", Probability: 0.6},
            Amount:   0,
        },
    }

    lotteries, _ := lottery.NewLotteriesPool(redPackets)

    // Draw
    result := lotteries.Draw()
    fmt.Printf("Congratulations! You got: %s\n", result)
}
```

## 🧪 Testing

```bash
# Run all tests
go test ./test/ -v

# Run performance tests
go test -bench=. -benchmem ./test/

# Data race detection
go test -race ./test/
```

## 📈 Detailed Performance Report

### Test Environment

- **CPU**: Apple M4
- **Architecture**: arm64
- **OS**: macOS
- **Go Version**: 1.21+

### Complete Benchmark Results

```
BenchmarkLockFree_SingleThread-10        78,868,198    12.74 ns/op    0 B/op    0 allocs/op
BenchmarkPool_SingleThread-10            63,218,757    20.08 ns/op    0 B/op    0 allocs/op
BenchmarkPool_Parallel-10               266,882,718     4.32 ns/op    0 B/op    0 allocs/op  ⚡
BenchmarkLockFree_LargeItems-10          14,587,876    81.11 ns/op    0 B/op    0 allocs/op
BenchmarkPool_LargeItems-10              11,683,934   103.50 ns/op    0 B/op    0 allocs/op
BenchmarkPool_LargeItems_Parallel-10     80,511,921    15.44 ns/op    0 B/op    0 allocs/op  ⚡
BenchmarkPool_RealWorldScenario-10      240,111,070     5.02 ns/op    0 B/op    0 allocs/op  ⚡
```

### Key Performance Indicators

| Metric | Value | Description |
|--------|-------|-------------|
| **Fastest Concurrent** | **4.32 ns/op** | sync.Pool concurrent (4 items) |
| **Max Throughput** | **231M ops/sec** | Concurrent scenario |
| **Memory Allocation** | **0 B/op** | Zero heap allocation |
| **GC Pressure** | **0 allocs/op** | No GC pressure |

## ❓ FAQ

### Q: How to choose between Lotteries and LotteriesPool?

**A:** 
- Single-thread or independent instance per goroutine → `NewLotteries()`
- High concurrency with shared instance → `NewLotteriesPool()` ⭐ **Recommended**

### Q: What if probabilities don't sum to 1.0?

**A:** The system will automatically detect and return an error:
```go
lotteries, err := lottery.NewLotteries(prizes)
if err != nil {
    // err: cannot find suitable mul value, probabilities may be too precise or sum is not 1.0
}
```

### Q: Can I dynamically modify probabilities?

**A:** Not supported in current version. Please reinitialize the lottery if you need to change probabilities.

### Q: How to verify randomness?

**A:** Run a large number of tests and verify the distribution:
```go
results := make(map[string]int)
for i := 0; i < 100000; i++ {
    result := lotteries.Draw()
    results[result]++
}
// Check distribution in results
```

### Q: Where are the performance bottlenecks?

**A:** 
1. Prize count (O(n) complexity)
2. Using lock-free version in concurrent scenarios (data race)
3. Frequently creating new instances (should reuse)

## 🤝 Contributing

Issues and Pull Requests are welcome!

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

Thanks to all contributors and users!

---

**⭐ If this project helps you, please give it a Star!**

**🔗 Links**:
- [中文文档](README.md)

