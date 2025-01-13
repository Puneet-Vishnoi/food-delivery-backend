package grpc

import (
	"context"

	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/adapter/grpc/pb"
	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/domain/entity/models/response"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.RestaurantServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url)
	if err != nil {
		return nil, err
	}

	c := pb.NewRestaurantServiceClient(conn)

	return &Client{
		conn,
		c,
	}, nil

}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateRestaurant(ctx context.Context, name string, location string) (string, error) {
	r, err := c.service.CreateRestaurant(
		ctx,
		&pb.CreateRestaurantRequest{
			Name:     name,
			Location: location,
		},
	)

	if err != nil {
		return "", err
	}
	return r.Id, nil
}

func (c *Client) GetRestaurant(ctx context.Context, id string) (*response.Restaurant, error) {
	r, err := c.service.GetRestaurant(
		ctx,
		&pb.GetRestaurantRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}

	resp := &response.Restaurant{
		ID:       r.Id,
		Name:     r.Name,
		Location: r.Location,
		Status:   r.Status,
	}

	return resp, nil
}

func (c *Client) ListRestaurants(ctx context.Context) (*[]response.Restaurant, error) {
	r, err := c.service.ListRestaurants(
		ctx,
		&pb.ListRestaurantsRequest{},
	)

	if err != nil {
		return nil, err
	}

	var resp []response.Restaurant
	for _, restaurant := range r.Restaurants {
		resp = append(resp, response.Restaurant{
			ID:       restaurant.Id,
			Name:     restaurant.Name,
			Location: restaurant.Location,
			Status:   restaurant.Status,
		})
	}

	return &resp, nil
}

func (c *Client) AddMenuItem(ctx context.Context, restaurant_id, name, description string, price float32, availability bool) (string, error) {
	r, err := c.service.AddMenuItem(
		ctx,
		&pb.AddMenuItemRequest{
			RestaurantId: restaurant_id,
			Name:         name,
			Description:  description,
			Price:        price,
			Availability: availability,
		},
	)

	if err != nil {
		return "", err
	}

	return r.Id, nil
}

func (c *Client) GetMenu(ctx context.Context, id string) (*[]response.MenuItem, error) {
	r, err := c.service.GetMenu(
		ctx,
		&pb.GetMenuRequest{
			RestaurantId: id,
		},
	)

	if err != nil {
		return nil, err
	}
	var resp []response.MenuItem
	for _, menuItem := range r.Items {
		resp = append(resp, response.MenuItem{
			ID:           menuItem.Id,
			Name:         menuItem.Name,
			Description:  menuItem.Description,
			Price:        float64(menuItem.Price),
			Availability: menuItem.Availability,
			RestaurantID: id,
		})
	}

	return &resp, nil
}
