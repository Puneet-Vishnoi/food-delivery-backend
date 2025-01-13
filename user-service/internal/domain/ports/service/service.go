package service

import "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/response"

type UserService interface {
	VerifyEmail(string) (string, error)
	VerifyOtp( string, int64) (string, error) 
	RegisterUser(string, string, string, []byte) (*response.User, string, error)

	AddAddressOfUser(string, string, string, string) error 

}
