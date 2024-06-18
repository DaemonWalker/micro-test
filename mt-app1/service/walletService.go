package service

import (
	"mt-app1/data"
	"mt-app1/optl"
	optlattr "mt-app1/optl/attr"

	"github.com/labstack/echo/v4"
)

const (
	sellServiceSpanName = "Wallet Service"
)

func MinusMoney(c echo.Context, userId int, money int) error {
	span := optl.StartTrace(c.Request().Context(), sellServiceSpanName,
		optlattr.MethodName("MinusMoney"))
	defer span.End()

	tx, err := getTransaction(c)
	if err != nil {
		return err
	}

	_, err = checkExecResult(data.UpdateMoney(userId, money, tx))

	return err
}
