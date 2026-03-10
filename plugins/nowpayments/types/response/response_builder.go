package response

import (
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/pkg/utils"
)

func BuildCreatePaymentLink(resp CreateInvoice) (*models.CreatePaymentLinkResponse, error) {
	createdAt := utils.ToUnixMilli(resp.CreatedAt)
	updatedAt := utils.ToUnixMilli(resp.UpdatedAt)
	if updatedAt == 0 {
		updatedAt = createdAt
	}
	return &models.CreatePaymentLinkResponse{
		InvoiceID:       resp.ID,
		Amount:          resp.PriceAmount.String(),
		Currency:        resp.PriceCurrency,
		ReceiveCurrency: utils.PtrValueOrDefault(resp.PayCurrency, ""),
		OrderID:         utils.PtrValueOrDefault(resp.OrderID, ""),
		CheckoutURL:     resp.InvoiceURL,
		WebhookURL:      utils.PtrValueOrDefault(resp.IPNCallbackURL, ""),
		SuccessURL:      utils.PtrValueOrDefault(resp.SuccessURL, ""),
		CancelURL:       utils.PtrValueOrDefault(resp.CancelURL, ""),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Raw:             resp,
	}, nil
}
