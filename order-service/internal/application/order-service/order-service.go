package orderservice

import (
	"fmt"
	"time"

	"github.com/MarNawar/food-delivery-backend/order-service/internal/adapter/repository"
	"github.com/MarNawar/food-delivery-backend/order-service/internal/domain/entity/models/constants"
	"github.com/MarNawar/food-delivery-backend/order-service/internal/domain/entity/models/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	orderRepo repository.OrderRepo
}

func (serv *OrderService) CreateOrder(order request.Order)(string, error) {
	order.Status = "pending"
	order.CreatedAt = time.Now().Unix()
	order.UpdatedAt = time.Now().Unix()
	var totalPrice float64
	for _, itm := range order.Items {
		totalPrice = totalPrice + itm.Price*float64(itm.Quantity)
	}

	order.TotalAmount = totalPrice
	id, err := serv.orderRepo.Insert(order, constants.OrderCollection)

	insertedID, ok := id.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return insertedID.Hex(), err
}

func (serv *OrderService) GetOrderById(id string)(*request.Order, error){
	order, err := serv.orderRepo.GetOrderById(id, constants.OrderCollection)

	if err != nil{
		return nil, err
	}

	return order, nil
}


func (serv *OrderService)UpdateOrderStatus(id, status string)(bool, error){
	_, err := serv.orderRepo.UpdateOrderStatus(id, status, constants.OrderCollection)

	if err != nil{
		return false, err
	}

	return true, nil
}

func (serv *OrderService) ListOrdersByUserId(userId string) ([]*request.Order, error) {
	orders, err := serv.orderRepo.ListOrdersByUserId(userId, constants.OrderCollection)

	if err != nil{
		return nil, err
	}

	return orders,  nil
}


func (serv *OrderService) ListOrdersByRestaurant(restaurantId string) ([]*request.Order, error) {
	orders, err := serv.orderRepo.ListOrdersByRestaurant(restaurantId, constants.OrderCollection)

	if err != nil{
		return nil, err
	}

	return orders,  nil
}
