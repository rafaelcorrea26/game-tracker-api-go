package handlers

import (
	"encoding/json"
	"net/http"

	"game-tracker-api-go/database"
	"game-tracker-api-go/models"
)

func GetGames(w http.ResponseWriter, r *http.Request) {

	rows, err := database.DB.Query(`
		SELECT id, name, status, year_completed, rating, notes
		FROM games
	`)

	if err != nil {
		http.Error(w, err.Error(), 500)
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

	json.NewEncoder(w).Encode(games)
}