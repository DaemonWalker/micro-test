package optl

import (
	"context"
	"fmt"
	"log"

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

var tracer = otel.Tracer("echo-server")

func InitOptl() *sdktrace.TracerProvider {
	conn, err := initConn()
	if err != nil {
		panic(err)
	}
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			semconv.ServiceNameKey.String("mt-app1"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	tp, err := initTracer(conn, res)
	if err != nil {
		log.Fatal(err)
	}

	return tp
}

func StartTrace(ctx context.Context, name string, attributes ...attribute.KeyValue) trace.Span {
	_, span := tracer.Start(ctx, name, trace.WithAttributes(attributes...))
	return span
}

func StartTraceDefer(ctx context.Context, name string, attributes ...attribute.KeyValue) func(...trace.SpanEndOption) {
	_, span := tracer.Start(ctx, name, trace.WithAttributes(attributes...))
	return span.End
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
