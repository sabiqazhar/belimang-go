package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/model"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/service"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
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
	adminRoute.POST("/merchants/:merchantId/item", h.addItem) // just for testing purpose
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
		Limit:    5,
		Offset:   0,
		SortAsc:  false,
		SortDesc: false,
	}

	if limit := c.Query("limit"); limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid limit",
			})
			return
		}
		params.Limit = int32(limitInt)
	}

	if offset := c.Query("offset"); offset != "" {
		offsetInt, err := strconv.ParseInt(offset, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid offset",
			})
			return
		}
		params.Offset = int32(offsetInt)
	}

	// Optional filters
	if merchantID := c.Query("merchantId"); merchantID != "" {
		merchantIDInt, err := strconv.ParseInt(merchantID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid merchantID",
			})
			return
		}
		params.MerchantID = pgtype.Int8(sql.NullInt64{Int64: merchantIDInt, Valid: true})
	}

	if name := c.Query("name"); name != "" {
		// Add wildcards for LIKE search
		params.Name = pgtype.Text(sql.NullString{String: "%" + name + "%", Valid: true})
	}

	if category := c.Query("merchantCategory"); category != "" {
		params.MerchantCategory = pgtype.Text(sql.NullString{String: category, Valid: true})
	}

	if sort := c.Query("createdAt"); sort == "asc" {
		params.SortAsc = true
	} else if sort == "desc" {
		params.SortDesc = true
	}

	merchants, err := h.merchantService.GetMerchants(ctx, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get merchants",
		})
		return
	}

	var response model.GetMerchantResponse
	if merchants == nil {
		response.Merchants = []model.MerchantListResponse{}
	}

	for _, m := range merchants {
		response.Merchants = append(response.Merchants, model.MerchantListResponse{
			MerchantID:       strconv.FormatInt(m.ID, 10),
			Name:             m.Name,
			MerchantCategory: m.MerchantCategory,
			ImageURL:         m.ImageUrl,
			Location: model.Location{
				Lat:  m.Latitude,
				Long: m.Longitude,
			},
			CreatedAt: m.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	response.Meta = model.Meta{
		Limit:  int(params.Limit),
		Offset: int(params.Offset),
		Total:  len(response.Merchants),
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response.Merchants,
		"meta": response.Meta,
	})
}

func (h *MerchantHandler) addItem(c *gin.Context) {
	ctx := c.Request.Context()
	merchantIdStr := c.Param("merchantId")
	if merchantIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "merchantId is required",
		})
		return
	}

	merchantId, err := strconv.ParseInt(merchantIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid merchantId",
		})
		return
	}

	var req model.AddItemRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	itemId, err := h.merchantService.AddItem(ctx, req, merchantId)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to add item")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to add item",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"itemId": itemId,
	})

}

func (h *MerchantHandler) GetItems(c *gin.Context) {
	ctx := c.Request.Context()

	params := db.GetItemListParams{
		Limit:    5,
		Offset:   0,
		SortAsc:  false,
		SortDesc: false,
	}

	if limit := c.Query("limit"); limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid limit",
			})
			return
		}
		params.Limit = int32(limitInt)
	}

	if offset := c.Query("offset"); offset != "" {
		offsetInt, err := strconv.ParseInt(offset, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid offset",
			})
			return
		}
		params.Offset = int32(offsetInt)
	}

	if name := c.Query("name"); name != "" {
		params.Name = pgtype.Text(sql.NullString{String: "%" + name + "%", Valid: true})
	}

	if itemID := c.Query("itemId"); itemID != "" {
		itemIDInt, err := strconv.ParseInt(itemID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid itemID",
			})
			return
		}
		params.ItemId = pgtype.Int8(sql.NullInt64{Int64: itemIDInt, Valid: true})
	}

	if category := c.Query("productCategory"); category != "" {
		params.ProductCategory = pgtype.Text(sql.NullString{String: category, Valid: true})
	}

	if sort := c.Query("createdAt"); sort == "asc" {
		params.SortAsc = true
	} else if sort == "desc" {
		params.SortDesc = true
	}

	items, err := h.merchantService.GetItems(ctx, params)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to get items")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get items",
		})
		return
	}

	if items == nil {
		c.JSON(http.StatusOK, gin.H{
			"data": []model.GetItemResponse{},
			"meta": model.Meta{
				Limit:  int(params.Limit),
				Offset: int(params.Offset),
				Total:  0,
			},
		})
		return
	}

	var response []model.GetItemResponse
	for _, item := range items {
		response = append(response, model.GetItemResponse{
			ItemId:          strconv.FormatInt(item.ID, 10),
			Name:            item.Name,
			Price:           item.Price.Int.Int64(),
			ImageURL:        item.ImageUrl,
			CreatedAt:       item.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			ProductCategory: item.ProductCategory,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"meta": model.Meta{
			Limit:  int(params.Limit),
			Offset: int(params.Offset),
			Total:  len(response),
		},
	})
}
