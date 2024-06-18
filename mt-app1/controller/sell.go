package controller

import (
	"mt-app1/model"
	"mt-app1/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Buy(c echo.Context) error {
	model := new(model.BuyItemModel)
	if err := c.Bind(model); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(model); err != nil {
		return err
	}
	if err := service.BuyItem(c, model); err != nil {
		log.Error(err)
		return err
	}
	return c.JSON(http.StatusOK, struct{}{})
}
