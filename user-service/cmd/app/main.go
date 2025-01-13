package main

import (
	"log"

	"github.com/MarNawar/food-delivery-backend/user-service/internal/adapter/grpc"
	repo "github.com/MarNawar/food-delivery-backend/user-service/internal/adapter/respository"
	db "github.com/MarNawar/food-delivery-backend/user-service/internal/adapter/respository/mongo"
	"github.com/MarNawar/food-delivery-backend/user-service/internal/adapter/respository/mongo/providers"
	userservice "github.com/MarNawar/food-delivery-backend/user-service/internal/application/user-service"
)

func main() {
	mongoClient := db.ConnectDB()
	dbHelper := providers.NewDbProvider(mongoClient.MongoClient)
	dashboardRepository := repo.NewUserRepo(dbHelper)

	servicePort := userservice.NewUserServiceProvider(dashboardRepository)
	port := 50051 // You can modify the port if needed

	log.Printf("Starting gRPC server on port %d...", port)
	if err := grpc.ListenGRPC(servicePort, port); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

}
