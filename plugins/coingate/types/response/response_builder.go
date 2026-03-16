package response

import (
	"strconv"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/pkg/utils"
)

func BuildCreatePaymentLink(resp CreateOrder) (*models.CreatePaymentLinkResponse, error) {
	createdAt := utils.ToUnixMilli(resp.CreatedAt)
	return &models.CreatePaymentLinkResponse{
		InvoiceID:       strconv.FormatInt(resp.ID, 10),
		Amount:          resp.PriceAmount,
		Currency:        strings.ToUpper(resp.PriceCurrency),
		ReceiveCurrency: resp.ReceiveCurrency,
		OrderID:         resp.OrderID,
		CheckoutURL:     resp.PaymentURL,
		CreatedAt:       createdAt,
		UpdatedAt:       createdAt,
		Raw:             resp,
	}, nil
}
