package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/MarNawar/food-delivery-backend/user-service/internal/adapter/grpc/pb"
	userservice "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/ports/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedUserServiceServer
	service userservice.UserService
}

func ListenGRPC(s userservice.UserService, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}
	serv := grpc.NewServer()
	pb.RegisterUserServiceServer(serv, &grpcServer{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) VerifyEmail(ctx context.Context, r *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	msg, err := s.service.VerifyEmail(r.Email)
	if err != nil {
		return nil, err
	}

	return &pb.VerifyEmailResponse{
		Message: msg,
	}, nil
}

func (s *grpcServer) VerifyOtp(ctx context.Context, r *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	msg, err := s.service.VerifyOtp(r.Email, r.Otp)
	if err != nil {
		return nil, err
	}

	return &pb.VerifyOTPResponse{
		Message: msg,
	}, nil
}

func (s *grpcServer) RegisterUser(ctx context.Context, user *pb.UserRequest) (*pb.UserResponse, error) {
	userResp, token, err := s.service.RegisterUser(user.User.Email, user.User.Name, user.User.Phone, user.User.Password)
	if err != nil {
		return nil, err
	}
	pbRespUser := pb.User{
		Id:        userResp.Id.Hex(),
		Name:      userResp.Name,
		Email:     userResp.Email,
		Phone:     userResp.Phone,
		Password:  []byte(userResp.Password),
		UserType:  userResp.UserType,
		CreatedAt: userResp.CreatedAt,
		UpdatedAt: userResp.UpdatedAt,
		// string id = 1; // ObjectID will be string in proto
		// string name = 2;
		// string email = 3;
		// string phone = 4;
		// bytes password = 5;
		// string user_type = 6;
		// int64 created_at = 7;
		// int64 updated_at = 8;

	}
	return &pb.UserResponse{
		Message: "User Registered Succesfully",
		Token:   token,
		User:    &pbRespUser,
	}, nil

}

func (s *grpcServer) AddAddressOfUser(ctx context.Context, req *pb.AddressRequest) (*pb.AddressResponse, error) {
	err := s.service.AddAddressOfUser(req.Address.Email, req.Address.Address_1, req.Address.City, req.Address.Country)
	if err != nil {
		return nil, err
	}

	return &pb.AddressResponse{
		Message: "address added succesfully",
	}, nil
}
