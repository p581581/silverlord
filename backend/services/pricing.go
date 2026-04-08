package services

import (
	"math"

	"silverlord/backend/cache"
	"silverlord/backend/models"
)

// CalcPrice 依規格公式計算單件商品售價（台幣整數）
//
// 步驟：
//  1. 現貨價(USD/oz) + 修正值 → 基準現貨價
//  2. 基準現貨價 * USD/TWD    → 每 troy oz 台幣價
//  3. 每 troy oz 台幣價 * 重量 → 白銀本體台幣價
//  4. * (1 + 溢價%)           → 含溢價台幣價
//  5. + 工錢                  → 最終售價
func CalcPrice(p *models.Product, cfg *models.SystemConfig) int {
	snap := cache.Get()

	// 1. 基準現貨價（USD/oz）
	baseUsd := snap.XagUsd + cfg.ManualOffset

	// 2. 匯率：ManualExRate > 0 時強制使用固定匯率
	usdTwd := snap.UsdTwd
	if cfg.ManualExRate > 0 {
		usdTwd = cfg.ManualExRate
	}

	// 3. 每 troy oz 換算為台幣
	pricePerOzTwd := baseUsd * usdTwd

	// 4. 乘上商品重量，套用溢價
	silverValue := pricePerOzTwd * p.Weight * (1 + p.PremiumPercent/100)

	// 5. 加工錢，四捨五入
	return int(math.Round(silverValue + float64(p.LaborCost)))
}
