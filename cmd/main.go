package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vimalrajliya/backend-assignment-app/cmd/api"
	"github.com/vimalrajliya/backend-assignment-app/cmd/auth"
	"github.com/vimalrajliya/backend-assignment-app/database"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

func main() {
	database.ConnectDB()
	database.ConnectRedis()
	fmt.Print(" Program started")
	router := gin.Default()
	router.POST("/auth", api.PostUser)
	router.POST("/auth/sign-in", api.SignInUser)
	protected := router.Group("/api")
	protected.Use(auth.AuthenticateToken())
	{
		protected.GET("/user", api.GetUserDetails)
		protected.GET("/token/refresh", api.RefreshToken)
		protected.POST("/user/log-out", api.LogOutUser)
	}
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
