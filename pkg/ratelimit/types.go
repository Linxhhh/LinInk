package ratelimit

import "context"

type Limiter interface {
	// key 就是限流对象
	// bool 代表是否限流，true 表示需要限流
	// err 限流器本身有咩有错误
	Limit(ctx context.Context, key string) (bool, error)
}