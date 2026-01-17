package models

import "time"

// CreatePaymentRequest is the unified API request payload for creating a payment.
type CreatePaymentRequest struct {
	Provider    string  `json:"provider"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Crypto      string  `json:"crypto"`
	Description string  `json:"description"`
	OrderID     string  `json:"order_id"`
	ReturnURL   string  `json:"return_url"`
	CancelURL   string  `json:"cancel_url"`
	NotifyURL   string  `json:"notify_url"`
}

// CreatePaymentResponse is the unified response for payment creation.
type CreatePaymentResponse struct {
	Success    bool      `json:"success"`
	PaymentID  string    `json:"payment_id"`
	PaymentURL string    `json:"payment_url"`
	QRCode     string    `json:"qr_code"`
	Amount     string    `json:"amount"`
	Currency   string    `json:"currency"`
	Status     string    `json:"status"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// PaymentStatusResponse represents the current status of a payment.
type PaymentStatusResponse struct {
	PaymentID   string    `json:"payment_id"`
	Status      string    `json:"status"`
	Amount      string    `json:"amount"`
	Currency    string    `json:"currency"`
	Transaction string    `json:"transaction"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RefundRequest is the unified refund request payload.
type RefundRequest struct {
	Provider  string  `json:"provider"`
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

// RefundResponse represents the unified refund response.
type RefundResponse struct {
	Success    bool   `json:"success"`
	RefundID   string `json:"refund_id"`
	PaymentID  string `json:"payment_id"`
	Amount     string `json:"amount"`
	Currency   string `json:"currency"`
	Status     string `json:"status"`
	RefundType string `json:"refund_type"`
}

// WebhookResult is returned after a provider webhook is processed.
type WebhookResult struct {
	Accepted bool   `json:"accepted"`
	Message  string `json:"message"`
}
