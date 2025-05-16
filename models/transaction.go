package models

import "time"

// id | user_id | amount | trx_type | trx_category |         created_at

type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"`
	Type      string    `json:"trx_type"`
	Category  string    `json:"trx_category"`
	CreatedAt time.Time `json:"created_at"`
}
