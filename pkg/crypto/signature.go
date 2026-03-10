package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func GenerateSignature(payload any, secret string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))

	switch v := payload.(type) {
	case nil:
	case []byte:
		_, _ = mac.Write(v)
	case string:
		_, _ = mac.Write([]byte(v))
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		_, _ = mac.Write(b)
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func VerifySignature(payload any, signature string, secret string) (bool, error) {
	expected, err := GenerateSignature(payload, secret)
	if err != nil {
		return false, err
	}
	return hmac.Equal([]byte(expected), []byte(signature)), nil
}
