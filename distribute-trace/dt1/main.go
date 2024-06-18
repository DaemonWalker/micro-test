package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

var tracer = otel.Tracer("dt1-tracer")

func main() {
	conn, err := initConn()
	if err != nil {
		panic(err)
	}
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			semconv.ServiceNameKey.String("dt1"),
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
	r.Use(otelecho.Middleware("dt1-trace"))

	r.GET("/", func(c echo.Context) error {
		_, span := tracer.Start(c.Request().Context(), "dt1-home")
		defer span.End()

		client := &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   10 * time.Second,
		}
		time.Sleep(3 * time.Second)
		_, span = tracer.Start(c.Request().Context(), "say hello", trace.WithAttributes(semconv.PeerService("dt1-http-client")))
		defer span.End()

		span.AddEvent("Start Http Request")

		req, _ := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, "http://localhost:8081", nil)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		msg := string(body)

		if span.IsRecording() {
			span.AddEvent("Get Http Request Result", trace.WithAttributes(attribute.String("Get response body", msg)))
		}

		time.Sleep(1 * time.Second)
		return c.String(200, "Get response from 8081: \""+msg+"\"")
	})
	log.Fatal(r.Start(":8080"))
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
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, nil
}
