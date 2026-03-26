package request

import (
	"crypto/sha256"
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

	description := strings.TrimSpace(req.Description)
	if n := len(req.Items); n > 0 {
		item := req.Items[0]
		if description != "" {
			description += " | "
		}
		description += fmt.Sprintf("%s x%d", item.Name, item.Quantity)
		if n > 1 {
			description += fmt.Sprintf(" (+%d more)", n-1)
		}
	}

	r := &CreateOrder{
		OrderID:         orderID,
		PriceAmount:     priceAmount,
		PriceCurrency:   strings.ToUpper(strings.TrimSpace(req.Currency)),
		ReceiveCurrency: strings.ToUpper(strings.TrimSpace(req.ReceiveCurrency)),
		Title:           strings.TrimSpace("Order " + orderID),
		Description:     description,
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
			r.Shopper.Email = req.CustomerEmail

			if payload.Shopper.CompanyDetails != nil {
				scd := *payload.Shopper.CompanyDetails
				r.Shopper.CompanyDetails = &scd
			}
		}
	}

	if r.Shopper == nil {
		r.Shopper = &CreateOrderShopper{Email: req.CustomerEmail}
	}

	signature, err := crypto.GenerateHMACSignature(r, sha256.New, callbackSecret, "hex")
	if err != nil {
		return nil, fmt.Errorf("generate ecdsa signature: %w", err)
	}
	r.Token = signature

	return r, nil
}
