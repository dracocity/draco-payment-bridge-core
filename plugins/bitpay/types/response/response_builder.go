package response

import (
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/consts"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/types"
)

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

func BuildGetPayment(resp RetrieveInvoice) (*models.GetPaymentResponse, error) {
	return &models.GetPaymentResponse{
		Status:         normalizeBitPayStatus(resp.Data.Status),
		ProviderStatus: strings.TrimSpace(resp.Data.Status),
		PaymentID:      strings.TrimSpace(resp.Data.ID),
		OrderID:        strings.TrimSpace(resp.Data.OrderID),
		Currency:       strings.ToUpper(strings.TrimSpace(resp.Data.Currency)),
		Amount:         resp.Data.Price.String(),
		CreatedAt:      resp.Data.InvoiceTime,
		Raw:            resp,
	}, nil
}

func BuildCreateRefund(resp CreateRefundRequest, fallbackPaymentID string) (*models.CreateRefundResponse, error) {
	status := strings.TrimSpace(resp.Data.Status)
	success := status != "" && !strings.EqualFold(status, "failed") && !strings.EqualFold(status, "canceled")

	paymentID := strings.TrimSpace(resp.Data.Invoice)
	if paymentID == "" {
		paymentID = strings.TrimSpace(fallbackPaymentID)
	}

	return &models.CreateRefundResponse{
		Success:    success,
		RefundID:   strings.TrimSpace(resp.Data.ID),
		PaymentID:  paymentID,
		Amount:     resp.Data.Amount.String(),
		Currency:   strings.ToUpper(strings.TrimSpace(resp.Data.Currency)),
		Status:     status,
		RefundType: "standard",
	}, nil
}

func normalizeBitPayStatus(status string) types.PaymentStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "new":
		return consts.StatusCreated
	case "paid":
		return consts.StatusPending
	case "confirmed":
		return consts.StatusConfirming
	case "complete":
		return consts.StatusCompleted
	case "expired":
		return consts.StatusExpired
	case "invalid":
		return consts.StatusFailed
	default:
		normalized := strings.ToUpper(strings.TrimSpace(status))
		if normalized == "" {
			return consts.StatusPending
		}
		return types.PaymentStatus(normalized)
	}
}
