package stripe

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/vietquan-37/payment-service/pkg/config"
	"github.com/vietquan-37/payment-service/pkg/pb"
	"github.com/vietquan-37/payment-service/pkg/provider"
)

type Stripe struct {
	cfg *config.Config
}

func NewPaymentProvider(cfg *config.Config) provider.PaymentProvider {
	return &Stripe{cfg: cfg}
}

func (s *Stripe) CreatePaymentLink(p *pb.PaymentLinkRequest) (string, error) {

	stripe.Key = s.cfg.StripeSecretKey
	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range p.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:          stripe.String(s.cfg.Currency),
				UnitAmountDecimal: stripe.Float64(float64(item.Price) * 100),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.ProductName),
				},
			},
			Quantity: stripe.Int64(item.Quantity),
		})
	}
	itemsJSON, err := json.Marshal(p.Items)
	if err != nil {
		return "", fmt.Errorf("failed to serialize items: %w", err)
	}
	customerJSON, err := json.Marshal(p.Customer)
	if err != nil {
		return "", fmt.Errorf("failed to serialize customer: %w", err)
	}
	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.cfg.SuccessURL),
		CancelURL:  stripe.String(s.cfg.CancelURL),
		LineItems:  items,
		Metadata: map[string]string{
			"order_id": fmt.Sprintf("%d", p.OrderId),
			"customer": string(customerJSON),
			"items":    string(itemsJSON),
		},
	}
	sess, err := session.New(params)
	if err != nil {
		log.Error().Err(err).Msg("stripe: failed to create checkout session")
		return "", err
	}

	log.Info().Str("url", sess.URL).Msg("stripe: created checkout session")
	return sess.URL, nil
}
