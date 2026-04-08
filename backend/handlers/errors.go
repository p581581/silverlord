package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorCode 業務錯誤碼，前端依此判斷錯誤類型
type ErrorCode string

const (
	// ── 認證 ────────────────────────────────────────────────────
	ErrAuthPasswordRequired ErrorCode = "AUTH_PASSWORD_REQUIRED"
	ErrAuthWrongPassword    ErrorCode = "AUTH_WRONG_PASSWORD"
	ErrAuthMissingToken     ErrorCode = "AUTH_MISSING_TOKEN"
	ErrAuthInvalidToken     ErrorCode = "AUTH_INVALID_TOKEN"
	ErrAuthForbidden        ErrorCode = "AUTH_FORBIDDEN"

	// ── 報價 / Cache ─────────────────────────────────────────────
	ErrPriceNotReady ErrorCode = "PRICE_NOT_READY" // Cache 尚未初始化

	// ── 商品 ────────────────────────────────────────────────────
	ErrProductNotFound     ErrorCode = "PRODUCT_NOT_FOUND"
	ErrProductNameRequired ErrorCode = "PRODUCT_NAME_REQUIRED"
	ErrProductInvalidWeight ErrorCode = "PRODUCT_INVALID_WEIGHT"
	ErrProductInactive     ErrorCode = "PRODUCT_INACTIVE"

	// ── 訂單 ────────────────────────────────────────────────────
	ErrOrderItemsRequired  ErrorCode = "ORDER_ITEMS_REQUIRED"
	ErrOrderInvalidQty     ErrorCode = "ORDER_INVALID_QTY"
	ErrOrderNotFound       ErrorCode = "ORDER_NOT_FOUND"
	ErrOrderAlreadyReported ErrorCode = "ORDER_ALREADY_REPORTED"
	ErrOrderInvalidStatus  ErrorCode = "ORDER_INVALID_STATUS"
	ErrOrderCreateFailed   ErrorCode = "ORDER_CREATE_FAILED"

	// ── 系統設定 ─────────────────────────────────────────────────
	ErrConfigTgNotConfigured ErrorCode = "CONFIG_TG_NOT_CONFIGURED"
	ErrConfigTgSendFailed    ErrorCode = "CONFIG_TG_SEND_FAILED"

	// ── 上傳 ────────────────────────────────────────────────────
	ErrUploadFileRequired ErrorCode = "UPLOAD_FILE_REQUIRED"
	ErrUploadWriteFailed  ErrorCode = "UPLOAD_WRITE_FAILED"

	// ── 系統 ────────────────────────────────────────────────────
	ErrInternal ErrorCode = "INTERNAL_ERROR"
)

// apiError 是統一的錯誤回應結構
type apiError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// fail 回傳標準錯誤 JSON
func fail(c echo.Context, status int, code ErrorCode, message string) error {
	return c.JSON(status, apiError{Code: code, Message: message})
}

// 常用快捷函式

func errBadRequest(c echo.Context, code ErrorCode, msg string) error {
	return fail(c, http.StatusBadRequest, code, msg)
}

func errNotFound(c echo.Context, code ErrorCode, msg string) error {
	return fail(c, http.StatusNotFound, code, msg)
}

func errConflict(c echo.Context, code ErrorCode, msg string) error {
	return fail(c, http.StatusConflict, code, msg)
}

func errUnauthorized(c echo.Context, code ErrorCode, msg string) error {
	return fail(c, http.StatusUnauthorized, code, msg)
}

func errServiceUnavailable(c echo.Context, code ErrorCode, msg string) error {
	return fail(c, http.StatusServiceUnavailable, code, msg)
}

func errInternal(c echo.Context, msg string) error {
	return fail(c, http.StatusInternalServerError, ErrInternal, msg)
}
