package webhook

import (
	"common/kafka/producer"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"github.com/vietquan-37/payment-service/pkg/config"
	"github.com/vietquan-37/payment-service/pkg/message"
	"io"
	"net/http"

	"strconv"
	"time"
)

type PaymentHttpHandler struct {
	cfg *config.Config
	p   *producer.Producer
}

func NewPaymentHttpHandler(cfg *config.Config, p *producer.Producer) *PaymentHttpHandler {
	return &PaymentHttpHandler{
		cfg: cfg,
		p:   p,
	}
}

func (h *PaymentHttpHandler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handleCheckoutWebhook)
}
func (h *PaymentHttpHandler) handleCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEventWithOptions(payload, signatureHeader, h.cfg.StripeSignature, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})

	if err != nil {
		log.Error().Msgf("Webhook signature verification failed. %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Error().Msgf("Error unmarshalling checkout session data: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		if session.PaymentStatus == "paid" {
			log.Printf("Payment successfull for %s", session.ID)
			orderID := session.Metadata["order_id"]
			items := session.Metadata["items"]
			Id, err := strconv.Atoi(orderID)
			if err != nil {
				log.Error().Msgf("Error converting order id to integer: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var orderItems []message.ItemPurchased
			if err := json.Unmarshal([]byte(items), &orderItems); err != nil {
				log.Error().Msgf("Error unmarshalling items metadata: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			payload, err := message.NewPaymentEnvelope("payment-service", "1", message.PaymentSucceededMessage{
				OrderID: int32(Id),
				Items:   orderItems,
			})
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = h.p.SendMessage(ctx, h.cfg.Topic, nil, payload)
			if err != nil {
				log.Error().Msgf("Error sending message: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Info().Msgf("Payment successful for %s", session.ID)

		}
	}

	w.WriteHeader(http.StatusOK)
}
