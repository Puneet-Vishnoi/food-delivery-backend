package service

import "github.com/MarNawar/food-delivery-backend/order-service/internal/domain/entity/models/request"

type OrderService interface {
	CreateOrder(request.Order) (string, error)
	GetOrderById(id string) (*request.Order, error)
	UpdateOrderStatus(string, string) (bool, error)
	ListOrdersByUserId(string) ([]*request.Order, error)
	ListOrdersByRestaurant(string) ([]*request.Order, error)
}
