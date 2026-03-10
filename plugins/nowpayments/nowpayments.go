package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/consts"
	"github.com/dracocity/draco-payment-bridge-core/internal/httpx"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/plugin"
	"github.com/dracocity/draco-payment-bridge-core/internal/types"
	"github.com/dracocity/draco-payment-bridge-core/plugins/nowpayments/types/request"
	"github.com/dracocity/draco-payment-bridge-core/plugins/nowpayments/types/response"
)

const (
	prodBaseURL    = "https://api.nowpayments.io/v1"
	sandboxBaseURL = "https://api-sandbox.nowpayments.io/v1"
)

type nowPaymentsPlugin struct {
	pg *nowPaymentsPG
}

type nowPaymentsPG struct {
	client    *httpx.Client
	ipnSecret string
}

func New() plugin.Plugin {
	return &nowPaymentsPlugin{}
}

func (p *nowPaymentsPlugin) Load(cfg config.PGConfig) error {
	var baseURL string
	switch strings.ToLower(cfg["mode"]) {
	case "sandbox":
		baseURL = sandboxBaseURL
	case "production":
		baseURL = prodBaseURL
	default:
		return fmt.Errorf("invalid mode: %s", cfg["mode"])
	}
	apiKey := strings.TrimSpace(cfg["api_key"])
	if apiKey == "" {
		return errors.New("nowpayments api_key is required")
	}
	ipnSecret := strings.TrimSpace(cfg["ipn_secret"])
	if ipnSecret == "" {
		return errors.New("nowpayments ipn_secret is required")
	}
	timeout := 7 * time.Second
	if timeoutStr := strings.TrimSpace(cfg["http_timeout"]); timeoutStr != "" {
		parsedTimeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid http_timeout: %w", err)
		}
		timeout = parsedTimeout
	}
	p.pg = &nowPaymentsPG{
		client: httpx.New(httpx.ClientOptions{
			Client:  &http.Client{Timeout: timeout},
			BaseURL: baseURL,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
				"x-api-key":    []string{apiKey},
			},
			ErrorPrefix: "nowpayments",
		}),
		ipnSecret: ipnSecret,
	}
	return nil
}

func (p *nowPaymentsPlugin) Unload() error {
	return nil
}

func (p *nowPaymentsPlugin) Name() string {
	return "nowpayments"
}

func (p *nowPaymentsPlugin) PaymentGateway() pg.PaymentGateway {
	return p.pg
}

func (b *nowPaymentsPG) Name() string {
	return "nowpayments"
}

func (b *nowPaymentsPG) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	body, err := request.BuildCreateInvoice(req)
	if err != nil {
		return nil, err
	}

	var resp response.CreateInvoice
	if err := b.client.Post(ctx, "/invoice", nil, body, &resp); err != nil {
		return nil, err
	}

	receiveCurrency := ""
	if resp.PayCurrency != nil {
		receiveCurrency = *resp.PayCurrency
	}

	orderID := req.OrderID
	if resp.OrderID != nil {
		orderID = *resp.OrderID
	}

	webhookURL := ""
	if resp.IPNCallbackURL != nil {
		webhookURL = *resp.IPNCallbackURL
	}

	successURL := ""
	if resp.SuccessURL != nil {
		successURL = *resp.SuccessURL
	}

	cancelURL := ""
	if resp.CancelURL != nil {
		cancelURL = *resp.CancelURL
	}

	createdAt := toUnixMilli(resp.CreatedAt)
	updatedAt := toUnixMilli(resp.UpdatedAt)
	if updatedAt == 0 {
		updatedAt = createdAt
	}
	return &models.CreatePaymentLinkResponse{
		InvoiceID:       resp.ID,
		Amount:          resp.PriceAmount.String(),
		Currency:        resp.PriceCurrency,
		ReceiveCurrency: receiveCurrency,
		OrderID:         orderID,
		CheckoutURL:     resp.InvoiceURL,
		WebhookURL:      webhookURL,
		SuccessURL:      successURL,
		CancelURL:       cancelURL,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Raw:             resp,
	}, nil
}

func (b *nowPaymentsPG) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	// TODO
	body, err := request.BuildCreateInvoicePayment(req)
	if err != nil {
		return nil, err
	}

	var resp response.CreateInvoicePayment
	if err := b.client.Post(ctx, "/invoice-payment", nil, body, &resp); err != nil {
		return nil, err
	}

	// TODO: CreatePaymentResponse 필드를 정리부터 해야..
	orderID := req.OrderID
	if resp.OrderID != nil {
		orderID = *resp.OrderID
	}

	createdAt := toUnixMilli(resp.CreatedAt)
	expiresAt := toUnixMilli(resp.ExpirationEstimateDate)
	if expiresAt == 0 {
		expiresAt = createdAt + 1800000 // 30 minutes
	}
	updatedAt := toUnixMilli(resp.UpdatedAt)
	if updatedAt == 0 {
		updatedAt = createdAt
	}

	return &models.CreatePaymentResponse{
		Status:       consts.StatusPending,
		PaymentID:    resp.PaymentID,
		OrderID:      orderID,
		Currency:     strings.ToUpper(resp.PayCurrency),
		Amount:       resp.PayAmount.String(),
		FiatCurrency: strings.ToUpper(resp.PriceCurrency),
		FiatAmount:   resp.PriceAmount.String(),
		Network:      resp.Network,
		// ExchangeRate: "",
		DepositAddress: resp.PayAddress,
		Memo:           resp.PayinExtraID,
		ExpiresAt:      expiresAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Raw:            resp,
	}, nil
}

