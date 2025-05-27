package models

import (
	"database/sql"
	"time"

	"github.com/khralenok/mywallet-api/database"
)

type Balance struct {
	UserID       int       `json:"user_id"`
	Balance      int       `json:"balance"`
	SnapshotDate time.Time `json:"snapshot_date"`
}

func CreateNewBalance(userID int) error {
	newBalance := Balance{
		UserID:       userID,
		Balance:      0,
		SnapshotDate: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	query := "INSERT INTO balances (user_id, balance, snapshot_date) VALUES ($1, $2, $3)"
	_, err := database.DB.Exec(query, newBalance.UserID, newBalance.Balance, newBalance.SnapshotDate)

	if err != nil {
		return err
	}

	return nil
}

func CalcBalance(userID int) error {
	var newBalance Balance

	query := "SELECT * FROM balances WHERE user_id=$1"

	err := database.DB.QueryRow(query, userID).Scan(&newBalance.UserID, &newBalance.Balance, &newBalance.SnapshotDate)

	if err != nil {
		return err
	}

	query = "SELECT * FROM transactions WHERE user_id=$1 AND created_at>$2"

	rows, err := database.DB.Query(query, newBalance.UserID, newBalance.SnapshotDate)

	if err != nil {
		return err
	}

	defer rows.Close()

	if err := UpdateBalance(rows, userID); err != nil {
		return err
	}

	return nil
}

func RecalcBalance(userID int) error {
	query := "SELECT * FROM transactions WHERE user_id=$1"

	rows, err := database.DB.Query(query, userID)

	if err != nil {
		return err
	}

	defer rows.Close()

	if err := UpdateBalance(rows, userID); err != nil {
		return err
	}

	return nil
}

func UpdateBalance(rawTrxs *sql.Rows, userID int) error {
	var newBalance Balance
	currentTime := time.Now()

	sumOfTransactions, err := SumOfTransactions(rawTrxs)

	if err != nil {
		return err
	}

	newBalance.UserID = userID
	newBalance.Balance += sumOfTransactions
	newBalance.SnapshotDate = currentTime

	query := "UPDATE balances SET balance=$1,snapshot_date=$2 WHERE user_id=$3"

	_, err = database.DB.Exec(query, newBalance.Balance, newBalance.SnapshotDate, newBalance.UserID)

	if err != nil {
		return err
	}

	return nil
}

func GetBalance(userID int) (int, error) {
	var curBalance int
	var lastSnapshot Balance

	query := "SELECT * FROM balances WHERE user_id=$1"

	err := database.DB.QueryRow(query, userID).Scan(&lastSnapshot.UserID, &lastSnapshot.Balance, &lastSnapshot.SnapshotDate)

	if err != nil {
		return curBalance, err
	}

	return lastSnapshot.Balance, nil
}
