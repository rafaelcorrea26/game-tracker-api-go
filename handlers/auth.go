package handlers

import (
	"net/http"

	"game-tracker-api-go/database"
	"game-tracker-api-go/models"
	"game-tracker-api-go/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// =========================
// INPUTS
// =========================

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// =========================
// REGISTER
// =========================

// Register godoc
// @Summary Criar usuário
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Dados do usuário"
// @Success 201 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {

	var input RegisterInput

	// 📥 validar JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}

	// 🔐 Hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro ao gerar hash da senha",
		})
		return
	}

	// 💾 Inserir no banco
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

// =========================
// LOGIN
// =========================

// Login godoc
// @Summary Login do usuário
// @Description Retorna um token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Credenciais"
// @Success 200 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {

	var input LoginInput

	// 📥 validar JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "dados inválidos",
		})
		return
	}

	var user models.User

	// 🔎 Buscar usuário no banco
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

	// 🔐 Comparar senha com hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "credenciais inválidas",
		})
		return
	}

	// 🎟️ Gerar token JWT
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
