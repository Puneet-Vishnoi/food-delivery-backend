package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/MarNawar/food-delivery-backend/order-service/internal/adapter/grpc/pb"
	orderservice "github.com/MarNawar/food-delivery-backend/order-service/internal/application/order-service"
	"github.com/MarNawar/food-delivery-backend/order-service/internal/domain/entity/models/request"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service orderservice.OrderService
}

func ListenGRPC(s orderservice.OrderService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{service: s})

	// func pb.RegisterRestaurantServiceServer(s grpc.ServiceRegistrar, srv pb.RestaurantServiceServer)

	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var items []request.OrderItem
	for _, itm := range req.Items {
		items = append(items, request.OrderItem{
			MenuItemID: itm.MenuItemId,
			Name:       itm.Name,
			Quantity:   int(itm.Quantity),
			Price:      float64(itm.Price),
		})
	}

	id, err := s.service.CreateOrder(request.Order{
		UserID:       req.UserId,
		RestaurantID: req.RestaurantId,
		Items:        items,
	})

	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
			Id: id,
		},
		nil

}

func (s *grpcServer) GetOrderById(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.service.GetOrderById(req.Id)

	if err != nil {
		return nil, err
	}

	pbItems := []*pb.OrderItem{}

	for _, itm := range order.Items {
		pbItems = append(pbItems, &pb.OrderItem{
			MenuItemId: itm.MenuItemID,
			Name:       itm.Name,
			Quantity:   int32(itm.Quantity),
			Price:      float32(itm.Price),
		})
	}

	return &pb.GetOrderResponse{
		Id:           order.ID,
		RestaurantId: order.RestaurantID,
		UserId:       order.UserID,
		Items:        pbItems,
		TotalAmount:  float32(order.TotalAmount),
		Status:       order.Status,
	}, nil
}

func (s *grpcServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	isUpdated, err := s.service.UpdateOrderStatus(req.Id, req.Status)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateOrderStatusResponse{
		Success: isUpdated,
	}, nil
}

func (s *grpcServer) ListOrdersByUser(ctx context.Context, req *pb.ListOrdersByUserRequest) (*pb.ListOrdersByUserResponse, error) {
	orders, err := s.service.ListOrdersByUserId(req.UserId)
	if err != nil {
		return nil, err
	}

	pbOrders := []*pb.GetOrderResponse{}

	for _, order := range orders {
		pbItems := []*pb.OrderItem{}

		for _, itm := range order.Items {
			pbItems = append(pbItems, &pb.OrderItem{
				MenuItemId: itm.MenuItemID,
				Name:       itm.Name,
				Quantity:   int32(itm.Quantity),
				Price:      float32(itm.Price),
			})
		}
		pbOrders = append(pbOrders, &pb.GetOrderResponse{
			Id:           order.ID,
			UserId:       order.UserID,
			RestaurantId: order.RestaurantID,
			Items:        pbItems,
			TotalAmount:  float32(order.TotalAmount),
			Status:       order.Status,
		})
	}

	return &pb.ListOrdersByUserResponse{
		Orders: pbOrders,
	}, nil

}


func (s *grpcServer) ListOrdersByRestaurant(ctx context.Context, req *pb.ListOrdersByRestaurantRequest) (*pb.ListOrdersByRestaurantResponse, error) {
	orders, err := s.service.ListOrdersByUserId(req.RestaurantId)
	if err != nil {
		return nil, err
	}

	pbOrders := []*pb.GetOrderResponse{}

	for _, order := range orders {
		pbItems := []*pb.OrderItem{}

		for _, itm := range order.Items {
			pbItems = append(pbItems, &pb.OrderItem{
				MenuItemId: itm.MenuItemID,
				Name:       itm.Name,
				Quantity:   int32(itm.Quantity),
				Price:      float32(itm.Price),
			})
		}
		pbOrders = append(pbOrders, &pb.GetOrderResponse{
			Id:           order.ID,
			UserId:       order.UserID,
			RestaurantId: order.RestaurantID,
			Items:        pbItems,
			TotalAmount:  float32(order.TotalAmount),
			Status:       order.Status,
		})
	}

	return &pb.ListOrdersByRestaurantResponse{
		Orders: pbOrders,
	}, nil

}
