package response

import "github.com/dracocity/draco-payment-bridge-core/internal/models"

func BuildCreatePaymentLink(resp CreateInvoice, webhookURL string) (*models.CreatePaymentLinkResponse, error) {
	receiveCurrency := ""
	if resp.Data.BuyerProvidedInfo != nil {
		receiveCurrency = resp.Data.BuyerProvidedInfo.SelectedTransactionCurrency
	}

	return &models.CreatePaymentLinkResponse{
		InvoiceID:       resp.Data.ID,
		Amount:          resp.Data.DisplayAmountPaid,
		Currency:        resp.Data.Currency,
		ReceiveCurrency: receiveCurrency,
		OrderID:         resp.Data.OrderID,
		Description:     resp.Data.ItemDesc,
		CheckoutURL:     resp.Data.URL,
		WebhookURL:      webhookURL,
		SuccessURL:      resp.Data.RedirectURL,
		CancelURL:       resp.Data.CloseURL,
		ExpiresAt:       &resp.Data.ExpirationTime,
		CreatedAt:       resp.Data.InvoiceTime,
		UpdatedAt:       resp.Data.InvoiceTime,
		Raw:             resp,
	}, nil
}
