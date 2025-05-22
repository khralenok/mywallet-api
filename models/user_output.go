package models

import "time"

type UserOutput struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	BalanceUSD   float64   `json:"balance_usd"`
	SnapshotDate time.Time `json:"snapshot_date"`
}
