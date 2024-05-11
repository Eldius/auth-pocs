package helper

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var (
	NotAuthorizedErr = errors.New("user lookup")
)

func HashPassword(password string) (string, error) {
	bPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("generating bcrypt hash: %w", err)
		return "", err
	}
	return string(bPass), nil
}

func ValidatePassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		err = fmt.Errorf("%w: %w", NotAuthorizedErr, err)
		return err
	}

	return nil
}
