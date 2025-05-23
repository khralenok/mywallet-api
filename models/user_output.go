package models

type UserOutput struct {
	ID         int     `json:"id"`
	Username   string  `json:"username"`
	BalanceUSD float64 `json:"balance_usd"`
}
