package handler

import (
	"context"

	"github.com/sabiqazhar/belimang-go/merchant-service/internal/service"
	pb "github.com/sabiqazhar/belimang-go/proto/merchant"
)

type MerchantGRPCHandler struct {
	pb.UnimplementedMerchantServiceServer
	merchantService service.MerchantService
}

func NewMerchantGRPCHandler(merchantService service.MerchantService) *MerchantGRPCHandler {
	return &MerchantGRPCHandler{
		merchantService: merchantService,
	}
}

func (h *MerchantGRPCHandler) GetMerchant(ctx context.Context, req *pb.GetMerchantRequest) (*pb.GetMerchantResponse, error) {
	merchant, err := h.merchantService.GetMerchantById(ctx, req.Id)
	if err != nil {
		return &pb.GetMerchantResponse{
			Error: true,
		}, nil
	}

	return &pb.GetMerchantResponse{
		Id:        merchant.ID,
		Name:      merchant.Name,
		Latitude:  merchant.Latitude,
		Longitude: merchant.Longitude,
		Error:     false,
	}, nil
}

func (h *MerchantGRPCHandler) GetItemsDetail(ctx context.Context, req *pb.GetItemsDetailRequest) (*pb.GetItemsDetailResponse, error) {
	itemDetail, err := h.merchantService.IsValidMerchantItem(ctx, req.MerchantId, req.ItemIds)
	if err != nil {
		return &pb.GetItemsDetailResponse{
			Error:        true,
			ErrorMessage: err.Error(),
		}, err
	}

	responseItems := make([]*pb.ItemDetail, 0, len(itemDetail))
	var errAppear bool

	for _, item := range itemDetail {
		priceAsFloat8, err := item.Price.Float64Value()
		if err != nil || !item.Price.Valid {
			errAppear = true
			break
		}

		responseItems = append(responseItems, &pb.ItemDetail{
			Id:    item.ID,
			Price: float32(priceAsFloat8.Float64),
			Name:  item.Name,
		})
	}
	return &pb.GetItemsDetailResponse{
		Items: responseItems,
		Error: errAppear,
	}, nil
}

func (h *MerchantGRPCHandler) GetItemDetail(ctx context.Context, req *pb.GetItemDetailRequest) (*pb.GetItemDetailResponse, error) {
	itemDetail, err := h.merchantService.GetItemByID(ctx, req.Id)
	if err != nil {
		return &pb.GetItemDetailResponse{
			Error: true,
		}, nil
	}

	priceAsFloat8, err := itemDetail.Price.Float64Value()
	if err != nil || !itemDetail.Price.Valid {
		return &pb.GetItemDetailResponse{
			Error: true,
		}, nil
	}

	return &pb.GetItemDetailResponse{
		Id:         itemDetail.ID,
		Name:       itemDetail.Name,
		Price:      priceAsFloat8.Float64,
		MerchantId: itemDetail.MerchantID.Int64,
		Error:      false,
	}, nil
}
