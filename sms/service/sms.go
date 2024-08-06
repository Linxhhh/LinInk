package service

import "context"

// Sms 短信服务接口
type Sms interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}