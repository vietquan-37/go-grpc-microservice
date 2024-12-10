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
		Path: "/v1/create_user", Method: http.MethodPost,
	},
	
	{
		Path: "/v1/products", Method: http.MethodGet,
	},
	{
		Path: "/v1/product/", Method: http.MethodGet,
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
		Path: "/v1/cart", Method: http.MethodGet, Roles: []string{"USER"},
	},
	{
		Path: "/v1/detail/", Method: http.MethodDelete, Roles: []string{"USER"},
	},
}
