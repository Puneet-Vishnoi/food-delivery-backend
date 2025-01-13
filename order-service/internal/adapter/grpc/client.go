package grpc

import (
	"context"
	"strconv"

	"github.com/MarNawar/food-delivery-backend/order-service/internal/adapter/grpc/pb"
	"github.com/MarNawar/food-delivery-backend/order-service/internal/domain/entity/models/request"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url)
	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)

	return &Client{
		conn,
		c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateOrder(ctx context.Context, user_id string, restaurant_id string, items []request.OrderItem) (string, error) {

	pbitems := []*pb.OrderItem{}

	for _, item := range items {
		pbitems = append(pbitems, &pb.OrderItem{
			MenuItemId: item.MenuItemID,
			Name:       item.Name,
			Quantity:   int32(item.Quantity),
			Price:      float32(item.Price),
		})
	}

	r, err := c.service.CreateOrder(
		ctx,
		&pb.CreateOrderRequest{
			UserId:       user_id,
			RestaurantId: restaurant_id,
			Items:        pbitems,
		},
	)

	if err != nil {
		return "", err
	}

	return r.Id, nil
}

func (c *Client) GetOrder(ctx context.Context, id string) (*request.Order, error) {
	r, err := c.service.GetOrder(
		ctx,
		&pb.GetOrderRequest{
			Id: id,
		},
	)

	if err != nil {
		return nil, err
	}

	respItems := []request.OrderItem{}

	for _, item := range r.Items {
		respItems = append(respItems, request.OrderItem{
			MenuItemID: item.MenuItemId,
			Name:       item.Name,
			Quantity:   int(item.Quantity),
			Price:      float64(item.Price),
		})
	}

	// Parse CreatedAt string to int64 (assuming it's a Unix timestamp string)
	createdAt, err := strconv.ParseInt(r.CreatedAt, 10, 64)
	if err != nil {
		return nil, err // Handle error if parsing fails
	}

	return &request.Order{
		ID:           r.Id,
		UserID:       r.UserId,
		RestaurantID: r.RestaurantId,
		Items:        respItems,
		TotalAmount:  float64(r.TotalAmount),
		Status:       r.Status,
		CreatedAt:    createdAt,
	}, nil
}

func (c *Client) UpdateOrderStatus(ctx context.Context, id, status string) (bool, error) {
	r, err := c.service.UpdateOrderStatus(
		ctx,
		&pb.UpdateOrderStatusRequest{
			Id:     id,
			Status: status,
		},
	)

	if err != nil {
		return false, err
	}

	return r.Success, nil
}

func (c *Client) ListOrdersByUser(ctx context.Context, userId string) ([]*request.Order, error) {
	r, err := c.service.ListOrdersByUser(
		ctx,
		&pb.ListOrdersByUserRequest{
			UserId: userId,
		},
	)

	if err != nil {
		return nil, err
	}

	var resp []*request.Order
	for _, userOrders := range r.Orders {
		var orderItems []request.OrderItem
		for _, item := range userOrders.Items {
			orderItems = append(orderItems, request.OrderItem{
				MenuItemID: item.MenuItemId,
				Name:       item.Name,
				Quantity:   int(item.Quantity),
				Price:      float64(item.Price),
			})
		}
		resp = append(resp, &request.Order{
			ID:           userOrders.Id,
			UserID:       userOrders.UserId,
			RestaurantID: userOrders.RestaurantId,
			Items:        orderItems,
			TotalAmount:  float64(userOrders.TotalAmount),
			Status:       userOrders.Status,
		})
	}
	
	return resp, nil
}

func (c *Client)ListOrdersByRestaurant(ctx context.Context, restaurantId string)([]*request.Order, error){
	r, err := c.service.ListOrdersByRestaurant(
		ctx,
		&pb.ListOrdersByRestaurantRequest{
			RestaurantId: restaurantId,
		},
	)

	if err != nil{
		return nil, err
	}

	var resp []*request.Order
	for _, userOrders := range r.Orders {
		var orderItems []request.OrderItem
		for _, item := range userOrders.Items {
			orderItems = append(orderItems, request.OrderItem{
				MenuItemID: item.MenuItemId,
				Name:       item.Name,
				Quantity:   int(item.Quantity),
				Price:      float64(item.Price),
			})
		}
		resp = append(resp, &request.Order{
			ID:           userOrders.Id,
			UserID:       userOrders.UserId,
			RestaurantID: userOrders.RestaurantId,
			Items:        orderItems,
			TotalAmount:  float64(userOrders.TotalAmount),
			Status:       userOrders.Status,
		})
	}

	return resp, nil
}
