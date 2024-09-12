package test

import (
	"context"
	"log"
	"testing"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/storage/postgres"
	"github.com/jackc/pgx/v5"

	"github.com/stretchr/testify/assert"
)

func TestProductRepo(t *testing.T) {
	db := createDBConnection(t)
	defer db.Close(context.Background())

	productRepo := postgres.NewProductRepo(db)

	t.Run("CreateProduct", func(t *testing.T) {
		req := &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Test Product",
				Description:   "Test Description",
				BasePrice:     10.0,
				CurrentPrice:  10.0,
				ImageUrl:      "https://example.com/image.jpg",
				StockQuantity: 100,
			},
		}

		product, err := productRepo.CreateProduct(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.NotEmpty(t, product.Id)

		// Cleanup
		defer deleteProduct(t, db, product.Id)
	})

	t.Run("GetProduct", func(t *testing.T) {
		// Create a product first
		createdProduct, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Get Product",
				Description:   "Get Product Description",
				BasePrice:     12.0,
				CurrentPrice:  12.0,
				ImageUrl:      "https://example.com/get_image.jpg",
				StockQuantity: 50,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdProduct)

		// Get the product
		product, err := productRepo.GetProduct(context.Background(), &product_service.GetProductRequest{Id: createdProduct.Id})
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, createdProduct.Id, product.Id)

		// Cleanup
		defer deleteProduct(t, db, createdProduct.Id)
	})

	t.Run("ListProducts", func(t *testing.T) {
		// Create a few test products
		productsToCreate := []*product_service.Product{
			{
				Name:          "List Product 1",
				Description:   "List Product Description 1",
				BasePrice:     15.0,
				CurrentPrice:  15.0,
				ImageUrl:      "https://example.com/list_image1.jpg",
				StockQuantity: 75,
			},
			{
				Name:          "List Product 2",
				Description:   "List Product Description 2",
				BasePrice:     20.0,
				CurrentPrice:  20.0,
				ImageUrl:      "https://example.com/list_image2.jpg",
				StockQuantity: 25,
			},
		}

		for _, product := range productsToCreate {
			_, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{Product: product})
			assert.NoError(t, err)
			defer deleteProduct(t, db, product.Id)
		}

		// Test ListProducts with no filters
		products, err := productRepo.ListProducts(context.Background(), &product_service.ListProductsRequest{
			Page:  1,
			Limit: 10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.GreaterOrEqual(t, len(products.Products), 2)

		// Test ListProducts with name filter
		products, err = productRepo.ListProducts(context.Background(), &product_service.ListProductsRequest{
			Page:  1,
			Limit: 10,
			Name:  "List Product 1",
		})
		log.Println(products)
		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Equal(t, 1, len(products.Products))

		// Test ListProducts with price filter
		products, err = productRepo.ListProducts(context.Background(), &product_service.ListProductsRequest{
			Page:     1,
			Limit:    10,
			MinPrice: 16.0,
			MaxPrice: 25.0,
		})
		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Equal(t, 1, len(products.Products))
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		// Create a product first
		createdProduct, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Update Product",
				Description:   "Update Product Description",
				BasePrice:     25.0,
				CurrentPrice:  25.0,
				ImageUrl:      "https://example.com/update_image.jpg",
				StockQuantity: 80,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdProduct)

		// Update the product
		createdProduct.Name = "Updated Product"
		createdProduct.CurrentPrice = 30.0

		updatedProduct, err := productRepo.UpdateProduct(context.Background(), &product_service.UpdateProductRequest{
			Product: createdProduct,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedProduct)
		assert.Equal(t, "Updated Product", updatedProduct.Name)
		assert.Equal(t, float32(30.0), updatedProduct.CurrentPrice)

		// Cleanup
		defer deleteProduct(t, db, createdProduct.Id)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		// Create a product first
		createdProduct, err := productRepo.CreateProduct(context.Background(), &product_service.CreateProductRequest{
			Product: &product_service.Product{
				Name:          "Delete Product",
				Description:   "Delete Product Description",
				BasePrice:     30.0,
				CurrentPrice:  30.0,
				ImageUrl:      "https://example.com/delete_image.jpg",
				StockQuantity: 90,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdProduct)

		// Delete the product
		_, err = productRepo.DeleteProduct(context.Background(), &product_service.DeleteProductRequest{Id: createdProduct.Id})
		assert.NoError(t, err)

		// Try to get the deleted product
		_, err = productRepo.GetProduct(context.Background(), &product_service.GetProductRequest{Id: createdProduct.Id})
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})
}

func deleteProduct(t *testing.T, db *pgx.Conn, productID string) {
	_, err := db.Exec(context.Background(), "DELETE FROM products WHERE id = $1", productID)
	if err != nil {
		t.Fatalf("Error deleting product: %v", err)
	}
}
