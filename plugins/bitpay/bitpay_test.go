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

	pgCfg, ok := cfg.Providers["bitpay"]
	if !ok {
		t.Fatal("bitpay config was not loaded from file")
	}

	p := &bitpayPlugin{}
	if err := p.Load(pgCfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	resp, err := p.pg.CreatePaymentLink(context.Background(), models.CreatePaymentLinkRequest{
		OrderID:      "SAMPLE",
		FiatAmount:   "3.21",
		FiatCurrency: "USD",
		WebhookURL:   "https://local.draco.city/api/v1/providers/bitpay/webhooks",
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

	pgCfg, ok := cfg.Providers["bitpay"]
	if !ok {
		t.Fatal("bitpay config was not loaded from file")
	}

	p := &bitpayPlugin{}
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
