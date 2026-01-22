package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
)

type coinbaseCommercePlugin struct {
	apiKey string
	pg     *coinbaseCommercePG
}

type coinbaseCommercePG struct {
	apiKey string
}

func New() plugin.Plugin {
	return &coinbaseCommercePlugin{}
}

func (p *coinbaseCommercePlugin) New() error {
	return nil
}

func (p *coinbaseCommercePlugin) Load() error {
	p.apiKey = os.Getenv("COINBASE_COMMERCE_API_KEY")
	p.pg = &coinbaseCommercePG{apiKey: p.apiKey}
	return nil
}

func (p *coinbaseCommercePlugin) Unload() error {
	return nil
}

func (p *coinbaseCommercePlugin) GetName() string {
	return "coinbasecommerce"
}

func (p *coinbaseCommercePlugin) PaymentGateway() pg.PaymentGateway {
	return p.pg
}

func (b *coinbaseCommercePG) Name() string {
	return "coinbasecommerce"
}

func (b *coinbaseCommercePG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	paymentID := fmt.Sprintf("cbc-%d", time.Now().UnixNano())
	return &models.CreatePaymentResponse{
		Success:    true,
		PaymentID:  paymentID,
		PaymentURL: "https://commerce.coinbase.com/charges/" + paymentID,
		QRCode:     "",
		Amount:     formatAmount(req.Amount),
		Currency:   req.Currency,
		Status:     "new",
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	}, nil
}

func (b *coinbaseCommercePG) GetStatus(ctx context.Context, paymentID string) (*models.PaymentStatusResponse, error) {
	return &models.PaymentStatusResponse{
		PaymentID:   paymentID,
		Status:      "pending",
		Amount:      "",
		Currency:    "",
		Transaction: "",
		UpdatedAt:   time.Now(),
	}, nil
}

func (b *coinbaseCommercePG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
	refundType := "full"
	amount := ""
	if req.Amount > 0 {
		refundType = "partial"
		amount = formatAmount(req.Amount)
	}
	return &models.RefundResponse{
		Success:    true,
		RefundID:   fmt.Sprintf("cbc-refund-%d", time.Now().UnixNano()),
		PaymentID:  req.PaymentID,
		Amount:     amount,
		Currency:   "",
		Status:     "pending",
		RefundType: refundType,
	}, nil
}

func (b *coinbaseCommercePG) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
