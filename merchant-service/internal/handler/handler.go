package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
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
	adminRoute.GET("/merchants", h.GetMerchants)
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

func (h *MerchantHandler) GetMerchants(c *gin.Context) {
	ctx := c.Request.Context()

	params := db.GetMerchantListParams{
		Limit:  5,
		Offset: 0,
	}

	if limit := c.Query("limit"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid limit",
			})
			return
		}
		params.Limit = int32(limitInt)
	}

	if offset := c.Query("offset"); offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid offset",
			})
			return
		}
		params.Offset = int32(offsetInt)
	}

	// For merchantID - need to handle non-existent case
	if merchantID := c.Query("merchantID"); merchantID != "" {
		merchantIDInt, err := strconv.Atoi(merchantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid merchantID"})
			return
		}
		params.MerchantID = int64(merchantIDInt)
		// Your service should return empty array if not found (200 status)
	}

	// For name - add wildcard wrapping
	if name := c.Query("name"); name != "" {
		params.Name = name
	}

	// For merchantCategory - validate enum
	if merchantCategory := c.Query("merchantCategory"); merchantCategory != "" {
		validCategories := []string{"SmallRestaurant", "MediumRestaurant", "LargeRestaurant", "MerchandiseRestaurant", "BoothKiosk", "ConvenienceStore"}
		isValid := false
		for _, v := range validCategories {
			if merchantCategory == v {
				isValid = true
				break
			}
		}
		if isValid {
			params.MerchantCategory = merchantCategory
		}
	}

	// For createdAt - ignore if wrong
	if createdAtSort := c.Query("createdAt"); createdAtSort != "" {
		if createdAtSort == "asc" {
			params.CreatedAtSortAsc = true
		} else if createdAtSort == "desc" {
			params.CreatedAtSortDesc = true
		}
	}

	merchants, err := h.merchantService.GetMerchants(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get merchants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"merchants": merchants,
	})
}
