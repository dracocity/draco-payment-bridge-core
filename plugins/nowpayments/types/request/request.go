package request

// CreateInvoice
type CreateInvoice struct {
	CreateInvoicePayload
	PriceAmount      float64 `json:"price_amount"`                // the amount that users have to pay for the order stated in fiat currency. In case you do not indicate the price in crypto, our system will automatically convert this fiat amount into its crypto equivalent. NOTE: Some of the assets (KISHU, NWC, FTT, CHR, XYM, SRK, KLV, SUPER, OM, XCUR, NOW, SHIB, SAND, MATIC, CTSI, MANA, FRONT, FTM, DAO, LGCY), have a maximum price limit of ~$2000;
	PriceCurrency    string  `json:"price_currency"`              // the fiat currency in which the price_amount is specified (usd, eur, etc);
	PayCurrency      string  `json:"pay_currency,omitempty"`      // the specified crypto currency (btc, eth, etc), or one of available fiat currencies if it's enabled for your account (USD, EUR, ILS, GBP, AUD, RON); If not specified, can be chosen on the invoice_url
	IPNCallbackURL   string  `json:"ipn_callback_url,omitempty"`  // url to receive callbacks, should contain "http" or "https", eg. "https://nowpayments.io";
	OrderID          string  `json:"order_id,omitempty"`          // internal store order ID, e.g. "RGDBP-21314"
	OrderDescription string  `json:"order_description,omitempty"` // internal store order description, e.g. "Apple Macbook Pro 2019 x 1"
	SuccessURL       string  `json:"success_url,omitempty"`       // url where the customer will be redirected after successful payment;
	CancelURL        string  `json:"cancel_url,omitempty"`        // url where the customer will be redirected after failed payment;
}

type CreateInvoicePayload struct {
	IsFixedRate     *bool `json:"is_fixed_rate,omitempty"`       // boolean, can be true or false. Required for fixed-rate exchanges
	IsFeePaidByUser *bool `json:"is_fee_paid_by_user,omitempty"` // boolean, can be true or false. Required for fixed-rate exchanges with all fees paid by users;
}

// CreateInvoicePayment
type CreateInvoicePayment struct {
	CreateInvoicePaymentPayload
	InvoiceID        string `json:"iid"`                         // invoice id. You can get invoice ID in response of POST Create_invoice method
	PayCurrency      string `json:"pay_currency"`                // the crypto currency in which the pay_amount is specified (btc, eth, etc). NOTE: some of the currencies require a Memo, Destination Tag, etc., to complete a payment (AVA, EOS, BNBMAINNET, XLM, XRP). This is unique for each payment. This ID is received in “payin_extra_id” parameter of the response. Payments made without "payin_extra_id" cannot be detected automatically;=
	OrderDescription string `json:"order_description,omitempty"` // inner store order description
	CustomerEmail    string `json:"customer_email,omitempty"`    // user email to which a notification about the successful completion of the payment will be sent
}

type CreateInvoicePaymentPayload struct {
	PayoutAddress string `json:"payout_address,omitempty"`  // usually the funds will go to the address you specify in your Personal account. In case you want to receive funds on another address, you can specify it in this parameter
	PayoutExtraID string `json:"payout_extra_id,omitempty"` // extra id or memo or tag for external payout_address
}
