package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func EncryptPassword(secretKey, iv, plainPassword string) (string, error) {
	if _, err := io.ReadFull(rand.Reader, []byte(iv)); err != nil {
		err = fmt.Errorf(ErrEncryptPassword, err.Error())
		return "", err
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		err = fmt.Errorf(ErrEncryptPassword, err.Error())
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		err = fmt.Errorf(ErrEncryptPassword, err.Error())
		return "", err
	}

	cipherText := aesGCM.Seal(nil, []byte(iv), []byte(plainPassword), nil)
	result := base64.StdEncoding.EncodeToString(cipherText)

	return result,
		nil
}

func DecryptPassword(secretKey, iv, hashPassword string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		err = fmt.Errorf(ErrDecryptPassword, err.Error())
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		err = fmt.Errorf(ErrDecryptPassword, err.Error())
		return "", err
	}

	decHashPassword, err := base64.StdEncoding.DecodeString(hashPassword)
	if err != nil {
		err = fmt.Errorf(ErrDecryptPassword, err.Error())
		return "", err
	}

	decryptedText, err := aesGCM.Open(nil, []byte(iv), decHashPassword, nil)
	if err != nil {
		err = fmt.Errorf(ErrEncryptPassword, err.Error())
		return "", err
	}

	return string(decryptedText),
		nil
}

func Encrypt(plainText string, secretKey string) (string, error) {
	if len(secretKey) != 32 {
		err := fmt.Errorf(ErrEncrypt, "secretKey must be 32 bytes long")
		return "", err
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		err = fmt.Errorf(ErrEncrypt, err.Error())
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		err = fmt.Errorf(ErrEncrypt, err.Error())
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plainText))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertextBase64 string, secretKey string) (string, error) {
	if len(secretKey) != 32 {
		err := fmt.Errorf(ErrDecrypt, "secretKey must be 32 bytes long")
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		err := fmt.Errorf(ErrDecrypt, err.Error())
		return "", err
	}

	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		err := fmt.Errorf(ErrDecrypt, err.Error())
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		err := fmt.Errorf(ErrDecrypt, "ciphertext too short")
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
