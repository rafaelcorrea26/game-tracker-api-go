package models

type Game struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	YearCompleted int    `json:"year_completed"`
	Rating        int    `json:"rating"`
	Notes         string `json:"notes"`
	Platform      string `json:"platform"`
}
