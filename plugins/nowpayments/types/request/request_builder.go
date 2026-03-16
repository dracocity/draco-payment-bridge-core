package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func BuildCreateInvoice(req models.CreatePaymentLinkRequest) (*CreateInvoice, error) {
	priceAmount, err := strconv.ParseFloat(strings.TrimSpace(req.Amount), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid fiat_amount: %w", err)
	}

	orderDescription := strings.TrimSpace(req.Description)
	if n := len(req.Items); n > 0 {
		item := req.Items[0]
		if orderDescription != "" {
			orderDescription += " | "
		}
		orderDescription += fmt.Sprintf("%s x%d", item.Name, item.Quantity)
		if n > 1 {
			orderDescription += fmt.Sprintf(" (+%d more)", n-1)
		}
	}

	r := &CreateInvoice{
		PriceAmount:      priceAmount,
		PriceCurrency:    strings.ToLower(strings.TrimSpace(req.Currency)),
		PayCurrency:      strings.ToLower(strings.TrimSpace(req.ReceiveCurrency)),
		IPNCallbackURL:   strings.TrimSpace(req.WebhookURL),
		OrderID:          strings.TrimSpace(req.OrderID),
		OrderDescription: orderDescription,
		SuccessURL:       strings.TrimSpace(req.SuccessURL),
		CancelURL:        strings.TrimSpace(req.CancelURL),
	}

	// TODO:
	// r.Items = make([]CreateInvoiceItem, 0, len(req.Items))
	// for _, item := range req.Items {
	// 	})
	// }

	if len(req.ProviderPayload) > 0 {
		var payload CreateInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}
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
