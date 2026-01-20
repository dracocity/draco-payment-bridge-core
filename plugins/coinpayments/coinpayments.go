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

type coinpaymentsPlugin struct {
	publicKey  string
	privateKey string
	bridge     *coinpaymentsBridge
}

type coinpaymentsBridge struct {
	publicKey  string
	privateKey string
}

func New() plugin.Plugin {
	return &coinpaymentsPlugin{}
}

func (p *coinpaymentsPlugin) New() error {
	return nil
}

func (p *coinpaymentsPlugin) Load() error {
	p.publicKey = os.Getenv("COINPAYMENTS_PUBLIC_KEY")
	p.privateKey = os.Getenv("COINPAYMENTS_PRIVATE_KEY")
	p.bridge = &coinpaymentsBridge{publicKey: p.publicKey, privateKey: p.privateKey}
	return nil
}

func (p *coinpaymentsPlugin) Unload() error {
	return nil
}

func (p *coinpaymentsPlugin) GetName() string {
	return "coinpayments"
}

func (p *coinpaymentsPlugin) Bridge() bridge.Bridge {
	return p.bridge
}

func (b *coinpaymentsBridge) Name() string {
	return "coinpayments"
}

func (b *coinpaymentsBridge) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	paymentID := fmt.Sprintf("cp-%d", time.Now().UnixNano())
	return &models.CreatePaymentResponse{
		Success:    true,
		PaymentID:  paymentID,
		PaymentURL: "https://www.coinpayments.net/index.php?cmd=_pay&reset=1&id=" + paymentID,
		QRCode:     "",
		Amount:     formatAmount(req.Amount),
		Currency:   req.Currency,
		Status:     "new",
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	}, nil
}

func (b *coinpaymentsBridge) GetStatus(ctx context.Context, paymentID string) (*models.PaymentStatusResponse, error) {
	return &models.PaymentStatusResponse{
		PaymentID:   paymentID,
		Status:      "pending",
		Amount:      "",
		Currency:    "",
		Transaction: "",
		UpdatedAt:   time.Now(),
	}, nil
}

func (b *coinpaymentsBridge) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
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

func (b *coinpaymentsBridge) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
