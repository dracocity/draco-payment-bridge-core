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

type bitpayPlugin struct {
	apiKey string
	bridge *bitpayBridge
}

type bitpayBridge struct {
	apiKey string
}

func New() plugin.Plugin {
	return &bitpayPlugin{}
}

func (p *bitpayPlugin) New() error {
	return nil
}

func (p *bitpayPlugin) Load() error {
	p.apiKey = os.Getenv("BITPAY_API_KEY")
	p.bridge = &bitpayBridge{apiKey: p.apiKey}
	return nil
}

func (p *bitpayPlugin) Unload() error {
	return nil
}

func (p *bitpayPlugin) GetName() string {
	return "bitpay"
}

func (p *bitpayPlugin) Bridge() bridge.Bridge {
	return p.bridge
}

func (b *bitpayBridge) Name() string {
	return "bitpay"
}

func (b *bitpayBridge) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	paymentID := fmt.Sprintf("bp-%d", time.Now().UnixNano())
	return &models.CreatePaymentResponse{
		Success:    true,
		PaymentID:  paymentID,
		PaymentURL: "https://bitpay.com/invoice?id=" + paymentID,
		QRCode:     "",
		Amount:     formatAmount(req.Amount),
		Currency:   req.Currency,
		Status:     "new",
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	}, nil
}

func (b *bitpayBridge) GetStatus(ctx context.Context, paymentID string) (*models.PaymentStatusResponse, error) {
	return &models.PaymentStatusResponse{
		PaymentID:   paymentID,
		Status:      "pending",
		Amount:      "",
		Currency:    "",
		Transaction: "",
		UpdatedAt:   time.Now(),
	}, nil
}

func (b *bitpayBridge) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
	refundType := "full"
	amount := ""
	if req.Amount > 0 {
		refundType = "partial"
		amount = formatAmount(req.Amount)
	}
	return &models.RefundResponse{
		Success:    true,
		RefundID:   fmt.Sprintf("bp-refund-%d", time.Now().UnixNano()),
		PaymentID:  req.PaymentID,
		Amount:     amount,
		Currency:   "",
		Status:     "pending",
		RefundType: refundType,
	}, nil
}

func (b *bitpayBridge) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
