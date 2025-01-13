package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/adapter/grpc/pb"
	restaurantservice "github.com/MarNawar/food-delivery-backend/restaurant-service/internal/application/restaurant-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedRestaurantServiceServer
	service restaurantservice.RestaurantService
}

func ListenGRPC(s restaurantservice.RestaurantService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	serv := grpc.NewServer()
	pb.RegisterRestaurantServiceServer(serv, &grpcServer{service: s})

	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) CreateRestaurant(ctx context.Context, req *pb.CreateRestaurantRequest) (*pb.CreateRestaurantResponse, error) {

	id, err := s.service.CreateRestaurantService(req.Name, req.Location)
	if err != nil {
		return nil, err
	}

	return &pb.CreateRestaurantResponse{Id: id}, nil
}

func (s *grpcServer)GetRestaurant(ctx context.Context, req *pb.GetRestaurantRequest)(*pb.GetRestaurantResponse, error){
	restaurant, err :=  s.service.GetRestaurantService(req.Id)
	if err != nil{
		return nil, err
	}
	return &pb.GetRestaurantResponse{
		Id: restaurant.ID,
		Name: restaurant.Name,
		Location: restaurant.Location,
		Status: restaurant.Status,
	}, nil
}

func (s *grpcServer)ListRestaurants(ctx context.Context, req *pb.ListRestaurantsRequest)(*pb.ListRestaurantsResponse, error){
	restaurants, err := s.service.ListRestaurantsService()
	if err != nil{
		return nil, err
	}

	pbRestaurants := []*pb.GetRestaurantResponse{}

	for _, restaurant := range *restaurants{
		pbRestaurants = append(pbRestaurants, &pb.GetRestaurantResponse{
			Id:       restaurant.ID,
			Name:     restaurant.Name,
			Location: restaurant.Location,
			Status:   restaurant.Status,
		} )
	}

	return &pb.ListRestaurantsResponse{
		Restaurants: pbRestaurants,
	}, nil
}

func (s *grpcServer)AddMenuItem(ctx context.Context, req *pb.AddMenuItemRequest)(*pb.AddMenuItemResponse, error){
	id, err := s.service.AddMenuItemService(req.RestaurantId, req.Name, req.Description, req.Price, req.Availability)

	if err != nil{
		return nil, err
	}

	
	return &pb.AddMenuItemResponse{Id: id}, nil
}

func (s *grpcServer) GetMenuItem(ctx context.Context, req *pb.GetMenuRequest)(*pb.GetMenuResponse, error){
	menus, err := s.service.GetMenu(req.RestaurantId)


	if err != nil{
		return nil, err
	}

	pbMenuItems := []*pb.MenuItem{}

	for _, menu := range menus{
		pbMenuItems = append(pbMenuItems, &pb.MenuItem{
			Id:       menu.ID,
			Name:     menu.Name,
			Description: menu.Description,
			Price:   float32(menu.Price),
			Availability: menu.Availability,
		} )
	}


	return &pb.GetMenuResponse{
		Items: pbMenuItems,
	}, nil
}