package response

import "github.com/shopspring/decimal"

// Create an Invoice
type CreateInvoice struct {
	Facade string `json:"facade"` // Invoice facade used by the token (for example: merchant).
	Data   struct {
		URL                 string           `json:"url"`                           // Hosted checkout URL for the invoice.
		PosData             string           `json:"posData,omitempty"`             // Merchant POS metadata echoed back by BitPay.
		Status              string           `json:"status"`                        // Current invoice status.
		Price               decimal.Decimal  `json:"price"`                         // Invoice amount in fiat currency.
		Currency            string           `json:"currency"`                      // Fiat currency code for `price` (ISO-4217).
		ItemDesc            string           `json:"itemDesc,omitempty"`            // Human-readable item description.
		OrderID             string           `json:"orderId,omitempty"`             // Merchant order identifier.
		InvoiceTime         int64            `json:"invoiceTime"`                   // Invoice creation time (Unix ms).
		ExpirationTime      int64            `json:"expirationTime"`                // Invoice expiration time (Unix ms).
		CurrentTime         int64            `json:"currentTime"`                   // BitPay server time when response is generated (Unix ms).
		GUID                string           `json:"guid,omitempty"`                // Internal GUID associated with the invoice.
		ID                  string           `json:"id"`                            // BitPay invoice identifier.
		LowFeeDetected      bool             `json:"lowFeeDetected"`                // Whether an underpaid/low-fee payment was detected.
		AmountPaid          decimal.Decimal  `json:"amountPaid"`                    // Total amount paid toward the invoice.
		DisplayAmountPaid   string           `json:"displayAmountPaid,omitempty"`   // Formatted amount paid for display purposes.
		ExceptionStatus     bool             `json:"exceptionStatus"`               // Whether the invoice is in an exception state.
		TargetConfirmations *int64           `json:"targetConfirmations,omitempty"` // Required blockchain confirmations for payment finality.
		Transactions        []map[string]any `json:"transactions,omitempty"`        // Payment transactions recorded for the invoice.
		TransactionSpeed    string           `json:"transactionSpeed,omitempty"`    // Requested payment speed policy (for example: high, medium, low).
		Buyer               *struct {
			Email string `json:"email"` // Buyer email address.
		} `json:"buyer,omitempty"` // Buyer contact object.
		RedirectURL     string `json:"redirectURL,omitempty"`  // URL to redirect buyer after successful payment.
		AutoRedirect    *bool  `json:"autoRedirect,omitempty"` // Whether checkout auto-redirects on completion.
		CloseURL        string `json:"closeURL,omitempty"`     // URL where buyer returns when checkout closes.
		RefundAddresses []map[string]struct {
			Type  string `json:"type"`  // Channel type used to collect refund address.
			Date  string `json:"date"`  // Collection timestamp (ISO-8601).
			Email string `json:"email"` // Email tied to the refund address request.
		} `json:"refundAddresses,omitempty"` // Refund address collection records.
		RefundAddressRequestPending bool   `json:"refundAddressRequestPending"`  // Whether refund address collection is pending.
		BuyerProvidedEmail          string `json:"buyerProvidedEmail,omitempty"` // Email supplied directly by buyer.
		BuyerProvidedInfo           *struct {
			EmailAddress                string `json:"emailAddress"`                // Buyer email captured during checkout.
			SelectedWallet              string `json:"selectedWallet"`              // Wallet selected by buyer.
			SelectedTransactionCurrency string `json:"selectedTransactionCurrency"` // Crypto currency selected by buyer.
		} `json:"buyerProvidedInfo,omitempty"` // Buyer-provided checkout information.
		PaymentSubtotals        map[string]int64                      `json:"paymentSubtotals,omitempty"`        // Subtotal per payment currency (usually in smallest unit).
		PaymentTotals           map[string]int64                      `json:"paymentTotals,omitempty"`           // Total per payment currency including fees.
		PaymentDisplayTotals    map[string]string                     `json:"paymentDisplayTotals,omitempty"`    // Formatted total per payment currency.
		PaymentDisplaySubTotals map[string]string                     `json:"paymentDisplaySubTotals,omitempty"` // Formatted subtotal per payment currency.
		ExchangeRates           map[string]map[string]decimal.Decimal `json:"exchangeRates,omitempty"`           // Exchange rates matrix used for invoice pricing.
		MinerFees               map[string]struct {
			SatoshisPerByte decimal.Decimal `json:"satoshisPerByte"` // Estimated satoshis-per-byte fee rate.
			TotalFee        decimal.Decimal `json:"totalFee"`        // Estimated network fee in transaction currency.
			FiatAmount      decimal.Decimal `json:"fiatAmount"`      // Estimated network fee converted to invoice fiat.
		} `json:"minerFees,omitempty"` // Miner fee estimates keyed by crypto currency.
		Shopper            any    `json:"shopper,omitempty"`          // Shopper profile object returned by BitPay.
		JsonPayProRequired bool   `json:"jsonPayProRequired"`         // Whether JSON Payment Protocol is required.
		MerchantName       string `json:"merchantName,omitempty"`     // Display name of merchant account.
		BitpayIdRequired   *bool  `json:"bitpayIdRequired,omitempty"` // Whether BitPay ID login is required.
		ItemizedDetails    []struct {
			Amount      int64  `json:"amount"`      // Line amount (usually in smallest unit expected by API response).
			Description string `json:"description"` // Line item description.
			IsFee       bool   `json:"isFee"`       // Whether line item represents a fee.
		} `json:"itemizedDetails,omitempty"` // Invoice line-item breakdown.
		SupportedTransactionCurrencies map[string]struct {
			Enabled bool   `json:"enabled"`          // Whether the payment currency is currently available.
			Reason  string `json:"reason,omitempty"` // Reason when a currency is unavailable.
		} `json:"supportedTransactionCurrencies,omitempty"` // Supported crypto currencies for this invoice.
		PaymentCodes   map[string]map[string]string `json:"paymentCodes,omitempty"` // Payment code values per chain/currency.
		UniversalCodes *struct {
			PaymentString string `json:"paymentString"` // Unified payment string for wallet deep links.
		} `json:"universalCodes,omitempty"` // Cross-wallet universal payment payload.
		Token string `json:"token,omitempty"` // Access token associated with this invoice response.
	} `json:"data"`
}

