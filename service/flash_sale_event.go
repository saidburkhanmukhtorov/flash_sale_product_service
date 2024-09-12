package service

import (
	"context"
	"fmt"
	"log"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/storage"
	"github.com/flash_sale/flash_sale_product_service/storage/redis"
)

// FlashSaleEventService implements the product_service.FlashSaleEventServiceServer interface.
type FlashSaleEventService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	product_service.UnimplementedFlashSaleEventServiceServer
}

// NewFlashSaleEventService creates a new FlashSaleEventService instance.
func NewFlashSaleEventService(storage storage.StorageI, redisClient *redis.Client) *FlashSaleEventService {
	return &FlashSaleEventService{
		storage:     storage,
		redisClient: redisClient,
	}
}

// CreateFlashSaleEvent creates a new flash sale event.
func (s *FlashSaleEventService) CreateFlashSaleEvent(ctx context.Context, req *product_service.CreateFlashSaleEventRequest) (*product_service.CreateFlashSaleEventResponse, error) {
	flashSaleEvent, err := s.storage.FlashSaleEvent().CreateFlashSaleEvent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create flash sale event: %w", err)
	}

	// Broadcast flash sale event creation notification
	notificationMessage := fmt.Sprintf(
		"New Flash Sale Event: %s! Starts on %s, ends on %s. Check it out!",
		flashSaleEvent.Name,
		flashSaleEvent.StartTime.AsTime().Format("2006-01-02 15:04:05"),
		flashSaleEvent.EndTime.AsTime().Format("2006-01-02 15:04:05"),
	)

	// Send notification for creation
	if err := s.redisClient.AddNotification(ctx, "brodacast", notificationMessage); err != nil {
		log.Printf("failed to send notification: %v", err)
	}

	return &product_service.CreateFlashSaleEventResponse{
		FlashSaleEvent: flashSaleEvent,
	}, nil
}

// GetFlashSaleEvent retrieves a flash sale event by its ID.
func (s *FlashSaleEventService) GetFlashSaleEvent(ctx context.Context, req *product_service.GetFlashSaleEventRequest) (*product_service.GetFlashSaleEventResponse, error) {
	flashSaleEvent, err := s.storage.FlashSaleEvent().GetFlashSaleEvent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get flash sale event: %w", err)
	}

	return &product_service.GetFlashSaleEventResponse{
		FlashSaleEvent: flashSaleEvent,
	}, nil
}

// UpdateFlashSaleEvent updates an existing flash sale event.
func (s *FlashSaleEventService) UpdateFlashSaleEvent(ctx context.Context, req *product_service.UpdateFlashSaleEventRequest) (*product_service.UpdateFlashSaleEventResponse, error) {
	flashSaleEvent, err := s.storage.FlashSaleEvent().UpdateFlashSaleEvent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update flash sale event: %w", err)
	}

	return &product_service.UpdateFlashSaleEventResponse{
		FlashSaleEvent: flashSaleEvent,
	}, nil
}

// DeleteFlashSaleEvent deletes a flash sale event by its ID.
func (s *FlashSaleEventService) DeleteFlashSaleEvent(ctx context.Context, req *product_service.DeleteFlashSaleEventRequest) (*product_service.DeleteFlashSaleEventResponse, error) {
	response, err := s.storage.FlashSaleEvent().DeleteFlashSaleEvent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete flash sale event: %w", err)
	}

	return response, nil
}

// ListFlashSaleEvents retrieves a list of flash sale events.
func (s *FlashSaleEventService) ListFlashSaleEvents(ctx context.Context, req *product_service.ListFlashSaleEventsRequest) (*product_service.ListFlashSaleEventsResponse, error) {
	response, err := s.storage.FlashSaleEvent().ListFlashSaleEvents(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list flash sale events: %w", err)
	}

	return response, nil
}
