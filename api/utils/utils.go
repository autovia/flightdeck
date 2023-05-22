package utils

import "crypto/rand"

func RandomString() (string, error) {
	bytes := make([]byte, 256)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
