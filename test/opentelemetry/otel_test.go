package opentelemetry

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
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
	serviceName    = "Go-Jaeger-Demo"
	jaegerEndpoint = "127.0.0.1:4318"
)

func TestServer(t *testing.T) {

	res, err := newResource(serviceName)
	require.NoError(t, err)

	// 设置传播器 ———— 在客户端和服务端之间传递 tracing 相关信息
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// 设置 trace provider ———— 用于在打点时，构建 trace
	tp, err := newTraceProvider(res)
	require.NoError(t, err)
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	// 启动服务器
	server := gin.Default()
	server.GET("/test", func(ginCtx *gin.Context) {

		// tracer name
		tracer := otel.Tracer("opentelemetry")

		// create top span
		var ctx context.Context = ginCtx
		ctx, span := tracer.Start(ctx, "top-span")
		defer span.End()

		// add event for top span
		span.AddEvent("test-event")
		time.Sleep(time.Second)

		// create sub span
		for i := 0; i < 10; i++ {
			_, subSpan := tracer.Start(ctx, fmt.Sprintf("sub-span-%d", i))
			time.Sleep(time.Millisecond * 300)

			// set attributes
			subSpan.SetAttributes(attribute.String("key1", "value1"))
			subSpan.End()
		}

		// http response
		ginCtx.String(http.StatusOK, "OK")
	})
	server.Run(":3333")
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
		trace.WithBatcher(exporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}