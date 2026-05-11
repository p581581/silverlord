package models

import "time"

type Order struct {
	Id              int64     `xorm:"pk autoincr" json:"id"`
	OrderNo         string    `xorm:"unique notnull" json:"order_no"`       // SN + YYYYMMDD + 4位流水
	TotalLockedPrice int      `xorm:"notnull" json:"total_locked_price"`    // 鎖定台幣總額
	Email           string    `json:"email"`                                 // 買家 Email
	Name            string    `json:"name"`                                  // 收件人姓名
	Phone           string    `json:"phone"`                                 // 收件人電話
	ShipMethod      string    `json:"ship_method"`                           // "mail" | "711"
	ShipAddress     string    `json:"ship_address"`                          // 郵寄地址
	ShipStore       string    `json:"ship_store"`                            // 7-11 店名
	ShipFee         int       `xorm:"default 0" json:"ship_fee"`             // 實際收取運費（TWD）
	RemitInfo       string    `json:"remit_info"`                            // 匯款末五碼
	Status          int       `xorm:"default 0" json:"status"`              // 0:待匯款 1:已通知 2:已確認/待出貨
	CreatedAt       time.Time `xorm:"created" json:"created_at"`
}

type OrderItem struct {
	Id              int64   `xorm:"pk autoincr" json:"id"`
	OrderId         int64   `xorm:"notnull index" json:"order_id"`
	ProductId       int64   `xorm:"notnull" json:"product_id"`
	ProductName     string  `json:"product_name"`                          // 快照
	ImageUrl        string  `json:"image_url"`                             // 快照
	Weight          float64 `json:"weight"`                                // 快照（troy oz）
	Quantity        int     `xorm:"notnull default 1" json:"quantity"`
	LockedUnitPrice int     `json:"locked_unit_price"`                     // 下單當下單價（台幣）
}
