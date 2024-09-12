package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FlashSaleRepo struct {
	db *pgx.Conn
}

func NewFlashSaleRepo(db *pgx.Conn) *FlashSaleRepo {
	return &FlashSaleRepo{
		db: db,
	}
}
func (r *FlashSaleRepo) CreateFlashSaleEvent(ctx context.Context, req *product_service.CreateFlashSaleEventRequest) (*product_service.FlashSaleEvent, error) {
	if req.FlashSaleEvent.Id == "" {
		req.FlashSaleEvent.Id = uuid.NewString()
	}

	query := `
		INSERT INTO flash_sale_events (
			id,
			name,
			description,
			start_time,
			end_time,
			status,
			event_type,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0
		) RETURNING id, created_at, updated_at
	`

	flashSaleEventModel := makeFlashSaleEventModel(req.FlashSaleEvent)

	err := r.db.QueryRow(ctx, query,
		flashSaleEventModel.Id,
		flashSaleEventModel.Name,
		flashSaleEventModel.Description,
		flashSaleEventModel.StartTime,
		flashSaleEventModel.EndTime,
		flashSaleEventModel.Status,
		flashSaleEventModel.EventType,
	).Scan(&flashSaleEventModel.Id, &flashSaleEventModel.CreatedAt, &flashSaleEventModel.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return makeFlashSaleEventProto(flashSaleEventModel), nil
}

func (r *FlashSaleRepo) GetFlashSaleEvent(ctx context.Context, req *product_service.GetFlashSaleEventRequest) (*product_service.FlashSaleEvent, error) {
	var flashSaleEventModel models.FlashSaleEvent

	query := `
		SELECT 
			id,
			name,
			description,
			start_time,
			end_time,
			status,
			event_type,
			created_at,
			updated_at,
			deleted_at
		FROM flash_sale_events
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&flashSaleEventModel.Id,
		&flashSaleEventModel.Name,
		&flashSaleEventModel.Description,
		&flashSaleEventModel.StartTime,
		&flashSaleEventModel.EndTime,
		&flashSaleEventModel.Status,
		&flashSaleEventModel.EventType,
		&flashSaleEventModel.CreatedAt,
		&flashSaleEventModel.UpdatedAt,
		&flashSaleEventModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeFlashSaleEventProto(flashSaleEventModel), nil
}

func (r *FlashSaleRepo) UpdateFlashSaleEvent(ctx context.Context, req *product_service.UpdateFlashSaleEventRequest) (*product_service.FlashSaleEvent, error) {
	query := `
		UPDATE flash_sale_events
		SET 
			name = $1,
			description = $2,
			start_time = $3,
			end_time = $4,
			status = $5,
			event_type = $6,
			updated_at = NOW()
		WHERE id = $7 AND deleted_at = 0 AND start_time > NOW()
		RETURNING id, name, description, start_time, end_time, status, event_type, created_at, updated_at
	`

	flashSaleEventModel := makeFlashSaleEventModel(req.FlashSaleEvent)

	err := r.db.QueryRow(ctx, query,
		flashSaleEventModel.Name,
		flashSaleEventModel.Description,
		flashSaleEventModel.StartTime,
		flashSaleEventModel.EndTime,
		flashSaleEventModel.Status,
		flashSaleEventModel.EventType,
		flashSaleEventModel.Id,
	).Scan(
		&flashSaleEventModel.Id,
		&flashSaleEventModel.Name,
		&flashSaleEventModel.Description,
		&flashSaleEventModel.StartTime,
		&flashSaleEventModel.EndTime,
		&flashSaleEventModel.Status,
		&flashSaleEventModel.EventType,
		&flashSaleEventModel.CreatedAt,
		&flashSaleEventModel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeFlashSaleEventProto(flashSaleEventModel), nil
}

func (r *FlashSaleRepo) DeleteFlashSaleEvent(ctx context.Context, req *product_service.DeleteFlashSaleEventRequest) (*product_service.DeleteFlashSaleEventResponse, error) {
	query := `
		UPDATE flash_sale_events
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at = 0
	`

	_, err := r.db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		return nil, err
	}

	return &product_service.DeleteFlashSaleEventResponse{
		Message: "Flash sale event deleted successfully",
	}, nil
}

func (r *FlashSaleRepo) ListFlashSaleEvents(ctx context.Context, req *product_service.ListFlashSaleEventsRequest) (*product_service.ListFlashSaleEventsResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			name,
			description,
			start_time,
			end_time,
			status,
			event_type,
			created_at,
			updated_at,
			deleted_at
		FROM 
			flash_sale_events
		WHERE deleted_at = 0
	`

	filter := ""

	if req.Name != "" {
		filter += fmt.Sprintf(" AND name ILIKE $%d", count)
		args = append(args, "%"+req.Name+"%")
		count++
	}

	if req.Description != "" {
		filter += fmt.Sprintf(" AND description ILIKE $%d", count)
		args = append(args, "%"+req.Description+"%")
		count++
	}

	if req.StartTime != nil {
		filter += fmt.Sprintf(" AND start_time >= $%d", count)
		args = append(args, req.StartTime.AsTime())
		count++
	}

	if req.EndTime != nil {
		filter += fmt.Sprintf(" AND end_time <= $%d", count)
		args = append(args, req.EndTime.AsTime())
		count++
	}

	if req.Status != "" {
		filter += fmt.Sprintf(" AND status = $%d", count)
		args = append(args, req.Status)
		count++
	}

	if req.EventType != "" {
		filter += fmt.Sprintf(" AND event_type = $%d", count)
		args = append(args, req.EventType)
		count++
	}

	query += filter

	// Handle invalid page or limit values
	if req.Page <= 0 {
		req.Page = 1 // Default to page 1
	}
	if req.Limit <= 0 {
		req.Limit = 10 // Default to a limit of 10
	}

	totalCountQuery := "SELECT count(*) FROM flash_sale_events WHERE 1=1 AND deleted_at = 0" + filter
	var totalCount int
	err := r.db.QueryRow(ctx, totalCountQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Add LIMIT and OFFSET for pagination using the proto fields
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", count, count+1)
	args = append(args, req.Limit, (req.Page-1)*req.Limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flashSaleEventList []*product_service.FlashSaleEvent

	for rows.Next() {
		var flashSaleEventModel models.FlashSaleEvent
		err = rows.Scan(
			&flashSaleEventModel.Id,
			&flashSaleEventModel.Name,
			&flashSaleEventModel.Description,
			&flashSaleEventModel.StartTime,
			&flashSaleEventModel.EndTime,
			&flashSaleEventModel.Status,
			&flashSaleEventModel.EventType,
			&flashSaleEventModel.CreatedAt,
			&flashSaleEventModel.UpdatedAt,
			&flashSaleEventModel.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		flashSaleEventList = append(flashSaleEventList, makeFlashSaleEventProto(flashSaleEventModel))
	}

	return &product_service.ListFlashSaleEventsResponse{
		FlashSaleEvents: flashSaleEventList,
		Total:           int32(totalCount),
	}, nil
}

// Convert db model to proto model
func makeFlashSaleEventProto(event models.FlashSaleEvent) *product_service.FlashSaleEvent {
	return &product_service.FlashSaleEvent{
		Id:          event.Id,
		Name:        event.Name,
		Description: event.Description,
		StartTime:   timestamppb.New(event.StartTime),
		EndTime:     timestamppb.New(event.EndTime),
		Status:      event.Status,
		EventType:   event.EventType,
		CreatedAt:   timestamppb.New(event.CreatedAt),
		UpdatedAt:   timestamppb.New(event.UpdatedAt),
	}
}

// Convert proto model to db model
func makeFlashSaleEventModel(event *product_service.FlashSaleEvent) models.FlashSaleEvent {
	return models.FlashSaleEvent{
		Id:          event.Id,
		Name:        event.Name,
		Description: event.Description,
		StartTime:   event.StartTime.AsTime(),
		EndTime:     event.EndTime.AsTime(),
		Status:      event.Status,
		EventType:   event.EventType,
		CreatedAt:   event.CreatedAt.AsTime(),
		UpdatedAt:   event.UpdatedAt.AsTime(),
	}
}
