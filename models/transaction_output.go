package models

import "time"

type TransactionOutput struct {
	ID        int       `json:"id"`
	AmountUSD float64   `json:"amount_usd"`
	Category  string    `json:"trx_category"`
	CreatedAt time.Time `json:"created_at"`
}
