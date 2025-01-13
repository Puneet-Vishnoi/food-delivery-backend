package router

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MarNawar/food-delivery-backend/payment-service/internal/adapter/http/handlers"
	"github.com/MarNawar/food-delivery-backend/payment-service/internal/domain/entity/constants"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
}

func NewServer(dashboardHandler *handlers.PaymentHandler) *Server {
	gin.SetMode(gin.DebugMode)

	server := gin.Default()

	server.Use(func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		log.Println("startTime: ", startTime)
		log.Println("path: ", path)
		log.Println("query: ", query)

		c.Next()
	})

	payment := server.Group("/api/payment")

	payment.GET("/create", dashboardHandler.CreatePaymentOrder)
	payment.GET("/validate", dashboardHandler.VerifiedPayment)

	return &Server{
		Engine: server,
	}
}

func (s *Server) Start() {
	if err := s.Engine.Run(os.Getenv(constants.GIN_PORT)); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}


