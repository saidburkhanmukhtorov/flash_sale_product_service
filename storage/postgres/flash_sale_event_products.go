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

type FlashSaleEventProductRepo struct {
	db *pgx.Conn
}

func NewFlashSaleEventProductRepo(db *pgx.Conn) *FlashSaleEventProductRepo {
	return &FlashSaleEventProductRepo{
		db: db,
	}
}

func (r *FlashSaleEventProductRepo) CreateFlashSaleEventProduct(ctx context.Context, req *product_service.CreateFlashSaleEventProductRequest) (*product_service.FlashSaleEventProduct, error) {
	err := r.checkFlashSaleEventStatus(ctx, req.FlashSaleEventProduct.EventId)
	if err != nil {
		return nil, err
	}
	if req.FlashSaleEventProduct.Id == "" {
		req.FlashSaleEventProduct.Id = uuid.NewString()
	}

	query := `
		INSERT INTO flash_sale_event_products (
			id,
			event_id,
			product_id,
			discount_percentage,
			sale_price,
			available_quantity,
			original_stock,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0
		) RETURNING id, created_at, updated_at
	`

	flashSaleEventProductModel := makeFlashSaleEventProductModel(req.FlashSaleEventProduct)

	err = r.db.QueryRow(ctx, query,
		flashSaleEventProductModel.Id,
		flashSaleEventProductModel.EventId,
		flashSaleEventProductModel.ProductId,
		flashSaleEventProductModel.DiscountPercentage,
		flashSaleEventProductModel.SalePrice,
		flashSaleEventProductModel.AvailableQuantity,
		flashSaleEventProductModel.OriginalStock,
	).Scan(
		&flashSaleEventProductModel.Id,
		&flashSaleEventProductModel.CreatedAt,
		&flashSaleEventProductModel.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return makeFlashSaleEventProductProto(flashSaleEventProductModel), nil
}

func (r *FlashSaleEventProductRepo) GetFlashSaleEventProduct(ctx context.Context, req *product_service.GetFlashSaleEventProductRequest) (*product_service.FlashSaleEventProduct, error) {
	var flashSaleEventProductModel models.FlashSaleEventProduct

	query := `
		SELECT 
			id,
			event_id,
			product_id,
			discount_percentage,
			sale_price,
			available_quantity,
			original_stock,
			created_at,
			updated_at,
			deleted_at
		FROM flash_sale_event_products
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&flashSaleEventProductModel.Id,
		&flashSaleEventProductModel.EventId,
		&flashSaleEventProductModel.ProductId,
		&flashSaleEventProductModel.DiscountPercentage,
		&flashSaleEventProductModel.SalePrice,
		&flashSaleEventProductModel.AvailableQuantity,
		&flashSaleEventProductModel.OriginalStock,
		&flashSaleEventProductModel.CreatedAt,
		&flashSaleEventProductModel.UpdatedAt,
		&flashSaleEventProductModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeFlashSaleEventProductProto(flashSaleEventProductModel), nil
}

func (r *FlashSaleEventProductRepo) UpdateFlashSaleEventProduct(ctx context.Context, req *product_service.UpdateFlashSaleEventProductRequest) (*product_service.FlashSaleEventProduct, error) {
	if req.FlashSaleEventProduct.EventId != "" {
		err := r.checkFlashSaleEventStatus(ctx, req.FlashSaleEventProduct.EventId)
		if err != nil {
			return nil, err
		}
	}
	query := `
		UPDATE flash_sale_event_products
		SET 
			event_id = $1,
			product_id = $2,
			discount_percentage = $3,
			sale_price = $4,
			available_quantity = $5,
			original_stock = $6,
			updated_at = NOW()
		WHERE id = $7 AND deleted_at = 0
		RETURNING id, event_id, product_id, discount_percentage, sale_price, available_quantity, original_stock, created_at, updated_at
	`

	flashSaleEventProductModel := makeFlashSaleEventProductModel(req.FlashSaleEventProduct)

	err := r.db.QueryRow(ctx, query,
		flashSaleEventProductModel.EventId,
		flashSaleEventProductModel.ProductId,
		flashSaleEventProductModel.DiscountPercentage,
		flashSaleEventProductModel.SalePrice,
		flashSaleEventProductModel.AvailableQuantity,
		flashSaleEventProductModel.OriginalStock,
		flashSaleEventProductModel.Id,
	).Scan(
		&flashSaleEventProductModel.Id,
		&flashSaleEventProductModel.EventId,
		&flashSaleEventProductModel.ProductId,
		&flashSaleEventProductModel.DiscountPercentage,
		&flashSaleEventProductModel.SalePrice,
		&flashSaleEventProductModel.AvailableQuantity,
		&flashSaleEventProductModel.OriginalStock,
		&flashSaleEventProductModel.CreatedAt,
		&flashSaleEventProductModel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeFlashSaleEventProductProto(flashSaleEventProductModel), nil
}

func (r *FlashSaleEventProductRepo) DeleteFlashSaleEventProduct(ctx context.Context, req *product_service.DeleteFlashSaleEventProductRequest) (*product_service.DeleteFlashSaleEventProductResponse, error) {
	query := `
		UPDATE flash_sale_event_products
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at = 0
	`

	_, err := r.db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		return nil, err
	}

	return &product_service.DeleteFlashSaleEventProductResponse{
		Message: "Flash sale event product deleted successfully",
	}, nil
}

func (r *FlashSaleEventProductRepo) ListFlashSaleEventProducts(ctx context.Context, req *product_service.ListFlashSaleEventProductsRequest) (*product_service.ListFlashSaleEventProductsResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			event_id,
			product_id,
			discount_percentage,
			sale_price,
			available_quantity,
			original_stock,
			created_at,
			updated_at,
			deleted_at
		FROM 
			flash_sale_event_products
		WHERE deleted_at = 0
	`

	filter := ""

	if req.EventId != "" {
		filter += fmt.Sprintf(" AND event_id = $%d", count)
		args = append(args, req.EventId)
		count++
	}

	if req.ProductId != "" {
		filter += fmt.Sprintf(" AND product_id = $%d", count)
		args = append(args, req.ProductId)
		count++
	}

	if req.MinDiscountPercentage > 0 {
		filter += fmt.Sprintf(" AND discount_percentage >= $%d", count)
		args = append(args, req.MinDiscountPercentage)
		count++
	}

	if req.MaxDiscountPercentage > 0 {
		filter += fmt.Sprintf(" AND discount_percentage <= $%d", count)
		args = append(args, req.MaxDiscountPercentage)
		count++
	}

	if req.MinSalePrice > 0 {
		filter += fmt.Sprintf(" AND sale_price >= $%d", count)
		args = append(args, req.MinSalePrice)
		count++
	}

	if req.MaxSalePrice > 0 {
		filter += fmt.Sprintf(" AND sale_price <= $%d", count)
		args = append(args, req.MaxSalePrice)
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

	totalCountQuery := "SELECT count(*) FROM flash_sale_event_products WHERE 1=1 AND deleted_at = 0" + filter
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

	var flashSaleEventProductList []*product_service.FlashSaleEventProduct

	for rows.Next() {
		var flashSaleEventProductModel models.FlashSaleEventProduct
		err = rows.Scan(
			&flashSaleEventProductModel.Id,
			&flashSaleEventProductModel.EventId,
			&flashSaleEventProductModel.ProductId,
			&flashSaleEventProductModel.DiscountPercentage,
			&flashSaleEventProductModel.SalePrice,
			&flashSaleEventProductModel.AvailableQuantity,
			&flashSaleEventProductModel.OriginalStock,
			&flashSaleEventProductModel.CreatedAt,
			&flashSaleEventProductModel.UpdatedAt,
			&flashSaleEventProductModel.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		flashSaleEventProductList = append(flashSaleEventProductList, makeFlashSaleEventProductProto(flashSaleEventProductModel))
	}

	return &product_service.ListFlashSaleEventProductsResponse{
		FlashSaleEventProducts: flashSaleEventProductList,
		Total:                  int32(totalCount),
	}, nil
}

// Convert db model to proto model
func makeFlashSaleEventProductProto(fsep models.FlashSaleEventProduct) *product_service.FlashSaleEventProduct {
	return &product_service.FlashSaleEventProduct{
		Id:                 fsep.Id,
		EventId:            fsep.EventId,
		ProductId:          fsep.ProductId,
		DiscountPercentage: fsep.DiscountPercentage,
		SalePrice:          fsep.SalePrice,
		AvailableQuantity:  fsep.AvailableQuantity,
		OriginalStock:      fsep.OriginalStock,
		CreatedAt:          timestamppb.New(fsep.CreatedAt),
		UpdatedAt:          timestamppb.New(fsep.UpdatedAt),
	}
}

// Convert proto model to db model
func makeFlashSaleEventProductModel(fsep *product_service.FlashSaleEventProduct) models.FlashSaleEventProduct {
	return models.FlashSaleEventProduct{
		Id:                 fsep.Id,
		EventId:            fsep.EventId,
		ProductId:          fsep.ProductId,
		DiscountPercentage: fsep.DiscountPercentage,
		SalePrice:          fsep.SalePrice,
		AvailableQuantity:  fsep.AvailableQuantity,
		OriginalStock:      fsep.OriginalStock,
		CreatedAt:          fsep.CreatedAt.AsTime(),
		UpdatedAt:          fsep.UpdatedAt.AsTime(),
	}
}

func (r *FlashSaleEventProductRepo) checkFlashSaleEventStatus(ctx context.Context, eventID string) error {
	query := `
		SELECT status, end_time
		FROM flash_sale_events
		WHERE id = $1 AND deleted_at = 0
	`

	var (
		status  string
		endTime time.Time
	)

	err := r.db.QueryRow(ctx, query, eventID).Scan(&status, &endTime)
	if err != nil {
		return fmt.Errorf("error checking flash sale event status: %w", err)
	}

	if status == "ENDED" {
		return fmt.Errorf("flash sale event has ended")
	}

	if endTime.Before(time.Now()) {
		return fmt.Errorf("flash sale event has expired")
	}

	return nil
}
