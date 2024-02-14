package utils

import (
	"assessment/pkg/database/mongodb/models"
	"errors"
)

func ValidateUser(user models.User) error {
	if err := ValidateUsername(user.Name); err != nil {
		return err
	}

	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("Username is required")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("Password is required")
	}
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	return nil
}
