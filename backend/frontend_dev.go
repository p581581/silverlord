//go:build !prod

package main

import "github.com/labstack/echo/v4"

// dev 模式：Vue 由 Vite dev server 提供，Go 只跑 API
func setupFrontend(_ *echo.Echo) {}
