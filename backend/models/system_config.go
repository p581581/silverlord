package models

type SystemConfig struct {
	Id           int     `xorm:"pk" json:"id"`
	ManualOffset float64 `xorm:"default 0.0" json:"manual_offset"`
	ManualExRate float64 `xorm:"default 0.0" json:"manual_ex_rate"`
	TgBotToken   string  `json:"tg_bot_token"`
	TgChatIds    string  `json:"tg_chat_ids"` // 逗號分隔，支援多組 Chat ID
	ShipFeeMail    int     `xorm:"default 0" json:"ship_fee_mail"`    // 郵寄運費（TWD）
	ShipFee711     int     `xorm:"default 0" json:"ship_fee_711"`     // 7-11 店到店運費（TWD）
	ShipFeePickup  int     `xorm:"default 0" json:"ship_fee_pickup"`  // 自取運費（TWD）
}
