package handlers

import (
	"net/http"
	"strconv"
	"time"

	"silverlord/backend/cache"
	"silverlord/backend/config"
	"silverlord/backend/database"
	"silverlord/backend/middleware"
	"silverlord/backend/models"
	"silverlord/backend/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// POST /api/v1/admin/login
func AdminLogin(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req struct {
			Password string `json:"password"`
		}
		if err := c.Bind(&req); err != nil || req.Password == "" {
			return errBadRequest(c, ErrAuthPasswordRequired, "password required")
		}
		if req.Password != cfg.AdminPass {
			return errUnauthorized(c, ErrAuthWrongPassword, "wrong password")
		}

		claims := &middleware.JWTClaims{
			Role: "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			return errInternal(c, "token sign failed")
		}

		return c.JSON(http.StatusOK, echo.Map{"token": signed})
	}
}

// ─── Product CRUD ────────────────────────────────────────────────────────────

type adminProductResponse struct {
	models.Product
	Price int `json:"price"`
}

// GET /api/admin/products
func AdminGetProducts(c echo.Context) error {
	var products []models.Product
	database.DB.Find(&products)

	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)

	snap := cache.Get()
	result := make([]adminProductResponse, 0, len(products))
	for _, p := range products {
		price := 0
		if snap.XagUsd > 0 {
			price = services.CalcPrice(&p, &cfg)
		}
		result = append(result, adminProductResponse{Product: p, Price: price})
	}
	return c.JSON(http.StatusOK, result)
}

// POST /api/admin/products
func AdminCreateProduct(c echo.Context) error {
	var p models.Product
	if err := c.Bind(&p); err != nil {
		return errBadRequest(c, ErrInternal, err.Error())
	}
	if p.Name == "" {
		return errBadRequest(c, ErrProductNameRequired, "name is required")
	}
	if p.Weight <= 0 {
		return errBadRequest(c, ErrProductInvalidWeight, "weight must be greater than 0")
	}
	if _, err := database.DB.Insert(&p); err != nil {
		return errInternal(c, err.Error())
	}
	return c.JSON(http.StatusCreated, p)
}

// PUT /api/admin/products/:id
func AdminUpdateProduct(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var p models.Product
	has, _ := database.DB.ID(id).Get(&p)
	if !has {
		return errNotFound(c, ErrProductNotFound, "product not found")
	}
	if err := c.Bind(&p); err != nil {
		return errBadRequest(c, ErrInternal, err.Error())
	}
	if p.Name == "" {
		return errBadRequest(c, ErrProductNameRequired, "name is required")
	}
	if p.Weight <= 0 {
		return errBadRequest(c, ErrProductInvalidWeight, "weight must be greater than 0")
	}
	p.Id = id
	database.DB.ID(id).AllCols().Update(&p)
	return c.JSON(http.StatusOK, p)
}

// PATCH /api/admin/products/:id/toggle
func AdminToggleProduct(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var p models.Product
	has, _ := database.DB.ID(id).Get(&p)
	if !has {
		return errNotFound(c, ErrProductNotFound, "product not found")
	}
	p.IsActive = !p.IsActive
	database.DB.ID(id).Cols("is_active").Update(&p)
	return c.JSON(http.StatusOK, echo.Map{"id": id, "is_active": p.IsActive})
}

// DELETE /api/admin/products/:id
func AdminDeleteProduct(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	has, _ := database.DB.ID(id).Exist(new(models.Product))
	if !has {
		return errNotFound(c, ErrProductNotFound, "product not found")
	}
	database.DB.ID(id).Delete(new(models.Product))
	return c.NoContent(http.StatusNoContent)
}

// ─── Order Management ────────────────────────────────────────────────────────

// GET /api/admin/orders?q=&status=
func AdminGetOrders(c echo.Context) error {
	q := c.QueryParam("q")
	statusStr := c.QueryParam("status")

	session := database.DB.OrderBy("created_at DESC")

	if q != "" {
		session = session.Where("order_no LIKE ?", "%"+q+"%")
	}
	if statusStr != "" {
		session = session.Where("status = ?", statusStr)
	}

	var orders []models.Order
	session.Find(&orders)
	return c.JSON(http.StatusOK, orders)
}

// DELETE /api/admin/orders/:id
func AdminDeleteOrder(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var order models.Order
	has, _ := database.DB.ID(id).Get(&order)
	if !has {
		return errNotFound(c, ErrOrderNotFound, "order not found")
	}

	// 若尚未標記不成立，刪除前歸還庫存品庫存
	if order.Status != 3 {
		var items []models.OrderItem
		database.DB.Where("order_id = ?", id).Find(&items)
		for _, item := range items {
			var p models.Product
			if ok, _ := database.DB.ID(item.ProductId).Get(&p); ok && p.ProductType == "in_stock" {
				p.StockQty += item.Quantity
				database.DB.ID(p.Id).Cols("stock_qty").Update(&p)
			}
		}
	}

	database.DB.Where("order_id = ?", id).Delete(new(models.OrderItem))
	database.DB.ID(id).Delete(new(models.Order))
	return c.NoContent(http.StatusNoContent)
}

