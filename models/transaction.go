package models

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"`
	Type      string    `json:"trx_type"`
	Category  string    `json:"trx_category"`
	CreatedAt time.Time `json:"created_at"`
}
