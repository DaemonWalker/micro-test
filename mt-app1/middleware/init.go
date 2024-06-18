package middleware

import "github.com/labstack/echo/v4"

func InitMiddleware(r *echo.Echo) {
	r.Use(transaction)
}