func (b *nowPaymentsPG) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
	var resp response.GetPayment
	if err := b.client.Get(ctx, "/payment/"+paymentID, nil, nil, &resp); err != nil {
		return nil, err
	}
	// TODO: GetPaymentResponse 필드를 정리부터 해야..
	orderID := ""
	if resp.OrderID != nil {
		orderID = *resp.OrderID
	}

	expectedAmount := resp.PayAmount
	if expectedAmount.IsZero() {
		expectedAmount = resp.PriceAmount
	}

	paidAmount := resp.ActuallyPaid
	remainingAmount := expectedAmount.Sub(paidAmount)
	if remainingAmount.IsNegative() {
		remainingAmount = remainingAmount.Abs()
	} else {
		remainingAmount = remainingAmount.Round(16)
	}

	overpaidAmount := paidAmount.Sub(expectedAmount)
	if overpaidAmount.IsNegative() {
		overpaidAmount = overpaidAmount.Neg()
	} else {
		overpaidAmount = overpaidAmount.Round(16)
	}

	if paidAmount.LessThanOrEqual(expectedAmount) {
		overpaidAmount = overpaidAmount.Sub(overpaidAmount)
	}

	if paidAmount.GreaterThanOrEqual(expectedAmount) {
		remainingAmount = remainingAmount.Sub(remainingAmount)
	}

	createdAt := toUnixMilli(resp.CreatedAt)
	updatedAt := toUnixMilli(resp.UpdatedAt)
	if updatedAt == 0 {
		updatedAt = createdAt
	}

	transactionHash := ""
	if resp.PayinHash != nil {
		transactionHash = strings.TrimSpace(*resp.PayinHash)
	}
	if transactionHash == "" && resp.PayoutHash != nil {
		transactionHash = strings.TrimSpace(*resp.PayoutHash)
	}

	currency := strings.ToUpper(strings.TrimSpace(resp.PayCurrency))
	if currency == "" {
		currency = strings.ToUpper(strings.TrimSpace(resp.OutcomeCurrency))
	}

	return &models.GetPaymentResponse{
		Status:          normalizeNowPaymentsStatus(resp.PaymentStatus),
		PaymentID:       resp.PaymentID,
		OrderID:         orderID,
		Currency:        currency,
		Amount:          expectedAmount.String(),
		FiatCurrency:    strings.ToUpper(strings.TrimSpace(resp.PriceCurrency)),
		FiatAmount:      resp.PriceAmount.String(),
		DepositAddress:  strings.TrimSpace(resp.PayAddress),
		Memo:            resp.PayinExtraID,
		TransactionHash: transactionHash,
		PaidAmount:      paidAmount.String(),
		RemainingAmount: remainingAmount.String(),
		OverpaidAmount:  overpaidAmount.String(),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Raw:             resp,
	}, nil
}

func (b *nowPaymentsPG) Refund(ctx context.Context, req models.RefundRequest) (*models.RefundResponse, error) {
	return nil, errors.New("nowpayments does not provide an API refund endpoint")
}

func (b *nowPaymentsPG) HandleWebhook(ctx context.Context, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	signature := headerValue(headers, "x-nowpayments-sig")
	if signature == "" {
		return nil, errors.New("missing x-nowpayments-sig header")
	}
	calculated, err := nowPaymentsSignature(payload, b.ipnSecret)
	if err != nil {
		return nil, err
	}
	if !hmac.Equal([]byte(strings.ToLower(signature)), []byte(calculated)) {
		return nil, errors.New("invalid nowpayments signature")
	}
	return &models.WebhookResult{
		Accepted: true,
		Message:  "webhook verified",
	}, nil
}

func nowPaymentsSignature(payload []byte, secret string) (string, error) {
	sortedJSON, err := sortJSONPayload(payload)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha512.New, []byte(secret))
	if _, err := mac.Write([]byte(sortedJSON)); err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func sortJSONPayload(payload []byte) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(payload, &params); err != nil {
		return "", fmt.Errorf("invalid webhook payload: %w", err)
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, key := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		keyJSON, err := marshalNoEscape(key)
		if err != nil {
			return "", err
		}
		valueJSON, err := marshalNoEscape(params[key])
		if err != nil {
			return "", err
		}
		buf.Write(keyJSON)
		buf.WriteByte(':')
		buf.Write(valueJSON)
	}
	buf.WriteByte('}')
	return buf.String(), nil
}

func marshalNoEscape(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

func toUnixMilli(datetime string) int64 {
	if datetime == "" {
		return 0
	}

	if ts, err := time.Parse(time.RFC3339, datetime); err == nil {
		return ts.UnixMilli()
	}

	return 0
}

func normalizeNowPaymentsStatus(status string) types.PaymentStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "new":
		return consts.StatusCreated
	case "waiting":
		return consts.StatusPending
	case "confirming", "sending":
		return consts.StatusConfirming
	case "partially_paid":
		return consts.StatusPartiallyPaid
	case "finished", "confirmed":
		return consts.StatusCompleted
	case "failed", "refunded":
		return consts.StatusFailed
	case "expired":
		return consts.StatusExpired
	default:
		normalized := strings.ToUpper(strings.TrimSpace(status))
		if normalized == "" {
			return consts.StatusPending
		}
		return types.PaymentStatus(normalized)
	}
}

func headerValue(headers map[string][]string, name string) string {
	for key, values := range headers {
		if strings.EqualFold(key, name) && len(values) > 0 {
			return strings.TrimSpace(values[0])
		}
	}
	return ""
}
