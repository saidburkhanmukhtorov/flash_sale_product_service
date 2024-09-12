package service

import (
	"context"
	"fmt"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/storage"
)

// FlashSaleEventProductService implements the product_service.FlashSaleEventProductServiceServer interface.
type FlashSaleEventProductService struct {
	storage storage.StorageI
	product_service.UnimplementedFlashSaleEventProductServiceServer
}

// NewFlashSaleEventProductService creates a new FlashSaleEventProductService instance.
func NewFlashSaleEventProductService(storage storage.StorageI) *FlashSaleEventProductService {
	return &FlashSaleEventProductService{
		storage: storage,
	}
}

// CreateFlashSaleEventProduct creates a new flash sale event product.
func (s *FlashSaleEventProductService) CreateFlashSaleEventProduct(ctx context.Context, req *product_service.CreateFlashSaleEventProductRequest) (*product_service.CreateFlashSaleEventProductResponse, error) {
	flashSaleEventProduct, err := s.storage.FlashSaleEventProduct().CreateFlashSaleEventProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create flash sale event product: %w", err)
	}

	return &product_service.CreateFlashSaleEventProductResponse{
		FlashSaleEventProduct: flashSaleEventProduct,
	}, nil
}

// GetFlashSaleEventProduct retrieves a flash sale event product by its ID.
func (s *FlashSaleEventProductService) GetFlashSaleEventProduct(ctx context.Context, req *product_service.GetFlashSaleEventProductRequest) (*product_service.GetFlashSaleEventProductResponse, error) {
	flashSaleEventProduct, err := s.storage.FlashSaleEventProduct().GetFlashSaleEventProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get flash sale event product: %w", err)
	}

	return &product_service.GetFlashSaleEventProductResponse{
		FlashSaleEventProduct: flashSaleEventProduct,
	}, nil
}

// UpdateFlashSaleEventProduct updates an existing flash sale event product.
func (s *FlashSaleEventProductService) UpdateFlashSaleEventProduct(ctx context.Context, req *product_service.UpdateFlashSaleEventProductRequest) (*product_service.UpdateFlashSaleEventProductResponse, error) {
	flashSaleEventProduct, err := s.storage.FlashSaleEventProduct().UpdateFlashSaleEventProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update flash sale event product: %w", err)
	}

	return &product_service.UpdateFlashSaleEventProductResponse{
		FlashSaleEventProduct: flashSaleEventProduct,
	}, nil
}

// DeleteFlashSaleEventProduct deletes a flash sale event product by its ID.
func (s *FlashSaleEventProductService) DeleteFlashSaleEventProduct(ctx context.Context, req *product_service.DeleteFlashSaleEventProductRequest) (*product_service.DeleteFlashSaleEventProductResponse, error) {
	response, err := s.storage.FlashSaleEventProduct().DeleteFlashSaleEventProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete flash sale event product: %w", err)
	}

	return response, nil
}

// ListFlashSaleEventProducts retrieves a list of flash sale event products.
func (s *FlashSaleEventProductService) ListFlashSaleEventProducts(ctx context.Context, req *product_service.ListFlashSaleEventProductsRequest) (*product_service.ListFlashSaleEventProductsResponse, error) {
	response, err := s.storage.FlashSaleEventProduct().ListFlashSaleEventProducts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list flash sale event products: %w", err)
	}

	return response, nil
}
