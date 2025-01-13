package application

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MarNawar/food-delivery-backend/payment-service/internal/domain/entity/request"
	"github.com/MarNawar/food-delivery-backend/payment-service/internal/domain/entity/response"
	"github.com/razorpay/razorpay-go"
	utils "github.com/razorpay/razorpay-go/utils"
)

type RazorPayService struct {
	Client *razorpay.Client
}

func NewRazorpayClient() *razorpay.Client {
	return razorpay.NewClient(os.Getenv("RAZORPAY_ID"), os.Getenv("RAZORPAY_SECRET_KEY"))
}

func NewRazorPayService() *RazorPayService {
	return &RazorPayService{
		Client: NewRazorpayClient(),
	}
}

func (serv *RazorPayService) CreatePaymentOrder(req *request.PaymentRequest) (*response.Response) {
	data := map[string]interface{}{
		"amount":   req.Amount * 100, // Razorpay requires amount in paise
		"currency": req.Currency,
		"receipt":  req.OrderID,
	}

	resp, err := serv.Client.Order.Create(data, nil)
	if err != nil {
		return &response.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create payment order",
			Data:    nil,
			Error:   err,
		}
	}

	res := &response.Response{
		Status:  http.StatusOK,
		Message: "Payment order created successfully",
		Data: response.PaymentResponse{
			Id:          resp["id"].(string),
			OrderID:     req.OrderID,
			RazorpayKey: os.Getenv("RAZORPAY_KEY"),
		},
		Error: nil,
	}

	return res
}

func (serv *RazorPayService) VerifyPayment(req *request.VerifyPaymentRequest) *response.Response {
	params := map[string]interface{}{
		"razorpay_order_id":   req.OrderID,
		"razorpay_payment_id": req.PaymentID,
	}

	isVerify := utils.VerifyPaymentSignature(params, req.Signature, os.Getenv("RAZORPAY_SECRET"))
	if !isVerify {
		return &response.Response{
			Status:  http.StatusInternalServerError,
			Message: "Payment verification failed",
			Data: nil,
			Error: fmt.Errorf("payment verification failed: %s", req.OrderID),
		}
	}

	return &response.Response{
		Status:  http.StatusOK,
		Message: "Payment verified successfully",
		Data: nil,
		Error: nil,
	}
}

// func (serv *RazorPayService) RefundPayment(paymentId string) {
// 	var refundRequest struct {
// 		PaymentID string `json:"payment_id"`
// 	}

// 	// client := services.NewRazorpayClient()
// 	data := map[string]interface{}{
// 		"payment_id": paymentId,
// 	}

// 	resp, err := serv.Client.Payment.Refund(refundRequest.PaymentID, data, nil)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate refund"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"refund": resp})
// }
