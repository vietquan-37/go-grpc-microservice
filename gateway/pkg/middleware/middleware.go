package middleware

import (
	"github.com/vietquan-37/gateway/pkg/auth/pb"
	"github.com/vietquan-37/gateway/pkg/routes"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader     = "Authorization"
	authorizationBearerType = "Bearer"
)

type AuthMiddleWareConfig struct {
	AuthService pb.AuthServiceClient
}

func NewAuthMiddleWareConfig(authService pb.AuthServiceClient) *AuthMiddleWareConfig {
	return &AuthMiddleWareConfig{
		AuthService: authService,
	}

}

func (a *AuthMiddleWareConfig) AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if PublicMethod(r.Method, r.URL.Path) {
			handler.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get(authorizationHeader)
		if len(authHeader) == 0 {
			http.Error(w, "Missing the authorization header", http.StatusUnauthorized)
			return
		}
		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}
		authorizationType := fields[0]
		if authorizationType != authorizationBearerType {
			http.Error(w, "Invalid authorization type", http.StatusUnauthorized)
			return
		}
		accessToken := fields[1]
		payload, err := a.AuthService.Validate(r.Context(), &pb.ValidateRequest{Token: accessToken})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		userRole := payload.Role
		if !HasPermission(r.Method, r.URL.Path, userRole) {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		}
		md := metadata.Pairs("id", strconv.Itoa(int(payload.UserId)), "role", payload.Role)
		ctx := metadata.NewOutgoingContext(r.Context(), md)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
func PublicMethod(method string, path string) bool {
	for _, apiRoute := range routes.ApiRoutes {
		if apiRoute.Method == method && apiRoute.Path == path && apiRoute.Roles == nil {
			return true

		}
	}
	return false
}
func HasPermission(method string, path string, userRole string) bool {
	for _, apiRoute := range routes.ApiRoutes {
		if apiRoute.Method == method && strings.Contains(path, apiRoute.Path) {
			for _, role := range apiRoute.Roles {
				if role == userRole {
					return true
				}
			}
		}
	}
	return false
}
