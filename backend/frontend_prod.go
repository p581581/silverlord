//go:build prod

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func init() {
	// 強制註冊 Vite/Vue 常用類型
	_ = mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
	_ = mime.AddExtensionType(".mjs", "application/javascript; charset=utf-8")
	_ = mime.AddExtensionType(".css", "text/css; charset=utf-8")
	_ = mime.AddExtensionType(".woff2", "font/woff2")
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
	_ = mime.AddExtensionType(".json", "application/json")
}

//go:embed all:web
var webFS embed.FS

// prod 模式：embed Vue dist，Go 同時 serve API + SPA
func setupFrontend(e *echo.Echo) {
	sub, _ := fs.Sub(webFS, "web")

	// 🔍 啟動時檢查：印出所有嵌入的檔案清單
    fs.WalkDir(sub, ".", func(path string, d fs.DirEntry, err error) error {
        if !d.IsDir() {
            fmt.Println("已嵌入檔案:", path)
        }
        return nil
    })

	e.RouteNotFound("/*", func(c echo.Context) error {
		if c.Request().Method != http.MethodGet {
			return echo.ErrMethodNotAllowed
		}

		path := strings.TrimPrefix(c.Request().URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		data, err := fs.ReadFile(sub, path)
		if err != nil {
			// 如果請求的是靜態資源但找不到，回傳 404
			if strings.HasPrefix(path, "assets/") || filepath.Ext(path) != "" {
				return echo.ErrNotFound
			}
			// SPA 路由回退
			data, err = fs.ReadFile(sub, "index.html")
			if err != nil {
				return echo.ErrNotFound
			}
			return c.HTMLBlob(http.StatusOK, data)
		}

		// 使用封裝好的函數獲取 MIME
		contentType := getMimeType(path, data)

		// --- 關鍵標頭設定 ---
		// 告訴瀏覽器：不要猜測內容類型，嚴格遵守我給你的 Content-Type
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		
		// 如果是 JS 檔案，明確加上 charset=utf-8 避免亂碼報錯
		return c.Blob(http.StatusOK, contentType, data)
	})
}

func getMimeType(path string, data []byte) string {
	ext := filepath.Ext(path)
	
	// 1. 優先從副檔名判斷 (最準確，因為我們已經在 init 註冊過)
	ct := mime.TypeByExtension(ext)
	if ct != "" {
		return ct
	}

	// 2. 如果是 .js 卻沒抓到（保險機制）
	if ext == ".js" || ext == ".mjs" {
		return "application/javascript; charset=utf-8"
	}

	// 3. 最後才嘗試內容嗅探 (DetectContentType 對於文本類型的判定較弱)
	return http.DetectContentType(data)
}