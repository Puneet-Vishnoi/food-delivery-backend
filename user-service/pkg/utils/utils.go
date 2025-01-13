package utils

import (
	"errors"
	"strconv"

	"github.com/MarNawar/food-delivery-backend/user-service/internal/domain/entity/models/request"
	"golang.org/x/crypto/bcrypt"
)

func GenPassHash(s string) []byte {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return nil
	}
	return bytes
}

func CheckUserValidation(u request.UserClient) error {
	if u.Email == "" {
		return errors.New("email can't be empty")
	}
	if u.Name == "" {
		return errors.New("name can't be empty")
	}
	if u.Phone == "" {
		return errors.New("phone can't be empty")
	}
	if u.Password == "" {
		return errors.New("password can't be empty")
	}
	return nil
}

func ConvertStringIntoInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val

}
