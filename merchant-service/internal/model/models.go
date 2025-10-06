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

type Meta struct {
	Limit  int `json:"limit" binding:"required,min=1,max=100"`
	Offset int `json:"offset" binding:"required,min=1,max=100"`
	Total  int `json:"total"`
}

type GetMerchantResponse struct {
	Merchants []MerchantListResponse `json:"data"`
	Meta      Meta                   `json:"meta"`
}

type MerchantListResponse struct {
	MerchantID       string   `json:"merchantId"`
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchantCategory"`
	ImageURL         string   `json:"imageUrl"`
	Location         Location `json:"location"`
	CreatedAt        string   `json:"createdAt"`
}

// Add Item Request and Response
type (
	AddItemRequest struct {
		Name            string `json:"name" binding:"required,min=1,max=30"`
		ProductCategory string `json:"productCategory" binding:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
		Price           int64  `json:"price" binding:"required,gt=1"`
		ImageURL        string `json:"imageUrl" binding:"required,url"`
	}
)

type GetItemResponse struct {
	ItemId          string `json:"itemId"`
	Name            string `json:"name"`
	Price           int64  `json:"price"`
	ImageURL        string `json:"imageUrl"`
	CreatedAt       string `json:"createdAt"`
	ProductCategory string `json:"productCategory"`
}
