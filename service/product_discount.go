package service

import (
	"context"
	"fmt"
	"log"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/storage"
)

// ProductDiscountService implements the product_service.ProductDiscountServiceServer interface.
type ProductDiscountService struct {
	storage storage.StorageI
	product_service.UnimplementedProductDiscountServiceServer
}

// NewProductDiscountService creates a new ProductDiscountService instance.
func NewProductDiscountService(storage storage.StorageI) *ProductDiscountService {
	return &ProductDiscountService{
		storage: storage,
	}
}

// CreateProductDiscount creates a new product discount.
func (s *ProductDiscountService) CreateProductDiscount(ctx context.Context, req *product_service.CreateProductDiscountRequest) (*product_service.CreateProductDiscountResponse, error) {
	productDiscount, err := s.storage.ProductDiscount().CreateProductDiscount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product discount: %w", err)
	}

	// Get product and discount details for notification
	product, err := s.storage.Product().GetProduct(ctx, &product_service.GetProductRequest{Id: productDiscount.ProductId})
	if err != nil {
		return nil, fmt.Errorf("failed to get product for notification: %w", err)
	}

	discount, err := s.storage.Discount().GetDiscount(ctx, &product_service.GetDiscountRequest{Id: productDiscount.DiscountId})
	if err != nil {
		return nil, fmt.Errorf("failed to get discount for notification: %w", err)
	}

	// Broadcast product discount creation notification
	notificationMessage := fmt.Sprintf(
		"Discount Alert! %s is now on sale with %s off! Original price: $%.2f, Sale price: $%.2f",
		product.Name,
		getDiscountDescription(discount),
		product.BasePrice,
		calculateDiscountedPrice(product.BasePrice, discount),
	)
	if err := s.storage.SendNotification().SendNotification(ctx, notificationMessage); err != nil {
		log.Printf("failed to send notification: %v", err)
	}

	return &product_service.CreateProductDiscountResponse{
		ProductDiscount: productDiscount,
	}, nil
}

// GetProductDiscount retrieves a product discount by its ID.
func (s *ProductDiscountService) GetProductDiscount(ctx context.Context, req *product_service.GetProductDiscountRequest) (*product_service.GetProductDiscountResponse, error) {
	productDiscount, err := s.storage.ProductDiscount().GetProductDiscount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get product discount: %w", err)
	}

	return &product_service.GetProductDiscountResponse{
		ProductDiscount: productDiscount,
	}, nil
}

// UpdateProductDiscount updates an existing product discount.
func (s *ProductDiscountService) UpdateProductDiscount(ctx context.Context, req *product_service.UpdateProductDiscountRequest) (*product_service.UpdateProductDiscountResponse, error) {
	productDiscount, err := s.storage.ProductDiscount().UpdateProductDiscount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product discount: %w", err)
	}

	return &product_service.UpdateProductDiscountResponse{
		ProductDiscount: productDiscount,
	}, nil
}

// DeleteProductDiscount deletes a product discount by its ID.
func (s *ProductDiscountService) DeleteProductDiscount(ctx context.Context, req *product_service.DeleteProductDiscountRequest) (*product_service.DeleteProductDiscountResponse, error) {
	response, err := s.storage.ProductDiscount().DeleteProductDiscount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete product discount: %w", err)
	}

	return response, nil
}

// ListProductDiscounts retrieves a list of product discounts.
func (s *ProductDiscountService) ListProductDiscounts(ctx context.Context, req *product_service.ListProductDiscountsRequest) (*product_service.ListProductDiscountsResponse, error) {
	response, err := s.storage.ProductDiscount().ListProductDiscounts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list product discounts: %w", err)
	}

	return response, nil
}

// Helper function to get a user-friendly discount description
func getDiscountDescription(discount *product_service.Discount) string {
	if discount.DiscountType == "PERCENTAGE" {
		return fmt.Sprintf("%.2f%%", discount.DiscountValue)
	} else if discount.DiscountType == "FIXED_AMOUNT" {
		return fmt.Sprintf("$%.2f", discount.DiscountValue)
	}
	return "Unknown discount type"
}

// Helper function to calculate the discounted price
func calculateDiscountedPrice(basePrice float32, discount *product_service.Discount) float32 {
	if discount.DiscountType == "PERCENTAGE" {
		return basePrice * (1 - discount.DiscountValue/100)
	} else if discount.DiscountType == "FIXED_AMOUNT" {
		return basePrice - discount.DiscountValue
	}
	return basePrice
}
