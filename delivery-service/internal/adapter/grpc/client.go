package grpc

import (
	"context"

	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/adapter/grpc/pb"
	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/request"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.DeliveryServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url)
	if err != nil {
		return nil, err
	}

	c := pb.NewDeliveryServiceClient(conn)

	return &Client{
		conn,
		c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) AssignOrder(ctx context.Context, orderID string, pickupLocation *request.Location, deliveryLocation *request.Location) (*request.DeliveryPersonal, error) {
	r, err := c.service.AssignOrder(
		ctx,
		&pb.AssignOrderRequest{
			OrderID: orderID,
			PickupLocation: &pb.Location{
				Lat: pickupLocation.Lat,
				Lng: pickupLocation.Lng,
			},
			DeliveryLocation: &pb.Location{
				Lat: deliveryLocation.Lat,
				Lng: deliveryLocation.Lng,
			},
		},
	)

	if err != nil {
		return nil, err
	}
	return &request.DeliveryPersonal{
		Id:     r.DeliveryPersonal.Id,
		Name:   r.DeliveryPersonal.Name,
		Phone:  r.DeliveryPersonal.Phone,
		Status: r.DeliveryPersonal.Status,
		CurrentLocation: request.Location{
			Lat: r.DeliveryPersonal.CurrentLocation.Lat,
			Lng: r.DeliveryPersonal.CurrentLocation.Lng,
		},
		Vehicle: r.DeliveryPersonal.Vehicle,
	}, nil
}

func (c *Client) GetDeliveryStatus(ctx context.Context, orderID string) (*request.DeliveryPersonal, error) {
	r, err := c.service.GetDeliveryStatus(
		ctx,
		&pb.GetDeliveryStatusRequest{
			OrderID: orderID,
		},
	)

	if err != nil {
		return nil, err
	}
	return &request.DeliveryPersonal{
		Id:     r.DeliveryPersonal.Id,
		// Name:   r.DeliveryPersonal.Name,
		// Phone:  r.DeliveryPersonal.Phone,
		Status: r.DeliveryPersonal.Status,
		// CurrentLocation: request.Location{
		// 	Lat: r.DeliveryPersonal.CurrentLocation.Lat,
		// 	Lng: r.DeliveryPersonal.CurrentLocation.Lng,
		// },
		// Vehicle: r.DeliveryPersonal.Vehicle,
	}, nil
}

func (c *Client) UpdatePersonnelLocation(ctx context.Context, id string, deliveryLocation *request.Location) (string, error) {
	r, err := c.service.UpdatePersonnelLocation(
		ctx,
		&pb.UpdateLocationRequest{
			Id: id,
			CurrentLocation: &pb.Location{
				Lat: deliveryLocation.Lat,
				Lng: deliveryLocation.Lng,
			},
		},
	)

	if err != nil {
		return "", err
	}
	return r.Status, nil
}
