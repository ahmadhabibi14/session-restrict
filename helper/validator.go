package helper

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	if len(validationErrors) == 0 {
		return nil
	}

	errMsg := fmt.Sprintf(
		"error when validating %s (%s): '%v'",
		validationErrors[0].Field(),
		validationErrors[0].ActualTag(),
		validationErrors[0].Value(),
	)

	return errors.New(errMsg)
}

func IsValidPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
