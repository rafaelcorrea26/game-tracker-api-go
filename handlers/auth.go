package handlers

import (
	"net/http"

	"game-tracker-api-go/database"
	"game-tracker-api-go/models"
	"game-tracker-api-go/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro ao gerar hash da senha",
		})
		return
	}

	_, err = database.DB.Exec(
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		input.Username,
		string(hash),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "usuário já existe",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "usuário criado com sucesso",
	})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}

	var user models.User

	err := database.DB.QueryRow(
		"SELECT id, username, password FROM users WHERE username=$1",
		input.Username,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "credenciais inválidas",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "credenciais inválidas",
		})
		return
	}

	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro ao gerar token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func ChangePassword(c *gin.Context) {
	username := c.GetString("username")

	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}

	if len(input.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "a nova senha deve ter pelo menos 6 caracteres",
		})
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, password FROM users WHERE username=$1",
		username,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "usuário não encontrado",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "senha atual incorreta",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro ao gerar hash da nova senha",
		})
		return
	}

	_, err = database.DB.Exec(
		"UPDATE users SET password=$1 WHERE id=$2",
		string(hash),
		user.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "não foi possível atualizar a senha",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "senha alterada com sucesso",
	})
}