// Retrieve an Invoice
type RetrieveInvoice struct {
	Facade string `json:"facade"` // Invoice facade used by the token (for example: merchant).
	Data   struct {
		URL                 string           `json:"url"`                           // Hosted checkout URL for the invoice.
		PosData             string           `json:"posData,omitempty"`             // Merchant POS metadata echoed back by BitPay.
		Status              string           `json:"status"`                        // Current invoice status.
		Price               decimal.Decimal  `json:"price"`                         // Invoice amount in fiat currency.
		Currency            string           `json:"currency"`                      // Fiat currency code for `price` (ISO-4217).
		OrderID             string           `json:"orderId,omitempty"`             // Merchant order identifier.
		InvoiceTime         int64            `json:"invoiceTime"`                   // Invoice creation time (Unix ms).
		ExpirationTime      int64            `json:"expirationTime"`                // Invoice expiration time (Unix ms).
		CurrentTime         int64            `json:"currentTime"`                   // BitPay server time when response is generated (Unix ms).
		GUID                string           `json:"guid,omitempty"`                // Internal GUID associated with the invoice.
		ID                  string           `json:"id"`                            // BitPay invoice identifier.
		LowFeeDetected      bool             `json:"lowFeeDetected"`                // Whether an underpaid/low-fee payment was detected.
		AmountPaid          decimal.Decimal  `json:"amountPaid"`                    // Total amount paid toward the invoice.
		DisplayAmountPaid   string           `json:"displayAmountPaid,omitempty"`   // Formatted amount paid for display purposes.
		ExceptionStatus     bool             `json:"exceptionStatus"`               // Whether the invoice is in an exception state.
		TargetConfirmations *int64           `json:"targetConfirmations,omitempty"` // Required blockchain confirmations for payment finality.
		Transactions        []map[string]any `json:"transactions,omitempty"`        // Payment transactions recorded for the invoice.
		TransactionSpeed    string           `json:"transactionSpeed,omitempty"`    // Requested payment speed policy (for example: high, medium, low).
		Buyer               *struct {
			Email string `json:"email"` // Buyer email address.
		} `json:"buyer,omitempty"` // Buyer contact object.
		RedirectURL     string `json:"redirectURL,omitempty"`  // URL to redirect buyer after successful payment.
		AutoRedirect    *bool  `json:"autoRedirect,omitempty"` // Whether checkout auto-redirects on completion.
		CloseURL        string `json:"closeURL,omitempty"`     // URL where buyer returns when checkout closes.
		RefundAddresses []map[string]struct {
			Type  string `json:"type"`  // Channel type used to collect refund address.
			Date  string `json:"date"`  // Collection timestamp (ISO-8601).
			Email string `json:"email"` // Email tied to the refund address request.
		} `json:"refundAddresses,omitempty"` // Refund address collection records.
		RefundAddressRequestPending bool   `json:"refundAddressRequestPending"`  // Whether refund address collection is pending.
		BuyerProvidedEmail          string `json:"buyerProvidedEmail,omitempty"` // Email supplied directly by buyer.
		BuyerProvidedInfo           *struct {
			EmailAddress                string `json:"emailAddress"`                // Buyer email captured during checkout.
			SMS                         string `json:"sms,omitempty"`               // Buyer phone number collected for notifications.
			SMSVerified                 bool   `json:"smsVerified"`                 // Whether buyer SMS number is verified.
			SelectedWallet              string `json:"selectedWallet"`              // Wallet selected by buyer.
			SelectedTransactionCurrency string `json:"selectedTransactionCurrency"` // Crypto currency selected by buyer.
		} `json:"buyerProvidedInfo,omitempty"` // Buyer-provided checkout information.
		PaymentSubtotals        map[string]decimal.Decimal            `json:"paymentSubtotals,omitempty"`        // Subtotal per payment currency.
		PaymentTotals           map[string]decimal.Decimal            `json:"paymentTotals,omitempty"`           // Total per payment currency including fees.
		PaymentDisplayTotals    map[string]string                     `json:"paymentDisplayTotals,omitempty"`    // Formatted total per payment currency.
		PaymentDisplaySubTotals map[string]string                     `json:"paymentDisplaySubTotals,omitempty"` // Formatted subtotal per payment currency.
		ExchangeRates           map[string]map[string]decimal.Decimal `json:"exchangeRates,omitempty"`           // Exchange rates matrix used for invoice pricing.
		MinerFees               map[string]struct {
			SatoshisPerByte decimal.Decimal `json:"satoshisPerByte"` // Estimated satoshis-per-byte fee rate.
			TotalFee        decimal.Decimal `json:"totalFee"`        // Estimated network fee in transaction currency.
			FiatAmount      decimal.Decimal `json:"fiatAmount"`      // Estimated network fee converted to invoice fiat.
		} `json:"minerFees,omitempty"` // Miner fee estimates keyed by crypto currency.
		Shopper            any    `json:"shopper,omitempty"`          // Shopper profile object returned by BitPay.
		JsonPayProRequired bool   `json:"jsonPayProRequired"`         // Whether JSON Payment Protocol is required.
		MerchantName       string `json:"merchantName,omitempty"`     // Display name of merchant account.
		BitpayIdRequired   *bool  `json:"bitpayIdRequired,omitempty"` // Whether BitPay ID login is required.
		IsCancelled        *bool  `json:"isCancelled,omitempty"`      // Whether the invoice has been canceled.
		ItemizedDetails    []struct {
			Amount      decimal.Decimal `json:"amount"`      // Line amount.
			Description string          `json:"description"` // Line item description.
			IsFee       bool            `json:"isFee"`       // Whether line item represents a fee.
		} `json:"itemizedDetails,omitempty"` // Invoice line-item breakdown.
		TransactionCurrency            string `json:"transactionCurrency,omitempty"` // Settled/selected transaction currency for payments.
		SupportedTransactionCurrencies map[string]struct {
			Enabled bool   `json:"enabled"`          // Whether the payment currency is currently available.
			Reason  string `json:"reason,omitempty"` // Reason when a currency is unavailable.
		} `json:"supportedTransactionCurrencies,omitempty"` // Supported crypto currencies for this invoice.
		PaymentCodes   map[string]map[string]string `json:"paymentCodes,omitempty"` // Payment code values per chain/currency.
		UniversalCodes *struct {
			PaymentString    string `json:"paymentString"`    // Unified payment string for wallet deep links.
			VerificationLink string `json:"verificationLink"` // Link to verify universal payment code.
		} `json:"universalCodes,omitempty"` // Cross-wallet universal payment payload.
		Token string `json:"token,omitempty"` // Access token associated with this invoice response.
	} `json:"data"`
}

