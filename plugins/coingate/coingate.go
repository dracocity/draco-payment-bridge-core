package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
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

func (p *coingatePlugin) Load(cfg config.PGConfig) error {
	p.apiKey = os.Getenv("COINGATE_API_TOKEN")
	p.pg = &coingatePG{apiKey: p.apiKey}
	return nil
}

func (p *coingatePlugin) Unload() error {
	return nil
}

func (p *coingatePlugin) Name() string {
	return "coingate"
}

func (p *coingatePlugin) PaymentGateway() pg.PaymentGateway {
	return p.pg
}

func (b *coingatePG) Name() string {
	return "coingate"
}

func (b *coingatePG) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	return nil, nil
}

func (b *coingatePG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	// paymentID := fmt.Sprintf("cg-%d", time.Now().UnixNano())
	createdAt := time.Now().UnixMilli()
	expiresAt := createdAt + 1800000 // 30 minutes
	return &models.CreatePaymentResponse{
		// PaymentID:  paymentID,
		// PaymentURL: "https://coingate.com/pay/" + paymentID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    "new",
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

func (b *coingatePG) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
	return &models.GetPaymentResponse{
		PaymentID: paymentID,
		Status:    "pending",
		Amount:    "",
		Currency:  "",
		// Transaction: "",
		UpdatedAt: time.Now().UnixMilli(),
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
