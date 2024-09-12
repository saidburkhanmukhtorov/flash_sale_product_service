package test

import (
	"context"
	"testing"
	"time"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/storage/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFlashSaleEventProductRepo(t *testing.T) {
	db := createDBConnection(t)
	defer db.Close(context.Background())

	productRepo := postgres.NewProductRepo(db)
	flashSaleEventRepo := postgres.NewFlashSaleRepo(db)
	flashSaleEventProductRepo := postgres.NewFlashSaleEventProductRepo(db)

	t.Run("CreateFlashSaleEventProduct", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Flash Sale",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a flash sale event
		flashSaleEvent, err := flashSaleEventRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Test Flash Sale Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 1, 0)),
				Status:      "UPCOMING",
				EventType:   "FLASH_SALE",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEvent)

		// Create a flash sale event product
		req := &product_service.CreateFlashSaleEventProductRequest{
			FlashSaleEventProduct: &product_service.FlashSaleEventProduct{
				EventId:            flashSaleEvent.Id,
				ProductId:          product.Id,
				DiscountPercentage: 10.0,
				SalePrice:          90.0,
				AvailableQuantity:  5,
				OriginalStock:      5,
			},
		}

		flashSaleEventProduct, err := flashSaleEventProductRepo.CreateFlashSaleEventProduct(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEventProduct)
		assert.NotEmpty(t, flashSaleEventProduct.Id)

		// Cleanup
		deleteFlashSaleEventProduct(t, db, flashSaleEventProduct.Id)
		deleteProduct(t, db, product.Id)
		deleteFlashSaleEvent(t, db, flashSaleEvent.Id)
	})

	t.Run("GetFlashSaleEventProduct", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Flash Sale",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a flash sale event
		flashSaleEvent, err := flashSaleEventRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Test Flash Sale Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 1, 0)),
				Status:      "UPCOMING",
				EventType:   "FLASH_SALE",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEvent)

		// Create a flash sale event product
		createdFlashSaleEventProduct, err := flashSaleEventProductRepo.CreateFlashSaleEventProduct(context.Background(), &product_service.CreateFlashSaleEventProductRequest{
			FlashSaleEventProduct: &product_service.FlashSaleEventProduct{
				EventId:            flashSaleEvent.Id,
				ProductId:          product.Id,
				DiscountPercentage: 10.0,
				SalePrice:          90.0,
				AvailableQuantity:  5,
				OriginalStock:      5,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdFlashSaleEventProduct)

		// Get the flash sale event product
		flashSaleEventProduct, err := flashSaleEventProductRepo.GetFlashSaleEventProduct(context.Background(), &product_service.GetFlashSaleEventProductRequest{Id: createdFlashSaleEventProduct.Id})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEventProduct)
		assert.Equal(t, createdFlashSaleEventProduct.Id, flashSaleEventProduct.Id)

		// Cleanup
		deleteFlashSaleEventProduct(t, db, createdFlashSaleEventProduct.Id)
		deleteProduct(t, db, product.Id)
		deleteFlashSaleEvent(t, db, flashSaleEvent.Id)
	})

	t.Run("UpdateFlashSaleEventProduct", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Flash Sale",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a flash sale event
		flashSaleEvent, err := flashSaleEventRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Test Flash Sale Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 1, 0)),
				Status:      "UPCOMING",
				EventType:   "FLASH_SALE",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEvent)

		// Create a flash sale event product
		createdFlashSaleEventProduct, err := flashSaleEventProductRepo.CreateFlashSaleEventProduct(context.Background(), &product_service.CreateFlashSaleEventProductRequest{
			FlashSaleEventProduct: &product_service.FlashSaleEventProduct{
				EventId:            flashSaleEvent.Id,
				ProductId:          product.Id,
				DiscountPercentage: 10.0,
				SalePrice:          90.0,
				AvailableQuantity:  5,
				OriginalStock:      5,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdFlashSaleEventProduct)

		// Update the flash sale event product
		createdFlashSaleEventProduct.DiscountPercentage = 20.0
		createdFlashSaleEventProduct.SalePrice = 80.0

		updatedFlashSaleEventProduct, err := flashSaleEventProductRepo.UpdateFlashSaleEventProduct(context.Background(), &product_service.UpdateFlashSaleEventProductRequest{
			FlashSaleEventProduct: createdFlashSaleEventProduct,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedFlashSaleEventProduct)
		assert.Equal(t, float32(20.0), updatedFlashSaleEventProduct.DiscountPercentage)
		assert.Equal(t, float32(80.0), updatedFlashSaleEventProduct.SalePrice)

		// Cleanup
		deleteFlashSaleEventProduct(t, db, createdFlashSaleEventProduct.Id)
		deleteProduct(t, db, product.Id)
		deleteFlashSaleEvent(t, db, flashSaleEvent.Id)
	})

	t.Run("DeleteFlashSaleEventProduct", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Flash Sale",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a flash sale event
		flashSaleEvent, err := flashSaleEventRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Test Flash Sale Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 1, 0)),
				Status:      "UPCOMING",
				EventType:   "FLASH_SALE",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEvent)

		// Create a flash sale event product
		createdFlashSaleEventProduct, err := flashSaleEventProductRepo.CreateFlashSaleEventProduct(context.Background(), &product_service.CreateFlashSaleEventProductRequest{
			FlashSaleEventProduct: &product_service.FlashSaleEventProduct{
				EventId:            flashSaleEvent.Id,
				ProductId:          product.Id,
				DiscountPercentage: 10.0,
				SalePrice:          90.0,
				AvailableQuantity:  5,
				OriginalStock:      5,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdFlashSaleEventProduct)

		// Delete the flash sale event product
		_, err = flashSaleEventProductRepo.DeleteFlashSaleEventProduct(context.Background(), &product_service.DeleteFlashSaleEventProductRequest{Id: createdFlashSaleEventProduct.Id})
		assert.NoError(t, err)

		// Try to get the deleted flash sale event product
		_, err = flashSaleEventProductRepo.GetFlashSaleEventProduct(context.Background(), &product_service.GetFlashSaleEventProductRequest{Id: createdFlashSaleEventProduct.Id})
		assert.ErrorIs(t, err, pgx.ErrNoRows)

		// 	// Cleanup
		// 	defer deleteProduct(t, db, product.Id)
		// 	deleteFlashSaleEvent(t, db, flashSaleEvent.Id)
	})

	t.Run("ListFlashSaleEventProducts", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Flash Sale",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a flash sale event
		flashSaleEvent, err := flashSaleEventRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Test Flash Sale Event",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 1, 0)),
				Status:      "UPCOMING",
				EventType:   "FLASH_SALE",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEvent)

		// Create a few test flash sale event products
		fsepToCreate := []*product_service.FlashSaleEventProduct{
			{
				EventId:            flashSaleEvent.Id,
				ProductId:          product.Id,
				DiscountPercentage: 10.0,
				SalePrice:          90.0,
				AvailableQuantity:  5,
				OriginalStock:      5,
			},
			{
				EventId:            flashSaleEvent.Id,
				ProductId:          product.Id,
				DiscountPercentage: 20.0,
				SalePrice:          80.0,
				AvailableQuantity:  3,
				OriginalStock:      3,
			},
		}

		for _, fsep := range fsepToCreate {
			_, err := flashSaleEventProductRepo.CreateFlashSaleEventProduct(context.Background(), &product_service.CreateFlashSaleEventProductRequest{FlashSaleEventProduct: fsep})
			assert.NoError(t, err)
			defer deleteFlashSaleEventProduct(t, db, fsep.Id) // Cleanup flash sale event products
		}

		// Test ListFlashSaleEventProducts
		flashSaleEventProducts, err := flashSaleEventProductRepo.ListFlashSaleEventProducts(context.Background(), &product_service.ListFlashSaleEventProductsRequest{
			Page:  1,
			Limit: 10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEventProducts)
		assert.GreaterOrEqual(t, len(flashSaleEventProducts.FlashSaleEventProducts), 2)

		// // Cleanup
		// deleteProduct(t, db, product.Id)
		// deleteFlashSaleEvent(t, db, flashSaleEvent.Id)
	})
}

func deleteFlashSaleEventProduct(t *testing.T, db *pgx.Conn, flashSaleEventProductID string) {
	_, err := db.Exec(context.Background(), "DELETE FROM flash_sale_event_products WHERE id = $1", flashSaleEventProductID)
	if err != nil {
		t.Fatalf("Error deleting flash sale event product: %v", err)
	}

}
