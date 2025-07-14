package provider

import (
	"github.com/vietquan-37/payment-service/pkg/pb"
)

type PaymentProvider interface {
	CreatePaymentLink(p *pb.PaymentLinkRequest) (string, error)
}
