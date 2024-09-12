package service

import (
	"context"
	"fmt"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/storage"
)

// DiscountService implements the product_service.DiscountServiceServer interface.
type DiscountService struct {
	storage storage.StorageI
	product_service.UnimplementedDiscountServiceServer
}

// NewDiscountService creates a new DiscountService instance.
func NewDiscountService(storage storage.StorageI) *DiscountService {
	return &DiscountService{
		storage: storage,
	}
}

// CreateDiscount creates a new discount.
func (s *DiscountService) CreateDiscount(ctx context.Context, req *product_service.CreateDiscountRequest) (*product_service.CreateDiscountResponse, error) {
	discount, err := s.storage.Discount().CreateDiscount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create discount: %w", err)
	}

	return &product_service.CreateDiscountResponse{
		Discount: discount,
	}, nil
}

// GetDiscount retrieves a discount by its ID.
func (s *DiscountService) GetDiscount(ctx context.Context, req *product_service.GetDiscountRequest) (*product_service.GetDiscountResponse, error) {
	discount, err := s.storage.Discount().GetDiscount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get discount: %w", err)
	}

	return &product_service.GetDiscountResponse{
		Discount: discount,
	}, nil
}

// UpdateDiscount updates an existing discount.
func (s *DiscountService) UpdateDiscount(ctx context.Context, req *product_service.UpdateDiscountRequest) (*product_service.UpdateDiscountResponse, error) {
	discount, err := s.storage.Discount().UpdateDiscount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update discount: %w", err)
	}

	return &product_service.UpdateDiscountResponse{
		Discount: discount,
	}, nil
}

// DeleteDiscount deletes a discount by its ID.
func (s *DiscountService) DeleteDiscount(ctx context.Context, req *product_service.DeleteDiscountRequest) (*product_service.DeleteDiscountResponse, error) {
	response, err := s.storage.Discount().DeleteDiscount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete discount: %w", err)
	}

	return response, nil
}

// ListDiscounts retrieves a list of discounts.
func (s *DiscountService) ListDiscounts(ctx context.Context, req *product_service.ListDiscountsRequest) (*product_service.ListDiscountsResponse, error) {
	response, err := s.storage.Discount().ListDiscounts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list discounts: %w", err)
	}

	return response, nil
}
