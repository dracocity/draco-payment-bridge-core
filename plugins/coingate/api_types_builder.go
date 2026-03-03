package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func buildCreateOrderRequest(req models.CreatePaymentLinkRequest, callbackToken string) (*createOrderRequest, error) {
	priceAmount, err := strconv.ParseFloat(strings.TrimSpace(req.FiatAmount), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid fiat_amount: %w", err)
	}

	return &createOrderRequest{
		OrderID:         strings.TrimSpace(req.OrderID),
		PriceAmount:     priceAmount,
		PriceCurrency:   strings.ToUpper(strings.TrimSpace(req.FiatCurrency)),
		ReceiveCurrency: strings.ToUpper(strings.TrimSpace(req.CryptoCurrency)),
		Title:           strings.TrimSpace(req.Description), // TODO: Title
		Description:     strings.TrimSpace(req.Description), // TODO: Description
		CallbackURL:     strings.TrimSpace(req.WebhookURL),
		CancelURL:       strings.TrimSpace(req.CancelURL),
		SuccessURL:      strings.TrimSpace(req.SuccessURL),
		Token:           strings.TrimSpace(callbackToken),
	}, nil
}
