package response

// CreateOrder
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
