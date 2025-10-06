package model

type LongLat struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"long" binding:"required"`
}

type (
	// Item item in order
	Item struct {
		ItemID   int64 `json:"itemId" binding:"required"`
		Quantity int64 `json:"quantity" binding:"required"`
	}

	// Order param
	Order struct {
		MerchantID      int64  `json:"merchantId" binding:"required"`
		IsStartingPoint bool   `json:"isStartingPoint" binding:"required"`
		Items           []Item `json:"items" binding:"required"`
	}

	// CreateEstimateRequest Post Estimate Request
	CreateEstimateRequest struct {
		UserLocation LongLat `json:"userLocation" binding:"required"`
		Orders       []Order `json:"orders" binding:"required"`
	}
)
