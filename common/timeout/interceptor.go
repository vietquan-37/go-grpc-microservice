package timeout

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func UnaryTimeoutInterceptor(duration time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, duration)
		defer cancel()
		return handler(ctx, req)
	}
}
