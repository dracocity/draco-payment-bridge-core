package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/pkg/crypto"
)

func BuildCreateOrder(req models.CreatePaymentLinkRequest, callbackSecret string) (*CreateOrder, error) {
	orderID := strings.TrimSpace(req.OrderID)

	priceAmount, err := strconv.ParseFloat(strings.TrimSpace(req.Amount), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	r := &CreateOrder{
		OrderID:         orderID,
		PriceAmount:     priceAmount,
		PriceCurrency:   strings.ToUpper(strings.TrimSpace(req.Currency)),
		ReceiveCurrency: strings.ToUpper(strings.TrimSpace(req.ReceiveCurrency)),
		Title:           strings.TrimSpace("Order - " + orderID),
		Description:     strings.TrimSpace(req.Description),
		CallbackURL:     strings.TrimSpace(req.WebhookURL),
		CancelURL:       strings.TrimSpace(req.CancelURL),
		SuccessURL:      strings.TrimSpace(req.SuccessURL),
	}

	if len(req.ProviderPayload) > 0 {
		var payload CreateOrderPayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}
		if payload.Shopper != nil {
			s := *payload.Shopper
			r.Shopper = &s

			if payload.Shopper.CompanyDetails != nil {
				scd := *payload.Shopper.CompanyDetails
				r.Shopper.CompanyDetails = &scd
			}
		}
	}

	signature, err := crypto.GenerateSignature(r, callbackSecret)
	if err != nil {
		return nil, err
	}
	r.Token = signature

	return r, nil
}
