package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/MarNawar/food-delivery-backend/order-service/internal/adapter/repository/mongo/providers"
	"github.com/MarNawar/food-delivery-backend/order-service/internal/domain/entity/models/constants"
	"github.com/MarNawar/food-delivery-backend/order-service/internal/domain/entity/models/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRepo struct {
	DBHelper *providers.DBHelper
}

func (repo *OrderRepo) Insert(data interface{}, collectionName string) (interface{}, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constants.Database).Collection(collectionName)
	result, err := orgCollection.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	log.Println(result)
	return result.InsertedID, nil
}

func (repo *OrderRepo) GetOrderById(id string, collectionName string) (*request.Order, error) {
	primitiveId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	orgCollection := repo.DBHelper.MongoClient.Database(constants.Database).Collection(collectionName)

	var order request.Order

	err = orgCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: primitiveId}}).Decode(&order)

	if err != nil {
		return nil, fmt.Errorf("failed to find order with given ID: %w", err)
	}

	return &order, nil
}

func (repo *OrderRepo) UpdateOrderStatus(id, status, collectionName string) (bool, error) {
	primitiveId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	orgCollection := repo.DBHelper.MongoClient.Database(constants.Database).Collection(collectionName)

	// Build the update query
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: status},
		}},
	}

	// Execute the update operation
	result, err := orgCollection.UpdateOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: primitiveId}},
		update,
	)

	if err != nil {
		return false, fmt.Errorf("failed to update order status: %w", err)
	}

	// Check if any document was updated
	if result.MatchedCount == 0 {
		return false, fmt.Errorf("no order found with the given ID")
	}

	return true, nil
}


func (repo *OrderRepo) ListOrdersByUserId(userId string, collectionName string) ([]*request.Order, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constants.Database).Collection(collectionName)

	cursor, err := orgCollection.Find(context.TODO(), bson.D{{Key: "user_id", Value: userId}})
	if err != nil {
		return []*request.Order{}, err
	}
	defer cursor.Close(context.TODO())

	var orders []*request.Order

	// Iterate through the cursor
	for cursor.Next(context.TODO()) {
		var order request.Order
		if err := cursor.Decode(&order); err != nil {
			log.Printf("Failed to decode order: %v", err)
			continue
		}
		orders = append(orders, &order)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered an error: %w", err)
	}

	return orders, nil
}

func (repo *OrderRepo) ListOrdersByRestaurant(restaurantId string, collectionName string)([]*request.Order, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constants.Database).Collection(collectionName)

	cursor, err := orgCollection.Find(context.TODO(), bson.D{{Key: "restaurent_id", Value: restaurantId}})
	if err != nil {
		return []*request.Order{}, err
	}
	defer cursor.Close(context.TODO())

	var orders []*request.Order

	// Iterate through the cursor
	for cursor.Next(context.TODO()) {
		var order request.Order
		if err := cursor.Decode(&order); err != nil {
			log.Printf("Failed to decode order: %v", err)
			continue
		}
		orders = append(orders, &order)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered an error: %w", err)
	}

	return orders, nil
}