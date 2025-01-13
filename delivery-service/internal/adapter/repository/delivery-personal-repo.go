package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/adapter/repository/mongo/providers"
	models "github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity"
	constant "github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/constants"
	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/request"
	"go.mongodb.org/mongo-driver/bson"
)

type DeliveryPersonalRepo struct {
	DBHelper *providers.DBHelper
}

func NewDeliveryPersonalRepo(DBHelper *providers.DBHelper) *DeliveryPersonalRepo {
	return &DeliveryPersonalRepo{
		DBHelper: DBHelper,
	}
}

func (repo *DeliveryPersonalRepo) Insert(data interface{}, collectionName string) (interface{}, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	result, err := orgCollection.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	log.Println(result)
	return result.InsertedID, nil
}

func (repo *DeliveryPersonalRepo) GetAvailablePersonnel(collectionName string) ([]*request.DeliveryPersonal, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	cursor, err := orgCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var personnels []*request.DeliveryPersonal

	// Iterate through the cursor
	for cursor.Next(context.TODO()) {
		var deliveryPersonal request.DeliveryPersonal
		if err := cursor.Decode(&deliveryPersonal); err != nil {
			log.Printf("Failed to decode deliveryPersonal: %v", err)
			continue
		}
		personnels = append(personnels, &deliveryPersonal)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor encountered an error: %w", err)
	}

	return personnels, nil
}

func (repo *DeliveryPersonalRepo) AssignOrderToPersonnel(orderID string, nearestPersonnelId string, pickupLocation *request.Location, deliveryLocation *request.Location, collectionName string) error {
	assignedOrder := models.AssignOrder{
		Id:               orderID,
		PersonnelID:      nearestPersonnelId,
		PickupLocation:   *pickupLocation,
		DeliveryLocation: *deliveryLocation,
	}

	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	_, err := orgCollection.InsertOne(context.TODO(), assignedOrder)
	if err != nil {
		return err
	}
	return nil
}

func (repo *DeliveryPersonalRepo) UpdatePersonnelStatus(personnelID string, statusUpdate string, collectionName string)(error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	update := bson.D{
		{Key:"$set", Value: bson.D{
			{Key: "status", Value: statusUpdate},
		}},
	}
	_, err := orgCollection.UpdateOne(context.TODO(), bson.D{{Key:"_id", Value: personnelID}}, update)
	if err != nil {
		return err
	}
	return nil
}

func (repo *DeliveryPersonalRepo)GetAssignedPersonnelOrderById(orderId string, collectionName string)(string, error){
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)

	var result models.AssignOrder
	err := orgCollection.FindOne(context.TODO(), bson.D{{Key:"_id", Value: orderId}}).Decode(result)
	if err != nil {
		return "", err
	}
	return result.PersonnelID, nil
}

func (repo *DeliveryPersonalRepo)GetAssignedOrderByPersonnelId(PersonnelID string, collectionName string)(*models.AssignOrder, error){
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)

	var result models.AssignOrder
	err := orgCollection.FindOne(context.TODO(), bson.D{{Key:"personnel_id", Value: PersonnelID}}).Decode(result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *DeliveryPersonalRepo) GetDeliveryStatus(personnelID string, collectionName string)(*request.DeliveryPersonal, error){
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)

	var result request.DeliveryPersonal
	err := orgCollection.FindOne(context.TODO(), bson.D{{Key:"_id", Value: personnelID}}).Decode(result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *DeliveryPersonalRepo) UpdatePersonnelLocation(personnelID string, currentLocation *request.Location,collectionName string)(error){
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	update := bson.D{
		{Key:"$set", Value: bson.D{
			{Key: "current_location", Value: currentLocation},
		}},
	}
	_, err := orgCollection.UpdateOne(context.TODO(), bson.D{{Key:"_id", Value: personnelID}}, update)
	if err != nil {
		return err
	}
	return nil
}