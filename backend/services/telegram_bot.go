package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"silverlord/backend/cache"
	"silverlord/backend/database"
	"silverlord/backend/models"
)

// ── Telegram Update 型別 ─────────────────────────────────────────────────────

type tgUpdate struct {
	UpdateId int64 `json:"update_id"`
	Message  *struct {
		Chat struct {
			Id int64 `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

var botOffset int64

// StartBot 啟動長輪詢 goroutine
func StartBot() {
	go func() {
		log.Println("[bot] Telegram bot polling started")
		for {
			var cfg models.SystemConfig
			if ok, _ := database.DB.ID(1).Get(&cfg); !ok || cfg.TgBotToken == "" {
				time.Sleep(30 * time.Second)
				continue
			}
			processBotUpdates(cfg.TgBotToken, cfg.TgChatIds)
			time.Sleep(2 * time.Second)
		}
	}()
}

func processBotUpdates(token, allowedIds string) {
	apiURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=10&allowed_updates=[\"message\"]",
		token, botOffset+1,
	)
	resp, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Ok     bool       `json:"ok"`
		Result []tgUpdate `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Ok {
		return
	}

	for _, upd := range result.Result {
		if upd.UpdateId > botOffset {
			botOffset = upd.UpdateId
		}
		if upd.Message == nil || upd.Message.Text == "" {
			continue
		}
		chatId := upd.Message.Chat.Id
		// 只允許設定中的 Chat ID 查詢
		if !isAllowedChat(chatId, allowedIds) {
			continue
		}
		handleBotCommand(token, chatId, upd.Message.Text)
	}
}

func isAllowedChat(chatId int64, allowedIds string) bool {
	for _, raw := range strings.Split(allowedIds, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err == nil && id == chatId {
			return true
		}
	}
	return false
}

// ── 指令處理 ─────────────────────────────────────────────────────────────────

func handleBotCommand(token string, chatId int64, text string) {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return
	}

	// 處理 /command@botname 格式
	cmd := strings.ToLower(parts[0])
	if idx := strings.Index(cmd, "@"); idx != -1 {
		cmd = cmd[:idx]
	}

	switch cmd {
	case "/start":
		botSend(token, chatId, botStartMsg())
	case "/help":
		botSend(token, chatId, botHelpMsg())
	case "/price":
		botSend(token, chatId, botPriceMsg())
	case "/query":
		if len(parts) < 2 {
			botSend(token, chatId, "❌ 請提供訂單編號\n範例：/query SN202504070001")
			return
		}
		botSend(token, chatId, botQueryOrder(parts[1]))
	case "/orders":
		if len(parts) < 2 {
			botSend(token, chatId, "❌ 請提供狀態\n範例：/orders 1\n或：/orders 1 2026-04-06")
			return
		}
		status, err := strconv.Atoi(parts[1])
		if err != nil || status < 0 || status > 5 {
			botSend(token, chatId, "❌ 狀態無效，請輸入 0~5")
			return
		}
		dateStr := time.Now().Format("2006-01-02")
		if len(parts) >= 3 {
			dateStr = parts[2]
		}
		botSend(token, chatId, botQueryByStatus(status, dateStr))
	}
}

// ── 訊息內容 ─────────────────────────────────────────────────────────────────

func botStartMsg() string {
	return `🥈 <b>白銀即時報價</b> — 訂單查詢機器人

輸入指令查詢訂單：
/price — 即時銀價
/query SN202504070001 — 查詢單筆訂單
/orders 1 — 今日已通知訂單

/help 查看所有指令`
}

func botHelpMsg() string {
	return `🥈 <b>白銀即時報價 指令列表</b>

<b>銀價查詢</b>
/price — 即時銀價與匯率

<b>單筆查詢</b>
/query &lt;訂單編號&gt;

<b>狀態查詢（預設今日）</b>
/orders 0 — ⏳ 待匯款
/orders 1 — 💰 已通知
/orders 2 — 📦 已確認/待出貨
/orders 3 — ❌ 不成立
/orders 4 — ⚠️ 異常
/orders 5 — ✅ 已完成

<b>指定日期</b>
/orders &lt;狀態&gt; YYYY-MM-DD`
}

func botStatusLabel(status int) string {
	switch status {
	case 0:
		return "⏳ 待匯款"
	case 1:
		return "💰 已通知"
	case 2:
		return "📦 已確認/待出貨"
	case 3:
		return "❌ 不成立"
	case 4:
		return "⚠️ 異常"
	case 5:
		return "✅ 已完成"
	default:
		return "❓"
	}
}

func botPriceMsg() string {
	snap := cache.Get()
	if snap.XagUsd == 0 {
		return "❌ 銀價資料尚未就緒，請稍後再試"
	}

	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)

	baseUsd := snap.XagUsd + cfg.ManualOffset
	usdTwd := snap.UsdTwd
	if cfg.ManualExRate > 0 {
		usdTwd = cfg.ManualExRate
	}
	pricePerOzTwd := int(baseUsd * usdTwd)

	updStr := snap.UpdatedAt.Format("2006-01-02 15:04")

	var sb strings.Builder
	fmt.Fprintf(&sb, "🥈 <b>即時銀價</b>\n\n")
	fmt.Fprintf(&sb, "XAG/USD：<b>$ %.2f</b>", snap.XagUsd)
	if cfg.ManualOffset != 0 {
		fmt.Fprintf(&sb, "（修正 %+.2f）", cfg.ManualOffset)
	}
	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "USD/TWD：<b>%.4f</b>", usdTwd)
	if cfg.ManualExRate > 0 {
		fmt.Fprintf(&sb, "（固定匯率）")
	}
	fmt.Fprintf(&sb, "\n\n")
	fmt.Fprintf(&sb, "銀價（每 oz）：<b>NT$ %s</b>\n\n", botFmtPrice(pricePerOzTwd))
	fmt.Fprintf(&sb, "<i>最後更新：%s</i>", updStr)
	return sb.String()
}

