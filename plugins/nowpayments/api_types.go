package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// CreateInvoice
type createInvoiceRequest struct {
	createInvoicePayload
	PriceAmount      float64 `json:"price_amount"`
	PriceCurrency    string  `json:"price_currency"`
	OrderID          string  `json:"order_id,omitempty"`
	OrderDescription string  `json:"order_description,omitempty"`
	IPNCallbackURL   string  `json:"ipn_callback_url,omitempty"`
	SuccessURL       string  `json:"success_url,omitempty"`
	CancelURL        string  `json:"cancel_url,omitempty"`
}

type createInvoicePayload struct {
	PayCurrency     string `json:"pay_currency,omitempty"`
	IsFixedRate     *bool  `json:"is_fixed_rate,omitempty"`
	IsFeePaidByUser *bool  `json:"is_fee_paid_by_user,omitempty"`
}

type createInvoiceResponse struct {
	ID               string          `json:"id"`
	TokenID          *string         `json:"token_id"`
	OrderID          *string         `json:"order_id"`
	OrderDescription *string         `json:"order_description"`
	PriceAmount      decimal.Decimal `json:"price_amount"`
	PriceCurrency    string          `json:"price_currency"`
	PayCurrency      *string         `json:"pay_currency"`
	IPNCallbackURL   *string         `json:"ipn_callback_url"`
	InvoiceURL       string          `json:"invoice_url"`
	SuccessURL       *string         `json:"success_url"`
	CancelURL        *string         `json:"cancel_url"`
	PartiallyPaidURL *string         `json:"partially_paid_url"`
	PayoutCurrency   *string         `json:"payout_currency,omitempty"`
	CreatedAt        string          `json:"created_at"`
	UpdatedAt        string          `json:"updated_at"`
	IsFixedRate      *bool           `json:"is_fixed_rate"`
	IsFeePaidByUser  *bool           `json:"is_fee_paid_by_user"`
}

// CreateInvoicePayment
type createInvoicePaymentRequest struct {
	createInvoicePaymentPayload
	InvoiceID        string `json:"iid"`                         // invoice id. You can get invoice ID in response of POST Create_invoice method
	PayCurrency      string `json:"pay_currency"`                // the crypto currency in which the pay_amount is specified (btc, eth, etc). NOTE: some of the currencies require a Memo, Destination Tag, etc., to complete a payment (AVA, EOS, BNBMAINNET, XLM, XRP). This is unique for each payment. This ID is received in “payin_extra_id” parameter of the response. Payments made without "payin_extra_id" cannot be detected automatically;=
	OrderDescription string `json:"order_description,omitempty"` // inner store order description
	CustomerEmail    string `json:"customer_email,omitempty"`    // user email to which a notification about the successful completion of the payment will be sent
	PayoutAddress    string `json:"payout_address,omitempty"`    // usually the funds will go to the address you specify in your Personal account. In case you want to receive funds on another address, you can specify it in this parameter
	PayoutExtraID    string `json:"payout_extra_id,omitempty"`   // extra id or memo or tag for external payout_address
}

type createInvoicePaymentPayload struct {
}

type createInvoicePaymentResponse struct {
	PaymentID              string           `json:"payment_id"`
	PaymentStatus          string           `json:"payment_status"`
	PayAddress             string           `json:"pay_address"`
	PriceAmount            decimal.Decimal  `json:"price_amount"`
	PriceCurrency          string           `json:"price_currency"`
	PayAmount              decimal.Decimal  `json:"pay_amount"`
	PayCurrency            string           `json:"pay_currency"`
	OrderID                *string          `json:"order_id"`
	OrderDescription       *string          `json:"order_description"`
	IPNCallbackURL         *string          `json:"ipn_callback_url"`
	CreatedAt              string           `json:"created_at"`
	UpdatedAt              string           `json:"updated_at"`
	PurchaseID             string           `json:"purchase_id"`
	AmountReceived         decimal.Decimal  `json:"amount_received"`
	PayinExtraID           *string          `json:"payin_extra_id"`
	SmartContract          *string          `json:"smart_contract"`
	Network                string           `json:"network"`
	NetworkPrecision       *string          `json:"network_precision"`
	TimeLimit              *string          `json:"time_limit"`
	BurningPercent         *decimal.Decimal `json:"burning_percent"`
	ExpirationEstimateDate string           `json:"expiration_estimate_date"`
	IsFixedRate            bool             `json:"is_fixed_rate"`
	IsFeePaidByUser        bool             `json:"is_fee_paid_by_user"`
	ValidUntil             string           `json:"valid_until"`
	Type                   string           `json:"type"`
	RedirectURL            *string          `json:"redirect_url"`
	Product                string           `json:"product"`
	Success                string           `json:"Success"`
}

