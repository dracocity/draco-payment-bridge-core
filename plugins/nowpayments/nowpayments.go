package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
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
	"github.com/dracocity/draco-payment-bridge-core/plugins/nowpayments/types/request"
	"github.com/dracocity/draco-payment-bridge-core/plugins/nowpayments/types/response"
)

const (
	prodBaseURL    = "https://api.nowpayments.io/v1"
	sandboxBaseURL = "https://api-sandbox.nowpayments.io/v1"
)

type nowPaymentsPlugin struct {
	pg *nowPaymentsPG
}

type nowPaymentsPG struct {
	client    *httpx.Client
	ipnSecret string
}

func New() plugin.Plugin {
	return &nowPaymentsPlugin{}
}

func (p *nowPaymentsPlugin) Load(cfg config.PGConfig) error {
	var baseURL string
	switch strings.ToLower(cfg["mode"]) {
	case "sandbox":
		baseURL = sandboxBaseURL
	case "production":
		baseURL = prodBaseURL
	default:
		return fmt.Errorf("invalid mode: %s", cfg["mode"])
	}
	apiKey := strings.TrimSpace(cfg["api_key"])
	if apiKey == "" {
		return errors.New("nowpayments api_key is required")
	}
	ipnSecret := strings.TrimSpace(cfg["ipn_secret"])
	if ipnSecret == "" {
		return errors.New("nowpayments ipn_secret is required")
	}
	timeout := 7 * time.Second
	if timeoutStr := strings.TrimSpace(cfg["http_timeout"]); timeoutStr != "" {
		parsedTimeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid http_timeout: %w", err)
		}
		timeout = parsedTimeout
	}
	p.pg = &nowPaymentsPG{
		client: httpx.New(httpx.ClientOptions{
			Client:  &http.Client{Timeout: timeout},
			BaseURL: baseURL,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				"x-api-key":    []string{apiKey},
			},
			ErrorPrefix: "nowpayments",
		}),
		ipnSecret: ipnSecret,
	}
	return nil
}

func (p *nowPaymentsPlugin) Unload() error {
	return nil
}

func (p *nowPaymentsPlugin) Name() string {
	return "nowpayments"
}

func (p *nowPaymentsPlugin) PaymentGateway() pg.PaymentGateway {
	return p.pg
}

func (b *nowPaymentsPG) Name() string {
	return "nowpayments"
}

func (b *nowPaymentsPG) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	body, err := request.BuildCreateInvoice(req)
	if err != nil {
		return nil, err
	}

	var resp response.CreateInvoice
	if err := b.client.Post(ctx, "/invoice", nil, body, &resp); err != nil {
		return nil, err
	}

	return response.BuildCreatePaymentLink(resp)
}

func (b *nowPaymentsPG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	body, err := request.BuildCreateInvoicePayment(req)
	if err != nil {
		return nil, err
	}

	var resp response.CreateInvoicePayment
	if err := b.client.Post(ctx, "/invoice-payment", nil, body, &resp); err != nil {
		return nil, err
	}

	return response.BuildCreatePayment(resp)
}

func (b *nowPaymentsPG) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
	var resp response.GetPayment
	if err := b.client.Get(ctx, "/payment/"+paymentID, nil, nil, &resp); err != nil {
		return nil, err
	}

	return response.BuildGetPayment(resp)
}

func (b *nowPaymentsPG) CreateRefund(ctx context.Context, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
	return nil, errors.New("nowpayments does not provide an API refund endpoint")
}

func (b *nowPaymentsPG) HandleWebhook(ctx context.Context, payload []byte, header http.Header) (*models.WebhookResult, error) {
	signature := header.Get("x-nowpayments-sig")
	if signature == "" {
		return nil, errors.New("missing x-nowpayments-sig header")
	}

	canonicalPayload, err := canonicalizePayload(payload)
	if err != nil {
		return nil, err
	}

	digest := hmac.New(sha512.New, []byte(b.ipnSecret))
	digest.Write(canonicalPayload)
	expected := hex.EncodeToString(digest.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(strings.ToLower(signature))) {
		return nil, errors.New("invalid nowpayments signature")
	}

	return &models.WebhookResult{
		Accepted: true,
		Message:  "verified",
	}, nil
}

func canonicalizePayload(payload []byte) ([]byte, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, fmt.Errorf("invalid webhook payload: %w", err)
	}

	canonicalPayload, err := marshalNoEscape(params)
	if err != nil {
		return nil, err
	}

	return canonicalPayload, nil
}

func marshalNoEscape(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		return nil, err
	}

	encoded := buf.Bytes()
	if n := len(encoded); n > 0 && encoded[n-1] == '\n' {
		encoded = encoded[:n-1]
	}
	return encoded, nil
}
