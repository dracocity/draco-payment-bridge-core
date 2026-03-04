package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func BuildCreateInvoice(req models.CreatePaymentLinkRequest) (*CreateInvoice, error) {
	priceAmount, err := strconv.ParseFloat(strings.TrimSpace(req.FiatAmount), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid fiat_amount: %w", err)
	}

	r := &CreateInvoice{
		PriceAmount:      priceAmount,
		PriceCurrency:    strings.ToLower(strings.TrimSpace(req.FiatCurrency)),
		OrderID:          strings.TrimSpace(req.OrderID),
		OrderDescription: strings.TrimSpace(req.Description),
		IPNCallbackURL:   strings.TrimSpace(req.WebhookURL),
		SuccessURL:       strings.TrimSpace(req.SuccessURL),
		CancelURL:        strings.TrimSpace(req.CancelURL),
	}

	if len(req.ProviderPayload) > 0 {
		var payload CreateInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}
		r.PayCurrency = payload.PayCurrency
		r.IsFixedRate = payload.IsFixedRate
		r.IsFeePaidByUser = payload.IsFeePaidByUser
	}

	return r, nil
}

func BuildCreateInvoicePayment(req models.CreatePaymentRequest) (*CreateInvoicePayment, error) {

	r := &CreateInvoicePayment{
		InvoiceID: req.InvoiceID,
	}

	if len(req.ProviderPayload) > 0 {
		var payload CreateInvoicePaymentPayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}
	}

	return r, nil
}
