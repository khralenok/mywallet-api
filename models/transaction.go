package models

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"`
	Type      string    `json:"trx_type"`
	Category  string    `json:"trx_category"`
	CreatedAt time.Time `json:"created_at"`
}

func SumOfTransactions(rawTrxs *sql.Rows) (int, error) {
	var sumOfTrxs int

	for rawTrxs.Next() {
		var newTransaction Transaction
		if err := rawTrxs.Scan(&newTransaction.ID, &newTransaction.UserID, &newTransaction.Amount, &newTransaction.Type, &newTransaction.Category, &newTransaction.CreatedAt); err != nil {
			return 0, err
		}

		if newTransaction.Type == "expense" {
			newTransaction.Amount *= -1
		}

		sumOfTrxs += newTransaction.Amount
	}

	return sumOfTrxs, nil
}