// GET /api/admin/orders/:id
func AdminGetOrderDetail(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var order models.Order
	has, _ := database.DB.ID(id).Get(&order)
	if !has {
		return errNotFound(c, ErrOrderNotFound, "order not found")
	}
	var items []models.OrderItem
	database.DB.Where("order_id = ?", id).Find(&items)
	return c.JSON(http.StatusOK, echo.Map{"order": order, "items": items})
}

// PATCH /api/admin/orders/:id/status
func AdminUpdateOrderStatus(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var order models.Order
	has, _ := database.DB.ID(id).Get(&order)
	if !has {
		return errNotFound(c, ErrOrderNotFound, "order not found")
	}
	var req struct {
		Status int `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return errBadRequest(c, ErrInternal, err.Error())
	}
	if req.Status < 0 || req.Status > 5 {
		return errBadRequest(c, ErrOrderInvalidStatus, "status must be 0–5")
	}

	// status 1（已通知）不可直接設為不成立，需先退回待匯款
	if req.Status == 3 && order.Status == 1 {
		return errBadRequest(c, ErrOrderInvalidStatus, "please reject payment report first before invalidating")
	}
	// 已完成（5）不可設為不成立
	if req.Status == 3 && order.Status == 5 {
		return errBadRequest(c, ErrOrderInvalidStatus, "completed orders cannot be invalidated")
	}
	// 已完成（5）只能從已確認/待出貨（2）設定
	if req.Status == 5 && order.Status != 2 {
		return errBadRequest(c, ErrOrderInvalidStatus, "order must be confirmed before marking as completed")
	}

	// 訂單不成立（→ 3）才歸還庫存品庫存
	if req.Status == 3 && order.Status != 3 {
		var items []models.OrderItem
		database.DB.Where("order_id = ?", id).Find(&items)
		for _, item := range items {
			var p models.Product
			if ok, _ := database.DB.ID(item.ProductId).Get(&p); ok && p.ProductType == "in_stock" {
				p.StockQty += item.Quantity
				database.DB.ID(p.Id).Cols("stock_qty").Update(&p)
			}
		}
	}

	order.Status = req.Status
	database.DB.ID(id).Cols("status").Update(&order)
	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}

// PATCH /api/admin/products/:id/stock
func AdminAdjustStock(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var p models.Product
	has, _ := database.DB.ID(id).Get(&p)
	if !has {
		return errNotFound(c, ErrProductNotFound, "product not found")
	}
	var req struct {
		StockQty int `json:"stock_qty"`
	}
	if err := c.Bind(&req); err != nil {
		return errBadRequest(c, ErrInternal, err.Error())
	}
	if req.StockQty < 0 {
		return errBadRequest(c, ErrInternal, "stock_qty cannot be negative")
	}
	p.StockQty = req.StockQty
	database.DB.ID(id).Cols("stock_qty").Update(&p)
	return c.JSON(http.StatusOK, echo.Map{"id": id, "stock_qty": p.StockQty})
}

// ─── System Config ───────────────────────────────────────────────────────────

// GET /api/v1/ship-fees（公開）
func GetShipFees(c echo.Context) error {
	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)
	return c.JSON(http.StatusOK, echo.Map{
		"mail":   cfg.ShipFeeMail,
		"711":    cfg.ShipFee711,
		"pickup": cfg.ShipFeePickup,
	})
}

// GET /api/v1/bank-info
func GetBankInfo(c echo.Context) error {
	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)
	return c.JSON(http.StatusOK, echo.Map{
		"bank_name":         cfg.BankName,
		"bank_account":      cfg.BankAccount,
		"bank_account_name": cfg.BankAccountName,
	})
}

// GET /api/admin/config
func AdminGetConfig(c echo.Context) error {
	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)
	return c.JSON(http.StatusOK, cfg)
}

// PUT /api/admin/config
func AdminUpdateConfig(c echo.Context) error {
	var cfg models.SystemConfig
	if err := c.Bind(&cfg); err != nil {
		return errBadRequest(c, ErrInternal, err.Error())
	}
	cfg.Id = 1
	database.DB.ID(1).AllCols().Update(&cfg)
	return c.JSON(http.StatusOK, cfg)
}

// POST /api/admin/config/test-tg
func AdminTestTelegram(c echo.Context) error {
	var cfg models.SystemConfig
	database.DB.ID(1).Get(&cfg)
	if cfg.TgBotToken == "" || cfg.TgChatIds == "" {
		return errBadRequest(c, ErrConfigTgNotConfigured, "telegram not configured")
	}
	services.SendTelegramMulti(cfg.TgBotToken, cfg.TgChatIds,
		"✅ <b>白銀平台</b> Telegram 通知測試成功！")
	return c.JSON(http.StatusOK, echo.Map{"message": "sent"})
}

// ─── Upload ──────────────────────────────────────────────────────────────────

// POST /api/admin/upload
func AdminUpload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return errBadRequest(c, ErrUploadFileRequired, "file required")
	}

	src, err := file.Open()
	if err != nil {
		return errInternal(c, err.Error())
	}
	defer src.Close()

	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	dst, err := createUploadFile("uploads/" + filename)
	if err != nil {
		return errInternal(c, err.Error())
	}
	defer dst.Close()

	if err := copyFile(src, dst); err != nil {
		return fail(c, http.StatusInternalServerError, ErrUploadWriteFailed, "file write failed")
	}

	return c.JSON(http.StatusOK, echo.Map{"url": "/uploads/" + filename})
}
