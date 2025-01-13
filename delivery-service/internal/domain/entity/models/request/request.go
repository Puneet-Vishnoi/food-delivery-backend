package request

type Location struct {
	Lat float32 `bson:"lat" json:"lat"`
	Lng float32 `bson:"lng" json:"lng"`
}

type DeliveryPersonal struct {
	Id              string   `bson:"_id" json:"_id"`
	Name            string   `bson:"name" json:"name"`
	Phone           string   `bson:"phone" json:"phone"`
	Status          string   `bson:"status" json:"status"`
	CurrentLocation Location `bson:"current_location" json:"current_location"`
	Vehicle         string   `bson:"vehicle" json:"vehicle"`
}


