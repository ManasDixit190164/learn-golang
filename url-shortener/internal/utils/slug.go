package utils

import (
	"crypto/rand"
	"math/big"
)

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		length = 7
	}

	result := make([]byte, length)
	max := big.NewInt(int64(len(base62Chars)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = base62Chars[n.Int64()]
	}

	return string(result), nil
}
