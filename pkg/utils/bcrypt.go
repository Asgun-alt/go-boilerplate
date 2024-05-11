package utils

import (
	"go-boilerplate/pkg/constant"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword),
		nil
}

func ComparePassword(password, confirmPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(password),
		[]byte(confirmPassword),
	); err != nil {
		if err != constant.ErrPasswordNotMatch {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
