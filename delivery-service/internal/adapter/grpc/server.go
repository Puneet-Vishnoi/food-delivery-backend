package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/adapter/grpc/pb"
	"github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/entity/models/request"
	deliveryPersonalService "github.com/MarNawar/food-delivery-backend/delivery-service/internal/domain/ports/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedDeliveryServiceServer
	service deliveryPersonalService.DeliveryPersonalService
}

func ListenGRPC (s deliveryPersonalService.DeliveryPersonalService, port int)error{
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil{
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	serv := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) AssignOrder(ctx context.Context, r *pb.AssignOrderRequest)(*pb.AssignOrderResponse, error){
	pickupLocation := request.Location{
		Lat : r.PickupLocation.Lat,
		Lng: r.PickupLocation.Lng,
	}
	deliveryLocation := request.Location{
		Lat: r.DeliveryLocation.Lat,
		Lng: r.DeliveryLocation.Lng,
	}
	deliveryPersonal, err := s.service.AssignOrder(r.OrderID, pickupLocation, deliveryLocation)
	

	if err != nil{
		return nil, err
	}

	return &pb.AssignOrderResponse{
		DeliveryPersonal: &pb.DeliveryPersonal{
			Id: deliveryPersonal.Id,
			Name: deliveryPersonal.Name,
			Phone: deliveryPersonal.Phone,
			Status: deliveryPersonal.Status,
			Vehicle: deliveryPersonal.Vehicle,
			CurrentLocation: &pb.Location{
				Lat : deliveryPersonal.CurrentLocation.Lat,
				Lng : deliveryPersonal.CurrentLocation.Lng,
			},
		},
	}, nil
}

func (s *grpcServer) UpdatePersonnelLocation(ctx context.Context, r *pb.UpdateLocationRequest)(*pb.UpdateLocationResponse, error){
	personnelLocation := request.Location{
		Lat: r.CurrentLocation.Lat,
		Lng: r.CurrentLocation.Lng,
	}
	err := s.service.UpdatePersonnelLocation(r.Id, personnelLocation)
	if err != nil{
		return nil, err
	}
	
	err = s.service.UpdateDeliveryStatus(r.Id, personnelLocation)
	if err != nil{
		return nil, err
	}

	return &pb.UpdateLocationResponse{
		Status: "Location Updated",
	}, nil
}

func (s *grpcServer) GetDeliveryStatus(ctx context.Context, r *pb.GetDeliveryStatusRequest)(*pb.GetDeliveryStatusResponse, error){
	// resp, err:= s.service.GetDeliveryStatus()
	deliveryPersonal, err := s.service.GetDeliveryStatus(r.OrderID)
	if err != nil{
		return nil, err
	}
	return &pb.GetDeliveryStatusResponse{
		DeliveryPersonal: &pb.DeliveryPersonal{
			Id: deliveryPersonal.Id,
			Status: deliveryPersonal.Status,	
		},
	}, nil
}
// (ctx context.Context, orderID string) (*request.DeliveryPersonal, error)