// GetPayment
type nowPaymentsGetPaymentResponse struct {
	PaymentID        string           `json:"payment_id"`
	InvoiceID        *string          `json:"invoice_id"`
	PaymentStatus    string           `json:"payment_status"`
	PayAddress       string           `json:"pay_address"`
	PayinExtraID     *string          `json:"payin_extra_id"`
	PriceAmount      decimal.Decimal  `json:"price_amount"`
	PriceCurrency    string           `json:"price_currency"`
	PayAmount        decimal.Decimal  `json:"pay_amount"`
	ActuallyPaid     decimal.Decimal  `json:"actually_paid"`
	PayCurrency      string           `json:"pay_currency"`
	OrderID          *string          `json:"order_id"`
	OrderDescription *string          `json:"order_description"`
	PurchaseID       string           `json:"purchase_id"`
	OutcomeAmount    decimal.Decimal  `json:"outcome_amount"`
	OutcomeCurrency  string           `json:"outcome_currency"`
	PayoutHash       *string          `json:"payout_hash"`
	PayinHash        *string          `json:"payin_hash"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
	BurningPercent   *decimal.Decimal `json:"burning_percent"`
	Type             string           `json:"type"`
	PaymentExtraIDs  []string         `json:"payment_extra_ids"`
}

func (r *nowPaymentsGetPaymentResponse) UnmarshalJSON(data []byte) error {
	type alias nowPaymentsGetPaymentResponse
	aux := struct {
		PaymentID       any   `json:"payment_id"`
		InvoiceID       any   `json:"invoice_id"`
		PurchaseID      any   `json:"purchase_id"`
		BurningPercent  any   `json:"burning_percent"`
		PaymentExtraIDs []any `json:"payment_extra_ids"`
		*alias
	}{
		alias: (*alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.PaymentID = toString(aux.PaymentID)
	r.InvoiceID = toPtrString(aux.InvoiceID)
	r.PurchaseID = toString(aux.PurchaseID)
	burningPercent, err := toPtrDecimal(aux.BurningPercent)
	if err != nil {
		return fmt.Errorf("invalid burning_percent (%v): %w", aux.BurningPercent, err)
	}
	r.BurningPercent = burningPercent

	for _, paymentExtraID := range aux.PaymentExtraIDs {
		id := toString(paymentExtraID)
		if id == "" {
			continue
		}
		r.PaymentExtraIDs = append(r.PaymentExtraIDs, id)
	}

	return nil
}

func toString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		if v == "" || strings.EqualFold(v, "null") {
			return ""
		}
		return strings.TrimSpace(v)
	case float64:
		return strings.TrimSpace(strconv.FormatFloat(v, 'f', -1, 64))
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}

func toPtrString(value any) *string {
	switch v := value.(type) {
	case nil:
		return nil
	default:
		s := toString(v)
		if s == "" {
			return nil
		}
		return &s
	}
}

func toDecimal(value any) (decimal.Decimal, error) {
	switch v := value.(type) {
	case nil:
		return decimal.Zero, nil
	case string:
		if v == "" || strings.EqualFold(v, "null") {
			return decimal.Zero, nil
		}
		return decimal.NewFromString(v)
	case float64:
		s := strings.TrimSpace(strconv.FormatFloat(v, 'f', -1, 64))
		if s == "" {
			return decimal.Zero, nil
		}
		return decimal.NewFromString(s)
	default:
		s := strings.TrimSpace(fmt.Sprint(v))
		if s == "" {
			return decimal.Zero, nil
		}
		return decimal.NewFromString(s)
	}
}

func toPtrDecimal(value any) (*decimal.Decimal, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	default:
		d, err := toDecimal(v)
		if err != nil {
			return nil, err
		}
		if d.Equal(decimal.Zero) {
			return nil, nil
		}
		return &d, nil

	}
}
