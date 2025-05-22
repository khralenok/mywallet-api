package models

type TransactionRequest struct {
	AmountUSD float64 `json:"amount_usd"`
	Category  string  `json:"trx_category"`
}
