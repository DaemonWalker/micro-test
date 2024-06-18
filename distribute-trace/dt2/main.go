package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

var tracer = otel.Tracer("dt2-tracer")

func main() {
	conn, err := initConn()
	if err != nil {
		panic(err)
	}
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			semconv.ServiceNameKey.String("dt2"),
		),
	)
	if err != nil {
		panic(err)
	}
	tp, err := initTracer(conn, res)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	r := echo.New()
	r.Use(otelecho.Middleware("dt2-trace"))

	r.GET("/", func(c echo.Context) error {
		_, span := tracer.Start(c.Request().Context(), "dt2-home")
		defer span.End()
		if span.IsRecording() {
			span.SetAttributes(attribute.String("result", "alive"))
		}
		time.Sleep(time.Second)
		return c.String(200, "alive")
	})
	log.Fatal(r.Start(":8081"))
}

func initConn() (*grpc.ClientConn, error) {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	// You may replace `localhost:4317` with your endpoint.
	conn, err := grpc.NewClient("localhost:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

func initTracer(conn *grpc.ClientConn, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
