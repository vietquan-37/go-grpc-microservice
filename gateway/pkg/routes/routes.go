package routes

import "net/http"

type APIRoutes struct {
	Path   string
	Method string
	Roles  []string
}

var ApiRoutes = []APIRoutes{
	{
		Path: "/v1/login", Method: http.MethodPost,
	},
	{
		Path: "/v1/register", Method: http.MethodPost,
	},
	{
		Path: "/v1/products", Method: http.MethodGet,
	},
	{
		Path: "/v1/products/{id}", Method: http.MethodGet,
	},
	{
		Path: "/v1/product", Method: http.MethodPost, Roles: []string{"ADMIN"},
	},
	{
		Path: "/v1/product", Method: http.MethodPut, Roles: []string{"ADMIN"},
	},
	{
		Path: "/v1/order", Method: http.MethodPost, Roles: []string{"USER"},
	},
	{
		Path: "/v1/cart/{user_id}", Method: http.MethodGet, Roles: []string{"USER"},
	},
	{
		Path: "/v1/detail/{id}", Method: http.MethodDelete, Roles: []string{"USER"},
	},
}
