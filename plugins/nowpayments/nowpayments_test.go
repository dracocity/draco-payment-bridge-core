package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestCreatePaymentLink_RequestsSandboxEndpointWhenModeIsSandbox(t *testing.T) {
	t.Parallel()

	cfgPath := "../../.vscode/tmp/config.toml"
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load config file: %v", err)
	}

	if err := logger.Init(cfg.Log); err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	pgCfg, ok := cfg.Providers["nowpayments"]
	if !ok {
		t.Fatal("nowpayments config was not loaded from file")
	}

	p := &nowPaymentsPlugin{}
	if err := p.Load(pgCfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	// var gotURL string
	// p.pg.client = &http.Client{
	// 	Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
	// 		gotURL = r.URL.String()
	// 		return &http.Response{
	// 			StatusCode: http.StatusOK,
	// 			Header:     make(http.Header),
	// 			Body: io.NopCloser(strings.NewReader(`{
	// 				"payment_id":"sandbox_1",
	// 				"price_amount":"3.21",
	// 				"price_currency":"usd",
	// 				"payment_status":"waiting"
	// 			}`)),
	// 		}, nil
	// 	}),
	// }

	resp, err := p.pg.CreatePaymentLink(context.Background(), models.CreatePaymentLinkRequest{
		Amount:      "3.21",
		Currency:    "USD",
		OrderID:     "SAMPLE",
		Description: "DESCRIPTION",
		Items: []models.Item{{
			ID:       "ItemID",
			Name:     "ItemName",
			Quantity: 1,
			Amount:   "3.21",
		}, {
			ID:       "ItemID",
			Name:     "ItemName",
			Quantity: 1,
			Amount:   "3.21",
		}},
		WebhookURL: "https://local.draco.city/api/v1/providers/nowpayments/webhooks",
	})
	if err != nil {
		t.Fatalf("CreatePaymentLink returned error: %v", err)
	}

	fmt.Println(resp)
}

func TestCreatePayment_RequestsSandboxEndpointWhenModeIsSandbox(t *testing.T) {
	t.Parallel()

	cfgPath := "../../.vscode/tmp/config.toml"
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load config file: %v", err)
	}

	if err := logger.Init(cfg.Log); err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	pgCfg, ok := cfg.Providers["nowpayments"]
	if !ok {
		t.Fatal("nowpayments config was not loaded from file")
	}

	p := &nowPaymentsPlugin{}
	if err := p.Load(pgCfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	resp, err := p.pg.CreatePayment(context.Background(), models.CreatePaymentRequest{
		InvoiceID:       "6201138408",
		ReceiveCurrency: "BTC",
		Description:     "DESCRIPTION",
		CustomerEmail:   "test@draco.city",
	})
	if err != nil {
		t.Fatalf("CreatePayment returned error: %v", err)
	}

	fmt.Println(resp)
}

func TestGetPayment_RequestsSandboxEndpointWhenModeIsSandbox(t *testing.T) {
	t.Parallel()

	cfgPath := "../../.vscode/tmp/config.toml"
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load config file: %v", err)
	}

	if err := logger.Init(cfg.Log); err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	pgCfg, ok := cfg.Providers["nowpayments"]
	if !ok {
		t.Fatal("nowpayments config was not loaded from file")
	}

	p := &nowPaymentsPlugin{}
	if err := p.Load(pgCfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	paymentID := "4635619774"
	resp, err := p.pg.GetPayment(context.Background(), paymentID)
	if err != nil {
		t.Fatalf("GetPayment returned error: %v", err)
	}

	fmt.Println(resp)
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
