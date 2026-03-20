package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/pkg/crypto"
)

func BuildRequestHeader(privKey *secp256k1.PrivateKey, baseURL, path string, query url.Values, payload any) (http.Header, error) {
	requestURL := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
	if encodedQuery := query.Encode(); encodedQuery != "" {
		requestURL += "?" + encodedQuery
	}

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

	message := requestURL + payloadMessage

	signature, err := crypto.GenerateECDSASignature(message, privKey)
	if err != nil {
		return nil, fmt.Errorf("generate ecdsa signature: %w", err)
	}

	return http.Header{
		"X-Signature": []string{signature},
	}, nil
}

func BuildCreateInvoice(req models.CreatePaymentLinkRequest, apiToken string) (*CreateInvoice, error) {
	price, err := strconv.ParseFloat(strings.TrimSpace(req.Amount), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid fiat_amount: %w", err)
	}
	itemDesc := strings.TrimSpace(req.Description)
	if n := len(req.Items); n > 0 {
		item := req.Items[0]
		if itemDesc != "" {
			itemDesc += " | "
		}
		itemDesc += fmt.Sprintf("%s x%d", item.Name, item.Quantity)
		if n > 1 {
			itemDesc += fmt.Sprintf(" (+%d more)", n-1)
		}
	}

	r := &CreateInvoice{
		Token:           apiToken,
		Price:           price,
		Currency:        strings.ToLower(strings.TrimSpace(req.Currency)),
		OrderID:         strings.TrimSpace(req.OrderID),
		ItemDesc:        itemDesc,
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
			r.Buyer.Email = req.BuyerEmail
		}
		r.JsonPayProRequired = payload.JsonPayProRequired
		r.AcceptanceWindow = payload.AcceptanceWindow
	}

	if r.Buyer == nil {
		r.Buyer = &CreateInvoiceBuyer{Email: req.BuyerEmail}
	}

	return r, nil
}

func BuildCreateRefundRequest(req models.CreateRefundRequest, apiToken string) *CreateRefundRequest {
	return &CreateRefundRequest{
		Amount:    req.Amount,
		InvoiceID: strings.TrimSpace(req.PaymentID),
		Token:     apiToken,
		Reference: strings.TrimSpace(req.Reason),
	}
}
