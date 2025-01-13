package apigateway

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"

	restaurant  "github.com/MarNawar/food-delivery-backend/restaurant-service"
)

type APIServer struct {
	restaurantClient *restaurant.Client
}

func NewAPIServer(restaurantURL string) (*APIServer, error) {
	client, err := restaurant.NewClient(restaurantURL)
	if err != nil {
		return nil, err
	}
	return &APIServer{restaurantClient: client}, nil
}

func (s *APIServer) Close() {
	s.restaurantClient.Close()
}

func (s *APIServer) SetupRouter() *gin.Engine {
	router := gin.Default()

	// Restaurant endpoints
	router.POST("/restaurants", s.CreateRestaurant)
	router.GET("/restaurants/:id", s.GetRestaurant)
	router.GET("/restaurants", s.ListRestaurants)

	// Menu endpoints
	router.POST("/restaurants/:id/menu", s.AddMenuItem)
	router.GET("/restaurants/:id/menu", s.GetMenu)

	return router
}

// CreateRestaurant - HTTP handler for creating a new restaurant
func (s *APIServer) CreateRestaurant(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Location string `json:"location" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.restaurantClient.CreateRestaurant(context.Background(), body.Name, body.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetRestaurant - HTTP handler for fetching a restaurant by ID
func (s *APIServer) GetRestaurant(c *gin.Context) {
	id := c.Param("id")

	restaurant, err := s.restaurantClient.GetRestaurant(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// ListRestaurants - HTTP handler for listing all restaurants
func (s *APIServer) ListRestaurants(c *gin.Context) {
	restaurants, err := s.restaurantClient.ListRestaurants(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// AddMenuItem - HTTP handler for adding a menu item to a restaurant
func (s *APIServer) AddMenuItem(c *gin.Context) {
	var body struct {
		Name         string  `json:"name" binding:"required"`
		Description  string  `json:"description"`
		Price        float64 `json:"price" binding:"required"`
		Availability bool    `json:"availability"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	price := float32(body.Price) // Convert to float32 for gRPC

	menuID, err := s.restaurantClient.AddMenuItem(context.Background(), id, body.Name, body.Description, price, body.Availability)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"menu_item_id": menuID})
}

// GetMenu - HTTP handler for fetching the menu of a restaurant
func (s *APIServer) GetMenu(c *gin.Context) {
	id := c.Param("id")

	menu, err := s.restaurantClient.GetMenu(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}
