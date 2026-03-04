package request

// CreateInvoice
type CreateInvoice struct {
	CreateInvoicePayload
	Token           string  `json:"token"`                     // The API token can be retrieved from the dashboard (limited to pos facade) or using the Tokens resource to get access to the merchant facade. This is described in the section Request an API token).
	Price           float64 `json:"price"`                     // Fixed price amount for the checkout, in the "currency" of the invoice object.
	Currency        string  `json:"currency"`                  // ISO 4217 3-character currency code. This is the currency associated with the price field, supported currencies are available via the Currencies resource.
	OrderID         string  `json:"orderId,omitempty"`         // Can be used by the merchant to assign their own internal Id to an invoice. If used, there should be a direct match between an orderId and an invoice id.
	ItemDesc        string  `json:"itemDesc,omitempty"`        // Invoice description - will be added as a line item on the BitPay checkout page, under the merchant name.
	NotificationURL string  `json:"notificationURL,omitempty"` // URL to which BitPay sends webhook notifications. HTTPS is mandatory.
	RedirectURL     string  `json:"redirectURL,omitempty"`     // The shopper will be redirected to this URL when clicking on the Return button after a successful payment or when clicking on the Close button if a separate closeURL is not specified. Be sure to include "http://" or "https://" in the url.
	CloseURL        string  `json:"closeURL,omitempty"`        // URL to redirect if the shopper does not pay the invoice and click on the Close button instead. Be sure to include "http://" or "https://" in the url.
}

type CreateInvoicePayload struct {
	BitpayIdRequired                       *bool                         `json:"bitpayIdRequired"`                                 // Forces the invoice to require BitPay ID to be completed, regardless of price.
	MerchantName                           string                        `json:"merchantName,omitempty"`                           // Display string for merchant identification (ex. Wal-Mart Store #1452, Bowling Green, KY).
	ForcedBuyerSelectedTransactionCurrency string                        `json:"forcedBuyerSelectedTransactionCurrency,omitempty"` // Merchant pre-selects transaction currency on behalf of buyer.
	ForcedBuyerSelectedWallet              string                        `json:"forcedBuyerSelectedWallet,omitempty"`              // Merchant pre-selects wallet on behalf of buyer.
	ItemCode                               string                        `json:"itemCode,omitempty"`                               // "bitcoindonation" for donations, otherwise do not include the field in the request.
	ItemizedDetails                        []CreateInvoiceItemizedDetail `json:"itemizedDetails,omitempty"`                        // Object containing line item details for display.
	NotificationEmail                      string                        `json:"notificationEmail,omitempty"`                      // Merchant email address for notification of invoice status change. It is also possible to configure this email via the account setting on the BitPay dashboard or disable the email notification
	AutoRedirect                           *bool                         `json:"autoRedirect,omitempty"`                           // Set to `false` by default, merchant can setup automatic redirect to their website by setting this parameter to `true`. (paid => RedirectURL, expires => CloseURL)
	PosData                                string                        `json:"posData,omitempty"`                                // A passthru variable provided by the merchant during invoice creation and designed to be used by the merchant to correlate the invoice with an order or other object in their system. This passthru variable can be a serialized object, e.g.: "posData": "\"{ \"ref\" : 711454, \"item\" : \"test_item\" }\"".
	GUID                                   string                        `json:"guid,omitempty"`                                   // A passthru variable provided by the merchant and designed to be used by the merchant to correlate the invoice with an order ID in their system, which can be used as a lookup variable in Retrieve Invoice by GUID.
	TransactionSpeed                       string                        `json:"transactionSpeed,omitempty"`                       // This is a risk mitigation parameter for the merchant to configure how they want to fulfill orders depending on the number of block confirmations for the transaction made by the consumer on the selected cryptocurrency. (high, medium, low)
	FullNotifications                      *bool                         `json:"fullNotifications,omitempty"`                      // This parameter is set to true by default, meaning all standard notifications are being sent for a payment made to an invoice. If you decide to set it to false instead, only 1 webhook will be sent for each invoice paid by the consumer. This webhook will be for the "confirmed" or "complete" invoice status, depending on the transactionSpeed selected.
	ExtendedNotifications                  *bool                         `json:"extendedNotifications,omitempty"`                  // Allows merchants to get access to additional webhooks. For instance when an invoice expires without receiving a payment or when it is refunded. If set to true, then fullNotifications is automatically set to true. When using the extendedNotifications parameter, the webhook also have a payload slightly different from the standard webhooks.
	Physical                               *bool                         `json:"physical,omitempty"`                               // Indicates whether items are physical goods. Alternatives include digital goods and services.
	BuyerSMS                               string                        `json:"buyerSms,omitempty"`                               // SMS phone number including country code i.e. +12223334444.
	Buyer                                  *CreateInvoiceBuyer           `json:"buyer,omitempty"`                                  // Allows merchant to pass buyer related information in the invoice object
	JsonPayProRequired                     string                        `json:"jsonPayProRequired,omitempty"`                     // If set to true, this means that the invoice will only accept payments from wallets which have implemented the BitPay JSON Payment Protocol
	AcceptanceWindow                       *int32                        `json:"acceptanceWindow,omitempty"`                       // Number of milliseconds that a user has to pay an invoice before it expires (0-900000). If not set, invoice will default to the account acceptanceWindow. If account acceptanceWindow is not set, invoice will default to 15 minutes (900,000 milliseconds).
}

type CreateInvoiceItemizedDetail struct {
	Amount      string `json:"amount"`      // Amount of currency for item
	Description string `json:"description"` // Display string for the item
	IsFee       bool   `json:"isFee"`       // Indicates whether or not the item is considered a fee/tax or part of the main purchase
}

type CreateInvoiceBuyer struct {
	Name       string `json:"name"`       // Buyer's name
	Address1   string `json:"address1"`   // Buyer's address
	Address2   string `json:"address2"`   // Buyer's appartment or suite number
	Locality   string `json:"locality"`   // Buyer's city or locality
	Region     string `json:"region"`     // Buyer's state or province
	PostalCode string `json:"postalCode"` // Buyer's Zip or Postal Code
	Country    string `json:"country"`    // Buyer's Country code. Format ISO 3166-1 alpha-2
	Email      string `json:"email"`      // Buyer's email address. If provided during invoice creation, this will bypass the email prompt for the consumer when opening the invoice.
	Phone      string `json:"phone"`      // Buyer's phone number
	Notify     bool   `json:"notify"`     // Indicates whether a BitPay email confirmation should be sent to the buyer once he has paid the invoice
}
