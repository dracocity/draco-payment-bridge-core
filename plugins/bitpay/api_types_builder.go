package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func buildCreateInvoiceRequest(req models.CreatePaymentLinkRequest, apiToken string) (createInvoiceRequest, error) {
	price, err := strconv.ParseFloat(strings.TrimSpace(req.FiatAmount), 64)
	if err != nil {
		return createInvoiceRequest{}, fmt.Errorf("invalid fiat_amount: %w", err)
	}

	r := createInvoiceRequest{
		Token:           apiToken,
		Price:           price,
		Currency:        strings.ToLower(strings.TrimSpace(req.FiatCurrency)),
		OrderID:         strings.TrimSpace(req.OrderID),
		ItemDesc:        strings.TrimSpace(req.Description),
		NotificationURL: strings.TrimSpace(req.WebhookURL),
		RedirectURL:     strings.TrimSpace(req.SuccessURL),
		CloseURL:        strings.TrimSpace(req.CancelURL),
	}

	if len(req.ProviderPayload) > 0 {
		var payload createInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return createInvoiceRequest{}, fmt.Errorf("invalid provider_payload: %w", err)
		}
		r.BitpayIdRequired = payload.BitpayIdRequired
		r.MerchantName = payload.MerchantName
		r.ForcedBuyerSelectedTransactionCurrency = payload.ForcedBuyerSelectedTransactionCurrency
		r.ForcedBuyerSelectedWallet = payload.ForcedBuyerSelectedWallet
		r.ItemCode = payload.ItemCode
		r.ItemizedDetails = payload.ItemizedDetails
		r.NotificationEmail = payload.NotificationEmail
		r.AutoRedirect = payload.AutoRedirect
		r.PosData = payload.PosData
		r.GUID = payload.GUID
		r.TransactionSpeed = payload.TransactionSpeed
		r.FullNotifications = payload.FullNotifications
		r.ExtendedNotifications = payload.ExtendedNotifications
		r.Physical = payload.Physical
		r.BuyerSMS = payload.BuyerSMS
		if payload.Buyer != nil {
			b := *payload.Buyer
			r.Buyer = &b
		}
		r.JsonPayProRequired = payload.JsonPayProRequired
		r.AcceptanceWindow = payload.AcceptanceWindow
	}

	return r, nil
}
