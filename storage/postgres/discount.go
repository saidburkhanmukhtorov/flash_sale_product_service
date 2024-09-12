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

type DiscountRepo struct {
	db *pgx.Conn
}

func NewDiscountuctRepo(db *pgx.Conn) *DiscountRepo {
	return &DiscountRepo{
		db: db,
	}
}
func (r *DiscountRepo) CreateDiscount(ctx context.Context, req *product_service.CreateDiscountRequest) (*product_service.Discount, error) {
	if req.Discount.Id == "" {
		req.Discount.Id = uuid.NewString()
	}

	query := `
		INSERT INTO discounts (
			id,
			name,
			description,
			discount_type,
			discount_value,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), 0
		) RETURNING id, created_at, updated_at
	`

	discountModel := makeDiscountModel(req.Discount)

	err := r.db.QueryRow(ctx, query,
		discountModel.Id,
		discountModel.Name,
		discountModel.Description,
		discountModel.DiscountType,
		discountModel.DiscountValue,
		discountModel.StartDate,
		discountModel.EndDate,
		discountModel.IsActive,
	).Scan(&discountModel.Id, &discountModel.CreatedAt, &discountModel.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return makeDiscountProto(discountModel), nil
}

func (r *DiscountRepo) GetDiscount(ctx context.Context, req *product_service.GetDiscountRequest) (*product_service.Discount, error) {
	var discountModel models.Discount

	query := `
		SELECT 
			id,
			name,
			description,
			discount_type,
			discount_value,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM discounts
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&discountModel.Id,
		&discountModel.Name,
		&discountModel.Description,
		&discountModel.DiscountType,
		&discountModel.DiscountValue,
		&discountModel.StartDate,
		&discountModel.EndDate,
		&discountModel.IsActive,
		&discountModel.CreatedAt,
		&discountModel.UpdatedAt,
		&discountModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeDiscountProto(discountModel), nil
}

func (r *DiscountRepo) UpdateDiscount(ctx context.Context, req *product_service.UpdateDiscountRequest) (*product_service.Discount, error) {
	query := `
		UPDATE discounts
		SET 
			name = $1,
			description = $2,
			discount_type = $3,
			discount_value = $4,
			start_date = $5,
			end_date = $6,
			is_active = $7,
			updated_at = NOW()
		WHERE id = $8 AND deleted_at = 0
		RETURNING id, name, description, discount_type, discount_value, start_date, end_date, is_active, created_at, updated_at
	`

	discountModel := makeDiscountModel(req.Discount)

	err := r.db.QueryRow(ctx, query,
		discountModel.Name,
		discountModel.Description,
		discountModel.DiscountType,
		discountModel.DiscountValue,
		discountModel.StartDate,
		discountModel.EndDate,
		discountModel.IsActive,
		discountModel.Id,
	).Scan(
		&discountModel.Id,
		&discountModel.Name,
		&discountModel.Description,
		&discountModel.DiscountType,
		&discountModel.DiscountValue,
		&discountModel.StartDate,
		&discountModel.EndDate,
		&discountModel.IsActive,
		&discountModel.CreatedAt,
		&discountModel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeDiscountProto(discountModel), nil
}

func (r *DiscountRepo) DeleteDiscount(ctx context.Context, req *product_service.DeleteDiscountRequest) (*product_service.DeleteDiscountResponse, error) {
	query := `
		UPDATE discounts
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at = 0
	`

	_, err := r.db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		return nil, err
	}

	return &product_service.DeleteDiscountResponse{
		Message: "Discount deleted successfully",
	}, nil
}
func (r *DiscountRepo) ListDiscounts(ctx context.Context, req *product_service.ListDiscountsRequest) (*product_service.ListDiscountsResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			name,
			description,
			discount_type,
			discount_value,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM 
			discounts
		WHERE 1=1 AND deleted_at = 0
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

	if req.DiscountType != "" {
		filter += fmt.Sprintf(" AND discount_type = $%d", count)
		args = append(args, req.DiscountType)
		count++
	}

	if req.MinDiscountValue > 0 {
		filter += fmt.Sprintf(" AND discount_value >= $%d", count)
		args = append(args, req.MinDiscountValue)
		count++
	}

	if req.MaxDiscountValue > 0 {
		filter += fmt.Sprintf(" AND discount_value <= $%d", count)
		args = append(args, req.MaxDiscountValue)
		count++
	}

	if req.StartDate != nil {
		filter += fmt.Sprintf(" AND start_date >= $%d", count)
		args = append(args, req.StartDate.AsTime())
		count++
	}

	if req.EndDate != nil {
		filter += fmt.Sprintf(" AND end_date <= $%d", count)
		args = append(args, req.EndDate.AsTime())
		count++
	}

	if req.IsActive {
		filter += fmt.Sprintf(" AND is_active = $%d", count)
		args = append(args, req.IsActive)
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

	totalCountQuery := "SELECT count(*) FROM discounts WHERE 1=1 AND deleted_at = 0" + filter
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

	var discountList []*product_service.Discount

	for rows.Next() {
		var discountModel models.Discount
		err = rows.Scan(
			&discountModel.Id,
			&discountModel.Name,
			&discountModel.Description,
			&discountModel.DiscountType,
			&discountModel.DiscountValue,
			&discountModel.StartDate,
			&discountModel.EndDate,
			&discountModel.IsActive,
			&discountModel.CreatedAt,
			&discountModel.UpdatedAt,
			&discountModel.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		discountList = append(discountList, makeDiscountProto(discountModel))
	}

	return &product_service.ListDiscountsResponse{
		Discounts: discountList,
		Total:     int32(totalCount),
	}, nil
}

// Convert db model to proto model
func makeDiscountProto(discount models.Discount) *product_service.Discount {
	return &product_service.Discount{
		Id:            discount.Id,
		Name:          discount.Name,
		Description:   discount.Description,
		DiscountType:  discount.DiscountType,
		DiscountValue: discount.DiscountValue,
		StartDate:     timestamppb.New(discount.StartDate),
		EndDate:       timestamppb.New(discount.EndDate),
		IsActive:      discount.IsActive,
		CreatedAt:     timestamppb.New(discount.CreatedAt),
		UpdatedAt:     timestamppb.New(discount.UpdatedAt),
	}
}

// Convert proto model to db model
func makeDiscountModel(discount *product_service.Discount) models.Discount {
	return models.Discount{
		Id:            discount.Id,
		Name:          discount.Name,
		Description:   discount.Description,
		DiscountType:  discount.DiscountType,
		DiscountValue: discount.DiscountValue,
		StartDate:     discount.StartDate.AsTime(),
		EndDate:       discount.EndDate.AsTime(),
		IsActive:      discount.IsActive,
		CreatedAt:     discount.CreatedAt.AsTime(),
		UpdatedAt:     discount.UpdatedAt.AsTime(),
	}
}
