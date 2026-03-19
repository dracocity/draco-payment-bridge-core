package response

import (
	"github.com/shopspring/decimal"
)

// CreateInvoice
type CreateInvoice struct {
	ID               string          `json:"id"`                        // Invoice ID
	TokenID          *string         `json:"token_id"`                  // Internal identifier
	OrderID          *string         `json:"order_id"`                  // Order ID specified in request
	OrderDescription *string         `json:"order_description"`         // Order description specified in request
	PriceAmount      decimal.Decimal `json:"price_amount"`              // Base price in fiat
	PriceCurrency    string          `json:"price_currency"`            // Ticker of base fiat currency
	PayCurrency      *string         `json:"pay_currency"`              // Currency your customer will pay with. If it's 'null' your customer can choose currency in web interface.
	IPNCallbackURL   *string         `json:"ipn_callback_url"`          // Link to your endpoint for IPN notifications catching
	InvoiceURL       string          `json:"invoice_url"`               // Link to the payment page that you can share with your customer
	SuccessURL       *string         `json:"success_url"`               // Customer will be redirected to this link once the payment is finished
	CancelURL        *string         `json:"cancel_url"`                // Customer will be redirected to this link if the payment fails
	PartiallyPaidURL *string         `json:"partially_paid_url"`        // Customer will be redirected to this link if the payment gets partially paid status
	PayoutCurrency   *string         `json:"payout_currency,omitempty"` // Ticker of payout currency
	CreatedAt        string          `json:"created_at"`                // Time of invoice creation
	UpdatedAt        string          `json:"updated_at"`                // Time of latest invoice information update
	IsFixedRate      *bool           `json:"is_fixed_rate"`             // This parameter is 'True' if Fixed Rate option is enabled and 'false' if it's disabled
	IsFeePaidByUser  *bool           `json:"is_fee_paid_by_user"`       // This parameter is 'True' if Fee Paid By User option is enabled and 'false' if it's disabled
}

// CreateInvoicePayment
type CreateInvoicePayment struct {
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
type GetPayment struct {
	PaymentID        int64            `json:"payment_id"`
	InvoiceID        *int64           `json:"invoice_id"`
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
	PurchaseID       int64            `json:"purchase_id"`
	OutcomeAmount    *decimal.Decimal `json:"outcome_amount"`
	OutcomeCurrency  string           `json:"outcome_currency"`
	PayoutHash       *string          `json:"payout_hash"`
	PayinHash        *string          `json:"payin_hash"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
	BurningPercent   string           `json:"burning_percent"`
	Type             string           `json:"type"`
	PaymentExtraIDs  []int64          `json:"payment_extra_ids"`
}
