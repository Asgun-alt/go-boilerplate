package utils

import (
	"crypto/rand"
	"go-boilerplate/pkg/constant"
	"math/big"
)

func GenerateRandomCode(length int) (string, error) {
	var result string
	alphabet := constant.CharsetLetters
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		result += string(alphabet[randomIndex.Int64()])
	}

	return result, nil
}

func GenerateRandomNumber(length int) (string, error) {
	numbering := constant.CharsetNumbers
	randomNumber := make([]byte, length)
	for i := range randomNumber {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(numbering))))
		if err != nil {
			return "", err
		}
		randomNumber[i] = numbering[randomIndex.Int64()]
	}

	return string(randomNumber), nil
}
