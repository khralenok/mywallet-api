package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/models"
	"github.com/khralenok/mywallet-api/utilities"
)

func GetMonthSummary(context *gin.Context) {
	userID := context.MustGet("userID").(int)
	monthInput := context.Query("month")

	var income int
	var expense int

	monthStart, err := time.Parse("2006-01", monthInput)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "month should be in 2006-01 format"})
		return
	}

	monthEnd := monthStart.AddDate(0, 1, 0)

	query := "SELECT * FROM transactions WHERE user_id=$1 and created_at BETWEEN $2 AND $3"

	rows, err := database.DB.Query(query, userID, monthStart, monthEnd)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Type, &transaction.Category, &transaction.CreatedAt); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if transaction.Type == "expense" {
			expense += transaction.Amount
			continue
		}

		income += transaction.Amount
	}

	balance, err := models.GetBalance(userID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	context.JSON(http.StatusOK, gin.H{"income_usd": utilities.ConvertToUSD(income), "expenses_usd": utilities.ConvertToUSD(expense), "balance_usd": utilities.ConvertToUSD(balance)})
}
