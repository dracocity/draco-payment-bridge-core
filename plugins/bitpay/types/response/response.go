package response

import "github.com/shopspring/decimal"

// Create an Invoice
type CreateInvoice struct {
	Facade string `json:"facade"`
	Data   struct {
		URL                         string          `json:"url"`
		PosData                     string          `json:"posData,omitempty"`
		Status                      string          `json:"status"`
		Price                       decimal.Decimal `json:"price"`
		Currency                    string          `json:"currency"`
		ItemDesc                    string          `json:"itemDesc,omitempty"`
		OrderID                     string          `json:"orderId,omitempty"`
		InvoiceTime                 int64           `json:"invoiceTime"`
		ExpirationTime              int64           `json:"expirationTime"`
		CurrentTime                 int64           `json:"currentTime"`
		GUID                        string          `json:"guid,omitempty"`
		ID                          string          `json:"id"`
		LowFeeDetected              bool            `json:"lowFeeDetected"`
		AmountPaid                  decimal.Decimal `json:"amountPaid"`
		DisplayAmountPaid           string          `json:"displayAmountPaid,omitempty"`
		ExceptionStatus             bool            `json:"exceptionStatus"`
		CloseURL                    string          `json:"closeURL,omitempty"`
		RedirectURL                 string          `json:"redirectURL,omitempty"`
		AutoRedirect                *bool           `json:"autoRedirect,omitempty"`
		RefundAddressRequestPending bool            `json:"refundAddressRequestPending"`
		BuyerProvidedInfo           *struct {
			EmailAddress                string `json:"emailAddress"`
			SelectedWallet              string `json:"selectedWallet"`
			SelectedTransactionCurrency string `json:"selectedTransactionCurrency"`
		} `json:"buyerProvidedInfo,omitempty"`
		PaymentSubtotals               map[string]int64                      `json:"paymentSubtotals,omitempty"`
		PaymentTotals                  map[string]int64                      `json:"paymentTotals,omitempty"`
		PaymentDisplayTotals           map[string]string                     `json:"paymentDisplayTotals,omitempty"`
		PaymentDisplaySubTotals        map[string]string                     `json:"paymentDisplaySubTotals,omitempty"`
		ExchangeRates                  map[string]map[string]decimal.Decimal `json:"exchangeRates,omitempty"`
		SupportedTransactionCurrencies map[string]struct {
			Enabled bool   `json:"enabled"`
			Reason  string `json:"reason,omitempty"`
		} `json:"supportedTransactionCurrencies,omitempty"`
		MinerFees map[string]struct {
			SatoshisPerByte decimal.Decimal `json:"satoshisPerByte"`
			TotalFee        decimal.Decimal `json:"totalFee"`
			FiatAmount      decimal.Decimal `json:"fiatAmount"`
		} `json:"minerFees,omitempty"`
		JsonPayProRequired bool `json:"jsonPayProRequired"`
		ItemizedDetails    []struct {
			Amount      int64  `json:"amount"`
			Description string `json:"description"`
			IsFee       bool   `json:"isFee"`
		} `json:"itemizedDetails,omitempty"`
		UniversalCodes *struct {
			PaymentString string `json:"paymentString"`
		} `json:"universalCodes,omitempty"`
		MerchantName string                       `json:"merchantName,omitempty"`
		PaymentCodes map[string]map[string]string `json:"paymentCodes,omitempty"`
		Token        string                       `json:"token,omitempty"`
	} `json:"data"`
}

// Retrieve an Invoice
type RetrieveInvoice struct {
	Facade string `json:"facade"`
	Data   struct {
		URL                         string          `json:"url"`
		PosData                     string          `json:"posData,omitempty"`
		Status                      string          `json:"status"`
		Price                       decimal.Decimal `json:"price"`
		Currency                    string          `json:"currency"`
		OrderID                     string          `json:"orderId,omitempty"`
		InvoiceTime                 int64           `json:"invoiceTime"`
		ExpirationTime              int64           `json:"expirationTime"`
		CurrentTime                 int64           `json:"currentTime"`
		GUID                        string          `json:"guid,omitempty"`
		ID                          string          `json:"id"`
		LowFeeDetected              bool            `json:"lowFeeDetected"`
		AmountPaid                  decimal.Decimal `json:"amountPaid"`
		DisplayAmountPaid           string          `json:"displayAmountPaid,omitempty"`
		ExceptionStatus             bool            `json:"exceptionStatus"`
		CloseURL                    string          `json:"closeURL,omitempty"`
		RedirectURL                 string          `json:"redirectURL,omitempty"`
		AutoRedirect                *bool           `json:"autoRedirect,omitempty"`
		RefundAddressRequestPending bool            `json:"refundAddressRequestPending"`
		BuyerProvidedInfo           *struct {
			EmailAddress                string `json:"emailAddress"`
			SMS                         string `json:"sms,omitempty"`
			SMSVerified                 bool   `json:"smsVerified"`
			SelectedWallet              string `json:"selectedWallet"`
			SelectedTransactionCurrency string `json:"selectedTransactionCurrency"`
		} `json:"buyerProvidedInfo,omitempty"`
		PaymentSubtotals               map[string]decimal.Decimal            `json:"paymentSubtotals,omitempty"`
		PaymentTotals                  map[string]decimal.Decimal            `json:"paymentTotals,omitempty"`
		PaymentDisplayTotals           map[string]string                     `json:"paymentDisplayTotals,omitempty"`
		PaymentDisplaySubTotals        map[string]string                     `json:"paymentDisplaySubTotals,omitempty"`
		ExchangeRates                  map[string]map[string]decimal.Decimal `json:"exchangeRates,omitempty"`
		SupportedTransactionCurrencies map[string]struct {
			Enabled bool   `json:"enabled"`
			Reason  string `json:"reason,omitempty"`
		} `json:"supportedTransactionCurrencies,omitempty"`
		MinerFees map[string]struct {
			SatoshisPerByte decimal.Decimal `json:"satoshisPerByte"`
			TotalFee        decimal.Decimal `json:"totalFee"`
			FiatAmount      decimal.Decimal `json:"fiatAmount"`
		} `json:"minerFees,omitempty"`
		JsonPayProRequired bool `json:"jsonPayProRequired"`
		ItemizedDetails    []struct {
			Amount      decimal.Decimal `json:"amount"`
			Description string          `json:"description"`
			IsFee       bool            `json:"isFee"`
		} `json:"itemizedDetails,omitempty"`
		TransactionCurrency string `json:"transactionCurrency,omitempty"`
		UniversalCodes      *struct {
			PaymentString string `json:"paymentString"`
		} `json:"universalCodes,omitempty"`
		MerchantName string                       `json:"merchantName,omitempty"`
		PaymentCodes map[string]map[string]string `json:"paymentCodes,omitempty"`
		Token        string                       `json:"token,omitempty"`
	} `json:"data"`
}
