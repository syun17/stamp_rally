package main

import (
	"log"

	controller "go_back/go_backend/go_back/Controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false,
	}))

	api := r.Group("/api")
	{
		api.GET("/spots", controller.GetSpots)
		api.POST("/stamps", controller.PostStamp)
		api.POST("/stamps/nfc", controller.PostStampByNfc)
		api.GET("/stamps/:user_id", controller.GetUserStamps)
		api.GET("/bingo/:user_id", controller.GetBingo)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
