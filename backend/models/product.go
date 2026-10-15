package models

type Product struct {
	Id              int64   `xorm:"pk autoincr" json:"id"`
	Name            string  `xorm:"notnull" json:"name"`
	Description     string  `xorm:"text" json:"description"`
	ImageUrl        string  `json:"image_url"`
	ImageUrl2       string  `json:"image_url_2"`
	ImageUrl3       string  `json:"image_url_3"`
	Weight          float64 `xorm:"notnull" json:"weight"`       // 重量（troy oz）
	LaborCost       int     `xorm:"default 0" json:"labor_cost"` // 台幣
	PremiumPercent  float64 `xorm:"default 5.0" json:"premium_percent"`
	IsActive        bool    `xorm:"index default true" json:"is_active"`
	ProductType     string  `xorm:"default 'in_stock'" json:"product_type"` // "in_stock" | "pre_order"
	StockQty        int     `xorm:"default 0" json:"stock_qty"`
	ExpectedArrival string  `xorm:"default ''" json:"expected_arrival"` // YYYY-MM-DD
	SortOrder       int64   `xorm:"index default 0" json:"sort_order"`
}
