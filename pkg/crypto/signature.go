package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
)

func GenerateSignature(payload any, secret string, encoding string) (string, error) {
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

	var sig string
	switch encoding {
	case "base64":
		sig = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	case "hex":
	default:
		sig = hex.EncodeToString(mac.Sum(nil))
	}
	return sig, nil
}

func VerifySignature(payload any, signature string, secret string, encoding string) (bool, error) {
	expected, err := GenerateSignature(payload, secret, encoding)
	if err != nil {
		return false, err
	}
	return hmac.Equal([]byte(expected), []byte(signature)), nil
}
