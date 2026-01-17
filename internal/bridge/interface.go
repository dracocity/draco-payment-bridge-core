package bridge

import (
	"context"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

// Bridge defines a payment provider integration.
type Bridge interface {
	Name() string
	CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error)
	GetStatus(ctx context.Context, paymentID string) (*models.PaymentStatusResponse, error)
	Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error)
	HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error)
}
