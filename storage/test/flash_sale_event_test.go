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

func TestFlashSaleEventRepo(t *testing.T) {
	db := createDBConnection(t)
	defer db.Close(context.Background())

	productRepo := postgres.NewFlashSaleRepo(db)

	t.Run("CreateFlashSaleEvent", func(t *testing.T) {
		req := &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Test Flash Sale",
				Description: "Test Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 1, 0)), // End date 1 month from now
				Status:      "UPCOMING",
				EventType:   "FLASH_SALE",
			},
		}

		flashSaleEvent, err := productRepo.CreateFlashSaleEvent(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEvent)
		assert.NotEmpty(t, flashSaleEvent.Id)

		// Cleanup
		defer deleteFlashSaleEvent(t, db, flashSaleEvent.Id)
	})

	t.Run("GetFlashSaleEvent", func(t *testing.T) {
		// Create a flash sale event first
		createdEvent, err := productRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Get Flash Sale",
				Description: "Get Flash Sale Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 2, 0)), // End date 2 months from now
				Status:      "ACTIVE",
				EventType:   "PROMOTION",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdEvent)

		// Get the flash sale event
		flashSaleEvent, err := productRepo.GetFlashSaleEvent(context.Background(), &product_service.GetFlashSaleEventRequest{Id: createdEvent.Id})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEvent)
		assert.Equal(t, createdEvent.Id, flashSaleEvent.Id)

		// Cleanup
		defer deleteFlashSaleEvent(t, db, createdEvent.Id)
	})

	t.Run("UpdateFlashSaleEvent", func(t *testing.T) {
		// Create a flash sale event first
		createdEvent, err := productRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Update Flash Sale",
				Description: "Update Flash Sale Description",
				StartTime:   timestamppb.New(time.Now().AddDate(0, 1, 0)),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 3, 0)), // End date 3 months from now
				Status:      "UPCOMING",
				EventType:   "FLASH_SALE",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdEvent)

		// Update the flash sale event
		createdEvent.Name = "Updated Flash Sale"
		createdEvent.Status = "ACTIVE"

		updatedEvent, err := productRepo.UpdateFlashSaleEvent(context.Background(), &product_service.UpdateFlashSaleEventRequest{
			FlashSaleEvent: createdEvent,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedEvent)
		assert.Equal(t, "Updated Flash Sale", updatedEvent.Name)
		assert.Equal(t, "ACTIVE", updatedEvent.Status)

		// Cleanup
		defer deleteFlashSaleEvent(t, db, createdEvent.Id)
	})

	t.Run("DeleteFlashSaleEvent", func(t *testing.T) {
		// Create a flash sale event first
		createdEvent, err := productRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{
			FlashSaleEvent: &product_service.FlashSaleEvent{
				Name:        "Delete Flash Sale",
				Description: "Delete Flash Sale Description",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 4, 0)), // End date 4 months from now
				Status:      "ENDED",
				EventType:   "PROMOTION",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdEvent)

		// Delete the flash sale event
		_, err = productRepo.DeleteFlashSaleEvent(context.Background(), &product_service.DeleteFlashSaleEventRequest{Id: createdEvent.Id})
		assert.NoError(t, err)

		// Try to get the deleted flash sale event
		_, err = productRepo.GetFlashSaleEvent(context.Background(), &product_service.GetFlashSaleEventRequest{Id: createdEvent.Id})
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})

	t.Run("ListFlashSaleEvents", func(t *testing.T) {
		// Create a few test flash sale events
		eventsToCreate := []*product_service.FlashSaleEvent{
			{
				Name:        "List Flash Sale 1",
				Description: "List Flash Sale Description 1",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 1, 0)),
				Status:      "UPCOMING",
				EventType:   "FLASH_SALE",
			},
			{
				Name:        "List Flash Sale 2",
				Description: "List Flash Sale Description 2",
				StartTime:   timestamppb.Now(),
				EndTime:     timestamppb.New(time.Now().AddDate(0, 2, 0)),
				Status:      "ACTIVE",
				EventType:   "PROMOTION",
			},
		}

		for _, event := range eventsToCreate {
			_, err := productRepo.CreateFlashSaleEvent(context.Background(), &product_service.CreateFlashSaleEventRequest{FlashSaleEvent: event})
			assert.NoError(t, err)
			defer deleteFlashSaleEvent(t, db, event.Id)
		}

		// Test ListFlashSaleEvents
		flashSaleEvents, err := productRepo.ListFlashSaleEvents(context.Background(), &product_service.ListFlashSaleEventsRequest{
			Page:  1,
			Limit: 10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, flashSaleEvents)
		assert.GreaterOrEqual(t, len(flashSaleEvents.FlashSaleEvents), 2)
	})
}

func deleteFlashSaleEvent(t *testing.T, db *pgx.Conn, eventID string) {
	_, err := db.Exec(context.Background(), "DELETE FROM flash_sale_events WHERE id = $1", eventID)
	if err != nil {
		t.Fatalf("Error deleting flash sale event: %v", err)
	}
}
