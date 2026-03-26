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
	PaymentID              string           `json:"payment_id"`               // Unique payment identifier
	PaymentStatus          string           `json:"payment_status"`           // Current payment status
	PayAddress             string           `json:"pay_address"`              // Destination address for payment
	PriceAmount            decimal.Decimal  `json:"price_amount"`             // Base price amount in fiat
	PriceCurrency          string           `json:"price_currency"`           // Ticker of base fiat currency
	PayAmount              decimal.Decimal  `json:"pay_amount"`               // Amount expected in pay currency
	PayCurrency            string           `json:"pay_currency"`             // Ticker of pay currency
	OrderID                *string          `json:"order_id"`                 // Order ID specified in request
	OrderDescription       *string          `json:"order_description"`        // Order description specified in request
	IPNCallbackURL         *string          `json:"ipn_callback_url"`         // IPN callback URL from request
	CreatedAt              string           `json:"created_at"`               // Time of payment creation
	UpdatedAt              string           `json:"updated_at"`               // Time of latest payment update
	PurchaseID             string           `json:"purchase_id"`              // Purchase identifier in NOWPayments
	AmountReceived         decimal.Decimal  `json:"amount_received"`          // Amount already received from payer
	PayinExtraID           *string          `json:"payin_extra_id"`           // Extra destination tag/memo for payment
	SmartContract          *string          `json:"smart_contract"`           // Token smart contract address (if applicable)
	Network                string           `json:"network"`                  // Blockchain network name
	NetworkPrecision       *string          `json:"network_precision"`        // Decimal precision supported by network
	TimeLimit              *string          `json:"time_limit"`               // Time limit for completing payment
	BurningPercent         *decimal.Decimal `json:"burning_percent"`          // Burn fee percent for specific assets
	ExpirationEstimateDate string           `json:"expiration_estimate_date"` // Estimated payment expiration timestamp
	IsFixedRate            bool             `json:"is_fixed_rate"`            // True when fixed-rate exchange is enabled
	IsFeePaidByUser        bool             `json:"is_fee_paid_by_user"`      // True when network/exchange fee is paid by user
	ValidUntil             string           `json:"valid_until"`              // Payment validity deadline
	Type                   string           `json:"type"`                     // Payment type
	RedirectURL            *string          `json:"redirect_url"`             // URL to redirect customer after processing
	Product                string           `json:"product"`                  // Product name/type for checkout flow
	Success                string           `json:"Success"`                  // Request result flag from API response
}

// GetPayment
type GetPayment struct {
	PaymentID        int64            `json:"payment_id"`        // Unique payment identifier
	InvoiceID        *int64           `json:"invoice_id"`        // Related invoice identifier (if payment belongs to invoice)
	PaymentStatus    string           `json:"payment_status"`    // Current payment status
	PayAddress       string           `json:"pay_address"`       // Destination address for payment
	PayinExtraID     *string          `json:"payin_extra_id"`    // Extra destination tag/memo for payment
	PriceAmount      decimal.Decimal  `json:"price_amount"`      // Base price amount in fiat
	PriceCurrency    string           `json:"price_currency"`    // Ticker of base fiat currency
	PayAmount        decimal.Decimal  `json:"pay_amount"`        // Amount expected in pay currency
	ActuallyPaid     decimal.Decimal  `json:"actually_paid"`     // Amount actually received from payer
	PayCurrency      string           `json:"pay_currency"`      // Ticker of pay currency
	OrderID          *string          `json:"order_id"`          // Order ID specified in request
	OrderDescription *string          `json:"order_description"` // Order description specified in request
	PurchaseID       int64            `json:"purchase_id"`       // Purchase identifier in NOWPayments
	OutcomeAmount    *decimal.Decimal `json:"outcome_amount"`    // Final payout amount after conversion/fees
	OutcomeCurrency  string           `json:"outcome_currency"`  // Ticker of payout currency
	PayoutHash       *string          `json:"payout_hash"`       // Transaction hash of payout transfer
	PayinHash        *string          `json:"payin_hash"`        // Transaction hash of incoming payment
	CreatedAt        string           `json:"created_at"`        // Time of payment creation
	UpdatedAt        string           `json:"updated_at"`        // Time of latest payment update
	BurningPercent   string           `json:"burning_percent"`   // Burn fee percent for specific assets
	Type             string           `json:"type"`              // Payment type
	PaymentExtraIDs  []int64          `json:"payment_extra_ids"` // Extra payment identifiers associated with this payment
}
