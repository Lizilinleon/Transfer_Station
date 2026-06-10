package util

import (
	crand "crypto/rand"
	"math/big"
)

const keyChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateAccessKey() string {
	key, err := generateRandomCharsKey(48)
	if err != nil {
		return "sk-fallback-key-generation-error"
	}
	return "sk-" + key
}

func generateRandomCharsKey(length int) (string, error) {
	chars := make([]byte, length)
	maxIndex := big.NewInt(int64(len(keyChars)))

	for index := range chars {
		number, err := crand.Int(crand.Reader, maxIndex)
		if err != nil {
			return "", err
		}
		chars[index] = keyChars[number.Int64()]
	}

	return string(chars), nil
}
