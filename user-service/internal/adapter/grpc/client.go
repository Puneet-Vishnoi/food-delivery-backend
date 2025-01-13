package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/MarNawar/food-delivery-backend/user-service/internal/adapter/grpc/pb"
	requestModels "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/request"
	responseModels "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/response"
	"google.golang.org/grpc"

	utils "github.com/MarNawar/food-delivery-backend/user-service/pkg/utils"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.UserServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url)
	if err != nil {
		return nil, err
	}

	c := pb.NewUserServiceClient(conn)

	return &Client{
		conn,
		c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
func (c *Client) VerifyEmail(ctx context.Context, email string) (string, error) {
	r, err := c.service.VerifyEmail(
		ctx,
		&pb.VerifyEmailRequest{
			Email: email,
		},
	)
	if err != nil {
		return "", err
	}
	return r.Message, nil
}
func (c *Client) VerifyOtp(ctx context.Context, email string, otp int64) (string, error) {
	r, err := c.service.VerifyOTP(
		ctx,
		&pb.VerifyOTPRequest{
			Email: email,
			Otp:   otp,
		},
	)

	if err != nil {
		return "", err
	}

	return r.Message, nil
}

func (c *Client) RegisterUser(ctx context.Context, user *requestModels.UserClient) (*responseModels.User, error) {
	err := utils.CheckUserValidation(*user)
	if err != nil {
		log.Println(err)
		// c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return nil, err
	}

	pbUser := &pb.UserClient{
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Password: utils.GenPassHash(user.Password),
	}
	r, err := c.service.RegisterUser(
		ctx,
		&pb.UserRequest{
			User: pbUser,
		},
	)

	if err != nil {
		return nil, err
	}
	// Check if r.User is not nil to prevent runtime panic
	if r.User == nil {
		return nil, errors.New("received nil user in response")
	}

	userResponse := &responseModels.User{
		Name:      r.User.Name,
		Email:     r.User.Email,
		Phone:     r.User.Phone,
		UserType:  r.User.UserType,
		CreatedAt: r.User.CreatedAt,
		UpdatedAt: r.User.UpdatedAt,
	}

	return userResponse, nil
}

func (c *Client) UserLogin(ctx context.Context, email string, password []byte) (string, error) {
	r, err := c.service.AuthenticateUser(ctx, &pb.AuthRequest{
		Login: &pb.Login{
			Email:    email,
			Password: password,
		},
	})

	if err != nil {
		return "", err
	}

	return r.Message, nil
}

func (c *Client) UserAddress(ctx context.Context, address *requestModels.AddressClient, email string) (string, error) {

	r, err := c.service.UserAddress(ctx, &pb.AddressRequest{
		Address: &pb.AddressClient{
			Email:     email,
			Address_1: address.Address1,
			UserId:    address.UserId,
			City:      address.City,
			Country:   address.Country,
		},
	})

	if err != nil {
		return "", err
	}

	return r.Message, nil
}
