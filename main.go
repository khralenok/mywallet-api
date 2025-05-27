package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/handlers"
	"github.com/khralenok/mywallet-api/utilities"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer database.DB.Close()

	router := gin.Default()

	//User Management
	router.POST("/signin", handlers.CreateUser)
	router.POST("/login", handlers.LoginUser)
	router.GET("/profile", utilities.AuthMiddleware(), handlers.GetProfile)

	//Transaction management
	router.POST("/add_income", utilities.AuthMiddleware(), handlers.AddIncome)
	router.POST("/add_expense", utilities.AuthMiddleware(), handlers.AddExpense)
	router.GET("/transactions/:id", utilities.AuthMiddleware(), handlers.GetTransactionById)
	router.GET("/transactions/date", utilities.AuthMiddleware(), handlers.GetTransactionByDate)
	router.PUT("/update_transaction/:id", utilities.AuthMiddleware(), handlers.UpdateTransaction)
	router.DELETE("/delete_transaction/:id", utilities.AuthMiddleware(), handlers.DeleteTransaction)

	//Balance management
	router.GET("/balance", utilities.AuthMiddleware(), handlers.GetBalance)
	router.GET("/transactions", utilities.AuthMiddleware(), handlers.GetTransactions)
	router.GET("/snapshot", utilities.AuthMiddleware(), handlers.TakeBalanceSnapshot)

	//Reports
	router.GET("/month_summary", utilities.AuthMiddleware(), handlers.GetMonthSummary)

	router.Run(":8080")
}
