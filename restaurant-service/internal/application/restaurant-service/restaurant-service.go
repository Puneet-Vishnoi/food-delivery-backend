package restaurantservice

import (
	"fmt"
	"time"

	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/adapter/respository"
	constant "github.com/MarNawar/food-delivery-backend/restaurant-service/internal/domain/entity/models/constants"
	"github.com/MarNawar/food-delivery-backend/restaurant-service/internal/domain/entity/models/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestaurantService struct {
	restaurentRepo *respository.RestaurantRepo
}

func NewRestaurantServiceProvider(restaurentRepo *respository.RestaurantRepo) RestaurantService {
	return RestaurantService{restaurentRepo: restaurentRepo}
}


func (serv *RestaurantService) CreateRestaurantService(name, location string) (string, error) {
	var req response.Restaurant

	req.Name = name
	req.Location = location
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	req.Status = "active"

	restaurantResp, err := serv.restaurentRepo.Insert(req, constant.RestaurantCollection)

	if err != nil {
		return "", err
	}

	// Ensure the ID is converted properly from MongoDB's ObjectID
	insertedID, ok := restaurantResp.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert inserted ID to ObjectID")
	}
	

	return insertedID.Hex(), err
}

func (serv *RestaurantService) GetRestaurantService(id string)(*response.Restaurant, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	restaurantResp, err := serv.restaurentRepo.GetSingleRestaurantById(objectID, constant.RestaurantCollection)

	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}

	return restaurantResp, nil
}

func (serv *RestaurantService)ListRestaurantsService()(*[]response.Restaurant, error){
	restaurantResp, err := serv.restaurentRepo.GetRestaurantList(constant.RestaurantCollection)

	if err != nil{
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}

	return restaurantResp, nil
}

func (serv *RestaurantService)AddMenuItemService( restaurant_id, name, description string, price float32, availability	bool )(string, error){
	var req response.MenuItem

	req.RestaurantID = restaurant_id
	req.Name = name
	req.Description = description
	req.Price = float64(price)
	req.Availability = availability

	menuResp, err := serv.restaurentRepo.Insert(req, constant.MenuCollection)
	if err != nil{
		return "", err
	}

	insertedID, ok := menuResp.(primitive.ObjectID)

	if !ok {
		return "", fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	return insertedID.Hex(), err
}

func (serv *RestaurantService)GetMenu( restaurantId string)([]*response.MenuItem, error){
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid menu ID: %w", err)
	// }

	menuResp, err := serv.restaurentRepo.GetMenuById(restaurantId, constant.MenuCollection)

	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}

	return menuResp, nil
}