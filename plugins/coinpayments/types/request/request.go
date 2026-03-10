package request

// CreateInvoice
type CreateInvoice struct {
	CreateInvoicePayload
	Currency    string                 `json:"currency"`    // Currency of the invoice items and amount
	Amount      *CreateInvoiceAmount   `json:"amount"`      // Total invoice amount with a breakdown that provides details such as total item amount, total tax amount, shipping, handling, insurance, and discounts, if any.
	Description string                 `json:"description"` // the purchase description, can be provided instead of a list of `items`
	Webhooks    []CreateInvoiceWebhook `json:"webhooks"`    // Contains an array of
	Payment     *CreateInvoicePayment  `json:"payment"`     // Request to create an invoice payment with predefined currency
	SuccessURL  string                 `json:"successUrl"`  // the url to redirect to once an invoice is successfully paid
	CancelURL   string                 `json:"cancelUrl"`   // the url to redirect to if payment of an invoice fails (e.g. expired) or is cancelled by the user
}

type CreateInvoicePayload struct {
	Items                      []CreateInvoiceItem           `json:"items,omitempty"`                      // Array of items that a buyer intends to purchase from the merchant
	IsEmailDelivery            *bool                         `json:"isEmailDelivery,omitempty"`            // indicates if invoice will be email delivered
	EmailDelivery              *CreateInvoiceEmailDelivery   `json:"emailDelivery,omitempty"`              // Email delivery options of a merchant invoice
	DueDate                    string                        `json:"dueDate,omitempty"`                    // optional due date to be shown on the invoice, Format: date-time
	InvoiceDate                string                        `json:"invoiceDate,omitempty"`                // optional custom invoice date if not the created date of the invoice, invoices with a future date will be scheduled, Format: date-time
	Draft                      *bool                         `json:"draft,omitempty"`                      // flag indicating whether this is a draft invoice
	ClientID                   string                        `json:"clientId,omitempty"`                   // the id of the client creating this invoice (optional)
	InvoiceID                  string                        `json:"invoiceId,omitempty"`                  // the optional API caller provided external invoice number. Appears in screens shown to the Buyer and emails sent.
	Buyer                      *CreateInvoiceBuyer           `json:"buyer,omitempty"`                      // the buyer information
	Shipping                   *CreateInvoiceShipping        `json:"shipping,omitempty"`                   // Invoice shipping address and method
	RequireBuyerNameAndEmail   *bool                         `json:"requireBuyerNameAndEmail,omitempty"`   // flag indicating whether a buyer name and email are required, they will be requested at checkout if not provider by the caller. The CoinPayments.Api.Models.Merchant.CreateMerchantInvoiceRequestV2Dto.BuyerDataCollectionMessage will be displayed to the buyer when prompted.
	BuyerDataCollectionMessage string                        `json:"buyerDataCollectionMessage,omitempty"` // the message to display when collecting buyer user data
	Notes                      string                        `json:"notes,omitempty"`                      // notes for the merchant only, these are not visible to the buyers
	NotesToRecipient           string                        `json:"notesToRecipient,omitempty"`           // any additional information to share with the buyer about the transaction
	TermsAndConditions         string                        `json:"termsAndConditions,omitempty"`         // any terms and conditions, e.g. a cancellation policy
	MerchantOptions            *CreateInvoiceMerchantOptions `json:"merchantOptions,omitempty"`            // Options to show/hide merchant information on an invoice, or include additional merchant information specific to an invoice
	CustomData                 any                           `json:"customData,omitempty"`                 // any custom data the caller wishes to attach to the invoice which will be sent back in notifications
	Metadata                   *CreateInvoiceMetadata        `json:"metadata,omitempty"`                   // Represents metadata information related to an invoice, including integration and hostname details.
	PoNumber                   string                        `json:"poNumber,omitempty"`                   // Merchant's invoice number. Сan store any string. If the field is not filled in, a sequence number will be generated. Must be unique per merchant
	PayoutOverrides            []CreateInvoicePayoutOverride `json:"payoutOverrides,omitempty"`            // Optionally specifies payout configs for this invoice
	UseCoinReservation         *bool                         `json:"useCoinReservation,omitempty"`         // Indicates whether the invoice will use coin reservation for FIAT currency during transaction processing.
	HideShoppingCart           *bool                         `json:"hideShoppingCart,omitempty"`           // Flag for hiding icon on the checkout app
	AffiliateID                string                        `json:"affiliateId,omitempty"`                // Identifier of the affiliate associated with the invoice, if applicable, Format: uuid
	IsSimpleQR                 *bool                         `json:"isSimpleQR,omitempty"`                 // If IsSimpleQR is true, the checkout app generates a QR code with only the address. Otherwise, the QR code includes the address, currency, and amount.
}

type CreateInvoiceItem struct {
	CustomID       string                 `json:"customId"`       // the API caller provided external ID for the item. Appears on the Merchant dashboard and reports only.
	SKU            string                 `json:"sku"`            // the stock keeping unit (SKU) of the item
	Name           string                 `json:"name"`           // (* required) the name or title of the item
	Description    string                 `json:"description"`    // the detailed description of the item
	Quantity       *CreateInvoiceQuantity `json:"quantity"`       // Represents the quantity details for a line item in the CoinPayments merchant system. This class defines both the value of the quantity and its associated type.
	OriginalAmount string                 `json:"originalAmount"` // the original total price of the item if CoinPayments.Services.Merchants.Dtos.LineItemV2Dto.Amount represents a discounted price
	Amount         string                 `json:"amount"`         // the subtotal price of the item (note: this is not a per unit price but the total price for the total quantity)
	Tax            string                 `json:"tax"`            // the total taxes charged on this item
}

