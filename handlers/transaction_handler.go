package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/models"
	"github.com/khralenok/mywallet-api/utilities"
)

func GetTransactions(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	query := "SELECT * FROM transactions WHERE user_id=$1"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var transactions []models.TransactionOutput

	for rows.Next() {
		var rawTRX models.Transaction
		var transaction models.TransactionOutput
		if err := rows.Scan(&rawTRX.ID, &rawTRX.UserID, &rawTRX.Amount, &rawTRX.Type, &rawTRX.Category, &rawTRX.CreatedAt); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		transaction.ID = rawTRX.ID
		transaction.AmountUSD = utilities.ConvertToUSD(rawTRX.Amount)
		if rawTRX.Type == "expense" {
			transaction.AmountUSD *= -1
		}

		transaction.Category = rawTRX.Category
		transaction.CreatedAt = rawTRX.CreatedAt

		transactions = append(transactions, transaction)
	}

	context.JSON(http.StatusOK, transactions)
}

func AddIncome(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var newTransactionRequest models.TransactionRequest
	var newTRX models.Transaction

	if err := context.BindJSON(&newTransactionRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	newTRX.UserID = userID
	newTRX.Amount = utilities.ConvertToCents(newTransactionRequest.AmountUSD)
	newTRX.Type = "income"
	newTRX.Category = newTransactionRequest.Category

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

	var newTransactionRequest models.TransactionRequest
	var newTRX models.Transaction

	if err := context.BindJSON(&newTransactionRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	newTRX.UserID = userID
	newTRX.Amount = utilities.ConvertToCents(newTransactionRequest.AmountUSD)

	curBalance, err := getCurBalance(newTRX.UserID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if newTRX.Amount > curBalance {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Not enough funds to make transaction"})
		return
	}

	newTRX.Type = "expense"
	newTRX.Category = newTransactionRequest.Category

	query := "INSERT INTO transactions (user_id, amount, trx_type, trx_category) VALUES ($1, $2, $3, $4) RETURNING id, created_at"
	err = database.DB.QueryRow(query, newTRX.UserID, newTRX.Amount, newTRX.Type, newTRX.Category).Scan(&newTRX.ID, &newTRX.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert transaction"})
		return
	}

	context.JSON(http.StatusCreated, newTRX)
}
