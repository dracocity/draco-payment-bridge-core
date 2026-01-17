package handlers

import (
	"net/http"

	"github.com/dracocity/draco-payment-bridge-core/internal/bridge"
	"github.com/dracocity/draco-payment-bridge-core/internal/logger"
	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/gin-gonic/gin"
)

// Handler provides HTTP handlers for the payment bridge API.
type Handler struct {
	registry *bridge.Registry
	logger   logger.Logger
}

func New(registry *bridge.Registry) *Handler {
	return &Handler{
		registry: registry,
		logger:   logger.WithModule("handlers"),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.handleHealth)
	router.POST("/api/payment", h.handleCreatePayment)
	router.GET("/api/payment/status", h.handlePaymentStatus)
	router.POST("/api/payment/refund", h.handleRefund)
	router.POST("/api/webhook", h.handleWebhook)
}

func (h *Handler) handleHealth(c *gin.Context) {
	resp := map[string]interface{}{
		"status":    "ok",
		"providers": h.registry.Names(),
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) handleCreatePayment(c *gin.Context) {
	var req models.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}
	b, ok := h.registry.Get(req.Provider)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}

	resp, err := b.CreatePayment(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("failed to create payment", "provider", req.Provider, "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) handlePaymentStatus(c *gin.Context) {
	provider := c.Query("provider")
	paymentID := c.Query("payment_id")
	if provider == "" || paymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider and payment_id are required"})
		return
	}
	b, ok := h.registry.Get(provider)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}

	resp, err := b.GetStatus(c.Request.Context(), paymentID)
	if err != nil {
		h.logger.Warn("failed to get payment status", "provider", provider, "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) handleRefund(c *gin.Context) {
	var req models.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Provider == "" || req.PaymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider and payment_id are required"})
		return
	}
	b, ok := h.registry.Get(req.Provider)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}

	resp, err := b.Refund(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("failed to refund payment", "provider", req.Provider, "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) handleWebhook(c *gin.Context) {
	provider := c.Query("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}
	b, ok := h.registry.Get(provider)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown provider"})
		return
	}

	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := b.HandleWebhook(c.Request.Context(), payload, c.Request.Header)
	if err != nil {
		h.logger.Warn("failed to process webhook", "provider", provider, "error", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "provider error"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
