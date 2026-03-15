package main

import (
	"fmt"
	"net/http"

	"game-tracker-api-go/database"
	"game-tracker-api-go/handlers"
)

func main() {

	database.Connect()

	http.HandleFunc("/games", handlers.GetGames)

	fmt.Println("Servidor rodando na porta 8080")

	http.ListenAndServe(":8080", nil)
}