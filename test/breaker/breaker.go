package breaker

import (
	"context"
	"log"
	"time"

	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
)

// breaker Interceptor
type InterceptorBuilder struct {
	breaker *gobreaker.CircuitBreaker[any]
}

func NewInterceptorBuilder(name string, timeout time.Duration, maxRequests uint32) *InterceptorBuilder {

	// 使用 OnStateChange 来设置熔断器状态变化的回调
	onStateChange := func(name string, from gobreaker.State, to gobreaker.State) {
		if from == gobreaker.StateClosed && to == gobreaker.StateOpen {
			// 熔断器触发
			log.Printf("Circuit breaker '%s' has tripped", name)
		}
	}

	breaker := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
		Name:          name,
		Timeout:       timeout,          // 熔断后多久转为半开状态
		MaxRequests:   maxRequests,      // 半开状态下的最大请求数
		Interval:      10 * time.Second, // 正常状态下多久清理一次计数
		OnStateChange: onStateChange,
	})

	return &InterceptorBuilder{
		breaker: breaker,
	}
}

func (b *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		// 使用熔断器
		resp, err = b.breaker.Execute(func() (any, error) {
			return handler(ctx, req)
		})
		return
	}
}
