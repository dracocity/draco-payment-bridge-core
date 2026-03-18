package response

// Create Order
type CreateOrder struct {
	ID                     int64  `json:"id"`
	Status                 string `json:"status"`
	Title                  string `json:"title"`
	DoNotConvert           bool   `json:"do_not_convert"`
	OrderableType          string `json:"orderable_type"`
	OrderableID            int64  `json:"orderable_id"`
	PriceCurrency          string `json:"price_currency"`
	PriceAmount            string `json:"price_amount"`
	LightningNetwork       bool   `json:"lightning_network"`
	ReceiveCurrency        string `json:"receive_currency"`
	ReceiveAmount          string `json:"receive_amount"`
	CreatedAt              string `json:"created_at"`
	OrderID                string `json:"order_id"`
	PaymentURL             string `json:"payment_url"`
	UnderpaidAmount        string `json:"underpaid_amount"`
	OverpaidAmount         string `json:"overpaid_amount"`
	IsRefundable           bool   `json:"is_refundable"`
	Refunds                []any  `json:"refunds"`
	Voids                  []any  `json:"voids"`
	Fees                   []any  `json:"fees"`
	Token                  string `json:"token"`
	BlockchainTransactions []any  `json:"blockchain_transactions"`
}

// Get Order
type GetOrder struct {
	ID                     int64                           `json:"id"`
	Status                 string                          `json:"status"`
	Title                  string                          `json:"title"`
	DoNotConvert           bool                            `json:"do_not_convert"`
	OrderableType          string                          `json:"orderable_type"`
	OrderableID            int64                           `json:"orderable_id"`
	PriceCurrency          string                          `json:"price_currency"`
	PriceAmount            string                          `json:"price_amount"`
	PayCurrency            string                          `json:"pay_currency"`
	PayAmount              string                          `json:"pay_amount"`
	LightningNetwork       bool                            `json:"lightning_network"`
	ReceiveCurrency        string                          `json:"receive_currency"`
	ReceiveAmount          string                          `json:"receive_amount"`
	CreatedAt              string                          `json:"created_at"`
	ExpireAt               string                          `json:"expire_at"`
	PaidAt                 string                          `json:"paid_at"`
	PaymentAddress         string                          `json:"payment_address"`
	PaymentGateway         string                          `json:"payment_gateway"`
	OrderID                string                          `json:"order_id"`
	PaymentURL             string                          `json:"payment_url"`
	PaymentRequestURI      string                          `json:"payment_request_uri"`
	UnderpaidAmount        string                          `json:"underpaid_amount"`
	OverpaidAmount         string                          `json:"overpaid_amount"`
	IsRefundable           bool                            `json:"is_refundable"`
	ConversionRate         string                          `json:"conversion_rate"`
	Refunds                []GetOrderRefund                `json:"refunds"`
	Voids                  []any                           `json:"voids"`
	Fees                   []GetOrderFee                   `json:"fees"`
	BlockchainTransactions []GetOrderBlockchainTransaction `json:"blockchain_transactions"`
}

type GetOrderRefund struct {
	ID             int64                 `json:"id"`
	RequestAmount  string                `json:"request_amount"`
	RefundAmount   string                `json:"refund_amount"`
	Address        string                `json:"address"`
	Status         string                `json:"status"`
	Memo           any                   `json:"memo"`
	CreatedAt      string                `json:"created_at"`
	Order          GetOrderRefundOrder   `json:"order"`
	RefundCurrency GetOrderCurrencyInfo  `json:"refund_currency"`
	Transactions   []any                 `json:"transactions"`
	LedgerAccount  GetOrderLedgerAccount `json:"ledger_account"`
}

type GetOrderRefundOrder struct {
	ID int64 `json:"id"`
}

type GetOrderCurrencyInfo struct {
	ID       int64            `json:"id"`
	Title    string           `json:"title"`
	Symbol   string           `json:"symbol"`
	Platform GetOrderPlatform `json:"platform"`
}

type GetOrderPlatform struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type GetOrderLedgerAccount struct {
	ID       string                 `json:"id"`
	Currency GetOrderLedgerCurrency `json:"currency"`
}

type GetOrderLedgerCurrency struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Symbol string `json:"symbol"`
}

type GetOrderFee struct {
	Type     string              `json:"type"`
	Amount   string              `json:"amount"`
	Currency GetOrderFeeCurrency `json:"currency"`
}

type GetOrderFeeCurrency struct {
	ID     int64  `json:"id"`
	Symbol string `json:"symbol"`
}

type GetOrderBlockchainTransaction struct {
	ID                   int64                `json:"id"`
	Txid                 string               `json:"txid"`
	Amount               string               `json:"amount"`
	Status               string               `json:"status"`
	NetworkConfirmations int64                `json:"network_confirmations"`
	Currency             GetOrderCurrencyInfo `json:"currency"`
}
