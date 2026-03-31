package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/services"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	if err := logger.Init(config.LogConfig{
		Level:  "error",
		Format: "json",
		Output: config.LogOutputConfig{
			Stdout: true,
		},
		Rotation: &config.LogRotationConfig{
			MaxSize:  10,
			MaxAge:   1,
			Compress: false,
		},
	}); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

type mockPaymentService struct {
	providers []string

	createPaymentLinkFn func(ctx context.Context, provider string, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error)
	createPaymentFn     func(ctx context.Context, provider string, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error)
	getPaymentFn        func(ctx context.Context, provider, paymentID string) (*models.GetPaymentResponse, error)
	createRefundFn      func(ctx context.Context, provider string, req models.CreateRefundRequest) (*models.CreateRefundResponse, error)
	handleWebhookFn     func(ctx context.Context, provider string, payload []byte, header http.Header) (*models.WebhookResult, error)
}

func (m *mockPaymentService) Providers() []string {
	return m.providers
}

func (m *mockPaymentService) CreatePaymentLink(ctx context.Context, provider string, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	return m.createPaymentLinkFn(ctx, provider, req)
}

func (m *mockPaymentService) CreatePayment(ctx context.Context, provider string, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	return m.createPaymentFn(ctx, provider, req)
}

func (m *mockPaymentService) GetPayment(ctx context.Context, provider, paymentID string) (*models.GetPaymentResponse, error) {
	return m.getPaymentFn(ctx, provider, paymentID)
}

func (m *mockPaymentService) CreateRefund(ctx context.Context, provider string, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
	return m.createRefundFn(ctx, provider, req)
}

func (m *mockPaymentService) HandleWebhook(ctx context.Context, provider string, payload []byte, header http.Header) (*models.WebhookResult, error) {
	return m.handleWebhookFn(ctx, provider, payload, header)
}

func newTestRouter(svc *mockPaymentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	New(svc).RegisterRoutes(r)
	return r
}

func decodeBody[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()
	var out T
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, body = %s", err, rr.Body.String())
	}
	return out
}

func TestPaymentHandler_Health(t *testing.T) {
	r := newTestRouter(&mockPaymentService{providers: []string{"nowpayments", "testpg"}})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status field = %v, want ok", body["status"])
	}

	providersAny, ok := body["providers"].([]any)
	if !ok {
		t.Fatalf("providers field type = %T, want []any", body["providers"])
	}
	if len(providersAny) != 2 {
		t.Fatalf("len(providers) = %d, want 2", len(providersAny))
	}
}

func TestPaymentHandler_CreatePaymentLink_InvalidJSON(t *testing.T) {
	r := newTestRouter(&mockPaymentService{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/nowpayments/payment-links", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rr.Code)
	}
	resp := decodeBody[map[string]string](t, rr)
	if resp["error"] != "invalid request body" {
		t.Fatalf("error = %q, want invalid request body", resp["error"])
	}
}

func TestPaymentHandler_CreatePaymentLink_ErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{name: "unknown provider", err: services.ErrUnknownProvider, wantStatus: http.StatusBadRequest, wantError: "unknown provider"},
		{name: "provider error", err: errors.New("timeout"), wantStatus: http.StatusBadGateway, wantError: "provider error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := newTestRouter(&mockPaymentService{
				createPaymentLinkFn: func(ctx context.Context, provider string, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
					return nil, tc.err
				},
			})

			req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/missing/payment-links", bytes.NewBufferString(`{"amount":"1","currency":"USD"}`))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tc.wantStatus)
			}
			resp := decodeBody[map[string]string](t, rr)
			if resp["error"] != tc.wantError {
				t.Fatalf("error = %q, want %q", resp["error"], tc.wantError)
			}
		})
	}
}

