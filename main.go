package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/handlers"
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

	router.GET("/users", handlers.GetUsers)
	router.GET("/transactions", handlers.GetTransactions)
	router.GET("/balances", handlers.GetBalances)

	router.POST("/signin", handlers.CreateUser)
	router.POST("/login", handlers.LoginUser)

	router.Run(":8080")
}
