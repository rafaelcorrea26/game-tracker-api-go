// @title Game Tracker API
// @version 1.0
// @description API para rastrear progresso de jogos
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"os"
	"strings"

	"game-tracker-api-go/database"
	"game-tracker-api-go/handlers"
	"game-tracker-api-go/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "game-tracker-api-go/docs"
)

func main() {
	database.Connect()
	database.RunMigrations()

	r := gin.Default()

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if origin == "http://localhost:5173" {
				return true
			}

			if origin == frontendURL {
				return true
			}

			if strings.HasPrefix(origin, "https://game-tracker-frontend-react-") &&
				strings.HasSuffix(origin, ".vercel.app") {
				return true
			}

			return false
		},
	}))

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())

	auth.GET("/games", handlers.GetGames)
	auth.POST("/games", handlers.CreateGame)
	auth.PUT("/games/:id", handlers.UpdateGame)
	auth.DELETE("/games/:id", handlers.DeleteGame)
	auth.PUT("/change-password", handlers.ChangePassword)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
