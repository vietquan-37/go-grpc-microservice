package routes

func AccessibleRoles() map[string][]string {
	return map[string][]string{
		"/orderpb.OrderService/AddProduct":   {"USER"},
		"/orderpb.OrderService/DeleteDetail": {"USER"},
		"/orderpb.OrderService/GetUserCart":  {"USER"},
	}
}