// Create a Refund Request
type CreateRefundRequest struct {
	Data struct {
		ID                     string          `json:"id"`                     // Unique refund request identifier.
		Invoice                string          `json:"invoice"`                // Related BitPay invoice ID.
		Reference              string          `json:"reference"`              // Merchant-provided reference for reconciliation.
		Status                 string          `json:"status"`                 // Current refund request status (for example: pending, success, canceled).
		Amount                 decimal.Decimal `json:"amount"`                 // Refund amount in `currency`.
		TransactionCurrency    string          `json:"transactionCurrency"`    // Cryptocurrency used for the refund transaction.
		TransactionAmount      decimal.Decimal `json:"transactionAmount"`      // Refund amount in `transactionCurrency`.
		TransactionRefundFee   decimal.Decimal `json:"transactionRefundFee"`   // Network/miner fee amount in `transactionCurrency`.
		Currency               string          `json:"currency"`               // Reference fiat currency for the refund request.
		LastRefundNotification string          `json:"lastRefundNotification"` // Last refund notification timestamp (ISO-8601).
		RefundFee              decimal.Decimal `json:"refundFee"`              // Fee applied to the refund in `currency`.
		Immediate              bool            `json:"immediate"`              // Whether funds are removed immediately from merchant ledger.
		BuyerPaysRefundFee     bool            `json:"buyerPaysRefundFee"`     // Whether the buyer is charged the refund fee.
		RequestDate            string          `json:"requestDate"`            // Refund request creation timestamp (ISO-8601).
	} `json:"data"`
}
