package state

import (
	"context"
	"share-profile-allocator/internal/grpc"
	"share-profile-allocator/internal/utils"
	"sync"
	"time"
)

const (
	shareDataCacheTimeout    = 3 * time.Hour
	requestTickerDataTimeout = 20 * time.Second
)

var createCache sync.Once
var cache ShareDataCache

func GetShareDataCache() *ShareDataCache {
	createCache.Do(func() {
		cache = ShareDataCache{
			data: make(map[string]shareEntry),
		}
		cache.cleanupCacheJob()
	})

	return &cache
}

type shareEntry struct {
	shareData grpc.WrappedShareData
	acquired  time.Time
}

type ShareDataCache struct {
	data map[string]shareEntry
	mu   sync.RWMutex
}

// GetShareData returns the data corresponding to the given ticker. This data may have been cached, or it may be fresh.
func (cache *ShareDataCache) GetShareData(ticker string) (*grpc.WrappedShareData, error) {
	cache.mu.RLock()
	data, ok := cache.data[ticker]
	cache.mu.RUnlock()
	if ok {
		utils.Log("55d88c8c").Info("Request for share data, found in cache", "ticker", ticker, "acquired", data.acquired)
		return &data.shareData, nil
	}

	utils.Log("a8be4b7a").Info("Request for share data, cache miss", "ticker", ticker)

	ctx, cancel := context.WithTimeout(context.Background(), requestTickerDataTimeout)
	defer cancel()

	shareData, err := grpc.RequestDataForTicker(ctx, ticker)
	if err != nil {
		utils.Log("3778cca3").Error("Unable to get share data, GRPC request failed")
		return &grpc.ZeroShareData, err
	}

	cache.mu.Lock()
	cache.data[ticker] = shareEntry{
		shareData: *shareData,
		acquired:  time.Now(),
	}
	cache.mu.Unlock()

	return shareData, nil
}

type ShareDataResult struct {
	Data *grpc.WrappedShareData
	Err  error
}

// BatchGetShareData gets the data for the given tickers asynchronously. This is good for if you want to
// get the data for a lot of tickers, so you don't want to wait for each request to complete
// before requesting another.
func (cache *ShareDataCache) BatchGetShareData(tickers ...string) []ShareDataResult {
	res := make([]ShareDataResult, len(tickers))
	var wg sync.WaitGroup

	for i, t := range tickers {
		wg.Add(1)

		go func(index int, ticker string) {
			data, err := cache.GetShareData(ticker)
			res[index] = ShareDataResult{
				Data: data,
				Err:  err,
			}

			wg.Done()
		}(i, t)
	}

	wg.Wait()

	return res
}

// cleanupCacheJob launches a goroutine to delete any cache entries that have expired
func (cache *ShareDataCache) cleanupCacheJob() {
	cleanUpFreq := shareDataCacheTimeout / 2
	utils.Log("6112612a").Info("Starting share data cache clean up job", "Interval", cleanUpFreq)

	go func() {
		ticker := time.NewTicker(cleanUpFreq)
		for range ticker.C {
			cache.mu.Lock()

			for id, entry := range cache.data {
				if time.Since(entry.acquired) > shareDataCacheTimeout {
					utils.Log("8d37c0f4").Info("Deleting expired cached share data", "id", id, "acquired", entry.acquired)
					delete(cache.data, id)
				}
			}

			cache.mu.Unlock()
		}
	}()
}
