package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/httpx"
	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func newTestNowPaymentsPG(t *testing.T, rt roundTripFunc) *nowPaymentsPG {
	t.Helper()

	if os.Getenv("NOWPAYMENTS_REAL_NETWORK") == "1" {
		cfgPath := os.Getenv("NOWPAYMENTS_CONFIG_PATH")
		if cfgPath == "" {
			cfgPath = "../../.vscode/tmp/config.toml"
		}

		cfg, err := config.Load(cfgPath)
		if err != nil {
			t.Fatalf("failed to load config file: %v", err)
		}

		pgCfg, ok := cfg.Providers["nowpayments"]
		if !ok {
			t.Fatal("nowpayments config was not loaded from file")
		}

		p := &nowPaymentsPlugin{}
		if err := p.Load(pgCfg); err != nil {
			t.Fatalf("Load returned error: %v", err)
		}
		return p.pg
	}

	if rt == nil {
		t.Fatal("roundTripFunc is required for mock mode")
	}

	return &nowPaymentsPG{
		client: httpx.New(httpx.ClientOptions{
			Client:  &http.Client{Transport: rt},
			BaseURL: sandboxBaseURL,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				"x-api-key":    []string{"test-key"},
			},
			ErrorPrefix: "nowpayments",
		}),
		ipnSecret: "test-secret",
	}
}

func TestMain(m *testing.M) {
	_ = logger.Init(config.LogConfig{
		Level:  "error",
		Format: "json",
		Output: config.LogOutputConfig{
			Stdout: true,
		},
		Rotation: &config.LogRotationConfig{},
	})
	code := m.Run()
	_ = logger.Sync()
	os.Exit(code)
}

func TestCreatePaymentLink_RequestsSandboxEndpointWhenModeIsSandbox(t *testing.T) {
	t.Parallel()

	var gotURL string
	var gotMethod string
	var gotBody map[string]any

	p := newTestNowPaymentsPG(t, func(r *http.Request) (*http.Response, error) {
		gotURL = r.URL.String()
		gotMethod = r.Method

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"id":"sandbox_invoice_1",
				"order_id":"SAMPLE",
				"order_description":"DESCRIPTION | ItemName x1 (+1 more)",
				"price_amount":"3.21",
				"price_currency":"usd",
				"pay_currency":"btc",
				"ipn_callback_url":"https://local.draco.city/api/v1/providers/nowpayments/webhooks",
				"invoice_url":"https://sandbox.nowpayments.test/invoice/sandbox_invoice_1",
				"created_at":"2026-03-27T00:00:00Z",
				"updated_at":"2026-03-27T00:01:00Z"
			}`)),
		}, nil
	})

	resp, err := p.CreatePaymentLink(context.Background(), models.CreatePaymentLinkRequest{
		Amount:          "3.21",
		Currency:        "USD",
		ReceiveCurrency: "BTC",
		OrderID:         "SAMPLE",
		Description:     "DESCRIPTION",
		Items: []models.Item{{
			ID:       "ItemID",
			Name:     "ItemName",
			Quantity: 1,
			Amount:   "3.21",
		}, {
			ID:       "ItemID2",
			Name:     "ItemName2",
			Quantity: 1,
			Amount:   "1.00",
		}},
		WebhookURL: "https://local.draco.city/api/v1/providers/nowpayments/webhooks",
	})
	if err != nil {
		t.Fatalf("CreatePaymentLink returned error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("expected method POST, got %s", gotMethod)
	}
	if gotURL != "https://api-sandbox.nowpayments.io/v1/invoice" {
		t.Fatalf("expected sandbox invoice endpoint, got %s", gotURL)
	}
	if gotBody["price_currency"] != "usd" {
		t.Fatalf("expected price_currency usd, got %#v", gotBody["price_currency"])
	}
	if gotBody["pay_currency"] != "btc" {
		t.Fatalf("expected pay_currency btc, got %#v", gotBody["pay_currency"])
	}

	if resp.InvoiceID != "sandbox_invoice_1" {
		t.Fatalf("unexpected invoice id: %s", resp.InvoiceID)
	}
	if resp.Currency != "USD" || resp.ReceiveCurrency != "BTC" {
		t.Fatalf("unexpected response currency mapping: %#v", resp)
	}
}

func TestCreatePayment_RequestsSandboxEndpointWhenModeIsSandbox(t *testing.T) {
	t.Parallel()

	var gotURL string
	var gotMethod string
	var gotBody map[string]any

	p := newTestNowPaymentsPG(t, func(r *http.Request) (*http.Response, error) {
		gotURL = r.URL.String()
		gotMethod = r.Method

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		if err := json.Unmarshal(raw, &gotBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"payment_id":"payment_1",
				"payment_status":"waiting",
				"price_amount":"3.21",
				"price_currency":"usd",
				"pay_amount":"0.00011",
				"amount_received":"0",
				"pay_currency":"btc",
				"order_id":"SAMPLE",
				"order_description":"DESCRIPTION",
				"created_at":"2026-03-27T00:00:00Z",
				"updated_at":"2026-03-27T00:01:00Z",
				"expiration_estimate_date":"2026-03-27T00:20:00Z"
			}`)),
		}, nil
	})

	resp, err := p.CreatePayment(context.Background(), models.CreatePaymentRequest{
		InvoiceID:       "6201138408",
		ReceiveCurrency: "BTC",
		Description:     "DESCRIPTION",
		CustomerEmail:   "test@draco.city",
	})
	if err != nil {
		t.Fatalf("CreatePayment returned error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("expected method POST, got %s", gotMethod)
	}
	if gotURL != "https://api-sandbox.nowpayments.io/v1/invoice-payment" {
		t.Fatalf("expected sandbox invoice-payment endpoint, got %s", gotURL)
	}
	if gotBody["iid"] != "6201138408" {
		t.Fatalf("expected iid 6201138408, got %#v", gotBody["iid"])
	}
	if gotBody["pay_currency"] != "btc" {
		t.Fatalf("expected normalized pay_currency btc, got %#v", gotBody["pay_currency"])
	}
	if gotBody["customer_email"] != "test@draco.city" {
		t.Fatalf("expected customer_email test@draco.city, got %#v", gotBody["customer_email"])
	}

	if resp.PaymentID != "payment_1" {
		t.Fatalf("unexpected payment id: %s", resp.PaymentID)
	}
	if resp.Currency != "USD" || resp.ReceiveCurrency != "BTC" {
		t.Fatalf("unexpected response currency mapping: %#v", resp)
	}
	if resp.EstimatedAmount != "0.00011" {
		t.Fatalf("unexpected estimated amount: %s", resp.EstimatedAmount)
	}
}

