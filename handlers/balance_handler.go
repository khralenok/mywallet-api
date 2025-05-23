package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/models"
	"github.com/khralenok/mywallet-api/utilities"
)

func GetBalance(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var balance models.Balance

	query := "SELECT * FROM balances WHERE user_id=$1"
	err := database.DB.QueryRow(query, userID).Scan(&balance.UserID, &balance.Balance, &balance.SnapshotDate)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"balance_usd": utilities.ConvertToUSD(balance.Balance), "snapshot_date": balance.SnapshotDate})
}

func TakeBalanceSnapshot(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var newBalance models.Balance
	currentTime := time.Now()

	query := "SELECT * FROM balances WHERE user_id=$1"

	err := database.DB.QueryRow(query, userID).Scan(&newBalance.UserID, &newBalance.Balance, &newBalance.SnapshotDate)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	query = "SELECT * FROM transactions WHERE user_id=$1 AND created_at>$2"

	rows, err := database.DB.Query(query, newBalance.UserID, newBalance.SnapshotDate)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var sumOfTransactions int

	for rows.Next() {
		var newTransaction models.Transaction
		if err := rows.Scan(&newTransaction.ID, &newTransaction.UserID, &newTransaction.Amount, &newTransaction.Type, &newTransaction.Category, &newTransaction.CreatedAt); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if newTransaction.Type == "expense" {
			newTransaction.Amount *= -1
		}

		sumOfTransactions += newTransaction.Amount
	}

	newBalance.Balance += sumOfTransactions
	newBalance.SnapshotDate = currentTime

	query = "UPDATE balances SET balance=$1,snapshot_date=$2 WHERE user_id=$3"

	_, err = database.DB.Exec(query, newBalance.Balance, newBalance.SnapshotDate, newBalance.UserID)

	if err != nil {
		context.JSON(http.StatusNotModified, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Snapshot taken successfuly"})
}

func getCurBalance(userID int) (int, error) {
	var curBalance int
	var lastSnapshot models.Balance

	query := "SELECT * FROM balances WHERE user_id=$1"

	err := database.DB.QueryRow(query, userID).Scan(&lastSnapshot.UserID, &lastSnapshot.Balance, &lastSnapshot.SnapshotDate)

	if err != nil {
		return curBalance, err
	}

	query = "SELECT * FROM transactions WHERE user_id=$1 AND created_at>$2"

	rows, err := database.DB.Query(query, lastSnapshot.UserID, lastSnapshot.SnapshotDate)

	if err != nil {
		return curBalance, err
	}

	defer rows.Close()

	for rows.Next() {
		var newTRX models.Transaction
		if err := rows.Scan(&newTRX.ID, &newTRX.UserID, &newTRX.Amount, &newTRX.Type, &newTRX.Category, &newTRX.CreatedAt); err != nil {
			return curBalance, err
		}

		if newTRX.Type == "expense" {
			newTRX.Amount *= -1
		}

		lastSnapshot.Balance += newTRX.Amount
	}

	curBalance = lastSnapshot.Balance

	return curBalance, nil
}

func createNewBalance(userID int) error {
	newBalance := models.Balance{
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
