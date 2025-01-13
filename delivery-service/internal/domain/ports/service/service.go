package service

import "github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/request"

type GeolocationWorkflowService interface {
	IsWithinGeoFence(request.Location, []request.Location) bool
	GetDistanceFromPickupToPersonnel(float32, float32, float32, float32) (string, error) 
	CalculateDistance(float32, float32, float32, float32) (string, error) 
}

type DeliveryPersonalService interface{
	AssignOrder(string, request.Location, request.Location) (*request.DeliveryPersonal, error)
	UpdatePersonnelLocation( string, request.Location) error 
	UpdateDeliveryStatus(string, request.Location) error
	GetDeliveryStatus(string) (*request.DeliveryPersonal, error) 
}