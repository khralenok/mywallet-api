package handlers

import (
	"net/http"

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
