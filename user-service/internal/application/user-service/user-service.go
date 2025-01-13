package userservice

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/MarNawar/food-delivery-backend/user-service/internal/application/auth"
	constant "github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/constants"
	"github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/request"
	"github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/response"
	"github.com/MarNawar/food-delivery-backend/user-service/internal/domain/ports/repository"
	"github.com/MarNawar/food-delivery-backend/user-service/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserServiceProvider(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (serv *UserService) VerifyEmail(email string) (string, error) {
	var req request.Verification

	if email == "" {
		return "", errors.New(constant.EmailValidationError)
	}

	resp := serv.userRepo.GetSingleRecordByEmail(email, constant.VerificationsCollection)

	if resp.Otp != 0 {
		expirationTime := resp.CreatedAt + constant.OtpValidation
		if expirationTime < time.Now().Unix() {
			req, checkEmail := utils.SendEmailSendGrid(req)
			if checkEmail != nil {
				log.Panicln(checkEmail)

				return "", errors.New(constant.EmailValidationError)
			}
			req.CreatedAt = time.Now().Unix()
			serv.userRepo.UpdateVerification(req, constant.VerificationsCollection)

			return "", errors.New("OTP sent successfully")
		}
		return "", errors.New(constant.OptAlreadySentError)
	}

	req, checkEmail := utils.SendEmailSendGrid(req)
	if checkEmail != nil {
		log.Println(checkEmail)
		return "", errors.New(constant.EmailValidationError)
	}

	req.CreatedAt = time.Now().Unix()
	serv.userRepo.Insert(req, constant.VerificationsCollection)

	return "OTP sent successfully", nil
}

func (serv *UserService) VerifyOtp(email string, otp int64) (string, error) {
	var req request.Verification

	// Check if email and OTP fields are provided in the request
	if email == "" {
		return "", errors.New(constant.EmailValidationError)
	}
	if otp <= 0 {
		return "", errors.New(constant.OtpValidationError)
	}

	req.Email = email
	req.Otp = otp

	// Fetch the OTP record associated with the given email

	resp := serv.userRepo.GetSingleRecordByEmail(email, constant.VerificationsCollection)

	// Check if the email has already been verified
	if resp.Status {
		return "", errors.New(constant.AlreadyVerifiedError)

	}

	// Verify the OTP and check if it is expired
	expirationTime := resp.CreatedAt + constant.OtpValidation
	if resp.Otp != req.Otp {
		// c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpValidationError})
		return "", errors.New(constant.OtpValidationError)
	}
	if expirationTime < time.Now().Unix() {
		// c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpExpiredValidationError})
		return "", errors.New(constant.OtpValidationError)
	}

	// Update the verification record to mark the email as verified
	req.Status = true
	req.CreatedAt = time.Now().Unix()
	err := serv.userRepo.UpdateEmailVerifiedStatus(req, constant.VerificationsCollection)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpValidationError})
		return "", errors.New(constant.OtpValidationError)
	}
	// c.JSON(http.StatusOK, gin.H{"error": false, "message": "Email verified successfully"})
	return "Email verified successfully", nil
}

func (serv *UserService) RegisterUser(name, email, phone string, password []byte) (*response.User, string, error) {
	var dbUser response.User

	// Check if the email is verified in the verification records
	verificationResp := serv.userRepo.GetSingleRecordByEmail(email, constant.VerificationsCollection)
	if !verificationResp.Status {
		return nil, "", errors.New(constant.EmailIsNotVerified)
	}

	userResp := serv.userRepo.GetSingleRecordByEmailForUser(email, constant.UserCollection)
	if userResp.Email != "" {
		return nil, "", errors.New(constant.AlreadyRegisterWithThisEmail)

	}

	dbUser.Email = email
	dbUser.Name = name
	dbUser.Phone = phone
	dbUser.UserType = constant.NormalUser
	dbUser.Password = string(password)
	dbUser.CreatedAt = time.Now().Unix()
	dbUser.UpdatedAt = time.Now().Unix()

	InsertedID, err := serv.userRepo.Insert(dbUser, constant.UserCollection)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	jwtWrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}
	userID := InsertedID.(primitive.ObjectID)
	token, err := jwtWrapper.GenrateToken(userID, email, constant.NormalUser)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	return &dbUser, token, nil
}

// UserLogin authenticates a user and generates a JWT token
func (serv *UserService) UserLogin(email string, password []byte) (string, error) {
	// Fetch the user record from the database using the email
	userResp := serv.userRepo.GetSingleRecordByEmailForUser(email, constant.UserCollection)
	if userResp.Email == "" {
		return "", errors.New(constant.NotRegisteredUser)
	}

	// Validate the user's password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(userResp.Password), password); err != nil {
		return "", errors.New(constant.PasswordNotMatchedError)
	}

	// Generate a JWT token for the authenticated user
	jwtWrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48,
	}
	token, err := jwtWrapper.GenrateToken(userResp.Id, userResp.Email, userResp.UserType)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (serv *UserService) AddAddressOfUser(email string, address1, city, country string) error {
	// email, ok := c.Get("email")
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotAuthorizedUserError})
	// 	return
	// }

	userDBResp := serv.userRepo.GetSingleRecordByEmail(email, constant.UserCollection)
	if userDBResp.Email == "" {
		// c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotRegisteredUser})
		return errors.New(constant.NotRegisteredUser)
	}

	// var addressReq request.AddressClient
	// err := c.BindJSON(&addressReq)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
	// 	return
	// }

	// userId, err := primitive.ObjectIDFromHex(addressReq.UserId)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
	// 	return
	// }

	var addressDB response.Address

	addressDB.Address1 = address1
	addressDB.UserId = userDBResp.ID
	addressDB.City = city
	addressDB.Country = country

	_, err := serv.userRepo.Insert(addressDB, constant.AddressCollection)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return err
	}

	// c.JSON(http.StatusOK, gin.H{"message": "success", "error": false})

	return nil

}
