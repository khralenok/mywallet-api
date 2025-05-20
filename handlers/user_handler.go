package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khralenok/mywallet-api/database"
	"github.com/khralenok/mywallet-api/models"
	"github.com/khralenok/mywallet-api/utilities"
)

func GetProfile(context *gin.Context) {
	userID := context.MustGet("userID").(int)

	var user models.User

	query := "SELECT * FROM users WHERE id=$1"
	err := database.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": user})
}

func CreateUser(context *gin.Context) {
	var newUser models.User

	if err := context.BindJSON(&newUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	var passwordHash string
	var err error

	if passwordHash, err = utilities.HashPassword(newUser.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
		return
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, created_at"

	err = database.DB.QueryRow(query, newUser.Username, passwordHash).Scan(&newUser.ID, &newUser.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user"})
		return
	}

	context.JSON(http.StatusCreated, newUser)
}

func LoginUser(context *gin.Context) {
	var input models.LoginInputs

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	var user models.User
	var err error

	if user, err = getUserByUsername(input.Username, context); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if !utilities.CheckPasswordHash(input.Password, user.Password) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
	}

	token, err := utilities.GenerateJWT(user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Success", "token": token})

}

func getUserByUsername(username string, context *gin.Context) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE username=$1"
	err := database.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return models.User{}, err
	}

	return user, nil
}
