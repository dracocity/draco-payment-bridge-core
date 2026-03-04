package response

// CreateInvoice
type CreateInvoice struct {
	Invoices []struct {
		ID           string `json:"id"`           // The id of the created invoic
		Link         string `json:"link"`         // The link to the created invoice
		CheckoutLink string `json:"checkoutLink"` // The link to the checkout app
		Payment      *struct {
			PaymentID         string `json:"paymentId"` // the id of the payment
			Expires           string `json:"expires"`   // the timestamp when the payment expires and new payments will no longer be accepted, Format: date-time
			PaymentCurrencies []struct {
				Currency *struct {
					ID     string `json:"id"`     // the unique id of the currency on the CoinPayments platform
					Type   string `json:"type"`   // Type of currency. Enum: crypto, token, fiat
					Symbol string `json:"symbol"` // ticker symbol for the currency
					Name   string `json:"name"`   // the name of the currency
					Logo   *struct {
						ImageURL string `json:"imageUrl"` // Link to a CoinPayments hosted image for a currency in a SVG format
					} `json:"logo"` // Contains the logo URLs for a currency
					DecimalPlaces int32    `json:"decimalPlaces"` // the number of digits after the decimal separator
					Rank          int32    `json:"rank"`          // relative ordering/ranking of this currency
					Status        string   `json:"status"`        // The status of the currency. Enum: active, underMaintenance, deleted
					Capabilities  []string `json:"capabilities"`  // capabilities of this currency on the CoinPayments platform
					URLs          *struct {
						Websites  []string `json:"websites"`  // the websites for the currency
						Explorers []string `json:"explorers"` // the explorers for the currency (if crypto or a token)
					} `json:"urls"` // Contains various URLs for a currency
					RequiredConfirmations int64 `json:"requiredConfirmations"` // required confirmations before transactions are treated as confirmed
					IsEnabledForPayment   *bool `json:"isEnabledForPayment"`   // whether currency is enabled as a payment option
				} `json:"currency"` // Represents money, in any form (e.g. crypto/fiat currencies and tokens), supported on the CoinPayments platform in some capacity (e.g. holding in a wallet and/or payment processing).
				IsDisabled               *bool  `json:"isDisabled"`               // flag indicating whether this currency is currently unavailable (e.g. node or services down)
				Amount                   string `json:"amount"`                   // the total amount due to pay the remainder of the invoices balance in this currency
				ApproximateNetworkAmount string `json:"approximateNetworkAmount"` // Represents the approximate amount of the payment in network-native currency terms. This property provides a network-based estimation converted into a string format.
				RemainingAmount          string `json:"remainingAmount"`          // The remaining amount to be paid for the current invoice in the specified CoinPayments.Api.Models.Merchant.InvoicePaymentCurrencyV2Dto.Currency.
			} `json:"paymentCurrencies"` // A collection of supported payment currencies, including relevant details such as currency information, amounts, and status for a specific invoice payment.
			RefundEmail string `json:"refundEmail"` // The email address associated with the recipient to which any refunds should be sent.
		} `json:"payment"` // Represents the payment details associated with an invoice, including payment ID, expiration information, supported currencies, and refund details.
		HotWallet *struct {
			Currency *struct {
				ID     string `json:"id"`     // the unique id of the currency on the CoinPayments platform
				Type   string `json:"type"`   // Type of currency. Enum: crypto, token, fiat
				Symbol string `json:"symbol"` // ticker symbol for the currency
				Name   string `json:"name"`   // the name of the currency
				Logo   *struct {
					ImageURL string `json:"imageUrl"` // Link to a CoinPayments hosted image for a currency in a SVG format
				} `json:"logo"` // Contains the logo URLs for a currency
				DecimalPlaces int32    `json:"decimalPlaces"` // the number of digits after the decimal separator
				Rank          int32    `json:"rank"`          // relative ordering/ranking of this currency
				Status        string   `json:"status"`        // The status of the currency. Enum: active, underMaintenance, deleted
				Capabilities  []string `json:"capabilities"`  // capabilities of this currency on the CoinPayments platform
				URLs          *struct {
					Websites  []string `json:"websites"`  // the websites for the currency
					Explorers []string `json:"explorers"` // the explorers for the currency (if crypto or a token)
				} `json:"urls"` // Contains various URLs for a currency
				RequiredConfirmations int64 `json:"requiredConfirmations"` // required confirmations before transactions are treated as confirmed
				IsEnabledForPayment   *bool `json:"isEnabledForPayment"`   // whether currency is enabled as a payment option
			} `json:"currency"` // Represents money, in any form (e.g. crypto/fiat currencies and tokens), supported on the CoinPayments platform in some capacity (e.g. holding in a wallet and/or payment processing).
			Amount          string `json:"amount"`          // the amount due to be paid in this currency
			RemainingAmount string `json:"remainingAmount"` // Represents the remaining amount to be paid for an invoice, expressed as a string.
			Addresses       *struct {
				Address string `json:"address"` // the raw payment address
				BIP21   string `json:"biP21"`   // the BIP21 payment code, if available
			} `json:"addresses"` // Encapsulates the addresses where payments to an invoice can be sent
			Expires string `json:"expires"` // the timestamp when the payment expires and new payments will no longer be accepted, Format: date-time
		} `json:"hotWallet"` // Represents the details of a currency payment for an invoice, including the currency, payment amount, remaining amount to be paid, payment addresses, and expiration information.
	} `json:"invoices"` // the id of the created invoice
}
