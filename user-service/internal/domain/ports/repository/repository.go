package repository

import (
	"github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/request"
	"github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Insert(interface{}, string) (interface{}, error)
	GetSingleRecordByEmail(string, string) *request.Verification
	UpdateVerification(data request.Verification, collectionName string) error
	UpdateEmailVerifiedStatus(req request.Verification, collectionName string) error
	GetSingleRecordByEmailForUser(email, collectionName string) *response.User
	GetSingleAddress(primitive.ObjectID, string) (response.Address, error)
	GetSingleUserByUserId(primitive.ObjectID, string) (response.User, error)
	UpdateUser(response.User, string) error
	
}
