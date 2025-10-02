package service

import (
	"context"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/model"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/repositories"
)

type MerchantSvcImpl struct {
	repo repositories.MerchantRepository
}

func NewMerchantService(repo repositories.MerchantRepository) MerchantService {
	return &MerchantSvcImpl{
		repo: repo,
	}
}

func (s *MerchantSvcImpl) CreateMerchant(ctx context.Context, req model.CreateMerchantRequest) error {

	_ := db.Merchants{
		Name:             req.Name,
		MerchantCategory: req.MerchantCategory,
		ImageUrl:         req.ImageURL,
		Longitude:        req.Location.Long,
		Latitude:         req.Location.Lat,
	}

	return nil
}
