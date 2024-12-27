package interceptor

import (
	"context"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	xForwardedForHeader = "x-forwarded-for"
)

func GrpcLoggerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(start)
	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}
	//clientIp, err := extractUserIp(ctx)
	//if err != nil {
	//	logger = log.Error().Err(err)
	//}
	logger.
		Str("protocol", "grpc").
		//Str("ip_address", clientIp).
		Dur("duration", duration).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Str("method", info.FullMethod).
		Msg("Received a grpc request")

	return result, err
}

//func extractUserIp(ctx context.Context) (Ip string, err error) {
//	if meta, ok := metadata.FromIncomingContext(ctx); ok {
//
//		if ClientIp := meta.Get(xForwardedForHeader); len(ClientIp) > 0 {
//			Ip = ClientIp[0]
//		}
//	}
//	if Ip == "" {
//		return "", errors.New("no client ip found")
//	}
//	return Ip, nil
//}
