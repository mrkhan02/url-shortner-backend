package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mrkhan02/url-shortner-api/routes"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = ":8000"
	}
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Adjust this based on your needs
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.New(config))
	routes.ResRoutes(router)
	routes.ShortRoutes(router)
	router.Run(port)
}
