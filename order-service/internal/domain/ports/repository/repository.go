package repository

import "github.com/MarNawar/food-delivery-backend/order-service/internal/domain/entity/models/request"

type OrderRepository interface {
	Insert(interface{}, string) (interface{}, error)
	GetOrderById(string, string) (*request.Order, error)
	UpdateOrderStatus(string, string, string) (bool, error)
	ListOrdersByUserId(string, string) ([]*request.Order, error)
	ListOrdersByRestaurant(string, string)([]*request.Order, error)
}
