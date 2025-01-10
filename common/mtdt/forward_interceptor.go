package mtdt

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ForwardMetadataUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//this receives metadata from when client call server
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "Metadata error")
		}
		//this for service to call another service
		//that is why metadata cannot receive it grpc service call to another grpc service
		outgoingContext := metadata.NewOutgoingContext(ctx, md)
		return handler(outgoingContext, req)
	}
}
