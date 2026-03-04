package response

import "github.com/shopspring/decimal"

// CreateInvoice
type CreateInvoice struct {
	Facade string `json:"facade"`
	Data   struct {
		URL                 string           `json:"url"`
		PosData             string           `json:"posData,omitempty"`
		Status              string           `json:"status"`
		Price               decimal.Decimal  `json:"price"`
		Currency            string           `json:"currency"`
		ItemDesc            string           `json:"itemDesc,omitempty"`
		OrderID             string           `json:"orderId,omitempty"`
		InvoiceTime         int64            `json:"invoiceTime"`
		ExpirationTime      int64            `json:"expirationTime"`
		CurrentTime         int64            `json:"currentTime"`
		GUID                string           `json:"guid,omitempty"`
		ID                  string           `json:"id"`
		LowFeeDetected      bool             `json:"lowFeeDetected"`
		AmountPaid          decimal.Decimal  `json:"amountPaid"`
		DisplayAmountPaid   string           `json:"displayAmountPaid,omitempty"`
		ExceptionStatus     bool             `json:"exceptionStatus"`
		TargetConfirmations *int64           `json:"targetConfirmations,omitempty"`
		Transactions        []map[string]any `json:"transactions,omitempty"`
		TransactionSpeed    string           `json:"transactionSpeed,omitempty"`
		Buyer               *struct {
			Email string `json:"email"`
		} `json:"buyer,omitempty"`
		RedirectURL     string `json:"redirectURL,omitempty"`
		AutoRedirect    *bool  `json:"autoRedirect,omitempty"`
		CloseURL        string `json:"closeURL,omitempty"`
		RefundAddresses []map[string]struct {
			Type  string `json:"type"`
			Date  string `json:"date"`
			Email string `json:"email"`
		} `json:"refundAddresses,omitempty"`
		RefundAddressRequestPending bool   `json:"refundAddressRequestPending"`
		BuyerProvidedEmail          string `json:"buyerProvidedEmail,omitempty"`
		BuyerProvidedInfo           *struct {
			EmailAddress                string `json:"emailAddress"`
			SelectedWallet              string `json:"selectedWallet"`
			SelectedTransactionCurrency string `json:"selectedTransactionCurrency"`
		} `json:"buyerProvidedInfo,omitempty"`
		PaymentSubtotals        map[string]int64                      `json:"paymentSubtotals,omitempty"`
		PaymentTotals           map[string]int64                      `json:"paymentTotals,omitempty"`
		PaymentDisplayTotals    map[string]string                     `json:"paymentDisplayTotals,omitempty"`
		PaymentDisplaySubTotals map[string]string                     `json:"paymentDisplaySubTotals,omitempty"`
		ExchangeRates           map[string]map[string]decimal.Decimal `json:"exchangeRates,omitempty"`
		MinerFees               map[string]struct {
			SatoshisPerByte decimal.Decimal `json:"satoshisPerByte"`
			TotalFee        decimal.Decimal `json:"totalFee"`
			FiatAmount      decimal.Decimal `json:"fiatAmount"`
		} `json:"minerFees,omitempty"`
		Shopper            any    `json:"shopper,omitempty"`
		JsonPayProRequired bool   `json:"jsonPayProRequired"`
		MerchantName       string `json:"merchantName,omitempty"`
		BitpayIdRequired   *bool  `json:"bitpayIdRequired,omitempty"`
		ItemizedDetails    []struct {
			Amount      int64  `json:"amount"`
			Description string `json:"description"`
			IsFee       bool   `json:"isFee"`
		} `json:"itemizedDetails,omitempty"`
		SupportedTransactionCurrencies map[string]struct {
			Enabled bool   `json:"enabled"`
			Reason  string `json:"reason,omitempty"`
		} `json:"supportedTransactionCurrencies,omitempty"`
		PaymentCodes   map[string]map[string]string `json:"paymentCodes,omitempty"`
		UniversalCodes *struct {
			PaymentString string `json:"paymentString"`
		} `json:"universalCodes,omitempty"`
		Token string `json:"token,omitempty"`
	} `json:"data"`
}
