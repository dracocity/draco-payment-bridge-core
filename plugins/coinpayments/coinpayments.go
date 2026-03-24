package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/consts"
	"github.com/dracocity/draco-payment-bridge-core/internal/httpx"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
	"github.com/dracocity/draco-payment-bridge-core/internal/types"
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
	baseURL      string
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
				"Content-Type":          []string{"application/json"},
				"X-CoinPayments-Client": []string{clientID},
			},
			ErrorPrefix: "coinpayments",
		}),
		baseURL:      baseURL,
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

func (b *coinpaymentsPG) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	body, err := request.BuildCreateInvoice(req)
	if err != nil {
		return nil, err
	}
	header, err := request.BuildRequestHeader(b.clientID, b.clientSecret, http.MethodPost, b.baseURL, "/invoices", nil, body)
	if err != nil {
		return nil, err
	}

	var resp response.CreateInvoice
	if err := b.client.Post(ctx, "/invoices", header, body, &resp); err != nil {
		return nil, err
	}

	return response.BuildCreatePaymentLink(resp, req.Amount, req.Currency, req.OrderID)
}

func (b *coinpaymentsPG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	return nil, errors.New("be using this feature on my self-hosted page(coinpayments)")
}

func (b *coinpaymentsPG) GetPayment(ctx context.Context, invoiceID string) (*models.GetPaymentResponse, error) {
	query := url.Values{"include_full_details": []string{"false"}}
	header, err := request.BuildRequestHeader(b.clientID, b.clientSecret, http.MethodGet, b.baseURL, "/invoices/"+invoiceID, query, nil)
	if err != nil {
		return nil, err
	}
	var resp response.GetInvoice
	if err := b.client.Get(ctx, "/invoices/"+invoiceID, query, header, &resp); err != nil {
		return nil, err
	}
	amount := ""
	if resp.Amount != nil {
		amount = strings.TrimSpace(resp.Amount.Total)
	}

	currency := ""
	if resp.Currency != nil {
		currency = strings.ToUpper(strings.TrimSpace(resp.Currency.Symbol))
	}

	return &models.GetPaymentResponse{
		Status:    normalizeCoinPaymentsStatus(resp.Status),
		PaymentID: strings.TrimSpace(resp.ID),
		OrderID:   strings.TrimSpace(resp.InvoiceID),
		Currency:  currency,
		Amount:    amount,
		CreatedAt: toUnixMilli(resp.Created),
		Raw:       resp,
	}, nil
}

func (b *coinpaymentsPG) CreateRefund(ctx context.Context, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
	return &models.CreateRefundResponse{}, nil
}

func (b *coinpaymentsPG) HandleWebhook(ctx context.Context, payload []byte, header http.Header) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

func toUnixMilli(datetime string) int64 {
	if datetime == "" {
		return 0
	}

	if ts, err := time.Parse(time.RFC3339, datetime); err == nil {
		return ts.UnixMilli()
	}

	return 0
}

func normalizeCoinPaymentsStatus(status string) types.PaymentStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "new":
		return consts.StatusCreated
	case "pending":
		return consts.StatusPending
	case "paid":
		return consts.StatusPending
	case "completed":
		return consts.StatusCompleted
	case "cancelled":
		return consts.StatusFailed
	case "timedout":
		return consts.StatusExpired
	default:
		normalized := strings.ToUpper(strings.TrimSpace(status))
		if normalized == "" {
			return consts.StatusPending
		}
		return types.PaymentStatus(normalized)
	}
}
