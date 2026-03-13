package response

import (
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func BuildCreatePaymentLink(resp CreateInvoice) (*models.CreatePaymentLinkResponse, error) {
	return &models.CreatePaymentLinkResponse{}, nil
	// createdAt := utils.ToUnixMilli(resp.CreatedAt)
	// return &models.CreatePaymentLinkResponse{
	// 	InvoiceID:       strconv.FormatInt(resp.ID, 10),
	// 	Amount:          resp.PriceAmount,
	// 	Currency:        strings.ToUpper(resp.PriceCurrency),
	// 	ReceiveCurrency: resp.ReceiveCurrency,
	// 	OrderID:         resp.OrderID,
	// 	CheckoutURL:     resp.PaymentURL,
	// 	CreatedAt:       createdAt,
	// 	UpdatedAt:       createdAt,
	// 	Raw:             resp,
	// }, nil
}
