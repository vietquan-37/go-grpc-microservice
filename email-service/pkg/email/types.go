package email

type Type string

const (
	TypeVerification      Type = "verification"
	TypeOrderConfirmation Type = "order_confirmation"
)

type Template struct {
	Subject      string
	TemplateName string
}

var Templates = map[Type]Template{
	TypeVerification: {
		Subject:      "Welcome! Please verify your account",
		TemplateName: "verification.html",
	},
	TypeOrderConfirmation: {
		Subject:      "Order Confirmation",
		TemplateName: "order_confirmation.html",
	},
}
