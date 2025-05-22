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

	router.GET("/profile", utilities.AuthMiddleware(), handlers.GetProfile)
	router.GET("/balance", utilities.AuthMiddleware(), handlers.GetMyBalance)
	router.GET("/my_trxs", utilities.AuthMiddleware(), handlers.GetTransactionByUserID)
	router.GET("/snapshot", utilities.AuthMiddleware(), handlers.TakeBalanceSnapshot)

	router.POST("/signin", handlers.CreateUser)
	router.POST("/login", handlers.LoginUser)
	router.POST("/add_income", utilities.AuthMiddleware(), handlers.AddIncome)
	router.POST("/add_expense", utilities.AuthMiddleware(), handlers.AddExpense)

	router.Run(":8080")
}
