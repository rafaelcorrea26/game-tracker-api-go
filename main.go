// @title Game Tracker API
// @version 1.0
// @description API para rastrear progresso de jogos
// @host localhost:8080
// @BasePath /

// 🔐 JWT config pro Swagger
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"game-tracker-api-go/database"
	"game-tracker-api-go/handlers"
	"game-tracker-api-go/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "game-tracker-api-go/docs"
)

func main() {

	// 🔌 conectar banco
	database.Connect()

	// 🧱 rodar migrations
	database.RunMigrations()

	// 🚀 criar servidor
	r := gin.Default()

	// =========================
	// 🔓 ROTAS PÚBLICAS
	// =========================

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Swagger (público)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// =========================
	// 🔒 ROTAS PROTEGIDAS
	// =========================

	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())

	auth.GET("/games", handlers.GetGames)
	auth.POST("/games", handlers.CreateGame)
	auth.DELETE("/games/:id", handlers.DeleteGame)

	// =========================
	// ▶️ START SERVER
	// =========================

	r.Run(":8080")
}
