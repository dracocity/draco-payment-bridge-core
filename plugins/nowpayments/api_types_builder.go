package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func buildNowPaymentsCreateInvoiceRequest(req models.CreatePaymentLinkRequest) (nowPaymentsCreateInvoiceRequest, error) {
	priceAmount, err := strconv.ParseFloat(strings.TrimSpace(req.FiatAmount), 64)
	if err != nil {
		return nowPaymentsCreateInvoiceRequest{}, fmt.Errorf("invalid amount: %w", err)
	}

	orderID := strings.TrimSpace(req.OrderID)
	r := nowPaymentsCreateInvoiceRequest{
		OrderID:       &orderID,
		PriceAmount:   priceAmount,
		PriceCurrency: strings.ToLower(strings.TrimSpace(req.FiatCurrency)),
	}

	if req.Description != nil {
		if v := strings.TrimSpace(*req.Description); v != "" {
			r.OrderDescription = &v
		}
	}
	if req.WebhookURL != nil {
		if v := strings.TrimSpace(*req.WebhookURL); v != "" {
			r.IPNCallbackURL = &v
		}
	}
	if req.SuccessURL != nil {
		if v := strings.TrimSpace(*req.SuccessURL); v != "" {
			r.SuccessURL = &v
		}
	}
	if req.CancelURL != nil {
		if v := strings.TrimSpace(*req.CancelURL); v != "" {
			r.CancelURL = &v
		}
	}

	if len(req.ProviderPayload) > 0 {
		var payload nowPaymentsCreateInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nowPaymentsCreateInvoiceRequest{}, fmt.Errorf("invalid provider_payload: %w", err)
		}
		r.PayCurrency = payload.PayCurrency
		r.IsFixedRate = payload.IsFixedRate
		r.IsFeePaidByUser = payload.IsFeePaidByUser
	}

	return r, nil
}

func buildNowPaymentsCreateInvoicePaymentRequest(req models.CreatePaymentRequest) (nowPaymentsCreateInvoicePaymentRequest, error) {

	r := nowPaymentsCreateInvoicePaymentRequest{
		InvoiceID: req.InvoiceID,
	}

	if len(req.ProviderPayload) > 0 {
		var payload nowPaymentsCreateInvoicePaymentPayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nowPaymentsCreateInvoicePaymentRequest{}, fmt.Errorf("invalid provider_payload: %w", err)
		}
	}

	return r, nil
}
