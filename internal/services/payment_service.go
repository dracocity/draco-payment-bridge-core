package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
)

var ErrUnknownProvider = errors.New("unknown provider")

// PaymentService handles payment-provider orchestration.
type PaymentService struct {
	pgRegistry *pg.Registry
	logger     logger.Logger
}

func NewPaymentService(pgRegistry *pg.Registry) *PaymentService {
	return &PaymentService{
		pgRegistry: pgRegistry,
		logger:     logger.WithModule("payment-service"),
	}
}

func (s *PaymentService) Providers() []string {
	return s.pgRegistry.Names()
}

func (s *PaymentService) CreatePaymentLink(ctx context.Context, provider string, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error) {
	p, ok := s.pgRegistry.Get(provider)
	if !ok {
		return nil, ErrUnknownProvider
	}

	resp, err := p.CreatePaymentLink(ctx, req)
	if err != nil {
		s.logger.Warn("failed to create payment-link", "provider", provider, "error", err)
		return nil, fmt.Errorf("create payment-link: %w", err)
	}
	return resp, nil
}

func (s *PaymentService) CreatePayment(ctx context.Context, provider string, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error) {
	p, ok := s.pgRegistry.Get(provider)
	if !ok {
		return nil, ErrUnknownProvider
	}

	resp, err := p.CreatePayment(ctx, req)
	if err != nil {
		s.logger.Warn("failed to create payment", "provider", provider, "error", err)
		return nil, fmt.Errorf("create payment: %w", err)
	}
	return resp, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, provider, paymentID string) (*models.GetPaymentResponse, error) {
	p, ok := s.pgRegistry.Get(provider)
	if !ok {
		return nil, ErrUnknownProvider
	}

	resp, err := p.GetPayment(ctx, paymentID)
	if err != nil {
		s.logger.Warn("failed to get payment status", "provider", provider, "error", err)
		return nil, fmt.Errorf("get payment: %w", err)
	}
	return resp, nil
}

func (s *PaymentService) CreateRefund(ctx context.Context, provider string, req models.CreateRefundRequest) (*models.CreateRefundResponse, error) {
	p, ok := s.pgRegistry.Get(provider)
	if !ok {
		return nil, ErrUnknownProvider
	}

	resp, err := p.CreateRefund(ctx, req)
	if err != nil {
		s.logger.Warn("failed to create refund", "provider", provider, "error", err)
		return nil, fmt.Errorf("create refund: %w", err)
	}
	return resp, nil
}

func (s *PaymentService) HandleWebhook(ctx context.Context, provider string, payload []byte, headers map[string][]string) (*models.WebhookResult, error) {
	p, ok := s.pgRegistry.Get(provider)
	if !ok {
		return nil, ErrUnknownProvider
	}

	resp, err := p.HandleWebhook(ctx, payload, headers)
	if err != nil {
		s.logger.Warn("failed to process webhook", "provider", provider, "error", err)
		return nil, fmt.Errorf("handle webhook: %w", err)
	}
	return resp, nil
}
