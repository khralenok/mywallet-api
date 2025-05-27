package handlers

import (
	"net/http"
	"time"

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
		if err := rows.Scan(&rawTRX.ID, &rawTRX.UserID, &rawTRX.Amount, &rawTRX.Type, &rawTRX.Category, &rawTRX.CreatedAt); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		transactions = append(transactions, prepareTrxToOutput(rawTRX))
	}

	context.JSON(http.StatusOK, transactions)
}

func GetTransactionById(context *gin.Context) {
	_ = context.MustGet("userID").(int)
	trxID := context.Param("id")

	var rawTRX models.Transaction

	query := "SELECT * FROM transactions WHERE id=$1"

	err := database.DB.QueryRow(query, trxID).Scan(&rawTRX.ID, &rawTRX.UserID, &rawTRX.Amount, &rawTRX.Type, &rawTRX.Category, &rawTRX.CreatedAt)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	context.JSON(http.StatusOK, prepareTrxToOutput(rawTRX))
}

func GetTransactionByDate(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	fromDate := context.Query("from_date")
	toDate := context.Query("to_date")

	from, err := time.Parse(time.DateOnly, fromDate)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "from_date in 2006-01-02 format needed"})
		return
	}
	to, err := time.Parse(time.DateOnly, toDate)

	if err != nil {
		to = time.Now()
	}

	var transactions []models.TransactionOutput

	query := "SELECT * FROM transactions WHERE user_id=$1 and created_at BETWEEN $2 AND $3"

	rows, err := database.DB.Query(query, userID, from, to)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var rawTRX models.Transaction
		if err := rows.Scan(&rawTRX.ID, &rawTRX.UserID, &rawTRX.Amount, &rawTRX.Type, &rawTRX.Category, &rawTRX.CreatedAt); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		transactions = append(transactions, prepareTrxToOutput(rawTRX))
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

	if newTransactionRequest.AmountUSD < 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
	}

	newTRX.UserID = userID
	newTRX.Amount = utilities.ConvertToCents(newTransactionRequest.AmountUSD)

	curBalance, err := models.GetBalance(newTRX.UserID)

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

func DeleteTransaction(context *gin.Context) {
	userID := context.MustGet("userID").(int)
	trxID := context.Param("id")

	query := "DELETE FROM transactions WHERE id=$1"

	_, err := database.DB.Exec(query, trxID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	models.RecalcBalance(userID)

	context.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

func UpdateTransaction(context *gin.Context) {
	userID := context.MustGet("userID").(int)
	trxID := context.Param("id")

	var updateTransactionRequest models.TransactionRequest

	if err := context.BindJSON(&updateTransactionRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	query := "UPDATE transactions SET amount=$1, trx_category=$2 WHERE id=$3"

	_, err := database.DB.Exec(query, utilities.ConvertToCents(updateTransactionRequest.AmountUSD), updateTransactionRequest.Category, trxID)

	if err != nil {
		context.JSON(http.StatusNotModified, gin.H{"error": err.Error()})
	}

	models.RecalcBalance(userID)

	context.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}

func prepareTrxToOutput(rawTRX models.Transaction) models.TransactionOutput {
	var trx models.TransactionOutput

	trx.ID = rawTRX.ID
	trx.AmountUSD = utilities.ConvertToUSD(rawTRX.Amount)
	if rawTRX.Type == "expense" {
		trx.AmountUSD *= -1
	}

	trx.Category = rawTRX.Category
	trx.CreatedAt = rawTRX.CreatedAt

	return trx
}
