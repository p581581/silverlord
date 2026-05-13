package main

import (
	"os"

	"silverlord/backend/config"
	"silverlord/backend/database"
	"silverlord/backend/handlers"
	mw "silverlord/backend/middleware"
	"silverlord/backend/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// 確保 uploads 目錄存在
	os.MkdirAll("uploads", 0755)

	// 初始化 DB
	database.Init(cfg)

	// 啟動報價 Cron
	services.StartPriceCron()

	// 啟動 Telegram Bot 長輪詢
	services.StartBot()

	// Echo
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// 靜態資源（商品圖片）
	e.Static("/uploads", "uploads")

	// 前端 SPA（prod embed，dev no-op；用 RouteNotFound 不影響 API）
	setupFrontend(e)

	// ── Public API ──────────────────────────────────────────────
	v1 := e.Group("/api/v1")
	v1.POST("/admin/login", handlers.AdminLogin(cfg))
	v1.GET("/products", handlers.GetProducts)
	v1.GET("/products/:id", handlers.GetProduct)
	v1.GET("/spot-price", handlers.GetSpotPrice)
	v1.GET("/ship-fees", handlers.GetShipFees)
	v1.GET("/bank-info", handlers.GetBankInfo)
	v1.GET("/site-config", handlers.GetSiteConfig)
	v1.POST("/checkout", handlers.Checkout(cfg))
	v1.POST("/report-payment", handlers.ReportPayment)
	v1.POST("/orders/:order_no/resend-email", handlers.ResendEmail(cfg))
	v1.GET("/orders/:order_no", handlers.GetOrderStatus)

	// ── Admin API（需 JWT）──────────────────────────────────────
	admin := e.Group("/api/admin", mw.RequireAdmin(cfg.JWTSecret))

	// 商品
	admin.GET("/products", handlers.AdminGetProducts)
	admin.POST("/products", handlers.AdminCreateProduct)
	admin.PUT("/products/:id", handlers.AdminUpdateProduct)
	admin.PATCH("/products/:id/toggle", handlers.AdminToggleProduct)
	admin.PATCH("/products/:id/stock", handlers.AdminAdjustStock)
	admin.DELETE("/products/:id", handlers.AdminDeleteProduct)

	// 訂單
	admin.GET("/orders", handlers.AdminGetOrders)
	admin.GET("/orders/:id", handlers.AdminGetOrderDetail)
	admin.PATCH("/orders/:id/status", handlers.AdminUpdateOrderStatus)
	admin.DELETE("/orders/:id", handlers.AdminDeleteOrder)

	// 系統設定
	admin.GET("/config", handlers.AdminGetConfig)
	admin.PUT("/config", handlers.AdminUpdateConfig)
	admin.POST("/config/test-tg", handlers.AdminTestTelegram)

	// 上傳
	admin.POST("/upload", handlers.AdminUpload)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
