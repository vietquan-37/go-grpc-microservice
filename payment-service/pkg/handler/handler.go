package handler

import (
	"context"

	"github.com/vietquan-37/payment-service/pkg/pb"
	"github.com/vietquan-37/payment-service/pkg/provider"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	payment provider.PaymentProvider
}

func NewPaymentHandler(payment provider.PaymentProvider) *PaymentHandler {
	return &PaymentHandler{
		payment: payment,
	}
}
func (h *PaymentHandler) CreatePaymentLink(ctx context.Context, req *pb.PaymentLinkRequest) (*pb.PaymentLinkResponse, error) {
	link, err := h.payment.CreatePaymentLink(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while creating payment link: %v", err)
	}
	return &pb.PaymentLinkResponse{
		Link: link,
	}, nil
}
