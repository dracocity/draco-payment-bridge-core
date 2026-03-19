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
	"github.com/dracocity/draco-payment-bridge-core/plugins/bitpay/types/request"
	"github.com/dracocity/draco-payment-bridge-core/plugins/bitpay/types/response"
)

const (
	prodBaseURL    = "https://bitpay.com"
	sandboxBaseURL = "https://test.bitpay.com"
)

type bitpayPlugin struct {
	pg *bitpayPG
}

type bitpayPG struct {
	client   *httpx.Client
	apiToken string
}

func New() plugin.Plugin {
	return &bitpayPlugin{}
}

func (p *bitpayPlugin) Load(cfg config.PGConfig) error {
	var baseURL string
	switch strings.ToLower(cfg["mode"]) {
	case "sandbox":
		baseURL = sandboxBaseURL
	case "production":
		baseURL = prodBaseURL
	default:
		return fmt.Errorf("invalid mode: %s", cfg["mode"])
	}
	apiToken := strings.TrimSpace(cfg["api_token"])
	if apiToken == "" {
		return errors.New("bitpay api_token is required")
	}
	timeout := 7 * time.Second
	if timeoutStr := strings.TrimSpace(cfg["http_timeout"]); timeoutStr != "" {
		parsedTimeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid http_timeout: %w", err)
		}
		timeout = parsedTimeout
	}
	p.pg = &bitpayPG{
		client: httpx.New(httpx.ClientOptions{
			Client:  &http.Client{Timeout: timeout},
			BaseURL: baseURL,
			Header: http.Header{
				"Content-Type":     []string{"application/json"},
				"X-Accept-Version": []string{"2.0.0"},
				"Accept":           []string{"application/json"},
			},
			ErrorPrefix: "bitpay",
		}),
		apiToken: apiToken,
	}
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
	body, err := request.BuildCreateInvoice(req, b.apiToken)
	if err != nil {
		return nil, err
	}

	var resp response.CreateInvoice
	if err := b.client.Post(ctx, "/invoices", nil, body, &resp); err != nil {
		return nil, err
	}

	return response.BuildCreatePaymentLink(resp, req.WebhookURL)
}

func (b *bitpayPG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	return nil, errors.New("be using this feature on my self-hosted page(bitpay)")
}

func (b *bitpayPG) GetPayment(ctx context.Context, invoiceID string) (*models.GetPaymentResponse, error) {
	var resp response.RetrieveInvoice
	if err := b.client.Get(ctx, "/invoices/"+invoiceID, url.Values{"token": []string{b.apiToken}}, nil, &resp); err != nil {
		return nil, err
	}
	return &models.GetPaymentResponse{
		Status:         normalizeBitPayStatus(resp.Data.Status),
		ProviderStatus: strings.TrimSpace(resp.Data.Status),
		PaymentID:      strings.TrimSpace(resp.Data.ID),
		OrderID:        strings.TrimSpace(resp.Data.OrderID),
		Currency:       strings.ToUpper(strings.TrimSpace(resp.Data.Currency)),
		Amount:         resp.Data.Price.String(),
		CreatedAt:      resp.Data.InvoiceTime,
		Raw:            resp,
	}, nil
}

func (b *bitpayPG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
	return &models.RefundResponse{}, nil
}

func (b *bitpayPG) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}

func normalizeBitPayStatus(status string) types.PaymentStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "new":
		return consts.StatusCreated
	case "paid":
		return consts.StatusPending
	case "confirmed":
		return consts.StatusConfirming
	case "complete":
		return consts.StatusCompleted
	case "expired":
		return consts.StatusExpired
	case "invalid":
		return consts.StatusFailed
	default:
		normalized := strings.ToUpper(strings.TrimSpace(status))
		if normalized == "" {
			return consts.StatusPending
		}
		return types.PaymentStatus(normalized)
	}
}
