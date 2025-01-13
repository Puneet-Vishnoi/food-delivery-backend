package main

import (
	"log"

	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/adapter/grpc"
	repo "github.com/MarNawar/food-delivery-backend/restaurant-service/internal/adapter/respository"
	db "github.com/MarNawar/food-delivery-backend/restaurant-service/internal/adapter/respository/mongo"
	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/adapter/respository/mongo/providers"
	restaurantservice "github.com/MarNawar/food-delivery-backend/restaurant-service/internal/application/restaurant-service"
)

func main() {
	mongoClient := db.ConnectDB()
	dbHelper := providers.NewDbProvider(mongoClient.MongoClient)
	dashboardRepository := repo.NewRestaurantRepo(dbHelper)

	servicePort := restaurantservice.NewRestaurantServiceProvider(dashboardRepository)
	port := 50051 // You can modify the port if needed

	log.Printf("Starting gRPC server on port %d...", port)
	if err := grpc.ListenGRPC(servicePort, port); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}