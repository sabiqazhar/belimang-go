package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sabiqazhar/belimang-go/order-service/internal/model"
	"github.com/sabiqazhar/belimang-go/order-service/internal/service"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
)

type Handler struct {
	engine  *gin.Engine
	service service.OrderService
}

func NewHandler(engine *gin.Engine, service service.OrderService) *Handler {
	return &Handler{
		engine:  engine,
		service: service,
	}
}

func (h *Handler) RegisterRoutes() {
	users := h.engine.Group("/users")
	users.POST("/estimate", h.CreateEstimate)
}

func (h *Handler) CreateEstimate(c *gin.Context) {
	ctx := c.Request.Context()

	var request model.CreateEstimateRequest
	err := c.ShouldBindJSON(&request)

	if err != nil {
		logger.Logger.Error().Err(err).Msg("error binding request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	orderID, err := h.service.CreateOrder(ctx, request, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"totalPrice":                     100.00,
		"estimatedDeliveryTimeInMinutes": 400,
		"calculatedEstimateId":           orderID,
	})
}
