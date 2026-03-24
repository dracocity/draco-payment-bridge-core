package crypto

import (
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func GenerateHMACSignature(payload any, hashFn func() hash.Hash, secret string, encoding string) (string, error) {
	var data []byte
	switch v := payload.(type) {
	case nil:
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		data = b
	}

	mac := hmac.New(hashFn, []byte(secret))
	_, _ = mac.Write(data)

	var sig string
	switch encoding {
	case "base64":
		sig = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	case "hex":
		sig = hex.EncodeToString(mac.Sum(nil))
	default:
		return "", fmt.Errorf("unsupported encoding: %s", encoding)
	}
	return sig, nil
}

func VerifyHMACSignature(payload any, signature string, hashFn func() hash.Hash, secret string, encoding string) (bool, error) {
	expected, err := GenerateHMACSignature(payload, hashFn, secret, encoding)
	if err != nil {
		return false, err
	}
	return hmac.Equal([]byte(expected), []byte(strings.ToLower(signature))), nil
}

func GenerateECDSASignature(payload any, privKey *secp256k1.PrivateKey) (string, error) {
	var data []byte
	switch v := payload.(type) {
	case nil:
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		data = b
	}

	hash := sha256.Sum256([]byte(data))
	sig, err := ecdsa.SignASN1(rand.Reader, privKey.ToECDSA(), hash[:])
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}
