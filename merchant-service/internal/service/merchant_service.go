package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/db"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/model"
	"github.com/sabiqazhar/belimang-go/merchant-service/internal/repositories"
	"github.com/sabiqazhar/belimang-go/pkg/logger"
	"github.com/uber/h3-go/v4"
)

type MerchantSvcImpl struct {
	repo repositories.MerchantRepository
}

func NewMerchantService(repo repositories.MerchantRepository) MerchantService {
	return &MerchantSvcImpl{
		repo: repo,
	}
}

func (s *MerchantSvcImpl) CreateMerchant(ctx context.Context, req model.CreateMerchantRequest) (int64, error) {
	h3Index, err := s.getH3Index(req.Location.Lat, req.Location.Long, 8)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to get h3 index")
		return 0, err
	}
	merchantParam := db.CreateMerchantParams{
		Name:             req.Name,
		MerchantCategory: req.MerchantCategory,
		ImageUrl:         req.ImageURL,
		Longitude:        req.Location.Long,
		Latitude:         req.Location.Lat,
		H3Index:          pgtype.Int8{Int64: h3Index, Valid: true},
	}

	merchant, err := s.repo.InsertMerchant(ctx, merchantParam)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to create merchant")
		return 0, err
	}

	return merchant.ID, nil
}

func (s *MerchantSvcImpl) getH3Index(lat, long float64, resolution int) (int64, error) {
	coordinate := h3.LatLng{Lat: lat, Lng: long}
	cell, err := h3.LatLngToCell(coordinate, resolution)
	if err != nil {
		return 0, err
	}
	return int64(cell), nil
}

func (s *MerchantSvcImpl) GetMerchants(ctx context.Context, param db.GetMerchantListParams) ([]db.GetMerchantListRow, error) {
	merchants, err := s.repo.GetMerchants(ctx, param)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("failed to get merchants")
		return []db.GetMerchantListRow{}, err
	}
	logger.Logger.Info().Interface("param", param).Msg("get merchants")
	logger.Logger.Info().Interface("merchants", merchants).Msg("merchants")
	return merchants, nil
}
