package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"silverlord/backend/cache"
	"silverlord/backend/config"
	"silverlord/backend/database"
	"silverlord/backend/models"
	"silverlord/backend/services"

	"github.com/labstack/echo/v4"
)

type cartItem struct {
	ProductId int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type checkoutRequest struct {
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Phone       string     `json:"phone"`
	Items       []cartItem `json:"items"`
	ShipMethod  string     `json:"ship_method"`
	ShipAddress string     `json:"ship_address"`
	ShipStore   string     `json:"ship_store"`
}

// POST /api/v1/checkout
func Checkout(appCfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req checkoutRequest
		if err := c.Bind(&req); err != nil || len(req.Items) == 0 {
			return errBadRequest(c, ErrOrderItemsRequired, "items required")
		}
		if req.Email == "" {
			return errBadRequest(c, ErrOrderItemsRequired, "email required")
		}
		if req.Name == "" || req.Phone == "" {
			return errBadRequest(c, ErrOrderItemsRequired, "name and phone required")
		}
		if req.ShipMethod != "mail" && req.ShipMethod != "711" && req.ShipMethod != "pickup" {
			return errBadRequest(c, ErrOrderItemsRequired, "ship_method must be mail, 711, or pickup")
		}
		if req.ShipMethod == "mail" && req.ShipAddress == "" {
			return errBadRequest(c, ErrOrderItemsRequired, "ship_address required for mail")
		}
		if req.ShipMethod == "711" && req.ShipStore == "" {
			return errBadRequest(c, ErrOrderItemsRequired, "ship_store required for 711")
		}
		for _, item := range req.Items {
			if item.Quantity <= 0 {
				return errBadRequest(c, ErrOrderInvalidQty, "quantity must be greater than 0")
			}
		}

		snap := cache.Get()
		if snap.XagUsd == 0 {
			return errServiceUnavailable(c, ErrPriceNotReady, "price not ready, please retry in a moment")
		}

		var dbCfg models.SystemConfig
		database.DB.ID(1).Get(&dbCfg)

		sess := database.DB.NewSession()
		defer sess.Close()

		if err := sess.Begin(); err != nil {
			return errInternal(c, "tx begin failed")
		}

		total := 0
		var orderItems []models.OrderItem

		for _, item := range req.Items {
			var p models.Product
			has, _ := sess.ID(item.ProductId).Get(&p)
			if !has {
				sess.Rollback()
				return errNotFound(c, ErrProductNotFound,
					fmt.Sprintf("product %d not found", item.ProductId))
			}
			if !p.IsActive {
				sess.Rollback()
				return errBadRequest(c, ErrProductInactive,
					fmt.Sprintf("product %d is no longer available", item.ProductId))
			}
			// 庫存品需檢查庫存
			if p.ProductType == "in_stock" && p.StockQty < item.Quantity {
				sess.Rollback()
				return errBadRequest(c, ErrProductInactive,
					fmt.Sprintf("product %d has insufficient stock", item.ProductId))
			}

			unitPrice := services.CalcPrice(&p, &dbCfg)
			total += unitPrice * item.Quantity

			orderItems = append(orderItems, models.OrderItem{
				ProductId:       p.Id,
				ProductName:     p.Name,
				ImageUrl:        p.ImageUrl,
				Weight:          p.Weight,
				Quantity:        item.Quantity,
				LockedUnitPrice: unitPrice,
			})

			// 庫存品扣庫存
			if p.ProductType == "in_stock" {
				newQty := p.StockQty - item.Quantity
				if _, err := sess.ID(p.Id).Cols("stock_qty").Update(&models.Product{StockQty: newQty}); err != nil {
					sess.Rollback()
					return errInternal(c, "stock update failed")
				}
			}
		}

		// 運費加進總額
		shipFee := 0
		if req.ShipMethod == "mail" {
			shipFee = dbCfg.ShipFeeMail
		} else if req.ShipMethod == "711" {
			shipFee = dbCfg.ShipFee711
		} else if req.ShipMethod == "pickup" {
			shipFee = dbCfg.ShipFeePickup
		}
		total += shipFee

		order := &models.Order{
			OrderNo:          generateOrderNo(),
			TotalLockedPrice: total,
			Email:            req.Email,
			Name:             req.Name,
			Phone:            req.Phone,
			ShipMethod:       req.ShipMethod,
			ShipAddress:      req.ShipAddress,
			ShipStore:        req.ShipStore,
			ShipFee:          shipFee,
			Status:           0,
		}
		if _, err := sess.Insert(order); err != nil {
			sess.Rollback()
			return errInternal(c, "create order failed")
		}

		for i := range orderItems {
			orderItems[i].OrderId = order.Id
		}
		if _, err := sess.Insert(orderItems); err != nil {
			sess.Rollback()
			return errInternal(c, "create order items failed")
		}

		if err := sess.Commit(); err != nil {
			return errInternal(c, "tx commit failed")
		}

		// Email 確認信（非同步，使用 appCfg 的 SMTP 設定）
		if req.Email != "" {
			go func() {
				if err := services.SendOrderConfirmEmail(appCfg, req.Email, services.OrderEmailData{
					OrderNo:          order.OrderNo,
					TotalLockedPrice: order.TotalLockedPrice,
					Items:            orderItems,
				}); err != nil {
					log.Printf("[email] send failed: %v", err)
				}
			}()
		}

		// Telegram 新訂單通知（非同步）
		if dbCfg.TgBotToken != "" && dbCfg.TgChatIds != "" {
			var sb strings.Builder
			fmt.Fprintf(&sb, "📋 <b>新訂單通知</b>\n訂單編號：<code>%s</code>\n", order.OrderNo)
			for _, item := range orderItems {
				fmt.Fprintf(&sb, "  • %s × %d = <b>%d</b> 元\n",
					item.ProductName, item.Quantity, item.LockedUnitPrice*item.Quantity)
			}
			fmt.Fprintf(&sb, "合計：<b>%d</b> 元", order.TotalLockedPrice)
			go services.SendTelegramMulti(dbCfg.TgBotToken, dbCfg.TgChatIds, sb.String())
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"order_no":           order.OrderNo,
			"total_locked_price": order.TotalLockedPrice,
		})
	}
}

