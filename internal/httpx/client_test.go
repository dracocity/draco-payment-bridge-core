package httpx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dracocity/draco-payment-bridge-core/internal/config"
	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
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

func TestClientDo_MergesHeadersAndQueryAndDecodesResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1/payments" {
			t.Fatalf("path = %s, want /v1/payments", r.URL.Path)
		}
		if got := r.URL.Query().Get("foo"); got != "bar" {
			t.Fatalf("query foo = %q, want bar", got)
		}
		if got := r.Header.Get("X-Base"); got != "base-header" {
			t.Fatalf("X-Base = %q, want base-header", got)
		}
		if got := r.Header.Get("X-Request"); got != "request-header" {
			t.Fatalf("X-Request = %q, want request-header", got)
		}

		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}
		body := strings.TrimSpace(string(rawBody))
		if !strings.Contains(body, "\"text\":\"a<b\"") {
			t.Fatalf("body = %s, want unescaped html content", body)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"id":"p-1"}`))
	}))
	defer ts.Close()

	client := New(ClientOptions{
		Client:      ts.Client(),
		BaseURL:     ts.URL,
		Header:      http.Header{"X-Base": []string{"base-header"}},
		ErrorPrefix: "nowpayments",
	})

	var out struct {
		OK bool   `json:"ok"`
		ID string `json:"id"`
	}

	err := client.Do(context.Background(), RequestOptions{
		Method: http.MethodPost,
		Path:   "/v1/payments",
		Query:  url.Values{"foo": []string{"bar"}},
		Header: http.Header{"X-Request": []string{"request-header"}},
		Body: map[string]string{
			"text": "a<b",
		},
	}, &out)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if !out.OK || out.ID != "p-1" {
		t.Fatalf("out = %#v, want ok=true and id=p-1", out)
	}
}

func TestClientDo_ReturnsProviderErrorOnNon2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad request"))
	}))
	defer ts.Close()

	client := New(ClientOptions{
		Client:      ts.Client(),
		BaseURL:     ts.URL,
		ErrorPrefix: "nowpayments",
	})

	err := client.Get(context.Background(), "/v1/payments", nil, nil, nil)
	if err == nil {
		t.Fatal("Get() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "nowpayments api error: status 400") {
		t.Fatalf("err = %v, want nowpayments status error", err)
	}
	if !strings.Contains(err.Error(), "bad request") {
		t.Fatalf("err = %v, want response body", err)
	}
}

func TestClientDo_ReturnsDecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not-json"))
	}))
	defer ts.Close()

	client := New(ClientOptions{
		Client:      ts.Client(),
		BaseURL:     ts.URL,
		ErrorPrefix: "nowpayments",
	})

	var out map[string]any
	err := client.Get(context.Background(), "/v1/payments", nil, nil, &out)
	if err == nil {
		t.Fatal("Get() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "failed to decode nowpayments response") {
		t.Fatalf("err = %v, want decode error", err)
	}
}

func TestMarshalNoEscape(t *testing.T) {
	b, err := marshalNoEscape(map[string]string{"text": "a<b"})
	if err != nil {
		t.Fatalf("marshalNoEscape() error = %v", err)
	}
	if strings.Contains(string(b), "\\u003c") {
		t.Fatalf("output = %s, want non-escaped '<'", string(b))
	}
	if !strings.Contains(string(b), "\"text\":\"a<b\"") {
		t.Fatalf("output = %s, want text field", string(b))
	}

	_, err = marshalNoEscape(make(chan int))
	if err == nil {
		t.Fatal("marshalNoEscape(chan) expected error")
	}
}
