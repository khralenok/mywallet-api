package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/models"
)

func GetTransactionByUserID(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	query := "SELECT * FROM transactions WHERE user_id=$1"
	rows, err := database.DB.Query(query, userID)
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

func AddIncome(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var newTRX models.Transaction

	if err := context.BindJSON(&newTRX); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	newTRX.UserID = userID
	newTRX.Type = "income"

	query := "INSERT INTO transactions (user_id, amount, trx_type, trx_category) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	err := database.DB.QueryRow(query, newTRX.UserID, newTRX.Amount, newTRX.Type, newTRX.Category).Scan(&newTRX.ID, &newTRX.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
		return
	}

	context.JSON(http.StatusCreated, newTRX)
}

func AddExpense(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var newTRX models.Transaction

	if err := context.BindJSON(&newTRX); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	newTRX.UserID = userID
	newTRX.Type = "expense"

	query := "INSERT INTO transactions (user_id, amount, trx_type, trx_category) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	err := database.DB.QueryRow(query, newTRX.UserID, newTRX.Amount, newTRX.Type, newTRX.Category).Scan(&newTRX.ID, &newTRX.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
		return
	}

	context.JSON(http.StatusCreated, newTRX)
}
