package main

import (
	"context"
	"log"
	"mt-app1/data"
	"mt-app1/middleware"
	"mt-app1/optl"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho" // nolint:staticcheck  // deprecated.
)

func initialize() {
	data.InitDb()
}

func main() {
	tp := optl.InitOptl()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	initialize()

	r := echo.New()
	bindValidator(r)
	bindRouter(r)
	middleware.InitMiddleware(r)

	r.Use(echoprometheus.NewMiddleware("mt_app1")) // adds middleware to gather metrics
	r.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	r.Use(otelecho.Middleware("mt-app1-middleware"))

	r.GET("/", func(c echo.Context) error {
		time.Sleep(1 * time.Second)
		return c.String(200, "alive")
	})
	_ = r.Start(":8080")
}
