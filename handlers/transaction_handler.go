package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/models"
)

func GetTransactions(context *gin.Context) {
	rows, err := database.DB.Query("SELECT id, user_id, amount, trx_type, trx_category, created_at FROM transactions")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Category, &transaction.Type, &transaction.CreatedAt); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transactions = append(transactions, transaction)
	}

	context.JSON(http.StatusOK, transactions)
}
