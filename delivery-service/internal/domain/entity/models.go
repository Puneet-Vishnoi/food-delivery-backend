package entity

import "github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/request"

type AssignOrder struct {
	Id             string `bson:"_id" json:"_id"`
	PersonnelID    string `bson:"personnel_id" json:"personnel_id"`
	PickupLocation request.Location `bson:"pickup_location" json:"pickup_location"`
	DeliveryLocation request.Location `bson:"delivery_location" json:"delivery_location"`
}

