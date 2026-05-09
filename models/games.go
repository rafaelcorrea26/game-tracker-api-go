package models

type Game struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Rating    int    `json:"rating"`
	Notes     string `json:"notes"`
	Platform  string `json:"platform"`
	ImageURL  string `json:"image_url"`
}
