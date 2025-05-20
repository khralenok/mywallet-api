package main

import (
	"log"
	"net/http"

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

	router.GET("/balances", handlers.GetBalances) // For testing. Should be deleted on release

	router.GET("/profile", utilities.AuthMiddleware(), handlers.GetProfile)
	router.GET("/my_trxs", utilities.AuthMiddleware(), handlers.GetTransactionByUserID)

	router.GET("/protected", utilities.AuthMiddleware(), func(context *gin.Context) {
		userID := context.MustGet("userID").(int)
		context.JSON(http.StatusOK, gin.H{"message": "Welcome!", "user_id": userID})
	}) // For testing JWT. Should be deleted on release

	router.POST("/signin", handlers.CreateUser)
	router.POST("/login", handlers.LoginUser)
	router.POST("/add_trx", utilities.AuthMiddleware(), handlers.CreateTransaction)

	router.Run(":8080")
}
