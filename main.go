package main

import (
	"game-tracker-api-go/database"
	"game-tracker-api-go/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()

	database.RunMigrations()

	r := gin.Default()

	r.GET("/games", handlers.GetGames)
	r.POST("/games", handlers.CreateGame)

	r.Run(":8080")
}
