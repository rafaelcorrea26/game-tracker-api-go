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
// @Success 201 {object} models.Game
// @Router /games [post]
func CreateGame(c *gin.Context) {
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

	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.QueryRow(`
		INSERT INTO games (name, status, year_completed, rating, notes, platform, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, status, year_completed, rating, notes, platform
	`,
		game.Name,
		game.Status,
		game.YearCompleted,
		game.Rating,
		game.Notes,
		game.Platform,
		userID,
	).Scan(
		&game.ID,
		&game.Name,
		&game.Status,
		&game.YearCompleted,
		&game.Rating,
		&game.Notes,
		&game.Platform,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, game)
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
		SELECT id, name, status, year_completed, rating, notes, COALESCE(platform, '') as platform
		FROM games
		WHERE user_id=$1
		ORDER BY id DESC
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
			&game.Platform,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		games = append(games, game)
	}

	if games == nil {
		games = []models.Game{}
	}

	c.JSON(http.StatusOK, games)
}

// =========================
// UPDATE GAME
// =========================

// UpdateGame godoc
// @Summary Atualiza um jogo
// @Description Atualiza nome, status, ano, nota, notas e plataforma de um jogo do usuário
// @Tags games
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID do jogo"
// @Param input body models.Game true "Jogo"
// @Success 200 {object} models.Game
// @Router /games/{id} [put]
func UpdateGame(c *gin.Context) {
	id := c.Param("id")
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

	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = database.DB.QueryRow(`
		UPDATE games
		SET name = $1,
			status = $2,
			year_completed = $3,
			rating = $4,
			notes = $5,
			platform = $6
		WHERE id = $7 AND user_id = $8
		RETURNING id, name, status, year_completed, rating, notes, platform
	`,
		game.Name,
		game.Status,
		game.YearCompleted,
		game.Rating,
		game.Notes,
		game.Platform,
		id,
		userID,
	).Scan(
		&game.ID,
		&game.Name,
		&game.Status,
		&game.YearCompleted,
		&game.Rating,
		&game.Notes,
		&game.Platform,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "jogo não encontrado"})
		return
	}

	c.JSON(http.StatusOK, game)
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
