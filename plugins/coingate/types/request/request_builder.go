package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func BuildCreateOrder(req models.CreatePaymentLinkRequest, callbackToken string) (*CreateOrder, error) {
	priceAmount, err := strconv.ParseFloat(strings.TrimSpace(req.Amount), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	r := &CreateOrder{
		OrderID:         strings.TrimSpace(req.OrderID),
		PriceAmount:     priceAmount,
		PriceCurrency:   strings.ToUpper(strings.TrimSpace(req.Currency)),
		ReceiveCurrency: strings.ToUpper(strings.TrimSpace(req.ReceiveCurrency)),
		Title:           strings.TrimSpace(req.Description), // TODO: Title
		Description:     strings.TrimSpace(req.Description), // TODO: Description
		CallbackURL:     strings.TrimSpace(req.WebhookURL),
		CancelURL:       strings.TrimSpace(req.CancelURL),
		SuccessURL:      strings.TrimSpace(req.SuccessURL),
	}
	if len(req.ProviderPayload) > 0 {
		var payload CreateOrderPayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}
		r.Token = payload.Token
	}

	return r, nil
}
