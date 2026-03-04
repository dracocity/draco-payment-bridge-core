package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
)

func buildRequestHeaders(clientID, clientSecret, method, requestURL string, payload any) (map[string]string, error) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	signature, err := generateSignature(clientID, clientSecret, timestamp, method, requestURL, payload)
	if err != nil {
		return nil, fmt.Errorf("generate signature: %w", err)
	}

	return map[string]string{
		"Content-Type":             "application/json",
		"X-CoinPayments-Client":    clientID,
		"X-CoinPayments-Timestamp": timestamp,
		"X-CoinPayments-Signature": signature,
	}, nil
}

func generateSignature(clientID, clientSecret, timestamp, method, requestURL string, payload any) (string, error) {
	if payload == nil {
		return "", nil
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(payload); err != nil {
		return "", err
	}
	payloadMessage := strings.TrimRight(buf.String(), "\n")

	signingValue := "\uFEFF" + strings.ToUpper(method) + requestURL + clientID + timestamp + payloadMessage

	mac := hmac.New(sha256.New, []byte(clientSecret))

	_, _ = mac.Write([]byte(signingValue))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func buildCreateInvoiceRequest(req models.CreatePaymentLinkRequest) (*createInvoiceRequest, error) {
	r := &createInvoiceRequest{}

	if len(req.ProviderPayload) > 0 {
		var payload createInvoicePayload
		if err := json.Unmarshal(req.ProviderPayload, &payload); err != nil {
			return nil, fmt.Errorf("invalid provider_payload: %w", err)
		}

	}

	return r, nil
}
