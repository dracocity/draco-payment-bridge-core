package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/httpx"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
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
	client     *httpx.Client
	baseURL    string
	apiToken   string
	apiPrivKey *secp256k1.PrivateKey
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
	var apiPrivKey *secp256k1.PrivateKey
	if privKeyStr := strings.TrimSpace(cfg["api_private_key"]); privKeyStr == "" {
		return errors.New("bitpay api_private_key is required")
	} else if privKeyBytes, err := hex.DecodeString(privKeyStr); err != nil {
		return fmt.Errorf("invalid api_private_key: %w", err)
	} else {
		apiPrivKey = secp256k1.PrivKeyFromBytes(privKeyBytes)
	}
	apiPubKeyHex := hex.EncodeToString(apiPrivKey.PubKey().SerializeCompressed())
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
				"X-Identity":       []string{apiPubKeyHex},
			},
			ErrorPrefix: "bitpay",
		}),
		baseURL:    baseURL,
		apiToken:   apiToken,
		apiPrivKey: apiPrivKey,
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
	header, err := request.BuildRequestHeader(b.apiPrivKey, b.baseURL, "/invoices", nil, body)
	if err != nil {
		return nil, err
	}

	var resp response.CreateInvoice
	if err := b.client.Post(ctx, "/invoices", header, body, &resp); err != nil {
		return nil, err
	}

	return response.BuildCreatePaymentLink(resp, req.WebhookURL)
}

func (b *bitpayPG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	return nil, errors.New("be using this feature on my self-hosted page(bitpay)")
}

func (b *bitpayPG) GetPayment(ctx context.Context, invoiceID string) (*models.GetPaymentResponse, error) {
	header, err := request.BuildRequestHeader(b.apiPrivKey, b.baseURL, "/invoices/"+invoiceID, url.Values{"token": []string{b.apiToken}}, nil)
	if err != nil {
		return nil, err
	}
	var resp response.RetrieveInvoice
	if err := b.client.Get(ctx, "/invoices/"+invoiceID, url.Values{"token": []string{b.apiToken}}, header, &resp); err != nil {
		return nil, err
	}
	return response.BuildGetPayment(resp)
}

func (b *bitpayPG) CreateRefund(ctx context.Context, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
	body := request.BuildCreateRefundRequest(req, b.apiToken)
	header, err := request.BuildRequestHeader(b.apiPrivKey, b.baseURL, "/refunds", nil, body)
	if err != nil {
		return nil, err
	}

	var resp response.CreateRefundRequest
	if err := b.client.Post(ctx, "/refunds", header, body, &resp); err != nil {
		return nil, err
	}

	return response.BuildCreateRefund(resp, req.PaymentID)
}

func (b *bitpayPG) HandleWebhook(ctx context.Context, payload []byte, header http.Header) (*models.WebhookResult, error) {
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook received",
	}, nil
}
