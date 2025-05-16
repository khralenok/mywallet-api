package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/models"
)

func GetUsers(context *gin.Context) {
	rows, err := database.DB.Query("SELECT id, username, password, created_at FROM users")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	context.JSON(http.StatusOK, users)
}

func CreateUser(context *gin.Context) {
	var newUser models.User

	if err := context.BindJSON(&newUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, created_at"

	err := database.DB.QueryRow(query, newUser.Username, newUser.Password).Scan(&newUser.ID, &newUser.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	context.JSON(http.StatusCreated, newUser)
}
