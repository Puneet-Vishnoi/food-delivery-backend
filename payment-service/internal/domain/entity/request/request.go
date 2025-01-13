package request

type PaymentRequest struct {
	OrderID  string  `json:"order_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}



type VerifyPaymentRequest struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature"`
}
