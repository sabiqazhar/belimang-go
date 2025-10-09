package model

type LongLat struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"long" binding:"required"`
}

type (
	// Item item in order
	Item struct {
		ItemID   string `json:"itemId" binding:"required"`
		Quantity int64  `json:"quantity" binding:"required,gt=0"`
	}

	// Order param
	Order struct {
		MerchantID      string `json:"merchantId" binding:"required"`
		IsStartingPoint *bool  `json:"isStartingPoint" binding:"required"`
		Items           []Item `json:"items" binding:"required,gt=0,dive"`
	}

	// CreateEstimateRequest Post Estimate Request
	CreateEstimateRequest struct {
		UserLocation LongLat `json:"userLocation" binding:"required"`
		Orders       []Order `json:"orders" binding:"required,gt=0,dive"`
	}

	CreateOrderResponse struct {
		CalculatedEstimateId           int64   `json:"calculatedEstimateId"`
		EstimatedDeliveryTimeInMinutes float64 `json:"estimatedDeliveryTimeInMinutes"`
		TotalPrice                     float64 `json:"totalPrice"`
	}
)

type (
	MerchantDetail struct {
		ID       string
		Location LongLat
	}
)

type (
	CreateOrderRequest struct {
		CalculatedEstimateId int64 `json:"calculatedEstimateId" binding:"required"`
	}
)
