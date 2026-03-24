package pg

import (
	"context"
	"net/http"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

// PaymentGateway defines a payment provider integration.
type PaymentGateway interface {
	Name() string
	CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error)
	CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error)
	GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error)
	CreateRefund(ctx context.Context, req models.CreateRefundRequest) (*models.CreateRefundResponse, error)
	HandleWebhook(ctx context.Context, payload []byte, header http.Header) (*models.WebhookResult, error)
}
