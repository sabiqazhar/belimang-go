package service

import "github.com/sabiqazhar/belimang-go/merchant-service/internal/repositories"

type MerchantSvcImpl struct {
	repo repositories.MerchantRepository
}

func NewMerchantService(repo repositories.MerchantRepository) MerchantService {
	return &MerchantSvcImpl{
		repo: repo,
	}
}
