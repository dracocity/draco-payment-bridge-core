package response

// Create an Invoice
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

// Get an Invoice
type GetInvoice struct {
	ID              string `json:"id"`              // the CoinPayments id for the invoice
	InvoiceID       string `json:"invoiceId"`       // the optional API caller provided external invoice number
	InvoiceIDSuffix string `json:"invoiceIdSuffix"` // the optional numeric suffix for the invoiceId
	Created         string `json:"created"`         // the timestamp when the invoice was created, Format: date-time
	InvoiceDate     string `json:"invoiceDate"`     // the date of the invoice, Format: date-time
	DueDate         string `json:"dueDate"`         // optional due date of the invoice, Format: date-time
	Confirmed       string `json:"confirmed"`       // the timestamp when the invoice was confirmed, Format: date-time
	Completed       string `json:"completed"`       // the timestamp when the invoice was completed, Format: date-time
	Cancelled       string `json:"cancelled"`       // the timestamp when the invoice was manually cancelled, Format: date-time
	Expires         string `json:"expires"`         // the timestamp when the invoice expires, Format: date-time
	Currency        *struct {
		ID     string `json:"id"`     // the unique id of the currency on the CoinPayments platform
		Symbol string `json:"symbol"` // ticker symbol for the currency
		Name   string `json:"name"`   // the name of the currency
		Token  *struct {
			Name            string `json:"name"`            // name for the token, if available
			Symbol          string `json:"symbol"`          // ticker symbol for the token, if available
			ContractAddress string `json:"contractAddress"` // the address of the contract
			DecimalPlaces   int32  `json:"decimalPlaces"`   // the number of digits after the decimal separator
		} `json:"token"` // token information
		Logo *struct {
			ImageURL  string `json:"imageUrl"`  // Link to a CoinPayments hosted image for a currency; default size is 64x64 and can be changed to 32, 64, 128, or 200.
			VectorURL string `json:"vectorUrl"` // If available then the link to a CoinPayments hosted vector image (SVG) for the currency.
		} `json:"logo"` // Contains the logo URLs for a currency
		DecimalPlaces int32 `json:"decimalPlaces"` // the number of digits after the decimal separator
	} `json:"currency"` // invoice currency information
	Merchant *struct {
		ID                 string `json:"id"`                 // the CoinPayments id of the merchant
		Name               string `json:"name"`               // the business name of the merchant
		UBOName            string `json:"uboName"`            // full name of the Ultimate Beneficiary Owner (UBO) of the business
		WebsiteURL         string `json:"websiteUrl"`         // the url to the merchants website
		Country            string `json:"country"`            // the merchant's business country of registration
		LogoURL            string `json:"logoUrl"`            // the url to the merchant's logo
		Email              string `json:"email"`              // the merchant's business email
		Address            string `json:"address"`            // the merchant's business address
		Phone              string `json:"phone"`              // the phone number of the business
		Description        string `json:"description"`        // the description of the merchant
		RegistrationNumber string `json:"registrationNumber"` // the business registration number
	} `json:"merchant"` // merchant information
	MerchantOptions *struct {
		ShowAddress            *bool  `json:"showAddress"`
		ShowEmail              *bool  `json:"showEmail"`
		ShowPhone              *bool  `json:"showPhone"`
		ShowRegistrationNumber *bool  `json:"showRegistrationNumber"`
		ShowTaxID              *bool  `json:"showTaxId"`
		AdditionalInfo         string `json:"additionalInfo"`
	} `json:"merchantOptions"` // merchant display options on the invoice
	Buyer *struct {
		CompanyName string `json:"companyName"` // the name of the buyer's company
		Name        *struct {
			FirstName string `json:"firstName"` // the given, or first, name
			LastName  string `json:"lastName"`  // the surname or family name
		} `json:"name"` // buyer name object
		EmailAddress string `json:"emailAddress"` // the email address of the buyer
		PhoneNumber  string `json:"phoneNumber"`  // the phone number
		Address      *struct {
			Address1         string `json:"address1"`         // the first line of the address
			Address2         string `json:"address2"`         // the second line of the address
			Address3         string `json:"address3"`         // the third line of the address
			ProvinceOrState  string `json:"provinceOrState"`  // province, state, or equivalent sub-division
			City             string `json:"city"`             // the city, town, or village
			SuburbOrDistrict string `json:"suburbOrDistrict"` // neighborhood, suburb, or district
			CountryCode      string `json:"countryCode"`      // two-character ISO-3166-1 country code
			PostalCode       string `json:"postalCode"`       // postal or ZIP code
		} `json:"address"` // buyer address object
	} `json:"buyer"` // buyer information
	Description string `json:"description"` // the purchase description
	Items       []struct {
		CustomID    string `json:"customId"`    // the API caller provided external ID for the item. Appears on the Merchant dashboard and reports only.
		SKU         string `json:"sku"`         // the stock keeping unit (SKU) of the item
		Name        string `json:"name"`        // the name or title of the item
		Description string `json:"description"` // the detailed description of the item
		Quantity    *struct {
			Value int64  `json:"value"` // the quantity of the item
			Type  string `json:"type"`  // quantity type. Enum: hours, quantity
		} `json:"quantity"` // Represents the quantity details for a line item in the CoinPayments merchant system. This class defines both the value of the quantity and its associated type.
		OriginalAmount string `json:"originalAmount"` // the original total price of the item if CoinPayments.Services.Merchants.Dtos.LineItemV2Dto.Amount represents a discounted price
		Amount         string `json:"amount"`         // the subtotal price of the item (note: this is not a per unit price but the total price for the total quantity)
		Tax            string `json:"tax"`            // the total taxes charged on this item
	} `json:"items"` // items and/or services included in the invoice
	Amount *struct {
		Breakdown *struct {
			Subtotal string `json:"subtotal"` // the subtotal for all items, required if items are included and must equal sum(items.amount)
			Shipping string `json:"shipping"` // the optional shipping fee for all items within the invoice
			Handling string `json:"handling"` // the optional handling fee for all items within the invoice
			TaxTotal string `json:"taxTotal"` // the total taxes charged on the invoice, required when item taxes are present and must equal sum(items.tax)
			Discount string `json:"discount"` // the optional discount for all items within the invoice
		} `json:"breakdown"` // no further information available
		Total string `json:"total"` // total amount associated with the invoice, including charges, fees, and adjustments
	} `json:"amount"` // total invoice amount breakdown
	Shipping *struct {
		Method      string `json:"method"`      // the shipping method
		CompanyName string `json:"companyName"` // the company name of the party to ship the items to
		Name        *struct {
			FirstName string `json:"firstName"` // the given, or first, name
			LastName  string `json:"lastName"`  // the surname or family name
		} `json:"name"` // shipping recipient name object
		EmailAddress string `json:"emailAddress"` // the email address of the party to ship the items to
		PhoneNumber  string `json:"phoneNumber"`  // the phone number
		Address      *struct {
			Address1         string `json:"address1"`         // the first line of the address
			Address2         string `json:"address2"`         // the second line of the address
			Address3         string `json:"address3"`         // the third line of the address
			ProvinceOrState  string `json:"provinceOrState"`  // province, state, or equivalent sub-division
			City             string `json:"city"`             // the city, town, or village
			SuburbOrDistrict string `json:"suburbOrDistrict"` // neighborhood, suburb, or district
			CountryCode      string `json:"countryCode"`      // two-character ISO-3166-1 country code
			PostalCode       string `json:"postalCode"`       // postal or ZIP code
		} `json:"address"` // shipping address object
		HasData bool `json:"hasData"` // indicates whether shipping data is present
	} `json:"shipping"` // shipping address and method
	CustomData                 any    `json:"customData"`                 // custom data attached to the invoice
	Status                     string `json:"status"`                     // invoice status
	RequireBuyerNameAndEmail   bool   `json:"requireBuyerNameAndEmail"`   // whether buyer name and email are required
	BuyerDataCollectionMessage string `json:"buyerDataCollectionMessage"` // message shown when collecting buyer data
	Notes                      string `json:"notes"`                      // notes for the merchant only
	NotesToRecipient           string `json:"notesToRecipient"`           // additional information shown to the buyer
	TermsAndConditions         string `json:"termsAndConditions"`         // terms and conditions text
	EmailDelivery              *struct {
		To  string `json:"to"`  // the email `to` field, multiple addresses separated by semicolons
		Cc  string `json:"cc"`  // the email `cc` field, multiple addresses separated by semicolons
		Bcc string `json:"bcc"` // the email `bcc` field, multiple addresses separated by semicolons
	} `json:"emailDelivery"` // email delivery options
	IsEmailDelivery bool `json:"isEmailDelivery"` // indicates if invoice was delivered by email or to be delivered by email
	Metadata        *struct {
		Integration string `json:"integration"` // the integration from which the invoice was created
		Hostname    string `json:"hostname"`    // the hostname on which the invoice was created
	} `json:"metadata"` // metadata related to the invoice
	PONumber      string `json:"poNumber"` // purchase order number associated with the invoice
	PayoutDetails *struct {
		PaidTransactions []struct {
			Hash   string `json:"hash"` // TxHash from the blockchain if the transaction is an external transaction.
			Amount *struct {
				DisplayValue   string  `json:"displayValue"`   // the value formatted for display
				Value          string  `json:"value"`          // the amount of money in the currency's base (smallest) monetary unit
				CurrencyID     string  `json:"currencyId"`     // the currency the monetary value is represented in
				ValueAsDecimal float64 `json:"valueAsDecimal"` // the monetary value represented as decimal for precise calculations and storage
			} `json:"amount"` // Monetary value (an amount with a currency).
			ConversionID int64 `json:"conversionId"` // Represents the unique identifier for a currency conversion associated with the transaction, if the transaction involves a currency conversion process. Format: int64
		} `json:"paidTransactions"` // paid transaction details
		PaidDate                string  `json:"paidDate"`                // The date and time when the payment was made. Format: date-time
		CompletedTxID           string  `json:"completedTxId"`           // The ID of the completed transaction.
		ExternalAddress         string  `json:"externalAddress"`         // The external address where the payout is deposited
		DestinationCurrencyID   string  `json:"destinationCurrencyId"`   // The currency ID of the destination for the payout
		ExpectedDisplayValue    string  `json:"expectedDisplayValue"`    // The expected display value of the payout.
		SourceCurrencyID        string  `json:"sourceCurrencyId"`        // The currency ID of the source for the payout
		DestinationWalletID     string  `json:"destinationWalletId"`     // The ID of the destination wallet for the payout
		IsConversion            bool    `json:"isConversion"`            // Indicates whether a currency conversion is involved in the payout
		ConversionProgress      float64 `json:"conversionProgress"`      // The progress status of the currency conversion. Format: double
		SettlementModeErrorCode int32   `json:"settlementModeErrorCode"` // Settlement mode error code for payout processing failures.
		DestinationAmount       *struct {
			Amount *struct {
				DisplayValue   string  `json:"displayValue"`   // the value formatted for display
				Value          string  `json:"value"`          // the amount of money in the currency's base (smallest) monetary unit
				CurrencyID     string  `json:"currencyId"`     // the currency the monetary value is represented in
				ValueAsDecimal float64 `json:"valueAsDecimal"` // the monetary value represented as decimal for precise calculations and storage
			} `json:"amount"` // Monetary value (an amount with a currency).
			NativeAmount *struct {
				DisplayValue   string  `json:"displayValue"`   // the value formatted for display
				Value          string  `json:"value"`          // the amount of money in the currency's base (smallest) monetary unit
				CurrencyID     string  `json:"currencyId"`     // the currency the monetary value is represented in
				ValueAsDecimal float64 `json:"valueAsDecimal"` // the monetary value represented as decimal for precise calculations and storage
			} `json:"nativeAmount"` // Monetary value (an amount with a currency).
		} `json:"destinationAmount"` // Payout destination amount details
		ReceivedBlockchainTxID string `json:"receivedBlockchainTxId"` // The blockchain transaction ID associated with the received payment for the invoice payout.
		Items                  []struct {
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
			MerchantFees *struct {
				TransactionFee *struct {
					CurrencyID     string  `json:"currencyId"`     // the currency the monetary value is represented in
					DisplayValue   string  `json:"displayValue"`   // the value formatted for display
					Value          string  `json:"value"`          // the amount of money in the currency's base (smallest) monetary unit
					ValueAsDecimal float64 `json:"valueAsDecimal"` // the monetary value represented as decimal for precise calculations and storage
				} `json:"transactionFee"` // No further information available.
				NetworkFee *struct {
					CurrencyID     string  `json:"currencyId"`     // the currency the monetary value is represented in
					DisplayValue   string  `json:"displayValue"`   // the value formatted for display
					Value          string  `json:"value"`          // the amount of money in the currency's base (smallest) monetary unit
					ValueAsDecimal float64 `json:"valueAsDecimal"` // the monetary value represented as decimal for precise calculations and storage
				} `json:"networkFee"` // No further information available.
				ConversionFee *struct {
					CurrencyID     string  `json:"currencyId"`     // the currency the monetary value is represented in
					DisplayValue   string  `json:"displayValue"`   // the value formatted for display
					Value          string  `json:"value"`          // the amount of money in the currency's base (smallest) monetary unit
					ValueAsDecimal float64 `json:"valueAsDecimal"` // the monetary value represented as decimal for precise calculations and storage
				} `json:"conversionFee"` // No further information available.
			} `json:"merchantFees"` // The amount for service fees in the merchant's accepted currency
			PayoutAmount *struct {
				CurrencyID     string  `json:"currencyId"`     // the currency the monetary value is represented in
				DisplayValue   string  `json:"displayValue"`   // the value formatted for display
				Value          string  `json:"value"`          // the amount of money in the currency's base (smallest) monetary unit
				ValueAsDecimal float64 `json:"valueAsDecimal"` // the monetary value represented as decimal for precise calculations and storage
			} `json:"payoutAmount"` // No further information available.
			PayoutAmountInInvoiceCurrency *struct {
				CurrencyID     string  `json:"currencyId"`     // the currency the monetary value is represented in
				DisplayValue   string  `json:"displayValue"`   // the value formatted for display
				Value          string  `json:"value"`          // the amount of money in the currency's base (smallest) monetary unit
				ValueAsDecimal float64 `json:"valueAsDecimal"` // the monetary value represented as decimal for precise calculations and storage
			} `json:"payoutAmountInInvoiceCurrency"` // No further information available.
			MerchantPayoutAddress string `json:"merchantPayoutAddress"` // the merchants payment output address at the time the hot wallet was created
			Created               string `json:"created"`               // the timestamp of when this payout was created (or scheduled), Format: date-time
			Sent                  string `json:"sent"`                  // the timestamp of when this payout was sent (e.g. broadcast on the blockchain), Format: date-time
			Expected              string `json:"expected"`              // the approximate timestamp of when this payout is expected to be sent (e.g. broadcast on the blockchain), Format: date-time
			Confirmed             string `json:"confirmed"`             // the timestamp of when this payout was confirmed (e.g. on the blockchain), Format: date-time
			State                 string `json:"state"`                 // payout state. Enum: scheduled, sending, sent, confirmed, waitingC
		} `json:"items"` // No further information available.
		Paging *struct {
			Cursors *struct {
				Before string `json:"before"` // No further information available.
				After  string `json:"after"`  // No further information available.
			} `json:"cursors"` // No further information available.
			Limit    int32  `json:"limit"`    // No further information available. Format: int32
			First    string `json:"first"`    // No further information available. Format: uri
			Next     string `json:"next"`     // No further information available. Format: uri
			Previous string `json:"previous"` // No further information available. Format: uri
			Last     string `json:"last"`     // No further information available. Format: uri
		} `json:"paging"` // No further information available.
	} `json:"payoutDetails"` // payout information for the merchant
	Payments []struct {
		PaymentCurrencyID     string `json:"paymentCurrencyId"`     // The currencyId selected by a buyer to pay the invoice
		PaymentCurrencySymbol string `json:"paymentCurrencySymbol"` // The currencySymbol selected by a buyer to pay the invoice (e.g., ETH, SOL, etc.)
		NativeCurrencyID      string `json:"nativeCurrencyId"`      // The currencyId selected by the merchant for native display values
		NativeCurrencySymbol  string `json:"nativeCurrencySymbol"`  // The currencySymbol selected by the merchant for native display values (e.g., USD, EUR, etc.)
		ExpectedAmount        string `json:"expectedAmount"`        // Expected amount to be paid by a buyer
		NativeExpectedAmount  string `json:"nativeExpectedAmount"`  // Expected amount in merchant-preferred native currency
		ActualAmount          string `json:"actualAmount"`          // Actual paid amount in the payment currency selected by a buyer
		NativeActualAmount    string `json:"nativeActualAmount"`    // Actual paid amount in merchant-preferred native currency
		PaymentAddress        string `json:"paymentAddress"`        // Blockchain payment address
		ErrorCode             string `json:"errorCode"`             // Settlement error code for payment processing
		Fees                  *struct {
			PaymentSubTotal          string `json:"paymentSubTotal"`          // The subtotal amount of the payment transaction before fees or adjustments.
			MerchantMarkupOrDiscount string `json:"merchantMarkupOrDiscount"` // The amount representing either a markup or discount applied by the merchant.
			BuyerFee                 *struct {
				CoinPaymentsFee string `json:"coinPaymentsFee"` // CoinPayments fee details.
				NetworkFee      string `json:"networkFee"`      // Network fee details.
				ConversionFee   string `json:"conversionFee"`   // Conversion fee details.
				Total           string `json:"total"`           // Total fee amount for the buyer.
			} `json:"buyerFee"` // Summary of fees applicable to the buyer.
			MerchantFee *struct {
				CoinPaymentsFee string `json:"coinPaymentsFee"` // CoinPayments fee details.
				NetworkFee      string `json:"networkFee"`      // Network fee details.
				ConversionFee   string `json:"conversionFee"`   // Conversion fee details.
				Total           string `json:"total"`           // Total fee amount for the buyer.
			} `json:"merchantFee"` // Summary of fees applicable to the merchant.
			Gross string `json:"gross"` // Gross total = payment subtotal + buyer fees + merchant fees.
		} `json:"fees"` // Payment fees summary in payment currency
		NativeFees *struct {
			PaymentSubTotal          string `json:"paymentSubTotal"`          // The subtotal amount of the payment transaction before fees or adjustments.
			MerchantMarkupOrDiscount string `json:"merchantMarkupOrDiscount"` // The amount representing either a markup or discount applied by the merchant.
			BuyerFee                 *struct {
				CoinPaymentsFee string `json:"coinPaymentsFee"` // CoinPayments fee details.
				NetworkFee      string `json:"networkFee"`      // Network fee details.
				ConversionFee   string `json:"conversionFee"`   // Conversion fee details.
				Total           string `json:"total"`           // Total fee amount for the buyer.
			} `json:"buyerFee"` // Summary of fees applicable to the buyer.
			MerchantFee *struct {
				CoinPaymentsFee string `json:"coinPaymentsFee"` // CoinPayments fee details.
				NetworkFee      string `json:"networkFee"`      // Network fee details.
				ConversionFee   string `json:"conversionFee"`   // Conversion fee details.
				Total           string `json:"total"`           // Total fee amount for the buyer.
			} `json:"merchantFee"` // Summary of fees applicable to the merchant.
			Gross string `json:"gross"` // Gross total = payment subtotal + buyer fees + merchant fees.
		} `json:"nativeFees"` // Payment fees summary in native currency
		Payout *struct {
			ScheduledAt           string  `json:"scheduledAt"`           // The date and time when the payout is scheduled to be processed.
			CompletedAt           string  `json:"completedAt"`           // The date and time when the payout has been completed.
			BlockchainTx          string  `json:"blockchainTx"`          // The TxHash on the blockchain associated with the payout.
			SpendRequestId        string  `json:"spendRequestId"`        // The identifier for the spend request linked to this payout.
			Address               string  `json:"address"`               // The blockchain address to which the payout is sent.
			WalletId              string  `json:"walletId"`              // The unique identifier for the wallet associated with the payout.
			SentAt                string  `json:"sentAt"`                // The date and time when the payout was sent.
			ExpectedExecutionDate string  `json:"expectedExecutionDate"` // The date and time when the payout execution is expected to be completed.
			ReceivedBlockchainTx  string  `json:"receivedBlockchainTx"`  // The identifier of the blockchain transaction confirming receipt, if available.
			IsBatched             bool    `json:"isBatched"`             // Indicates whether the payout is part of a batch transaction.
			CurrencySymbol        string  `json:"currencySymbol"`        // The symbol representing the currency (e.g., "BTC", "ETH", "TRX").
			CurrencyId            string  `json:"currencyId"`            // The currency the monetary value is represented in.
			DisplayValue          string  `json:"displayValue"`          // The value formatted for display.
			Value                 string  `json:"value"`                 // The amount in the currency's base (smallest) monetary unit.
			ValueAsDecimal        float64 `json:"valueAsDecimal"`        // Represents the monetary value converted to a decimal type.
		} `json:"payout"` // Payout summary for the merchant
		NativePayout          string `json:"nativePayout"`          // Final payout info in merchant-preferred currency
		RefundEmail           string `json:"refundEmail"`           // Refund notification email
		State                 string `json:"state"`                 // Payment state
		IsActive              bool   `json:"isActive"`              // Indicates whether the invoice payment is active
		PendingAt             string `json:"pendingAt"`             // Timestamp when the payment entered pending state, Format: date-time
		ConfirmedAt           string `json:"confirmedAt"`           // Timestamp when required blockchain confirmations were reached, Format: date-time
		CompletedAt           string `json:"completedAt"`           // Timestamp when the payment was completed, Format: date-time
		Confirmations         int32  `json:"confirmations"`         // Current number of blockchain confirmations, Format: int32
		RequiredConfirmations int32  `json:"requiredConfirmations"` // Required number of blockchain confirmations, Format: int32
	} `json:"payments"` // payment summaries related to the invoice
	IsLifeTimeFinished bool   `json:"isLifeTimeFinished"` // invoice in finished state more than 90 days
	HideShoppingCart   bool   `json:"hideShoppingCart"`   // flag for hiding icon on the checkout app
	SuccessURL         string `json:"successUrl"`         // redirect URL after successful payment
	CancelURL          string `json:"cancelUrl"`          // redirect URL after canceling payment
	PayoutOverrides    []struct {
		FromCurrency string `json:"fromCurrency"` // The buyer selected currency.
		ToCurrency   string `json:"toCurrency"`   // Currency of the payout wallet or address.
		Address      string `json:"address"`      // External address to pay out to.
		Frequency    string `json:"frequency"`    // Payout frequency: normal, asSoonAsPossible, hourly, nightly, weekly.
	} `json:"payoutOverrides"` // configuration for overriding default payout settings
}
