package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sabiqazhar/belimang-go/order-service/internal/model"
	"github.com/sabiqazhar/belimang-go/order-service/internal/service"
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
	users.POST("/orders", h.CreateOrder)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()
	var req model.CreateOrderRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	//userID := c.GetString("userID")
	userID := 1
	err = h.service.UpdateOrderStatusOrdered(ctx, req.CalculatedEstimateId, int64(userID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orderId": req.CalculatedEstimateId})
}

func (h *Handler) CreateEstimate(c *gin.Context) {
	ctx := c.Request.Context()

	var request model.CreateEstimateRequest
	err := c.ShouldBindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	orderDetail, err := h.service.CreateOrder(ctx, request, 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusOK, orderDetail)
}
