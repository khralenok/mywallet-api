package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/models"
)

// user_id | balance | snapshot_date

func GetBalances(context *gin.Context) {
	rows, err := database.DB.Query("SELECT user_id, balance, snapshot_date FROM balances")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var balances []models.Balance

	for rows.Next() {
		var balance models.Balance
		if err := rows.Scan(&balance.UserID, &balance.Balance, &balance.SnapshotDate); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		balances = append(balances, balance)
	}

	context.JSON(http.StatusOK, balances)
}

func TakeBalanceSnapshot(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var newBalance models.Balance
	currentTime := time.Now()

	query := "SELECT * FROM balances WHERE user_id=$1"

	row := database.DB.QueryRow(query, userID)

	err := row.Scan(&newBalance.UserID, &newBalance.Balance, &newBalance.SnapshotDate)

	if err != nil {
		newBalance = models.Balance{
			UserID:       userID,
			Balance:      0,
			SnapshotDate: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		}
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
		if err := rows.Scan(&newTransaction.ID, &newTransaction.UserID, &newTransaction.Amount, &newTransaction.Category, &newTransaction.Type, &newTransaction.CreatedAt); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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
