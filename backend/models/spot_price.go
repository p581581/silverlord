package models

import "time"

type SpotPrice struct {
	Id        int64     `xorm:"pk autoincr" json:"id"`
	XagUsd    float64   `json:"xag_usd"` // 白銀現貨美金價（每盎司）
	UsdTwd    float64   `json:"usd_twd"` // 台幣匯率
	CreatedAt time.Time `xorm:"created index" json:"created_at"`
}
