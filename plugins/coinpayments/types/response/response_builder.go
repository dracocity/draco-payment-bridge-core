package response

import (
	"errors"
	"strings"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/pkg/utils"
)

func BuildCreatePaymentLink(resp CreateInvoice, amount, currency, orderID string) (*models.CreatePaymentLinkResponse, error) {
	if len(resp.Invoices) == 0 {
		return nil, errors.New("could not retrieve invoices")
	}
	invoice := resp.Invoices[0]
	paymentCurrency := invoice.Payment.PaymentCurrencies[0]
	expiresAt := utils.ToUnixMilli(invoice.Payment.Expires)
	createdAt := time.Now().UnixMilli()
	return &models.CreatePaymentLinkResponse{
		InvoiceID:       invoice.ID,
		Amount:          amount,
		Currency:        currency,
		ReceiveCurrency: strings.ToUpper(paymentCurrency.Currency.Symbol),
		OrderID:         orderID,
		CheckoutURL:     invoice.CheckoutLink,
		ExpiresAt:       &expiresAt,
		CreatedAt:       createdAt,
		UpdatedAt:       createdAt,
		Raw:             resp,
	}, nil
}
