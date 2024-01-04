package utils

import (
	"crypto/rand"
	"encoding/base64"
)

const idLength = 10
const secretCodeLength = 8

func GenerateID() string {
	idBytes := make([]byte, idLength)
	_, err := rand.Read(idBytes) //common error syntax

	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(idBytes)[:idLength]
}

func GenerateSecretCode() string {
	secretCodeBytes := make([]byte, secretCodeLength)
	_, err := rand.Read(secretCodeBytes)

	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(secretCodeBytes)[:secretCodeLength]
}
