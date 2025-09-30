package handler

import (
	"github.com/gin-gonic/gin"
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
}
