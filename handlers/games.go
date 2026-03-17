package handlers

import (
	"fmt"
	"game-tracker-api-go/database"
	"game-tracker-api-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGames(c *gin.Context) {

	status := c.Query("status")
	year := c.Query("year")
	name := c.Query("name")

	query := `
	SELECT id, name, status, year_completed, rating, notes
	FROM games
	WHERE 1=1
	`

	args := []interface{}{}
	index := 1

	if status != "" {
		query += fmt.Sprintf(" AND status=$%d", index)
		args = append(args, status)
		index++
	}

	if year != "" {
		query += fmt.Sprintf(" AND year_completed=$%d", index)
		args = append(args, year)
		index++
	}

	if name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", index)
		args = append(args, "%"+name+"%")
		index++
	}

	rows, err := database.DB.Query(query, args...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()

	var games []models.Game

	for rows.Next() {

		var g models.Game

		rows.Scan(
			&g.ID,
			&g.Name,
			&g.Status,
			&g.YearCompleted,
			&g.Rating,
			&g.Notes,
		)

		games = append(games, g)
	}

	c.JSON(http.StatusOK, games)
}

func CreateGame(c *gin.Context) {

	var game models.Game

	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := `
	INSERT INTO games (name, status, year_completed, rating, notes)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id
	`

	err := database.DB.QueryRow(
		query,
		game.Name,
		game.Status,
		game.YearCompleted,
		game.Rating,
		game.Notes,
	).Scan(&game.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, game)
}
