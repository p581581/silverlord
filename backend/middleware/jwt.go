package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JWTClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type authError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RequireAdmin 驗證 Authorization: Bearer <token>
func RequireAdmin(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, authError{
					Code:    "AUTH_MISSING_TOKEN",
					Message: "missing or malformed token",
				})
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{},
				func(t *jwt.Token) (any, error) {
					return []byte(secret), nil
				})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, authError{
					Code:    "AUTH_INVALID_TOKEN",
					Message: "invalid or expired token",
				})
			}

			claims, ok := token.Claims.(*JWTClaims)
			if !ok || claims.Role != "admin" {
				return c.JSON(http.StatusForbidden, authError{
					Code:    "AUTH_FORBIDDEN",
					Message: "forbidden",
				})
			}

			c.Set("claims", claims)
			return next(c)
		}
	}
}
