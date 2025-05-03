package store_test

import (
	"fmt"
	"sync"
	"testing"

	"rate-limiter/internal/limiter/store"
)

func BenchmarkRateLimiter_SingleKey(b *testing.B) {
	mockStore := store.NewMockStore(false)
	key := "benchmark-key"
	limit := 1000
	blockSeconds := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockStore.Increment(key, limit, blockSeconds)
	}
}

func BenchmarkRateLimiter_MultipleKeys(b *testing.B) {
	mockStore := store.NewMockStore(false)
	limit := 1000
	blockSeconds := 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("benchmark-key-%d", i%100) // Use 100 different keys
		mockStore.Increment(key, limit, blockSeconds)
	}
}

func BenchmarkRateLimiter_ConcurrentRequests(b *testing.B) {
	mockStore := store.NewMockStore(false)
	key := "concurrent-benchmark-key"
	limit := 10000
	blockSeconds := 1

	b.ResetTimer()

	// Run b.N operations spread across multiple goroutines
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mockStore.Increment(key, limit, blockSeconds)
		}
	})
}

func BenchmarkRateLimiter_HighConcurrency(b *testing.B) {
	concurrencyLevels := []int{10, 100, 1000}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			mockStore := store.NewMockStore(false)
			key := "high-concurrency-key"
			limit := 10000
			blockSeconds := 1

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var wg sync.WaitGroup
				for j := 0; j < concurrency; j++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						mockStore.Increment(key, limit, blockSeconds)
					}()
				}
				wg.Wait()
			}
		})
	}
}

func BenchmarkRateLimiter_MixedOperations(b *testing.B) {
	mockStore := store.NewMockStore(false)
	limit := 1000
	blockSeconds := 1

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			key := fmt.Sprintf("mixed-ops-key-%d", counter%10)
			counter++

			// Mix different operations
			switch counter % 3 {
			case 0:
				mockStore.Increment(key, limit, blockSeconds)
			case 1:
				mockStore.IsBlocked(key)
			case 2:
				mockStore.GetRequestCount(key)
			}
		}
	})
}

func BenchmarkRateLimiter_BurstTraffic(b *testing.B) {
	mockStore := store.NewMockStore(false)
	key := "burst-traffic-key"
	limit := 5000
	blockSeconds := 1

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate burst traffic with 100 concurrent requests
		var wg sync.WaitGroup
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mockStore.Increment(key, limit, blockSeconds)
			}()
		}
		wg.Wait()
	}
}
