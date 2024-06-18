package main

import (
	"mt-app1/controller"

	"github.com/labstack/echo/v4"
)

func bindRouter(r *echo.Echo) {
	r.POST("/buy", controller.Buy)
}
