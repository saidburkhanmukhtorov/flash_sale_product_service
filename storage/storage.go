package storage

import (
	"context"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
)

// StorageI defines the interface for interacting with the storage layer.
type StorageI interface {
	Product() ProductI
	Discount() DiscountI
	FlashSaleEvent() FlashSaleEventI
	ProductDiscount() ProductDiscountI
	FlashSaleEventProduct() FlashSaleEventProductI
}

// ProductI defines methods for interacting with product data.
type ProductI interface {
	CreateProduct(ctx context.Context, req *product_service.CreateProductRequest) (*product_service.Product, error)
	GetProduct(ctx context.Context, req *product_service.GetProductRequest) (*product_service.Product, error)
	ListProducts(ctx context.Context, req *product_service.ListProductsRequest) (*product_service.ListProductsResponse, error)
	UpdateProduct(ctx context.Context, req *product_service.UpdateProductRequest) (*product_service.Product, error)
	DeleteProduct(ctx context.Context, req *product_service.DeleteProductRequest) (*product_service.DeleteProductResponse, error)
}

// DiscountI defines methods for interacting with discount data.
type DiscountI interface {
	CreateDiscount(ctx context.Context, req *product_service.CreateDiscountRequest) (*product_service.Discount, error)
	GetDiscount(ctx context.Context, req *product_service.GetDiscountRequest) (*product_service.Discount, error)
	UpdateDiscount(ctx context.Context, req *product_service.UpdateDiscountRequest) (*product_service.Discount, error)
	DeleteDiscount(ctx context.Context, req *product_service.DeleteDiscountRequest) (*product_service.DeleteDiscountResponse, error)
	ListDiscounts(ctx context.Context, req *product_service.ListDiscountsRequest) (*product_service.ListDiscountsResponse, error)
}

// FlashSaleEventI defines methods for interacting with flash sale event data.
type FlashSaleEventI interface {
	CreateFlashSaleEvent(ctx context.Context, req *product_service.CreateFlashSaleEventRequest) (*product_service.FlashSaleEvent, error)
	GetFlashSaleEvent(ctx context.Context, req *product_service.GetFlashSaleEventRequest) (*product_service.FlashSaleEvent, error)
	UpdateFlashSaleEvent(ctx context.Context, req *product_service.UpdateFlashSaleEventRequest) (*product_service.FlashSaleEvent, error)
	DeleteFlashSaleEvent(ctx context.Context, req *product_service.DeleteFlashSaleEventRequest) (*product_service.DeleteFlashSaleEventResponse, error)
	ListFlashSaleEvents(ctx context.Context, req *product_service.ListFlashSaleEventsRequest) (*product_service.ListFlashSaleEventsResponse, error)
}

// ProductDiscountI defines methods for interacting with product discount data.
type ProductDiscountI interface {
	CreateProductDiscount(ctx context.Context, req *product_service.CreateProductDiscountRequest) (*product_service.ProductDiscount, error)
	GetProductDiscount(ctx context.Context, req *product_service.GetProductDiscountRequest) (*product_service.ProductDiscount, error)
	UpdateProductDiscount(ctx context.Context, req *product_service.UpdateProductDiscountRequest) (*product_service.ProductDiscount, error)
	DeleteProductDiscount(ctx context.Context, req *product_service.DeleteProductDiscountRequest) (*product_service.DeleteProductDiscountResponse, error)
	ListProductDiscounts(ctx context.Context, req *product_service.ListProductDiscountsRequest) (*product_service.ListProductDiscountsResponse, error)
}

// FlashSaleEventProductI defines methods for interacting with flash sale event product data.
type FlashSaleEventProductI interface {
	CreateFlashSaleEventProduct(ctx context.Context, req *product_service.CreateFlashSaleEventProductRequest) (*product_service.FlashSaleEventProduct, error)
	GetFlashSaleEventProduct(ctx context.Context, req *product_service.GetFlashSaleEventProductRequest) (*product_service.FlashSaleEventProduct, error)
	UpdateFlashSaleEventProduct(ctx context.Context, req *product_service.UpdateFlashSaleEventProductRequest) (*product_service.FlashSaleEventProduct, error)
	DeleteFlashSaleEventProduct(ctx context.Context, req *product_service.DeleteFlashSaleEventProductRequest) (*product_service.DeleteFlashSaleEventProductResponse, error)
	ListFlashSaleEventProducts(ctx context.Context, req *product_service.ListFlashSaleEventProductsRequest) (*product_service.ListFlashSaleEventProductsResponse, error)
}
