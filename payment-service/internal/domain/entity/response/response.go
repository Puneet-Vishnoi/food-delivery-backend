package response

type Response struct{
	Status int64
	Message string
	Data interface{}
	Error error
}
type PaymentResponse struct {
	Id string `json:"id"`
	OrderID     string `json:"order_id"`
	RazorpayKey string `json:"razorpay_key"`
}