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

type coingatePlugin struct {
	apiKey string
	pg     *coingatePG
}

type coingatePG struct {
	apiKey string
}

func New() plugin.Plugin {
	return &coingatePlugin{}
}

func (p *coingatePlugin) New() error {
	return nil
}

func (p *coingatePlugin) Load() error {
	p.apiKey = os.Getenv("COINGATE_API_TOKEN")
	p.pg = &coingatePG{apiKey: p.apiKey}
	return nil
}

func (p *coingatePlugin) Unload() error {
	return nil
}

func (p *coingatePlugin) GetName() string {
	return "coingate"
}

func (p *coingatePlugin) PaymentGateway() pg.PaymentGateway {
	return p.pg
}

func (b *coingatePG) Name() string {
	return "coingate"
}

func (b *coingatePG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	paymentID := fmt.Sprintf("cg-%d", time.Now().UnixNano())
	return &models.CreatePaymentResponse{
		Success:    true,
		PaymentID:  paymentID,
		PaymentURL: "https://coingate.com/pay/" + paymentID,
		QRCode:     "",
		Amount:     formatAmount(req.Amount),
		Currency:   req.Currency,
		Status:     "new",
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	}, nil
}

func (b *coingatePG) GetStatus(ctx context.Context, paymentID string) (*models.PaymentStatusResponse, error) {
	return &models.PaymentStatusResponse{
		PaymentID:   paymentID,
		Status:      "pending",
		Amount:      "",
		Currency:    "",
		Transaction: "",
		UpdatedAt:   time.Now(),
	}, nil
}

func (b *coingatePG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
	refundType := "full"
	amount := ""
	if req.Amount > 0 {
		refundType = "partial"
		amount = formatAmount(req.Amount)
	}
	return &models.RefundResponse{
		Success:    true,
		RefundID:   fmt.Sprintf("cg-refund-%d", time.Now().UnixNano()),
		PaymentID:  req.PaymentID,
		Amount:     amount,
		Currency:   "",
		Status:     "pending",
		RefundType: refundType,
	}, nil
}

func (b *coingatePG) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
