package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
)

type mockGateway struct {
	name string

	createPaymentLinkFn func(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error)
	createPaymentFn     func(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error)
	getPaymentFn        func(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error)
	createRefundFn      func(ctx context.Context, req models.CreateRefundRequest) (*models.CreateRefundResponse, error)
	handleWebhookFn     func(ctx context.Context, payload []byte, header http.Header) (*models.WebhookResult, error)
}

func (m *mockGateway) Name() string {
	return m.name
}

func (m *mockGateway) CreatePaymentLink(ctx context.Context, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	return m.createPaymentLinkFn(ctx, req)
}

func (m *mockGateway) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	return m.createPaymentFn(ctx, req)
}

func (m *mockGateway) GetPayment(ctx context.Context, paymentID string) (*models.GetPaymentResponse, error) {
	return m.getPaymentFn(ctx, paymentID)
}

func (m *mockGateway) CreateRefund(ctx context.Context, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
	return m.createRefundFn(ctx, req)
}

func (m *mockGateway) HandleWebhook(ctx context.Context, payload []byte, header http.Header) (*models.WebhookResult, error) {
	return m.handleWebhookFn(ctx, payload, header)
}

func TestPaymentService_Providers(t *testing.T) {
	r := pg.NewRegistry()
	r.Register(&mockGateway{name: "nowpayments"})
	r.Register(&mockGateway{name: "another"})

	svc := NewPaymentService(r)
	providers := svc.Providers()

	if len(providers) != 2 {
		t.Fatalf("len(providers) = %d, want 2", len(providers))
	}

	seen := map[string]bool{}
	for _, name := range providers {
		seen[name] = true
	}
	if !seen["nowpayments"] || !seen["another"] {
		t.Fatalf("providers = %v, want nowpayments and another", providers)
	}
}

func TestPaymentService_UnknownProvider(t *testing.T) {
	svc := NewPaymentService(pg.NewRegistry())
	ctx := context.Background()

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "create payment-link",
			call: func() error {
				_, err := svc.CreatePaymentLink(ctx, "missing", models.CreatePaymentLinkRequest{})
				return err
			},
		},
		{
			name: "create payment",
			call: func() error {
				_, err := svc.CreatePayment(ctx, "missing", models.CreatePaymentRequest{})
				return err
			},
		},
		{
			name: "get payment",
			call: func() error {
				_, err := svc.GetPayment(ctx, "missing", "payment-id")
				return err
			},
		},
		{
			name: "create refund",
			call: func() error {
				_, err := svc.CreateRefund(ctx, "missing", models.CreateRefundRequest{})
				return err
			},
		},
		{
			name: "handle webhook",
			call: func() error {
				_, err := svc.HandleWebhook(ctx, "missing", []byte("{}"), http.Header{})
				return err
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if !errors.Is(err, ErrUnknownProvider) {
				t.Fatalf("err = %v, want ErrUnknownProvider", err)
			}
		})
	}
}

func TestPaymentService_CreatePayment_WrapsProviderError(t *testing.T) {
	providerErr := errors.New("provider timeout")
	r := pg.NewRegistry()
	r.Register(&mockGateway{
		name: "nowpayments",
		createPaymentFn: func(ctx context.Context, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
			return nil, providerErr
		},
	})

	svc := NewPaymentService(r)
	_, err := svc.CreatePayment(context.Background(), "nowpayments", models.CreatePaymentRequest{})
	if err == nil {
		t.Fatal("CreatePayment() error = nil, want non-nil")
	}

	if !errors.Is(err, providerErr) {
		t.Fatalf("err = %v, want wrapped provider error", err)
	}
	if !strings.Contains(err.Error(), "create payment") {
		t.Fatalf("err = %v, want message containing create payment", err)
	}
}

func TestPaymentService_HandleWebhook_Success(t *testing.T) {
	expected := &models.WebhookResult{Accepted: true, Message: "ok"}
	r := pg.NewRegistry()
	r.Register(&mockGateway{
		name: "nowpayments",
		handleWebhookFn: func(ctx context.Context, payload []byte, header http.Header) (*models.WebhookResult, error) {
			if string(payload) != "{\"event\":\"paid\"}" {
				t.Fatalf("payload = %s, want payment event json", string(payload))
			}
			if header.Get("X-Signature") != "sig" {
				t.Fatalf("X-Signature = %q, want sig", header.Get("X-Signature"))
			}
			return expected, nil
		},
	})

	svc := NewPaymentService(r)
	resp, err := svc.HandleWebhook(
		context.Background(),
		"nowpayments",
		[]byte("{\"event\":\"paid\"}"),
		http.Header{"X-Signature": []string{"sig"}},
	)
	if err != nil {
		t.Fatalf("HandleWebhook() error = %v", err)
	}
	if resp != expected {
		t.Fatalf("resp = %#v, want %#v", resp, expected)
	}
}
