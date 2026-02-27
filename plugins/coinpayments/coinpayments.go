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

type coinpaymentsPlugin struct {
	publicKey  string
	privateKey string
	pg         *coinpaymentsPG
}

type coinpaymentsPG struct {
	publicKey  string
	privateKey string
}

func New() plugin.Plugin {
	return &coinpaymentsPlugin{}
}

func (p *coinpaymentsPlugin) Load(cfg config.PGConfig) error {
	p.publicKey = os.Getenv("COINPAYMENTS_PUBLIC_KEY")
	p.privateKey = os.Getenv("COINPAYMENTS_PRIVATE_KEY")
	p.pg = &coinpaymentsPG{publicKey: p.publicKey, privateKey: p.privateKey}
	return nil
}

func (p *coinpaymentsPlugin) Unload() error {
	return nil
}

func (p *coinpaymentsPlugin) Name() string {
	return "coinpayments"
}

func (p *coinpaymentsPlugin) PaymentGateway() pg.PaymentGateway {
	return p.pg
}

func (b *coinpaymentsPG) Name() string {
	return "coinpayments"
}

func (b *coinpaymentsPG) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	return nil, nil
}

func (b *coinpaymentsPG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	createdAt := time.Now().UnixMilli()
	expiresAt := createdAt + 1800000 // 30 minutes
	return &models.CreatePaymentResponse{
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    "new",
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

func (b *coinpaymentsPG) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
	return &models.GetPaymentResponse{
		PaymentID: paymentID,
		Status:    "pending",
		Amount:    "",
		Currency:  "",
		// Transaction: "",
		UpdatedAt: time.Now().UnixMilli(),
	}, nil
}

func (b *coinpaymentsPG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
	refundType := "full"
	amount := ""
	if req.Amount > 0 {
		refundType = "partial"
		amount = formatAmount(req.Amount)
	}
	return &models.RefundResponse{
		Success:    true,
		RefundID:   fmt.Sprintf("cp-refund-%d", time.Now().UnixNano()),
		PaymentID:  req.PaymentID,
		Amount:     amount,
		Currency:   "",
		Status:     "pending",
		RefundType: refundType,
	}, nil
}

func (b *coinpaymentsPG) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
