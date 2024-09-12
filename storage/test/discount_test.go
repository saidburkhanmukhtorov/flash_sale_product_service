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

func TestDiscountRepo(t *testing.T) {
	db := createDBConnection(t)
	defer db.Close(context.Background())

	productRepo := postgres.NewDiscountuctRepo(db)

	t.Run("CreateDiscount", func(t *testing.T) {
		req := &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Test Discount",
				Description:   "Test Description",
				DiscountType:  "PERCENTAGE",
				DiscountValue: 10.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 1, 0)), // End date 1 month from now
				IsActive:      true,
			},
		}

		discount, err := productRepo.CreateDiscount(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, discount)
		assert.NotEmpty(t, discount.Id)

		// Cleanup
		defer deleteDiscount(t, db, discount.Id)
	})

	t.Run("GetDiscount", func(t *testing.T) {
		// Create a discount first
		createdDiscount, err := productRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Get Discount",
				Description:   "Get Discount Description",
				DiscountType:  "FIXED_AMOUNT",
				DiscountValue: 5.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 2, 0)), // End date 2 months from now
				IsActive:      true,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdDiscount)

		// Get the discount
		discount, err := productRepo.GetDiscount(context.Background(), &product_service.GetDiscountRequest{Id: createdDiscount.Id})
		assert.NoError(t, err)
		assert.NotNil(t, discount)
		assert.Equal(t, createdDiscount.Id, discount.Id)

		// Cleanup
		defer deleteDiscount(t, db, createdDiscount.Id)
	})

	t.Run("UpdateDiscount", func(t *testing.T) {
		// Create a discount first
		createdDiscount, err := productRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Update Discount",
				Description:   "Update Discount Description",
				DiscountType:  "PERCENTAGE",
				DiscountValue: 15.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 3, 0)), // End date 3 months from now
				IsActive:      false,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdDiscount)

		// Update the discount
		createdDiscount.Name = "Updated Discount"
		createdDiscount.DiscountValue = 20.0
		createdDiscount.IsActive = true

		updatedDiscount, err := productRepo.UpdateDiscount(context.Background(), &product_service.UpdateDiscountRequest{
			Discount: createdDiscount,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedDiscount)
		assert.Equal(t, "Updated Discount", updatedDiscount.Name)
		assert.Equal(t, float32(20.0), updatedDiscount.DiscountValue)
		assert.True(t, updatedDiscount.IsActive)

		// Cleanup
		defer deleteDiscount(t, db, createdDiscount.Id)
	})

	t.Run("DeleteDiscount", func(t *testing.T) {
		// Create a discount first
		createdDiscount, err := productRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{
			Discount: &product_service.Discount{
				Name:          "Delete Discount",
				Description:   "Delete Discount Description",
				DiscountType:  "FIXED_AMOUNT",
				DiscountValue: 8.0,
				StartDate:     timestamppb.Now(),
				EndDate:       timestamppb.New(time.Now().AddDate(0, 4, 0)), // End date 4 months from now
				IsActive:      true,
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdDiscount)

		// Delete the discount
		_, err = productRepo.DeleteDiscount(context.Background(), &product_service.DeleteDiscountRequest{Id: createdDiscount.Id})
		assert.NoError(t, err)

		// Try to get the deleted discount
		_, err = productRepo.GetDiscount(context.Background(), &product_service.GetDiscountRequest{Id: createdDiscount.Id})
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})

	t.Run("ListDiscounts", func(t *testing.T) {
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

		for _, discount := range discountsToCreate {
			_, err := productRepo.CreateDiscount(context.Background(), &product_service.CreateDiscountRequest{Discount: discount})
			assert.NoError(t, err)
			defer deleteDiscount(t, db, discount.Id)
		}

		// Test ListDiscounts
		discounts, err := productRepo.ListDiscounts(context.Background(), &product_service.ListDiscountsRequest{
			Page:  1,
			Limit: 10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, discounts)
		assert.GreaterOrEqual(t, len(discounts.Discounts), 2)
	})
}

func deleteDiscount(t *testing.T, db *pgx.Conn, discountID string) {
	_, err := db.Exec(context.Background(), "DELETE FROM discounts WHERE id = $1", discountID)
	if err != nil {
		t.Fatalf("Error deleting discount: %v", err)
	}
}