type CreateInvoiceQuantity struct {
	Value int64  `json:"value"` // the quantity of the item. Must be greater than 0 and less than 999,999,999. defaults to 1 if not provided.
	Type  string `json:"type"`  // Defines the types of quantities applicable to a line item in the CoinPayments merchant system. It categorizes the quantity as either a measure of time or a specific count. Enum: hours, quantity
}

type CreateInvoiceAmount struct {
	Breakdown *CreateInvoiceBreakdown `json:"breakdown"` // No further information available.
	Total     string                  `json:"total"`     // The total amount associated with the invoice, encompassing all charges, fees, and adjustments.
}

type CreateInvoiceBreakdown struct {
	Subtotal string `json:"subtotal"` // the subtotal for all items, required if the request includes `invoice.items[]` and must be equal to the sum of all `(invoice.items[].amount)` for all items.
	Shipping string `json:"shipping"` // the optional shipping fee for all items within the invoice
	Handling string `json:"handling"` // the optional handling fee for all items within the invoice
	TaxTotal string `json:"taxTotal"` // the total taxes charged on the invoice, required if the request includes `invoice.items[]` with taxes and must be equal to the sum of all `(invoice.items[].tax)` for all items.
	Discount string `json:"discount"` // the optional discount for all items within the invoice
}

type CreateInvoiceEmailDelivery struct {
	To  string `json:"to"`  // the email `to` field, multiple addresses separated by semicolons
	Cc  string `json:"cc"`  // the email `cc` field, multiple addresses separated by semicolons
	Bcc string `json:"bcc"` // the email `bcc` field, multiple addresses separated by semicolons
}

type CreateInvoiceBuyer struct {
	CompanyName  string                `json:"companyName"`  // the name of the buyer's company
	Name         *CreateInvoiceName    `json:"name"`         // No further information available.
	EmailAddress string                `json:"emailAddress"` // the email address of the buyer
	PhoneNumber  string                `json:"phoneNumber"`  // the phone number
	Address      *CreateInvoiceAddress `json:"address"`      // No further information available.
}

type CreateInvoiceName struct {
	FirstName string `json:"firstName"` // the given, or first, name
	LastName  string `json:"lastName"`  // the surname or family name. Required when the party is a person.
}

type CreateInvoiceAddress struct {
	Address1         string `json:"address1"`         // the first line of the address. For example, number or street.
	Address2         string `json:"address2"`         // the second line of the address. For example, suite or apartment number.
	Address3         string `json:"address3"`         // the third line of the address, if needed.
	ProvinceOrState  string `json:"provinceOrState"`  // the highest level sub-division in a country, usually a province or state.
	City             string `json:"city"`             // the city, town or village.
	SuburbOrDistrict string `json:"suburbOrDistrict"` // the neighborhood, suburb or district.
	CountryCode      string `json:"countryCode"`      // the two-character ISO-3166-1 country code.
	PostalCode       string `json:"postalCode"`       // the postal code, zip code, or equivalent.
}

type CreateInvoiceShipping struct {
	Method       string                `json:"method"`       // the shipping method
	CompanyName  string                `json:"companyName"`  // the company name of the party to ship the items to
	Name         *CreateInvoiceName    `json:"name"`         // shipping recipient name
	EmailAddress string                `json:"emailAddress"` // the email address of the party to ship the items to
	PhoneNumber  string                `json:"phoneNumber"`  // the phone number
	Address      *CreateInvoiceAddress `json:"address"`      // shipping address
	HasData      *bool                 `json:"hasData"`      // indicates shipping data is present
}

type CreateInvoiceMerchantOptions struct {
	ShowAddress            *bool  `json:"showAddress"`            // Indicates whether the address should be shown on the invoice.
	ShowEmail              *bool  `json:"showEmail"`              // Indicates whether the email should be shown on the invoice.
	ShowPhone              *bool  `json:"showPhone"`              // Indicates whether the phone should be shown on the invoice.
	ShowRegistrationNumber *bool  `json:"showRegistrationNumber"` // Indicates whether the business registration number should be shown on the invoice.
	ShowTaxID              *bool  `json:"showTaxId"`              // Specifies whether the tax ID should be displayed on the invoice.
	AdditionalInfo         string `json:"additionalInfo"`         // Miscellaneous invoice-specific merchant information.
}

type CreateInvoiceMetadata struct {
	Integration string `json:"integration"` // the integration from which the invoice was created
	Hostname    string `json:"hostname"`    // the hostname on which the invoice was created
}

type CreateInvoiceWebhook struct {
	NotificationsURL string   `json:"notificationsUrl"` // the url to which to POST webhook notifications to
	Notifications    []string `json:"notifications"`    // the types of notifications to send to this endpoint
}
type CreateInvoicePayoutOverride struct {
	FromCurrency string `json:"fromCurrency"` // The buyer selected currency
	ToCurrency   string `json:"toCurrency"`   // Currency of the payout wallet or address
	Address      string `json:"address"`      // External address to pay out to
	Frequency    string `json:"frequency"`    // Enum: normal, asSoonAsPossible, hourly, nightly, weekly
}

type CreateInvoicePayment struct {
	PaymentCurrency string `json:"paymentCurrency"` // Create payment address for currency.
	RefundEmail     string `json:"refundEmail"`     // Email for refund instructions if there is a payment issue.
}
