package handlers

import (
	"net/http"

	"game-tracker-api-go/database"
	"game-tracker-api-go/models"

	"github.com/gin-gonic/gin"
)

// =========================
// CREATE GAME
// =========================

// CreateGame godoc
// @Summary Cria um jogo
// @Tags games
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.Game true "Jogo"
// @Success 200 {object} map[string]string
// @Router /games [post]
func CreateGame(c *gin.Context) {

	// 🔐 pegar usuário do JWT
	username := c.GetString("username")

	var userID int

	err := database.DB.QueryRow(
		"SELECT id FROM users WHERE username=$1",
		username,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não encontrado"})
		return
	}

	var game models.Game

	// 📥 receber JSON
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 💾 inserir no banco
	_, err = database.DB.Exec(`
		INSERT INTO games (name, status, year_completed, rating, notes, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		game.Name,
		game.Status,
		game.YearCompleted,
		game.Rating,
		game.Notes,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "jogo criado com sucesso",
	})
}

// =========================
// GET GAMES
// =========================

// GetGames godoc
// @Summary Lista jogos do usuário
// @Tags games
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Game
// @Router /games [get]
func GetGames(c *gin.Context) {

	// 🔐 pegar usuário do JWT
	username := c.GetString("username")

	var userID int

	err := database.DB.QueryRow(
		"SELECT id FROM users WHERE username=$1",
		username,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não encontrado"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, name, status, year_completed, rating, notes
		FROM games
		WHERE user_id=$1
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var games []models.Game

	for rows.Next() {
		var game models.Game

		err := rows.Scan(
			&game.ID,
			&game.Name,
			&game.Status,
			&game.YearCompleted,
			&game.Rating,
			&game.Notes,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		games = append(games, game)
	}

	c.JSON(http.StatusOK, games)
}

// =========================
// DELETE GAME
// =========================
// DeleteGame godoc
// @Summary Deleta um jogo
// @Description Remove um jogo pelo ID
// @Tags games
// @Security BearerAuth
// @Param id path int true "ID do jogo"
// @Success 200 {object} map[string]string
// @Router /games/{id} [delete]
func DeleteGame(c *gin.Context) {

	id := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM games WHERE id=$1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "jogo deletado com sucesso",
	})
}
