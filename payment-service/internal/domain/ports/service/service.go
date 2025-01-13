package service

import (
	"github.com/MarNawar/food-delivery-backend/payment-service/internal/domain/entity/request"
	"github.com/MarNawar/food-delivery-backend/payment-service/internal/domain/entity/response"
)

type ServicePort interface {
	CreatePaymentOrder(req *request.PaymentRequest) (*response.Response)
	VerifyPayment(req *request.VerifyPaymentRequest) (*response.Response) 

}
