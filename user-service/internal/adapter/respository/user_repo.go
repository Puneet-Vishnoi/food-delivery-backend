package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/MarNawar/food-delivery-backend/user-service/internal/adapter/respository/mongo/providers"
	constant "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/constants"
	request "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/request"
	response "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepo struct {
	DBHelper *providers.DBHelper
}

func NewUserRepo(DBHelper *providers.DBHelper) *UserRepo {
	return &UserRepo{
		DBHelper: DBHelper,
	}
}


func (repo *UserRepo) Insert(data interface{}, collectionName string) (interface{}, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	result, err := orgCollection.InsertOne(context.TODO(), data)
	if err != nil {
		return nil, err
	}
	log.Println(result)
	return result.InsertedID, nil
}

func (repo *UserRepo) GetSingleRecordByEmail(email string, collectionName string) *request.Verification {
	resp := &request.Verification{}
	filter := bson.D{{Key: "email", Value: email}}
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	err := orgCollection.FindOne(context.TODO(), filter).Decode(&resp)
	fmt.Println(err)
	return resp
}

func (repo *UserRepo) UpdateVerification(data request.Verification, collectionName string) error {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "email", Value: data.Email}}
	update := bson.D{{Key: "$set", Value: data}}
	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (repo *UserRepo) UpdateEmailVerifiedStatus(req request.Verification, collectionName string) error {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "email", Value: req.Email}}
	update := bson.D{{Key: "$set", Value: req}}
	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (repo *UserRepo) GetSingleRecordByEmailForUser(email, collectionName string) *response.User {
	resp := &response.User{}
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "email", Value: email}}
	_ = orgCollection.FindOne(context.TODO(), filter).Decode(&resp)
	return resp
}

func (repo *UserRepo) GetSingleAddress(id primitive.ObjectID, collectionName string) (response.Address, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "user_id", Value: id}}
	var address response.Address
	err := orgCollection.FindOne(context.TODO(), filter).Decode(&address)
	return address, err
}

func (repo *UserRepo) GetSingleUserByUserId(id primitive.ObjectID, collectionName string) (response.User, error) {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "_id", Value: id}}
	var user response.User
	err := orgCollection.FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

func (repo *UserRepo) UpdateUser(u response.User, collectionName string) error {
	orgCollection := repo.DBHelper.MongoClient.Database(constant.Database).Collection(collectionName)
	filter := bson.D{{Key: "_id", Value: u.Id}}
	update := bson.D{{Key: "$set", Value: u}}
	_, err := orgCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
