package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"silverlord/backend/cache"
	"silverlord/backend/database"
	"silverlord/backend/models"

	"github.com/robfig/cron/v3"
)

// StartPriceCron 每 60 秒抓一次外部報價並更新 Cache
func StartPriceCron() {
	c := cron.New()
	c.AddFunc("@every 60s", fetchAndCachePrice)
	c.Start()

	// 啟動時立即執行一次，避免 Cache 為空
	go fetchAndCachePrice()

	log.Println("Price cron started (every 60s).")
}

func fetchAndCachePrice() {
	xagUsd, err := fetchYahooPrice("SI=F")
	if err != nil {
		log.Printf("[cron] fetch XAG/USD error: %v", err)
		return
	}

	usdTwd, err := fetchYahooPrice("USDTWD=X")
	if err != nil {
		log.Printf("[cron] fetch USD/TWD error: %v", err)
		return
	}

	// 更新記憶體 Cache
	cache.Set(xagUsd, usdTwd)

	// 寫入 DB（歷史紀錄，CreatedAt 由 XORM created tag 自動填入）
	if _, err := database.DB.Insert(&models.SpotPrice{
		XagUsd: xagUsd,
		UsdTwd: usdTwd,
	}); err != nil {
		log.Printf("[cron] db insert error: %v", err)
	}

	log.Printf("[cron] updated — XAG/USD=%.4f, USD/TWD=%.4f", xagUsd, usdTwd)
}

// yahooResponse 只取需要的欄位
type yahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}

// fetchYahooPrice 從 Yahoo Finance 非官方 API 抓取即時報價
// symbol 範例："XAGUSD=X"、"USDTWD=X"、"SI=F"
func fetchYahooPrice(symbol string) (float64, error) {
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1m&range=1d",
		symbol,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	// Yahoo 需要 User-Agent，否則可能回傳 401
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("yahoo finance status %d for %s", resp.StatusCode, symbol)
	}

	var result yahooResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Chart.Error != nil {
		return 0, fmt.Errorf("yahoo finance error for %s: %v", symbol, result.Chart.Error)
	}
	if len(result.Chart.Result) == 0 {
		return 0, fmt.Errorf("yahoo finance empty result for %s", symbol)
	}

	price := result.Chart.Result[0].Meta.RegularMarketPrice
	if price == 0 {
		return 0, fmt.Errorf("yahoo finance returned 0 for %s", symbol)
	}
	return price, nil
}
