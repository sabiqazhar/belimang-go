package model

// Location represents the nested "location" object in the JSON
type Location struct {
	// binding:"required" ensures these fields are not zero
	Lat  float64 `json:"lat" binding:"required,latitude"`
	Long float64 `json:"long" binding:"required,longitude"`
}

// CreateMerchantRequest is the main request body with validation
type CreateMerchantRequest struct {
	// binding:"required,min=2,max=30" enforces length constraints
	Name string `json:"name" binding:"required,min=2,max=30"`

	// binding:"required,oneof=..." ensures the value is one of the enum options
	MerchantCategory string `json:"merchantCategory" binding:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`

	// binding:"required,url" checks for a valid URL format
	ImageURL string `json:"imageUrl" binding:"required,url"`

	// binding:"required" ensures the location object itself is present
	Location Location `json:"location" binding:"required"`
}
