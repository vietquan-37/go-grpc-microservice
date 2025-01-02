package loggers

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

const (
	xForwardedForHeader = "x-forwarded-for"
	ipUnknown           = "unknown"
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
	clientIp := extractUserIp(ctx)
	if clientIp == ipUnknown {
		log.Warn().Err(err).Msg("Failed to extract client IP")
	}
	logger.
		Str("protocol", "grpc").
		Str("ip_address", clientIp).
		Dur("duration", duration).
		Str("status_text", statusCode.String()).
		Str("method", info.FullMethod).
		Msg("Received a grpc request")

	return result, err
}

func extractUserIp(ctx context.Context) string {
	// from mtdt in grpc
	if meta, ok := metadata.FromIncomingContext(ctx); ok {
		if clientIp := meta.Get(xForwardedForHeader); len(clientIp) > 0 {
			return clientIp[0]
		}
	}

	//if p, ok := peer.FromContext(ctx); ok {
	//	return p.Addr.String()
	//}

	return ipUnknown
}
