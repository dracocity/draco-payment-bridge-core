package models

type Item struct {
	ID       string `json:"itemId"`
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
	Amount   string `json:"amount"`
}