func TestPaymentHandler_CreatePayment_Success(t *testing.T) {
	r := newTestRouter(&mockPaymentService{
		createPaymentFn: func(ctx context.Context, provider string, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
			if provider != "nowpayments" {
				t.Fatalf("provider = %q, want nowpayments", provider)
			}
			return &models.CreatePaymentResponse{PaymentID: "pay-1", Currency: "USD", Amount: "1"}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/nowpayments/payments", bytes.NewBufferString(`{"invoice_id":"inv-1","receive_currency":"btc"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	resp := decodeBody[models.CreatePaymentResponse](t, rr)
	if resp.PaymentID != "pay-1" {
		t.Fatalf("payment_id = %q, want pay-1", resp.PaymentID)
	}
}

func TestPaymentHandler_GetPayment_ErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{name: "unknown provider", err: services.ErrUnknownProvider, wantStatus: http.StatusBadRequest, wantError: "unknown provider"},
		{name: "provider error", err: errors.New("upstream"), wantStatus: http.StatusBadGateway, wantError: "provider error"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := newTestRouter(&mockPaymentService{
				getPaymentFn: func(ctx context.Context, provider, paymentID string) (*models.GetPaymentResponse, error) {
					if paymentID != "p-1" {
						t.Fatalf("paymentID = %q, want p-1", paymentID)
					}
					return nil, tc.err
				},
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/missing/payments/p-1", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, tc.wantStatus)
			}
			resp := decodeBody[map[string]string](t, rr)
			if resp["error"] != tc.wantError {
				t.Fatalf("error = %q, want %q", resp["error"], tc.wantError)
			}
		})
	}
}

func TestPaymentHandler_CreateRefund_ValidationAndErrorMapping(t *testing.T) {
	t.Run("missing payment_id", func(t *testing.T) {
		r := newTestRouter(&mockPaymentService{})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/nowpayments/refunds", bytes.NewBufferString(`{"amount":1.2}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
		resp := decodeBody[map[string]string](t, rr)
		if resp["error"] != "payment_id are required" {
			t.Fatalf("error = %q, want payment_id are required", resp["error"])
		}
	})

	t.Run("unknown provider", func(t *testing.T) {
		r := newTestRouter(&mockPaymentService{
			createRefundFn: func(ctx context.Context, provider string, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
				return nil, services.ErrUnknownProvider
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/missing/refunds", bytes.NewBufferString(`{"payment_id":"p-1","amount":1.2}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("provider error", func(t *testing.T) {
		r := newTestRouter(&mockPaymentService{
			createRefundFn: func(ctx context.Context, provider string, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
				return nil, errors.New("provider down")
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/nowpayments/refunds", bytes.NewBufferString(`{"payment_id":"p-1","amount":1.2}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadGateway {
			t.Fatalf("status = %d, want 502", rr.Code)
		}
	})
}

func TestPaymentHandler_Webhook_ErrorMappingAndSuccess(t *testing.T) {
	t.Run("unknown provider returns 404", func(t *testing.T) {
		r := newTestRouter(&mockPaymentService{
			handleWebhookFn: func(ctx context.Context, provider string, payload []byte, header http.Header) (*models.WebhookResult, error) {
				return nil, services.ErrUnknownProvider
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/missing/webhooks", bytes.NewBufferString(`{"event":"paid"}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rr.Code)
		}
	})

	t.Run("provider error returns 502", func(t *testing.T) {
		r := newTestRouter(&mockPaymentService{
			handleWebhookFn: func(ctx context.Context, provider string, payload []byte, header http.Header) (*models.WebhookResult, error) {
				return nil, errors.New("upstream")
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/nowpayments/webhooks", bytes.NewBufferString(`{"event":"paid"}`))
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadGateway {
			t.Fatalf("status = %d, want 502", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		r := newTestRouter(&mockPaymentService{
			handleWebhookFn: func(ctx context.Context, provider string, payload []byte, header http.Header) (*models.WebhookResult, error) {
				if provider != "nowpayments" {
					t.Fatalf("provider = %q, want nowpayments", provider)
				}
				if string(payload) != `{"event":"paid"}` {
					t.Fatalf("payload = %s, want event payload", string(payload))
				}
				return &models.WebhookResult{Accepted: true, Message: "verified"}, nil
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/nowpayments/webhooks", bytes.NewBufferString(`{"event":"paid"}`))
		req.Header.Set("X-Test", "1")
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		resp := decodeBody[models.WebhookResult](t, rr)
		if !resp.Accepted || resp.Message != "verified" {
			t.Fatalf("resp = %#v, want accepted verified", resp)
		}
	})
}
