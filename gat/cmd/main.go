package main

import (
	"log"

	routing "github.com/MarNawar/food-delivery-system/api-gateway"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routing.SetupRoutes(r)

	port := "8000"
	err := r.Run(":" + port)
	if err != nil {
		log.Fatal("failed to start a server")
	}
}
