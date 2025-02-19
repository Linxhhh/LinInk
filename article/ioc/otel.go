package ioc

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	jaeger "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

/*
	批量创建 span，并将 trace 数据发送到 Jaeger
*/

const (
	serviceName    = "LinInk-HTTP-service"
	jaegerEndpoint = "127.0.0.1:4318"
)

func InitOtel() func(ctx context.Context) {

	res, err := newResource(serviceName)
	if err != nil {
		panic(err)
	}

	// 设置传播器 ———— 在客户端和服务端之间传递 tracing 相关信息
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// 设置 trace provider ———— 用于在打点时，构建 trace
	tp, err := newTraceProvider(res)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)

	return func(ctx context.Context) {
		tp.Shutdown(ctx)
	}
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newResource(serviceName string) (*resource.Resource, error) {
	return resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("environment", "env"),
			attribute.Int64("ID", 2),
		),
	)
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := jaeger.New(
		context.Background(),
		jaeger.WithEndpoint(jaegerEndpoint),
		jaeger.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		// Default is 5s. Set to 1s for demonstrative purposes.
		trace.WithBatcher(exporter,	trace.WithBatchTimeout(time.Second)),  
	)
	return traceProvider, nil
}
