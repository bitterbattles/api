package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bitterbattles/api/pkg/crypto"
)

// NewHS256 creates a new JWT signed using HS256
func NewHS256(payload interface{}, secret string) (string, error) {
	header := Header{
		Algorithm: HS256,
		Type:      JWT,
	}
	encodedHeader, err := encode(header)
	if err != nil {
		return "", err
	}
	encodedPayload, err := encode(payload)
	if err != nil {
		return "", err
	}
	signature, err := signHS256(encodedHeader, encodedPayload, secret)
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("%s.%s.%s", encodedHeader, encodedPayload, signature)
	return token, nil
}

// DecodeHS256 decodes a JWT signed using HS256
func DecodeHS256(token string, secret string, payload interface{}) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("malformed token")
	}
	header := Header{}
	err := decode(parts[0], &header)
	if err != nil {
		return err
	}
	if header.Algorithm != HS256 {
		return errors.New("unsupported JWT algorithm")
	}
	verifyHS256(parts[0], parts[1], parts[2], secret)
	err = decode(parts[1], payload)
	if err != nil {
		return err
	}
	return nil
}

func encode(obj interface{}) (string, error) {
	json, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString([]byte(json))
	return encoded, nil
}

func decode(part string, obj interface{}) error {
	partBytes, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return err
	}
	return json.Unmarshal(partBytes, obj)
}

func signHS256(encodedHeader string, encodedPayload string, secret string) (string, error) {
	value := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)
	signatureBytes, err := crypto.HS256(value, secret)
	if err != nil {
		return "", err
	}
	signature := base64.RawURLEncoding.EncodeToString(signatureBytes)
	return signature, nil
}

func verifyHS256(encodedHeader string, encodedPayload string, signature string, secret string) error {
	expectedSignature, err := signHS256(encodedHeader, encodedPayload, secret)
	if err != nil {
		return err
	}
	if expectedSignature != signature {
		return errors.New("signature does not match")
	}
	return nil
}
