package service

import (
	"context"
	"fmt"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/storage"
)

// ProductService implements the product_service.ProductServiceServer interface.
type ProductService struct {
	storage storage.StorageI
	product_service.UnimplementedProductServiceServer
}

// NewProductService creates a new ProductService instance.
func NewProductService(storage storage.StorageI) *ProductService {
	return &ProductService{
		storage: storage,
	}
}

// CreateProduct creates a new product.
func (s *ProductService) CreateProduct(ctx context.Context, req *product_service.CreateProductRequest) (*product_service.CreateProductResponse, error) {
	product, err := s.storage.Product().CreateProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &product_service.CreateProductResponse{
		Product: product,
	}, nil
}

// GetProduct retrieves a product by its ID.
func (s *ProductService) GetProduct(ctx context.Context, req *product_service.GetProductRequest) (*product_service.GetProductResponse, error) {
	product, err := s.storage.Product().GetProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product_service.GetProductResponse{
		Product: product,
	}, nil
}

// UpdateProduct updates an existing product.
func (s *ProductService) UpdateProduct(ctx context.Context, req *product_service.UpdateProductRequest) (*product_service.UpdateProductResponse, error) {
	product, err := s.storage.Product().UpdateProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &product_service.UpdateProductResponse{
		Product: product,
	}, nil
}

// DeleteProduct deletes a product by its ID.
func (s *ProductService) DeleteProduct(ctx context.Context, req *product_service.DeleteProductRequest) (*product_service.DeleteProductResponse, error) {
	response, err := s.storage.Product().DeleteProduct(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}

	return response, nil
}

// ListProducts retrieves a list of products.
func (s *ProductService) ListProducts(ctx context.Context, req *product_service.ListProductsRequest) (*product_service.ListProductsResponse, error) {
	response, err := s.storage.Product().ListProducts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return response, nil
}
