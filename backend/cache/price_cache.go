package cache

import (
	"sync"
	"time"
)

// PriceSnapshot 是每次 Cron 抓到的行情快照
type PriceSnapshot struct {
	XagUsd    float64   `json:"xag_usd"`
	UsdTwd    float64   `json:"usd_twd"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	mu       sync.RWMutex
	snapshot PriceSnapshot
)

// Set 由 Cron Job 呼叫，更新全域 Cache
func Set(xagUsd, usdTwd float64) {
	mu.Lock()
	defer mu.Unlock()
	snapshot = PriceSnapshot{
		XagUsd:    xagUsd,
		UsdTwd:    usdTwd,
		UpdatedAt: time.Now(),
	}
}

// Get 由 API Handler 呼叫，讀取最新快照
func Get() PriceSnapshot {
	mu.RLock()
	defer mu.RUnlock()
	return snapshot
}
