package consts

import "github.com/dracocity/draco-payment-bridge-core/internal/types"

const (
	// Payment has been created but no payment detected yet
	StatusCreated types.PaymentStatus = "CREATED"

	// Waiting for customer to send funds
	StatusPending types.PaymentStatus = "PENDING"

	// Transaction detected, waiting for required confirmations
	StatusConfirming types.PaymentStatus = "CONFIRMING"

	// Required confirmations reached, payment completed
	StatusCompleted types.PaymentStatus = "COMPLETED"

	// Payment failed due to error
	StatusFailed types.PaymentStatus = "FAILED"

	// Payment expired before funds were received
	StatusExpired types.PaymentStatus = "EXPIRED"

	// Customer sent less than required amount
	StatusUnderpaid types.PaymentStatus = "UNDERPAID"

	// Customer sent more than required amount
	StatusOverpaid types.PaymentStatus = "OVERPAID"

	// Partial payment received
	StatusPartiallyPaid types.PaymentStatus = "PARTIALLY_PAID"
)
