package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/bridge"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
)

type coinbaseCommercePlugin struct {
	apiKey string
	bridge *coinbaseCommerceBridge
}

type coinbaseCommerceBridge struct {
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
	p.bridge = &coinbaseCommerceBridge{apiKey: p.apiKey}
	return nil
}

func (p *coinbaseCommercePlugin) Unload() error {
	return nil
}

func (p *coinbaseCommercePlugin) GetName() string {
	return "coinbasecommerce"
}

func (p *coinbaseCommercePlugin) Bridge() bridge.Bridge {
	return p.bridge
}

func (b *coinbaseCommerceBridge) Name() string {
	return "coinbasecommerce"
}

func (b *coinbaseCommerceBridge) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
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

func (b *coinbaseCommerceBridge) GetStatus(ctx context.Context, paymentID string) (*models.PaymentStatusResponse, error) {
	return &models.PaymentStatusResponse{
		PaymentID:   paymentID,
		Status:      "pending",
		Amount:      "",
		Currency:    "",
		Transaction: "",
		UpdatedAt:   time.Now(),
	}, nil
}

func (b *coinbaseCommerceBridge) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
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

func (b *coinbaseCommerceBridge) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
