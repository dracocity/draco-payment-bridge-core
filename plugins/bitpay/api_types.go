package main

import "github.com/shopspring/decimal"

// CreateInvoice
type createInvoiceRequest struct {
	Token                                  string                   `json:"token"`
	Price                                  float64                  `json:"price"`
	Currency                               string                   `json:"currency"`
	BitpayIdRequired                       *bool                    `json:"bitpayIdRequired"`
	MerchantName                           string                   `json:"merchantName,omitempty"`
	ForcedBuyerSelectedTransactionCurrency string                   `json:"forcedBuyerSelectedTransactionCurrency,omitempty"`
	ForcedBuyerSelectedWallet              string                   `json:"forcedBuyerSelectedWallet,omitempty"`
	OrderID                                string                   `json:"orderId,omitempty"`
	ItemDesc                               string                   `json:"itemDesc,omitempty"`
	ItemCode                               string                   `json:"itemCode,omitempty"`
	ItemizedDetails                        []itemizedDetailsRequest `json:"itemizedDetails,omitempty"`
	NotificationEmail                      string                   `json:"notificationEmail,omitempty"`
	NotificationURL                        string                   `json:"notificationURL,omitempty"`
	RedirectURL                            string                   `json:"redirectURL,omitempty"`
	CloseURL                               string                   `json:"closeURL,omitempty"`
	AutoRedirect                           *bool                    `json:"autoRedirect,omitempty"`
	PosData                                string                   `json:"posData,omitempty"`
	GUID                                   string                   `json:"guid,omitempty"`
	TransactionSpeed                       string                   `json:"transactionSpeed,omitempty"`
	FullNotifications                      *bool                    `json:"fullNotifications,omitempty"`
	ExtendedNotifications                  *bool                    `json:"extendedNotifications,omitempty"`
	Physical                               *bool                    `json:"physical,omitempty"`
	BuyerSMS                               string                   `json:"buyerSms,omitempty"`
	Buyer                                  *buyerRequest            `json:"buyer,omitempty"`
	JsonPayProRequired                     string                   `json:"jsonPayProRequired,omitempty"`
	AcceptanceWindow                       *int32                   `json:"acceptanceWindow,omitempty"`
}

type itemizedDetailsRequest struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
	IsFee       bool   `json:"isFee"`
}

type buyerRequest struct {
	Name       string `json:"name"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	Locality   string `json:"locality"`
	Region     string `json:"region"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Notify     bool   `json:"notify"`
}

type createInvoicePayload struct {
	BitpayIdRequired                       *bool                    `json:"bitpayIdRequired"`
	MerchantName                           string                   `json:"merchantName"`
	ForcedBuyerSelectedTransactionCurrency string                   `json:"forcedBuyerSelectedTransactionCurrency"`
	ForcedBuyerSelectedWallet              string                   `json:"forcedBuyerSelectedWallet"`
	ItemCode                               string                   `json:"itemCode"`
	ItemizedDetails                        []itemizedDetailsRequest `json:"itemizedDetails"`
	NotificationEmail                      string                   `json:"notificationEmail"`
	AutoRedirect                           *bool                    `json:"autoRedirect"`
	PosData                                string                   `json:"posData"`
	GUID                                   string                   `json:"guid"`
	TransactionSpeed                       string                   `json:"transactionSpeed"`
	FullNotifications                      *bool                    `json:"fullNotifications"`
	ExtendedNotifications                  *bool                    `json:"extendedNotifications"`
	Physical                               *bool                    `json:"physical"`
	BuyerSMS                               string                   `json:"buyerSms"`
	Buyer                                  *buyerRequest            `json:"buyer"`
	JsonPayProRequired                     string                   `json:"jsonPayProRequired"`
	AcceptanceWindow                       *int32                   `json:"acceptanceWindow"`
}

type createInvoiceResponse struct {
	Facade string       `json:"facade"`
	Data   dataResponse `json:"data"`
}

