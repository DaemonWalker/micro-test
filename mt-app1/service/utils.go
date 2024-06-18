package service

import (
	"database/sql"
	"errors"
	"mt-app1/constance"
	"mt-app1/data"
	"mt-app1/optl"

	"github.com/labstack/echo/v4"
)

var errFoo = errors.New("no rows effected")

func getTransaction(c echo.Context) (*sql.Tx, error) {
	var tx *sql.Tx
	var err error
	value := c.Get(constance.TxKey)
	if value != nil {
		tx = value.(*sql.Tx)
	} else {
		tx, err = data.Begin()
		c.Set(constance.TxKey, tx)
		span := optl.StartTrace(c.Request().Context(), "Transaction")
		c.Set(constance.TraceKey, &span)
	}
	return tx, err
}

func checkExecResult(re sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	r, e := re.RowsAffected()
	if e != nil {
		return 0, e
	}
	if r == 0 {
		return 0, errFoo
	}
	return r, e
}
