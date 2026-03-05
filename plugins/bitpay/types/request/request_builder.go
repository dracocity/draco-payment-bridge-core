package request

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func BuildCreateInvoice(req models.CreatePaymentLinkRequest, apiToken string) (*CreateInvoice, error) {
	price, err := strconv.ParseFloat(strings.TrimSpace(req.Amount), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid fiat_amount: %w", err)
	}

	r := &CreateInvoice{
		Token:           apiToken,
		Price:           price,
		Currency:        strings.ToLower(strings.TrimSpace(req.Currency)),
		OrderID:         strings.TrimSpace(req.OrderID),
		ItemDesc:        strings.TrimSpace(req.Description),
		NotificationURL: strings.TrimSpace(req.WebhookURL),
		RedirectURL:     strings.TrimSpace(req.SuccessURL),
		CloseURL:        strings.TrimSpace(req.CancelURL),
	}

	if len(req.ProviderPayload) > 0 {
		var payload CreateInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
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
