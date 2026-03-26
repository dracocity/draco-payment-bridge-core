package request

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/pkg/crypto"
)

func BuildRequestHeader(clientID, clientSecret, method, baseURL, path string, query url.Values, payload any) (http.Header, error) {
	requestURL := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
	if encodedQuery := query.Encode(); encodedQuery != "" {
		requestURL += "?" + encodedQuery
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05") // time.RFC3339

	var buf bytes.Buffer
	if payload != nil {
		// Preserve previous behavior: skip body serialization when payload is an empty string.
		if s, ok := payload.(string); !ok || s != "" {
			enc := json.NewEncoder(&buf)
			enc.SetEscapeHTML(false)
			if err := enc.Encode(payload); err != nil {
				return nil, err
			}
		}
	}
	payloadMessage := strings.TrimRight(buf.String(), "\n")

	message := "\ufeff" + strings.ToUpper(method) + requestURL + clientID + timestamp + payloadMessage

	signature, err := crypto.GenerateHMACSignature(message, sha256.New, clientSecret, "base64")
	if err != nil {
		return nil, fmt.Errorf("generate hmac signature: %w", err)
	}

	return http.Header{
		"X-CoinPayments-Timestamp": []string{timestamp},
		"X-CoinPayments-Signature": []string{signature},
	}, nil
}

func BuildCreateInvoice(req models.CreatePaymentLinkRequest) (*CreateInvoice, error) {
	description := "Order " + req.OrderID
	if req.Description != "" {
		description = description + " | " + req.Description
	}
	r := &CreateInvoice{
		Currency:    req.Currency,
		Amount:      &CreateInvoiceAmount{Total: req.Amount},
		Description: description,
		Webhooks: []CreateInvoiceWebhook{{
			NotificationsURL: req.WebhookURL,
			Notifications: []string{ // TODO: 기본값은 남겨두고 나머지만 옵션 헝태로 고민 필요..
				"invoiceCreated", "invoicePending", "invoicePaid",
				"invoiceCompleted", "invoiceCancelled", "invoiceTimedOut",
				"invoicePaymentCreated", "invoicePaymentTimedOut",
			},
		}},
		Payment: &CreateInvoicePayment{
			PaymentCurrency: req.ReceiveCurrency,
			RefundEmail:     req.CustomerEmail,
		},
		SuccessURL: req.SuccessURL,
		CancelURL:  req.CancelURL,
	}
	if len(req.Items) == 0 {
		req.Items = append(req.Items, models.Item{
			ID:       req.OrderID,
			Name:     "Order " + req.OrderID,
			Quantity: 1,
			Amount:   req.Amount,
		})
	}
	r.Items = make([]CreateInvoiceItem, 0, len(req.Items))
	for _, item := range req.Items {
		r.Items = append(r.Items, CreateInvoiceItem{
			CustomID: item.ID,
			Name:     item.Name,
			Quantity: &CreateInvoiceQuantity{Value: item.Quantity, Type: "quantity"},
			Amount:   item.Amount,
		})
	}

	if len(req.ProviderPayload) > 0 {
		var payload CreateInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
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
		if len(payload.PayoutOverrides) > 0 {
			r.PayoutOverrides = append(r.PayoutOverrides, payload.PayoutOverrides...)
		}
		r.UseCoinReservation = payload.UseCoinReservation
		r.HideShoppingCart = payload.HideShoppingCart
		r.AffiliateID = payload.AffiliateID
		r.IsSimpleQR = payload.IsSimpleQR
	}

	return r, nil
}
