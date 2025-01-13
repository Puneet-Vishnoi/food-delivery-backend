package respository

import (
	"context"
	"fmt"
	"log"

	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/adapter/respository/mongo/providers"
	constant "github.com/MarNawar/food-delivery-backend/restaurant-service/internal/domain/entity/models/constants"
	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/domain/entity/models/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantRepo struct {
	DBHelper *providers.DBHelper
}

func NewRestaurantRepo(DBHelper *providers.DBHelper) *RestaurantRepo {
	return &RestaurantRepo{
		DBHelper: DBHelper,
	}
}

func (repo *RestaurantRepo) Insert(data interface{}, collectionName string) (interface{}, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	result, err := orgCollection.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	log.Println(result)
	return result.InsertedID, nil
}

func (repo *RestaurantRepo) GetSingleRestaurantById(id primitive.ObjectID, collectionName string) (*response.Restaurant, error) {
	resp := &response.Restaurant{}

	filter := bson.D{{Key: "_id", Value: id}}
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)

	err := orgCollection.FindOne(context.TODO(), filter).Decode(&resp)

	return resp, err
}

func (repo *RestaurantRepo) GetRestaurantList(collectionName string) (*[]response.Restaurant, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)

	cursor, err := orgCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch restaurant list: %w", err)
	}

	defer cursor.Close(context.TODO())

	var restaurants []response.Restaurant

	// Iterate through the cursor
	for cursor.Next(context.TODO()) {
		var restaurant response.Restaurant
		if err := cursor.Decode(&restaurant); err != nil {
			log.Printf("Failed to decode restaurant: %v", err)
			continue
		}
		restaurants = append(restaurants, restaurant)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered an error: %w", err)
	}

	return &restaurants, nil
}

func (repo *RestaurantRepo) GetMenuById(restaurantId, collectionName string) ([]*response.MenuItem, error) {
	menus := []*response.MenuItem{}

	filter := bson.D{{Key: "restaurant_id", Value: restaurantId}}
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)

	cursor, err := orgCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch menus from restaurant %s: %w", restaurantId, err)
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var menu response.MenuItem
		if err := cursor.Decode(&menu); err != nil {
			log.Printf("Failed to decode menu: %v", err)
			continue
		}
		menus = append(menus, &menu)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered an error: %w", err)
	}

	return menus, err
}
