package handlers

import (
	"net/http"

	"github.com/MarNawar/food-delivery-backend/payment-service/internal/domain/entity/request"
	"github.com/MarNawar/food-delivery-backend/payment-service/internal/domain/ports/service"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service service.ServicePort
}

func NewHandlerProvider(service service.ServicePort) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (handler *PaymentHandler) CreatePaymentOrder(c *gin.Context) {
	var req request.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := handler.service.CreatePaymentOrder(&req)

	c.JSON(int(resp.Status), resp)
}

func (handler *PaymentHandler) VerifiedPayment(c *gin.Context){
	var req request.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := handler.service.CreatePaymentOrder(&req)

	c.JSON(int(resp.Status), resp)
}