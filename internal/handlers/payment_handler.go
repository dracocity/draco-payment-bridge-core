package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/services"
	"github.com/gin-gonic/gin"
)

type PaymentService interface {
	Providers() []string
	CreatePaymentLink(ctx context.Context, provider string, req models.CreatePaymentLinkRequest) (*models.CreatePaymentLinkResponse, error)
	CreatePayment(ctx context.Context, provider string, req models.CreatePaymentRequest) (*models.CreatePaymentResponse, error)
	GetPayment(ctx context.Context, provider, paymentID string) (*models.GetPaymentResponse, error)
	CreateRefund(ctx context.Context, provider string, req models.CreateRefundRequest) (*models.CreateRefundResponse, error)
	HandleWebhook(ctx context.Context, provider string, payload []byte, header http.Header) (*models.WebhookResult, error)
}

// PaymentHandler provides HTTP handlers for the payment API.
type PaymentHandler struct {
	paymentService PaymentService
	logger         logger.Logger
}

func New(paymentService PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		logger:         logger.WithModule("payment-handler"),
	}
}

func (h *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/health", h.handleHealth)
	router.POST("/api/v1/providers/:provider/payment-links", h.handleCreatePaymentLink)
	router.POST("/api/v1/providers/:provider/payments", h.handleCreatePayment)
	router.GET("/api/v1/providers/:provider/payments/:paymentId", h.handleGetPayment)
	router.POST("/api/v1/providers/:provider/refunds", h.handleCreateRefund)
	router.POST("/api/v1/providers/:provider/webhooks", h.handleWebhook)
}

func (h *PaymentHandler) handleHealth(c *gin.Context) {
	resp := map[string]interface{}{
		"status":    "ok",
		"providers": h.paymentService.Providers(),
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) handleCreatePaymentLink(c *gin.Context) {
	provider := c.Param("provider")
	var req models.CreatePaymentLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.paymentService.CreatePaymentLink(c.Request.Context(), provider, req)
	if err != nil {
		if errors.Is(err, services.ErrUnknownProvider) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) handleCreatePayment(c *gin.Context) {
	provider := c.Param("provider")
	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.paymentService.CreatePayment(c.Request.Context(), provider, req)
	if err != nil {
		if errors.Is(err, services.ErrUnknownProvider) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) handleGetPayment(c *gin.Context) {
	provider := c.Param("provider")
	paymentID := c.Param("paymentId")

	resp, err := h.paymentService.GetPayment(c.Request.Context(), provider, paymentID)
	if err != nil {
		if errors.Is(err, services.ErrUnknownProvider) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) handleCreateRefund(c *gin.Context) {
	provider := c.Param("provider")

	var req models.CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.PaymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_id are required"})
		return
	}

	resp, err := h.paymentService.CreateRefund(c.Request.Context(), provider, req)
	if err != nil {
		if errors.Is(err, services.ErrUnknownProvider) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) handleWebhook(c *gin.Context) {
	provider := c.Param("provider")

	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	h.logger.Debug("webhook received", "provider", provider, "payload", strings.TrimSpace(string(payload)))

	resp, err := h.paymentService.HandleWebhook(c.Request.Context(), provider, payload, c.Request.Header)
	if err != nil {
		if errors.Is(err, services.ErrUnknownProvider) {
			c.JSON(http.StatusNotFound, gin.H{"error": "unknown provider"})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
