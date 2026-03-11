package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/httpx"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
	"github.com/dracocity/draco-payment-bridge-core/plugins/coinpayments/types/request"
	"github.com/dracocity/draco-payment-bridge-core/plugins/coinpayments/types/response"
)

const (
	aBaseURL = "https://a-api.coinpayments.net/api/v2/merchant"
	bBaseURL = "https://b-api.coinpayments.net/api/v2/merchant"
)

type coinpaymentsPlugin struct {
	pg *coinpaymentsPG
}

type coinpaymentsPG struct {
	client       *httpx.Client
	clientID     string
	clientSecret string
}

func New() plugin.Plugin {
	return &coinpaymentsPlugin{}
}

func (p *coinpaymentsPlugin) Load(cfg config.PGConfig) error {
	var baseURL string
	switch strings.ToLower(cfg["api_type"]) {
	case "a-api":
		baseURL = aBaseURL
	case "b-api":
		baseURL = bBaseURL
	default:
		return fmt.Errorf("invalid api_type: %s", cfg["api_type"])
	}
	clientID := strings.TrimSpace(cfg["client_id"])
	if clientID == "" {
		return errors.New("coinpayments client_id is required")
	}
	clientSecret := strings.TrimSpace(cfg["client_secret"])
	if clientSecret == "" {
		return errors.New("coinpayments client_secret is required")
	}
	timeout := 7 * time.Second
	if timeoutStr := strings.TrimSpace(cfg["http_timeout"]); timeoutStr != "" {
		parsedTimeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid http_timeout: %w", err)
		}
		timeout = parsedTimeout
	}
	p.pg = &coinpaymentsPG{
		client: httpx.New(httpx.ClientOptions{
			Client:  &http.Client{Timeout: timeout},
			BaseURL: baseURL,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			ErrorPrefix: "coinpayments",
		}),
		clientID:     clientID,
		clientSecret: clientSecret,
	}
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

// TODO:
func (b *coinpaymentsPG) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	body, err := request.BuildCreateInvoice(req)
	if err != nil {
		return nil, err
	}
	header, err := request.BuildRequestHeaders(b.clientID, b.clientSecret, "POST", "", body)
	if err != nil {
		return nil, err
	}

	var resp response.CreateInvoice
	if err := b.client.Post(ctx, "/invoices", header, body, &resp); err != nil {
		return nil, err
	}

	return &models.CreatePaymentLinkResponse{
		// InvoiceID:    resp.Data.ID,
		// OrderID:      resp.Data.OrderID,
		// FiatAmount:   resp.Data.DisplayAmountPaid,
		// FiatCurrency: resp.Data.Currency,
		// CheckoutURL:  resp.Data.RedirectURL,
		// CreatedAt:    resp.Data.InvoiceTime,
		// UpdatedAt:    resp.Data.InvoiceTime,
		Raw: resp,
	}, nil
}

func (b *coinpaymentsPG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	return &models.CreatePaymentResponse{}, nil
}

func (b *coinpaymentsPG) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
	return &models.GetPaymentResponse{}, nil
}

func (b *coinpaymentsPG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
	return &models.RefundResponse{}, nil
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
