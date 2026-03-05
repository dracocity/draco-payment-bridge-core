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
	"github.com/dracocity/draco-payment-bridge-core/internal/httpx"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
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
	client        *httpx.Client
	callbackToken string
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
	callbackToken := strings.TrimSpace(cfg["callback_token"])
	if callbackToken == "" {
		return errors.New("coingate callback_token is required")
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
		callbackToken: callbackToken,
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
	body, err := request.BuildCreateOrder(req, b.callbackToken)
	if err != nil {
		return nil, err
	}

	var resp response.CreateOrder
	if err := b.client.Post(ctx, "/orders", nil, body, &resp); err != nil {
		return nil, err
	}

	orderID := req.OrderID
	if resp.OrderID != "" {
		orderID = resp.OrderID
	}

	createdAt := toUnixMilli(resp.CreatedAt)
	return &models.CreatePaymentLinkResponse{
		InvoiceID:   strconv.FormatInt(resp.ID, 10),
		OrderID:     orderID,
		Amount:      resp.PriceAmount,
		Currency:    strings.ToUpper(resp.PriceCurrency),
		CheckoutURL: resp.PaymentURL,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
		Raw:         resp,
	}, nil
}

func (b *coingatePG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	return &models.CreatePaymentResponse{}, nil
}

func (b *coingatePG) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
	return &models.GetPaymentResponse{}, nil
}

func (b *coingatePG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
	return &models.RefundResponse{}, nil
}

func (b *coingatePG) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
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
