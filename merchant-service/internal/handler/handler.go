package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/model"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/service"
)

type MerchantHandler struct {
	engine          *gin.Engine
	merchantService service.MerchantService
}

func NewMerchantHandler(engine *gin.Engine, merchantService service.MerchantService) *MerchantHandler {
	return &MerchantHandler{
		engine:          engine,
		merchantService: merchantService,
	}
}

func (h *MerchantHandler) RegisterRoutes() {
	h.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	adminRoute := h.engine.Group("/admin")
	adminRoute.POST("/merchants", h.CreateMerchant)
}

func (h *MerchantHandler) CreateMerchant(c *gin.Context) {
	var req model.CreateMerchantRequest
	ctx := c.Request.Context()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	merchantID, err := h.merchantService.CreateMerchant(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create merchant",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"merchantId": merchantID,
	})
}
