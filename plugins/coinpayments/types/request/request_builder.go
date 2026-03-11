package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/pkg/crypto"
)

func BuildRequestHeader(clientID, clientSecret, method, requestURL string, payload any) (http.Header, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05") // time.RFC3339

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(payload); err != nil {
		return nil, err
	}
	payloadMessage := strings.TrimRight(buf.String(), "\n")

	message := "\ufeff" + strings.ToUpper(method) + requestURL + clientID + timestamp + payloadMessage

	signature, err := crypto.GenerateSignature(message, clientSecret, "base64")
	if err != nil {
		return nil, fmt.Errorf("generate signature: %w", err)
	}

	return http.Header{
		"X-CoinPayments-Timestamp": []string{timestamp},
		"X-CoinPayments-Signature": []string{signature},
	}, nil
}

// TODO: Items required
func BuildCreateInvoice(req models.CreatePaymentLinkRequest) (*CreateInvoice, error) {
	description := "Order " + req.OrderID
	if req.Description != "" {
		description = description + " | " + req.Description
	}
	r := &CreateInvoice{
		Currency:    req.Currency,
		Description: description,
		SuccessURL:  req.SuccessURL,
		CancelURL:   req.CancelURL,
	}
	r.Amount = &CreateInvoiceAmount{Total: req.Amount}
	r.Webhooks = append(r.Webhooks, CreateInvoiceWebhook{
		NotificationsURL: req.WebhookURL,
		Notifications: []string{ // TODO: 기본값은 남겨두고 나머지만 옵션 헝태로 고민 필요..
			"invoiceCreated", "invoicePending", "invoicePaid",
			"invoiceCompleted", "invoiceCancelled", "invoiceTimedOut",
			"invoicePaymentCreated", "invoicePaymentTimedOut",
		},
	})
	r.Payment = &CreateInvoicePayment{
		PaymentCurrency: req.ReceiveCurrency,
	}

	if len(req.ProviderPayload) > 0 {
		var payload CreateInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}
		r.Items = payload.Items
		if payload.Amount != nil && payload.Amount.Breakdown != nil {
			b := *payload.Amount.Breakdown
			r.Amount.Breakdown = &b
		}
		r.IsEmailDelivery = payload.IsEmailDelivery
		if payload.EmailDelivery != nil {
			e := *payload.EmailDelivery
			r.EmailDelivery = &e
		}
		r.DueDate = payload.DueDate
		r.InvoiceDate = payload.InvoiceDate
		r.Draft = payload.Draft
		r.ClientID = payload.ClientID
		r.InvoiceID = payload.InvoiceID
		if payload.Buyer != nil {
			b := *payload.Buyer
			r.Buyer = &b
		}
		if payload.Shipping != nil {
			s := *payload.Shipping
			r.Shipping = &s
		}
		r.RequireBuyerNameAndEmail = payload.RequireBuyerNameAndEmail
		r.BuyerDataCollectionMessage = payload.BuyerDataCollectionMessage
		r.Notes = payload.Notes
		r.NotesToRecipient = payload.NotesToRecipient
		r.TermsAndConditions = payload.TermsAndConditions
		if payload.MerchantOptions != nil {
			m := *payload.MerchantOptions
			r.MerchantOptions = &m
		}
		r.CustomData = payload.CustomData
		if payload.Metadata != nil {
			m := *payload.Metadata
			r.Metadata = &m
		}
		r.PoNumber = payload.PoNumber
		if len(payload.Webhooks) > 0 {
			r.Webhooks = append(r.Webhooks, payload.Webhooks...)
		}
		if len(payload.PayoutOverrides) > 0 {
			r.PayoutOverrides = append(r.PayoutOverrides, payload.PayoutOverrides...)
		}
		r.UseCoinReservation = payload.UseCoinReservation
		if payload.Payment != nil && payload.Payment.RefundEmail != "" {
			r.Payment.RefundEmail = payload.Payment.RefundEmail
		}
		r.HideShoppingCart = payload.HideShoppingCart
		r.AffiliateID = payload.AffiliateID
		r.IsSimpleQR = payload.IsSimpleQR
	}

	return r, nil
}
