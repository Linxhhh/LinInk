package ratelimit

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	limiter Limiter
	key     string

	// 服务降级 -> 返回默认值（json、默认函数）
	// defaultValueMap map[string]string
	// defaultValueFuncMap map[string]func() any 
}

func NewInterceptorBuilder(limiter Limiter, key string) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key}
}

func (b *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			log.Println("判定限流出现问题")
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")

			// 激进策略
			// return handler(ctx, req)
		}
		
		if limited {
			// defVal, ok := b.defaultValueMap[info.FullMethod]
			// if ok {
			// 	   err = json.Unmarshal([]byte(defVal), &resp)
			//	   return defVal, err
			// }
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		return handler(ctx, req)
	}
}