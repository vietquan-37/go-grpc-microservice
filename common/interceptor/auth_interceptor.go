package interceptor

import (
	"common/client"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextType string

const (
	BearTokenType                   = "Bearer"
	AuthorizationHeader             = "authorization"
	UserContextKey      contextType = "user"
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
		newCtx, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		//capture ctx to handle downstream
		return handler(newCtx, req)
	}

}
func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	accessibleRoles, ok := interceptor.accessRole[method]
	if !ok {

		return ctx, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md[AuthorizationHeader]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	if !strings.HasPrefix(values[0], BearTokenType) {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not of type Bearer")
	}
	accessToken := strings.TrimPrefix(values[0], "Bearer ")
	claims, err := interceptor.authClient.Validate(ctx, accessToken)

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.User.Role {
			return context.WithValue(ctx, UserContextKey, claims), nil
		}
	}

	return nil, status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
