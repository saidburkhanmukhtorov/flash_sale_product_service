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

func TestProductDiscountRepo(t *testing.T) {
	db := createDBConnection(t)
	defer db.Close(context.Background())

	productRepo := postgres.NewProductRepo(db)
	discountRepo := postgres.NewDiscountuctRepo(db)
	productDiscountRepo := postgres.NewProductDiscountuctRepo(db)

	t.Run("CreateProductDiscount", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Discount",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a discount
		discount, err := discountRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Test Discount for Product",
				Description:   "Test Description",
				DiscountType:  "PERCENTAGE",
				DiscountValue: 10.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 1, 0)),
				IsActive:      true,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, discount)

		// Create a product discount
		req := &product_service.CreateProductDiscountRequest{
			ProductDiscount: &product_service.ProductDiscount{
				ProductId:  product.Id,
				DiscountId: discount.Id,
			},
		}

		productDiscount, err := productDiscountRepo.CreateProductDiscount(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, productDiscount)
		assert.NotEmpty(t, productDiscount.Id)

		// Cleanup
		deleteProductDiscount(t, db, productDiscount.Id, product.Id, discount.Id)
	})

	t.Run("GetProductDiscount", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Discount",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a discount
		discount, err := discountRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Test Discount for Product",
				Description:   "Test Description",
				DiscountType:  "PERCENTAGE",
				DiscountValue: 10.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 1, 0)),
				IsActive:      true,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, discount)

		// Create a product discount
		createdProductDiscount, err := productDiscountRepo.CreateProductDiscount(context.Background(), &product_service.CreateProductDiscountRequest{
			ProductDiscount: &product_service.ProductDiscount{
				ProductId:  product.Id,
				DiscountId: discount.Id,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdProductDiscount)

		// Get the product discount
		productDiscount, err := productDiscountRepo.GetProductDiscount(context.Background(), &product_service.GetProductDiscountRequest{Id: createdProductDiscount.Id})
		assert.NoError(t, err)
		assert.NotNil(t, productDiscount)
		assert.Equal(t, createdProductDiscount.Id, productDiscount.Id)

		// Cleanup
		deleteProductDiscount(t, db, createdProductDiscount.Id, product.Id, discount.Id)
	})
	t.Run("UpdateProductDiscount", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Discount",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a discount
		discount, err := discountRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Test Discount for Product",
				Description:   "Test Description",
				DiscountType:  "PERCENTAGE",
				DiscountValue: 10.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 1, 0)),
				IsActive:      true,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, discount)

		// Create another discount for updating
		discount2, err := discountRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Test Discount 2 for Product",
				Description:   "Test Description 2",
				DiscountType:  "FIXED_AMOUNT",
				DiscountValue: 5.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 2, 0)),
				IsActive:      true,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, discount2)

		// Create a product discount
		createdProductDiscount, err := productDiscountRepo.CreateProductDiscount(context.Background(), &product_service.CreateProductDiscountRequest{
			ProductDiscount: &product_service.ProductDiscount{
				ProductId:  product.Id,
				DiscountId: discount.Id,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdProductDiscount)

		// Update the product discount with the new discount ID
		createdProductDiscount.DiscountId = discount2.Id

		updatedProductDiscount, err := productDiscountRepo.UpdateProductDiscount(context.Background(), &product_service.UpdateProductDiscountRequest{
			ProductDiscount: createdProductDiscount,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedProductDiscount)
		assert.Equal(t, createdProductDiscount.DiscountId, updatedProductDiscount.DiscountId)

		// Cleanup
		deleteProductDiscount(t, db, createdProductDiscount.Id, product.Id, discount2.Id) // Delete with the updated discount ID
		deleteDiscount(t, db, discount.Id)                                                // Delete the first discount as well
	})
	t.Run("DeleteProductDiscount", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Discount",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a discount
		discount, err := discountRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Test Discount for Product",
				Description:   "Test Description",
				DiscountType:  "PERCENTAGE",
				DiscountValue: 10.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 1, 0)),
				IsActive:      true,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, discount)

		// Create a product discount
		createdProductDiscount, err := productDiscountRepo.CreateProductDiscount(context.Background(), &product_service.CreateProductDiscountRequest{
			ProductDiscount: &product_service.ProductDiscount{
				ProductId:  product.Id,
				DiscountId: discount.Id,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdProductDiscount)

		// Delete the product discount
		_, err = productDiscountRepo.DeleteProductDiscount(context.Background(), &product_service.DeleteProductDiscountRequest{Id: createdProductDiscount.Id})
		assert.NoError(t, err)

		// Try to get the deleted product discount
		_, err = productDiscountRepo.GetProductDiscount(context.Background(), &product_service.GetProductDiscountRequest{Id: createdProductDiscount.Id})
		assert.ErrorIs(t, err, pgx.ErrNoRows)

		// // Cleanup
		// defer deleteProduct(t, db, product.Id)
		// defer deleteDiscount(t, db, discount.Id)
	})

	t.Run("ListProductDiscounts", func(t *testing.T) {
		// Create a product first
		product, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Product for Discount",
				Description:   "Description",
				BasePrice:     100.0,
				CurrentPrice:  100.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 10,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, product)

		// Create a few test discounts
		discountsToCreate := []*product_service.Discount{
			{
				Name:          "List Discount 1",
				Description:   "List Discount Description 1",
				DiscountType:  "PERCENTAGE",
				DiscountValue: 10.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 1, 0)),
				IsActive:      true,
			},
			{
				Name:          "List Discount 2",
				Description:   "List Discount Description 2",
				DiscountType:  "FIXED_AMOUNT",
				DiscountValue: 5.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 2, 0)),
				IsActive:      true,
			},
		}

		// Create product discounts for each test discount
		for _, discount := range discountsToCreate {
			createdDiscount, err := discountRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{Discount: discount})
			assert.NoError(t, err)

			_, err = productDiscountRepo.CreateProductDiscount(context.Background(), &product_service.CreateProductDiscountRequest{
				ProductDiscount: &product_service.ProductDiscount{
					ProductId:  product.Id,
					DiscountId: createdDiscount.Id,
				},
			})
			assert.NoError(t, err)

			// defer deleteDiscount(t, db, createdDiscount.Id) // Cleanup discounts
		}

		// Test ListProductDiscounts
		productDiscounts, err := productDiscountRepo.ListProductDiscounts(context.Background(), &product_service.ListProductDiscountsRequest{
			Page:  1,
			Limit: 10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, productDiscounts)
		assert.GreaterOrEqual(t, len(productDiscounts.ProductDiscounts), 2)

		// Cleanup
		// defer deleteProduct(t, db, product.Id)
	})
}

func deleteProductDiscount(t *testing.T, db *pgx.Conn, productDiscountID, productId, discountId string) {
	_, err := db.Exec(context.Background(), "DELETE FROM product_discounts WHERE id = $1", productDiscountID)
	if err != nil {
		t.Fatalf("Error deleting product discount: %v", err)
	}
	deleteProduct(t, db, productId)
	deleteDiscount(t, db, discountId)
}
