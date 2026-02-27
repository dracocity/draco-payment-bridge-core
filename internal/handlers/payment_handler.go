package handlers

import (
	"errors"
	"net/http"

	"github.com/dracocity/draco-payment-bridge-core/internal/models"
	"github.com/dracocity/draco-payment-bridge-core/internal/pg"
	"github.com/dracocity/draco-payment-bridge-core/internal/services"
	"github.com/gin-gonic/gin"
)

// PaymentHandler provides HTTP handlers for the payment API.
type PaymentHandler struct {
	paymentService *services.PaymentService
}

func New(pgRegistry *pg.Registry) *PaymentHandler {
	return &PaymentHandler{
		paymentService: services.NewPaymentService(pgRegistry),
	}
}

func (h *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/api/v1/health", h.handleHealth)
	router.POST("/api/v1/providers/:provider/payment-link", h.handleCreatePaymentLink)
	router.POST("/api/v1/providers/:provider/payments", h.handleCreatePayment)
	router.GET("/api/v1/providers/:provider/payments/:paymentId", h.handleGetPayment)
	router.POST("/api/v1/providers/:provider/payments/refund", h.handleRefund)
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

func (h *PaymentHandler) handleRefund(c *gin.Context) {
	var req models.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.Provider == "" || req.PaymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider and payment_id are required"})
		return
	}

	resp, err := h.paymentService.Refund(c.Request.Context(), req)
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