type dataResponse struct {
	URL                            string                                    `json:"url"`
	PosData                        string                                    `json:"posData,omitempty"`
	Status                         string                                    `json:"status"`
	Price                          decimal.Decimal                           `json:"price"`
	Currency                       string                                    `json:"currency"`
	ItemDesc                       string                                    `json:"itemDesc,omitempty"`
	OrderID                        string                                    `json:"orderId,omitempty"`
	InvoiceTime                    int64                                     `json:"invoiceTime"`
	ExpirationTime                 int64                                     `json:"expirationTime"`
	CurrentTime                    int64                                     `json:"currentTime"`
	GUID                           string                                    `json:"guid,omitempty"`
	ID                             string                                    `json:"id"`
	LowFeeDetected                 bool                                      `json:"lowFeeDetected"`
	AmountPaid                     decimal.Decimal                           `json:"amountPaid"`
	DisplayAmountPaid              string                                    `json:"displayAmountPaid,omitempty"`
	ExceptionStatus                bool                                      `json:"exceptionStatus"`
	TargetConfirmations            *int64                                    `json:"targetConfirmations,omitempty"`
	Transactions                   []map[string]any                          `json:"transactions,omitempty"`
	TransactionSpeed               string                                    `json:"transactionSpeed,omitempty"`
	Buyer                          *buyerResponse                            `json:"buyer,omitempty"`
	RedirectURL                    string                                    `json:"redirectURL,omitempty"`
	AutoRedirect                   *bool                                     `json:"autoRedirect,omitempty"`
	CloseURL                       string                                    `json:"closeURL,omitempty"`
	RefundAddresses                []map[string]refundAddressInfoResponse    `json:"refundAddresses,omitempty"`
	RefundAddressRequestPending    bool                                      `json:"refundAddressRequestPending"`
	BuyerProvidedEmail             string                                    `json:"buyerProvidedEmail,omitempty"`
	BuyerProvidedInfo              *buyerProvidedInfoResponse                `json:"buyerProvidedInfo,omitempty"`
	PaymentSubtotals               map[string]int64                          `json:"paymentSubtotals,omitempty"`
	PaymentTotals                  map[string]int64                          `json:"paymentTotals,omitempty"`
	PaymentDisplayTotals           map[string]string                         `json:"paymentDisplayTotals,omitempty"`
	PaymentDisplaySubTotals        map[string]string                         `json:"paymentDisplaySubTotals,omitempty"`
	ExchangeRates                  map[string]map[string]decimal.Decimal     `json:"exchangeRates,omitempty"`
	MinerFees                      map[string]minerFeeResponse               `json:"minerFees,omitempty"`
	Shopper                        any                                       `json:"shopper,omitempty"`
	JsonPayProRequired             bool                                      `json:"jsonPayProRequired"`
	MerchantName                   string                                    `json:"merchantName,omitempty"`
	BitpayIdRequired               *bool                                     `json:"bitpayIdRequired,omitempty"`
	ItemizedDetails                []itemizedDetailsResponse                 `json:"itemizedDetails,omitempty"`
	SupportedTransactionCurrencies map[string]supportedCurrencyStateResponse `json:"supportedTransactionCurrencies,omitempty"`
	PaymentCodes                   map[string]map[string]string              `json:"paymentCodes,omitempty"`
	UniversalCodes                 *universalCodesResponse                   `json:"universalCodes,omitempty"`
	Token                          string                                    `json:"token,omitempty"`
}

type buyerResponse struct {
	Email string `json:"email"`
}

type refundAddressInfoResponse struct {
	Type  string `json:"type"`
	Date  string `json:"date"`
	Email string `json:"email"`
}

type buyerProvidedInfoResponse struct {
	EmailAddress                string `json:"emailAddress"`
	SelectedWallet              string `json:"selectedWallet"`
	SelectedTransactionCurrency string `json:"selectedTransactionCurrency"`
}

type minerFeeResponse struct {
	SatoshisPerByte decimal.Decimal `json:"satoshisPerByte"`
	TotalFee        decimal.Decimal `json:"totalFee"`
	FiatAmount      decimal.Decimal `json:"fiatAmount"`
}

type itemizedDetailsResponse struct {
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	IsFee       bool   `json:"isFee"`
}

type supportedCurrencyStateResponse struct {
	Enabled bool   `json:"enabled"`
	Reason  string `json:"reason,omitempty"`
}

type universalCodesResponse struct {
	PaymentString string `json:"paymentString"`
}
