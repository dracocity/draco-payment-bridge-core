package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
)

type ClientOptions struct {
	Client      *http.Client
	BaseURL     string
	Header      http.Header
	ErrorPrefix string
}

type RequestOptions struct {
	Method string
	Path   string
	Query  url.Values
	Header http.Header
	Body   any
}

type Client struct {
	client      *http.Client
	baseURL     string
	header      http.Header
	errorPrefix string
	logger      logger.Logger
}

func New(opts ClientOptions) *Client {
	client := opts.Client
	if client == nil {
		client = http.DefaultClient
	}
	header := opts.Header.Clone()
	if header == nil {
		header = http.Header{}
	}
	return &Client{
		client:      client,
		baseURL:     opts.BaseURL,
		header:      header,
		errorPrefix: opts.ErrorPrefix,
		logger:      logger.WithModule("httpx"),
	}
}

func (c *Client) Get(ctx context.Context, path string, query url.Values, header http.Header, out any) error {
	return c.Do(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   path,
		Query:  query,
		Header: header,
	}, out)
}

func (c *Client) Post(ctx context.Context, path string, header http.Header, body any, out any) error {
	return c.Do(ctx, RequestOptions{
		Method: http.MethodPost,
		Path:   path,
		Header: header,
		Body:   body,
	}, out)
}

func (c *Client) Do(ctx context.Context, opts RequestOptions, out any) error {
	u, err := url.JoinPath(c.baseURL, opts.Path)
	if err != nil {
		return err
	}
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	if len(opts.Query) > 0 {
		parsedURL.RawQuery = opts.Query.Encode()
	}
	c.logger.Info("http request", "method", opts.Method, "url", parsedURL.String())

	var body io.Reader
	if opts.Body != nil {
		b, err := marshalNoEscape(opts.Body)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
		c.logger.Info("http request", "body", string(b))
	}

	req, err := http.NewRequestWithContext(ctx, opts.Method, parsedURL.String(), body)
	if err != nil {
		return err
	}

	for key, values := range c.header {
		req.Header[key] = append([]string(nil), values...)
	}
	for key, values := range opts.Header {
		req.Header[key] = append([]string(nil), values...)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.logger.Info("http response", "status", resp.StatusCode, "payload", strings.TrimSpace(string(raw)))

	errorPrefix := c.errorPrefix
	if errorPrefix == "" {
		errorPrefix = "http"
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%s api error: status %d, body %s", errorPrefix, resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("failed to decode %s response: %w", errorPrefix, err)
	}
	return nil
}

func marshalNoEscape(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}
