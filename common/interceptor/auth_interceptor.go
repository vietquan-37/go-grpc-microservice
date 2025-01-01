package interceptor

import (
	"common/client"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authClient client.AuthClient
	accessRole map[string][]string
}

func NewAuthInterceptor(authClient client.AuthClient, accessibleRole map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{authClient, accessibleRole}
}
func (interceptor *AuthInterceptor) UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}

}
func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) error {
	accessibleRoles, ok := interceptor.accessRole[method]
	if !ok {

		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := interceptor.authClient.Validate(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.User.Role {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
