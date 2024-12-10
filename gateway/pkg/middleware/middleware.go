package middleware

import (
	"encoding/json"
	"github.com/vietquan-37/gateway/pkg/auth/pb"
	"github.com/vietquan-37/gateway/pkg/routes"
	"google.golang.org/grpc/metadata"
	"log"
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

func writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(map[string]string{"message": message})
	if err != nil {
		log.Printf("Failed to write JSON response: %v", err)
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
	}
}

func (a *AuthMiddleWareConfig) AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !RouteExists(r.Method, r.URL.Path) {
			writeJSONResponse(w, http.StatusNotFound, "Route not found")
			return
		}
		if PublicMethod(r.Method, r.URL.Path) {
			handler.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get(authorizationHeader)
		if len(authHeader) == 0 {
			writeJSONResponse(w, http.StatusUnauthorized, "Missing the authorization header")
			return
		}
		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			writeJSONResponse(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}
		authorizationType := fields[0]
		if authorizationType != authorizationBearerType {
			writeJSONResponse(w, http.StatusUnauthorized, "Invalid authorization type")
			return
		}
		accessToken := fields[1]
		payload, err := a.AuthService.Validate(r.Context(), &pb.ValidateRequest{Token: accessToken})
		if err != nil {
			writeJSONResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		userRole := payload.Role
		if !HasPermission(r.Method, r.URL.Path, userRole) {
			writeJSONResponse(w, http.StatusForbidden, "Permission denied")
			return
		}
		md := metadata.Pairs("id", strconv.Itoa(int(payload.UserId)), "role", userRole)
		ctx := metadata.NewOutgoingContext(r.Context(), md)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RouteExists(method string, path string) bool {
	for _, apiRoute := range routes.ApiRoutes {
		if apiRoute.Method == method && strings.HasPrefix(path, apiRoute.Path) {
			return true
		}
	}
	return false
}

func PublicMethod(method string, path string) bool {
	for _, apiRoute := range routes.ApiRoutes {
		if apiRoute.Method == method && strings.HasPrefix(path, apiRoute.Path) && apiRoute.Roles == nil {
			return true
		}
	}
	return false
}

func HasPermission(method string, path string, userRole string) bool {
	for _, apiRoute := range routes.ApiRoutes {
		if apiRoute.Method == method && strings.HasPrefix(path, apiRoute.Path) {
			for _, role := range apiRoute.Roles {
				if role == userRole {
					return true
				}
			}
		}
	}
	return false
}
