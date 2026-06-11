package util

import (
	crand "crypto/rand"
	"math/big"
)

const keyChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateAccessKey creates a client-facing key with an sk- prefix.
func GenerateAccessKey() string {
	key, err := generateRandomCharsKey(48)
	if err != nil {
		return "sk-fallback-key-generation-error"
	}
	return "sk-" + key
}

// generateRandomCharsKey uses crypto/rand so generated keys are hard to guess.
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
