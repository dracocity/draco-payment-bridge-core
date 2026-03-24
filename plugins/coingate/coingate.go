package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/consts"
	"github.com/dracocity/draco-payment-bridge-core/internal/httpx"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
	"github.com/dracocity/draco-payment-bridge-core/internal/types"
	"github.com/dracocity/draco-payment-bridge-core/plugins/coingate/types/request"
	"github.com/dracocity/draco-payment-bridge-core/plugins/coingate/types/response"
)

const (
	prodBaseURL    = "https://api.coingate.com/v2"
	sandboxBaseURL = "https://api-sandbox.coingate.com/v2"
)

type coingatePlugin struct {
	pg *coingatePG
}

type coingatePG struct {
	client         *httpx.Client
	callbackSecret string
}

func New() plugin.Plugin {
	return &coingatePlugin{}
}

func (p *coingatePlugin) Load(cfg config.PGConfig) error {
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
		return errors.New("coingate api_token is required")
	}
	callbackSecret := strings.TrimSpace(cfg["callback_secret"])
	if callbackSecret == "" {
		return errors.New("coingate callback_secret is required")
	}
	timeout := 7 * time.Second
	if timeoutStr := strings.TrimSpace(cfg["http_timeout"]); timeoutStr != "" {
		parsedTimeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid http_timeout: %w", err)
		}
		timeout = parsedTimeout
	}

	p.pg = &coingatePG{
		client: httpx.New(httpx.ClientOptions{
			Client:  &http.Client{Timeout: timeout},
			BaseURL: baseURL,
			Header: http.Header{
				"Content-Type":  []string{"application/json"},
				"Accept":        []string{"application/json"},
				"Authorization": []string{"Token " + apiToken},
			},
			ErrorPrefix: "coingate",
		}),
		callbackSecret: callbackSecret,
	}
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
	body, err := request.BuildCreateOrder(req, b.callbackSecret)
	if err != nil {
		return nil, err
	}

	var resp response.CreateOrder
	if err := b.client.Post(ctx, "/orders", nil, body, &resp); err != nil {
		return nil, err
	}

	return response.BuildCreatePaymentLink(resp)
}

func (b *coingatePG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	return nil, errors.New("be using this feature on my self-hosted page(coingate)")
}

func (b *coingatePG) GetPayment(ctx context.Context, orderID string) (*models.GetPaymentResponse, error) {
	var resp response.GetOrder
	if err := b.client.Get(ctx, "/orders/"+orderID, nil, nil, &resp); err != nil {
		return nil, err
	}

	amount := strings.TrimSpace(resp.PayAmount)
	if amount == "" {
		amount = strings.TrimSpace(resp.PriceAmount)
	}
	currency := strings.ToUpper(strings.TrimSpace(resp.PayCurrency))
	if currency == "" {
		currency = strings.ToUpper(strings.TrimSpace(resp.PriceCurrency))
	}

	return &models.GetPaymentResponse{
		Status:    normalizeCoinGateStatus(resp.Status),
		PaymentID: strconv.FormatInt(resp.ID, 10),
		OrderID:   strings.TrimSpace(resp.OrderID),
		Currency:  currency,
		Amount:    amount,
		CreatedAt: toUnixMilli(resp.CreatedAt),
		Raw:       resp,
	}, nil
}

func (b *coingatePG) CreateRefund(ctx context.Context, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
	return &models.CreateRefundResponse{}, nil
}

func (b *coingatePG) HandleWebhook(ctx context.Context, payload []byte, header http.Header) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
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

func normalizeCoinGateStatus(status string) types.PaymentStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "new":
		return consts.StatusCreated
	case "pending":
		return consts.StatusPending
	case "confirming":
		return consts.StatusConfirming
	case "paid":
		return consts.StatusCompleted
	case "partially_paid":
		return consts.StatusPartiallyPaid
	case "invalid":
		return consts.StatusFailed
	case "expired":
		return consts.StatusExpired
	default:
		normalized := strings.ToUpper(strings.TrimSpace(status))
		if normalized == "" {
			return consts.StatusPending
		}
		return types.PaymentStatus(normalized)
	}
}
