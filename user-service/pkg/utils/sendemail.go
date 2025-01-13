package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	constant "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/constants"
	"github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/request"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/exp/rand"

	
)

func SendEmailSendGrid(req request.Verification) (request.Verification, error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return req, errors.New("SENDGRID_API_KEY environment variable is not set")
	}

	// Create a SendGrid client
	client := sendgrid.NewSendClient(apiKey)

	// Set up the email message
	from := mail.NewEmail("Sender Name", constant.Sender)
	to := mail.NewEmail("Recipient Name", req.Email)
	subject := "OTP verification mail"

	otp := Randomnum()
	req.Otp = int64(otp)
	htmlContent := "<p>This is a test otp forverification <strong>" + strconv.Itoa(otp) + "</strong> </p>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	// Send the email message
	_, err := client.Send(message)
	if err != nil {
		return req, err
	}

	return req, nil
}

func Randomnum() int {
	rng := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	return rng.Intn(1000) + 1000 // OTP length of 4 digits
}
