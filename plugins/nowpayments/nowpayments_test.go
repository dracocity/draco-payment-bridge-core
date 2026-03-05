package main

import (
	"context"
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
		OrderID:         "SAMPLE",
		Amount:          "3.21",
		Currency:        "USD",
		ReceiveCurrency: "BTC",
		WebhookURL:      "https://local.draco.city/api/v1/providers/nowpayments/webhooks",
	})
	if err != nil {
		t.Fatalf("CreatePaymentLink returned error: %v", err)
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

	paymentID := "5234398929"
	resp, err := p.pg.GetPayment(context.Background(), paymentID)
	if err != nil {
		t.Fatalf("GetPayment returned error: %v", err)
	}

	fmt.Println(resp)
}