func TestGetPayment_RequestsSandboxEndpointWhenModeIsSandbox(t *testing.T) {
	t.Parallel()

	var gotURL string
	var gotMethod string

	p := newTestNowPaymentsPG(t, func(r *http.Request) (*http.Response, error) {
		gotURL = r.URL.String()
		gotMethod = r.Method

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"payment_id":4635619774,
				"payment_status":"finished",
				"price_amount":"3.21",
				"price_currency":"usd",
				"pay_amount":"0.00011",
				"actually_paid":"0.00011",
				"pay_currency":"btc",
				"order_id":"SAMPLE",
				"order_description":"DESCRIPTION",
				"created_at":"2026-03-27T00:00:00Z",
				"updated_at":"2026-03-27T00:01:00Z"
			}`)),
		}, nil
	})

	resp, err := p.GetPayment(context.Background(), "4635619774")
	if err != nil {
		t.Fatalf("GetPayment returned error: %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Fatalf("expected method GET, got %s", gotMethod)
	}
	if gotURL != "https://api-sandbox.nowpayments.io/v1/payment/4635619774" {
		t.Fatalf("expected sandbox payment endpoint, got %s", gotURL)
	}
	if resp.PaymentID != "4635619774" {
		t.Fatalf("unexpected payment id: %s", resp.PaymentID)
	}
	if resp.Currency != "USD" || resp.ReceiveCurrency != "BTC" {
		t.Fatalf("unexpected response currency mapping: %#v", resp)
	}
}

func TestWebhook_ReceiveProcess(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"actually_paid":0,"actually_paid_at_fiat":0,"fee":{"currency":"eth","depositFee":0,"serviceFee":0,"withdrawalFee":0},"invoice_id":6070266182,"order_description":"DESCRIPTION | ItemName x1 (+1 more)","order_id":"SAMPLE","outcome_amount":0.001408,"outcome_currency":"eth","parent_payment_id":null,"pay_address":"MEEim9nyiAWcDbMAqVZj6tE5xktqbSDrCJ","pay_amount":0.05716941,"pay_currency":"ltc","payin_extra_id":null,"payment_extra_ids":null,"payment_id":4492847470,"payment_status":"finished","price_amount":3.21,"price_currency":"usd","purchase_id":"6105448077","updated_at":1774390831880}`)
	p := &nowPaymentsPG{ipnSecret: "test-secret"}

	canonicalPayload, err := canonicalizePayload(payload)
	if err != nil {
		t.Fatalf("canonicalPayload returned error: %v", err)
	}

	digest := hmac.New(sha512.New, []byte(p.ipnSecret))
	digest.Write(canonicalPayload)
	signature := hex.EncodeToString(digest.Sum(nil))

	header := make(http.Header)
	header.Set("x-nowpayments-sig", signature)
	result, err := p.HandleWebhook(context.Background(), payload, header)
	if err != nil {
		t.Fatalf("HandleWebhook returned error: %v", err)
	}
	if result == nil || !result.Accepted || result.Message != "verified" {
		t.Fatalf("unexpected webhook result: %#v", result)
	}
}
