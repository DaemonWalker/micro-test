package middleware

import (
	"database/sql"
	"mt-app1/constance"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

func transaction(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		err := next(c)

		v := c.Get(constance.TxKey)
		if v == nil {
			return err
		}

		tx := v.(*sql.Tx)
		v = c.Get(constance.TraceKey)
		trace := *(v.(*trace.Span))
		defer trace.End()
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		return err
	})
}