func botQueryOrder(orderNo string) string {
	var order models.Order
	has, _ := database.DB.Where("order_no = ?", orderNo).Get(&order)
	if !has {
		return fmt.Sprintf(
			"❌ 找不到訂單 <code>%s</code>\n\n請確認編號是否正確，格式範例：\nSN202504070001",
			orderNo,
		)
	}

	var items []models.OrderItem
	database.DB.Where("order_id = ?", order.Id).Find(&items)

	var sb strings.Builder
	fmt.Fprintf(&sb, "📋 <b>訂單查詢結果</b>\n\n")
	fmt.Fprintf(&sb, "訂單編號：<code>%s</code>\n", order.OrderNo)
	fmt.Fprintf(&sb, "狀態：%s\n", botStatusLabel(order.Status))
	fmt.Fprintf(&sb, "收件人：%s\n", order.Name)
	fmt.Fprintf(&sb, "電話：%s\n", order.Phone)
	fmt.Fprintf(&sb, "Email：%s\n", order.Email)

	if order.ShipMethod == "711" {
		fmt.Fprintf(&sb, "出貨：7-11（%s）\n", order.ShipStore)
	} else {
		fmt.Fprintf(&sb, "出貨：郵寄\n  %s\n", order.ShipAddress)
	}
	fmt.Fprintf(&sb, "下單時間：%s\n", order.CreatedAt.Format("2006-01-02 15:04"))

	fmt.Fprintf(&sb, "\n<b>商品明細</b>\n")
	for _, item := range items {
		fmt.Fprintf(&sb, "  • %s × %d　NT$ %s\n",
			item.ProductName, item.Quantity,
			botFmtPrice(item.LockedUnitPrice*item.Quantity))
	}

	subtotal := order.TotalLockedPrice - order.ShipFee
	fmt.Fprintf(&sb, "\n商品小計：NT$ %s\n", botFmtPrice(subtotal))
	if order.ShipFee > 0 {
		shipLabel := "郵寄"
		if order.ShipMethod == "711" {
			shipLabel = "7-11"
		}
		fmt.Fprintf(&sb, "運費（%s）：NT$ %s\n", shipLabel, botFmtPrice(order.ShipFee))
	} else {
		fmt.Fprintf(&sb, "運費：免運\n")
	}
	fmt.Fprintf(&sb, "━━━━━━━━━━\n")
	fmt.Fprintf(&sb, "訂單總計：<b>NT$ %s</b>", botFmtPrice(order.TotalLockedPrice))

	return sb.String()
}

func botQueryByStatus(status int, dateStr string) string {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "❌ 日期格式錯誤，請使用 YYYY-MM-DD"
	}

	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	var orders []models.Order
	database.DB.
		Where("status = ? AND created_at >= ? AND created_at < ?", status, start, end).
		OrderBy("created_at DESC").
		Find(&orders)

	if len(orders) == 0 {
		return fmt.Sprintf("📭 %s 無「%s」訂單", dateStr, botStatusLabel(status))
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%s — %s\n\n", botStatusLabel(status), dateStr)

	for i, o := range orders {
		shipInfo := "郵寄"
		if o.ShipMethod == "711" {
			shipInfo = "7-11"
		}
		fmt.Fprintf(&sb, "%d. <code>%s</code>\n   %s｜NT$ %s｜%s\n\n",
			i+1, o.OrderNo, o.Name, botFmtPrice(o.TotalLockedPrice), shipInfo)
	}
	fmt.Fprintf(&sb, "共 %d 筆", len(orders))
	return sb.String()
}

// botFmtPrice 加千分位符號
func botFmtPrice(n int) string {
	s := strconv.Itoa(n)
	var out []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out)
}

// botSend 發送 HTML 訊息
func botSend(token string, chatId int64, text string) {
	if err := SendTelegram(token, chatId, text); err != nil {
		log.Printf("[bot] send to %d failed: %v", chatId, err)
	}
}