// POST /api/v1/orders/:order_no/resend-email
func ResendEmail(appCfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		orderNo := c.Param("order_no")
		var order models.Order
		has, _ := database.DB.Where("order_no = ?", orderNo).Get(&order)
		if !has {
			return errNotFound(c, ErrOrderNotFound, "order not found")
		}
		if order.Email == "" {
			return errBadRequest(c, ErrOrderItemsRequired, "order has no email")
		}
		var items []models.OrderItem
		database.DB.Where("order_id = ?", order.Id).Find(&items)
		go func() {
			if err := services.SendOrderConfirmEmail(appCfg, order.Email, services.OrderEmailData{
				OrderNo:          order.OrderNo,
				TotalLockedPrice: order.TotalLockedPrice,
				Items:            items,
			}); err != nil {
				log.Printf("[email] resend failed: %v", err)
			}
		}()
		return c.JSON(http.StatusOK, echo.Map{"message": "sent"})
	}
}

// GET /api/v1/orders/:order_no
func GetOrderStatus(c echo.Context) error {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		return errBadRequest(c, ErrOrderNotFound, "order_no required")
	}

	var order models.Order
	has, _ := database.DB.Where("order_no = ?", orderNo).Get(&order)
	if !has {
		return errNotFound(c, ErrOrderNotFound, "order not found")
	}

	var items []models.OrderItem
	database.DB.Where("order_id = ?", order.Id).Find(&items)

	// 舊訂單無 image_url 快照時，補抓當前商品圖
	for i := range items {
		if items[i].ImageUrl == "" {
			var p models.Product
			if ok, _ := database.DB.ID(items[i].ProductId).Get(&p); ok {
				items[i].ImageUrl = p.ImageUrl
			}
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"order_no":           order.OrderNo,
		"total_locked_price": order.TotalLockedPrice,
		"status":             order.Status,
		"created_at":         order.CreatedAt,
		"email":              order.Email,
		"ship_method":        order.ShipMethod,
		"ship_fee":           order.ShipFee,
		"items":              items,
	})
}

// POST /api/v1/report-payment
func ReportPayment(c echo.Context) error {
	var req struct {
		OrderNo   string `json:"order_no"`
		RemitInfo string `json:"remit_info"`
	}
	if err := c.Bind(&req); err != nil || req.OrderNo == "" || req.RemitInfo == "" {
		return errBadRequest(c, ErrOrderItemsRequired, "order_no and remit_info required")
	}

	var order models.Order
	has, _ := database.DB.Where("order_no = ?", req.OrderNo).Get(&order)
	if !has {
		return errNotFound(c, ErrOrderNotFound, "order not found")
	}
	if order.Status != 0 {
		return errConflict(c, ErrOrderAlreadyReported, "order already reported")
	}

	order.RemitInfo = req.RemitInfo
	order.Status = 1
	database.DB.ID(order.Id).Cols("remit_info", "status").Update(&order)

	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)
	if cfg.TgBotToken != "" && cfg.TgChatIds != "" {
		msg := fmt.Sprintf(
			"💰 <b>新匯款通知</b>\n訂單編號：<code>%s</code>\n匯款末五碼：<code>%s</code>\n金額：<b>%d</b> 元",
			order.OrderNo, order.RemitInfo, order.TotalLockedPrice,
		)
		go services.SendTelegramMulti(cfg.TgBotToken, cfg.TgChatIds, msg)
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "reported"})
}

// generateOrderNo 產生 SN + YYYYMMDD + 4位流水號
func generateOrderNo() string {
	today := time.Now().Format("20060102")
	prefix := "SN" + today + "%"
	count, _ := database.DB.Where("order_no LIKE ?", prefix).Count(new(models.Order))
	return fmt.Sprintf("SN%s%04d", today, count+1)
}
