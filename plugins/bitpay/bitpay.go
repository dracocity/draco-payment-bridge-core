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

type bitpayPlugin struct {
	apiKey string
	pg     *bitpayPG
}

type bitpayPG struct {
	apiKey string
}

func New() plugin.Plugin {
	return &bitpayPlugin{}
}

func (p *bitpayPlugin) Load(cfg config.PGConfig) error {
	p.apiKey = os.Getenv("BITPAY_API_KEY")
	p.pg = &bitpayPG{apiKey: p.apiKey}
	return nil
}

func (p *bitpayPlugin) Unload() error {
	return nil
}

func (p *bitpayPlugin) Name() string {
	return "bitpay"
}

func (p *bitpayPlugin) PaymentGateway() pg.PaymentGateway {
	return p.pg
}

func (b *bitpayPG) Name() string {
	return "bitpay"
}

func (b *bitpayPG) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	return nil, nil
}

func (b *bitpayPG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	// paymentID := fmt.Sprintf("bp-%d", time.Now().UnixNano())
	return &models.CreatePaymentResponse{
		// PaymentID: paymentID,
		// PaymentURL: "https://bitpay.com/invoice?id=" + paymentID,
		// QRCode:     "",
		// Amount:     req.Amount,
		// Currency:   req.Currency,
		// Status:     "new",
		// ExpiresAt:  time.Now().Add(30 * time.Minute),
	}, nil
}

func (b *bitpayPG) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
	return &models.GetPaymentResponse{
		PaymentID: paymentID,
		Status:    "pending",
		Amount:    "",
		Currency:  "",
		// Transaction: "",
		UpdatedAt: time.Now().UnixMilli(),
	}, nil
}

func (b *bitpayPG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
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

func (b *bitpayPG) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
