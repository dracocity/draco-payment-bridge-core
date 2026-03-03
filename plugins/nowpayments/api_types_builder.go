package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func buildCreateInvoiceRequest(req models.CreatePaymentLinkRequest) (*createInvoiceRequest, error) {
	priceAmount, err := strconv.ParseFloat(strings.TrimSpace(req.FiatAmount), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid fiat_amount: %w", err)
	}

	r := &createInvoiceRequest{
		PriceAmount:      priceAmount,
		PriceCurrency:    strings.ToLower(strings.TrimSpace(req.FiatCurrency)),
		OrderID:          strings.TrimSpace(req.OrderID),
		OrderDescription: strings.TrimSpace(req.Description),
		IPNCallbackURL:   strings.TrimSpace(req.WebhookURL),
		SuccessURL:       strings.TrimSpace(req.SuccessURL),
		CancelURL:        strings.TrimSpace(req.CancelURL),
	}

	if len(req.ProviderPayload) > 0 {
		var payload createInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}
		r.PayCurrency = payload.PayCurrency
		r.IsFixedRate = payload.IsFixedRate
		r.IsFeePaidByUser = payload.IsFeePaidByUser
	}

	return r, nil
}

func buildCreateInvoicePaymentRequest(req models.CreatePaymentRequest) (*createInvoicePaymentRequest, error) {

	r := &createInvoicePaymentRequest{
		InvoiceID: req.InvoiceID,
	}

	if len(req.ProviderPayload) > 0 {
		var payload createInvoicePaymentPayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}
	}

	return r, nil
}
