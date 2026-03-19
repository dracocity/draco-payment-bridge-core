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

	pgCfg, ok := cfg.Providers["coinpayments"]
	if !ok {
		t.Fatal("coinpayments config was not loaded from file")
	}

	p := &coinpaymentsPlugin{}
	if err := p.Load(pgCfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	resp, err := p.pg.CreatePaymentLink(context.Background(), models.CreatePaymentLinkRequest{
		Amount:   "3.21",
		Currency: "USD",
		OrderID:  "SAMPLE",
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
		BuyerEmail: "test@draco.com",
		WebhookURL: "https://local.draco.city/api/v1/providers/coinpayments/webhooks",
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

	pgCfg, ok := cfg.Providers["coinpayments"]
	if !ok {
		t.Fatal("coinpayments config was not loaded from file")
	}

	p := &coinpaymentsPlugin{}
	if err := p.Load(pgCfg); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	invoiceID := "1585dcc0-fa23-4dbb-a992-b8e00b81a4e6"
	resp, err := p.pg.GetPayment(context.Background(), invoiceID)
	if err != nil {
		t.Fatalf("GetPayment returned error: %v", err)
	}

	fmt.Println(resp)
	// GEThttps://a-api.coinpayments.net/api/v2/merchant/invoices/1585dcc0-fa23-4dbb-a992-b8e00b81a4e6?include_full_details=false18c25a000c854730bd3627840eb72af82026-03-18T20:44:15
}
