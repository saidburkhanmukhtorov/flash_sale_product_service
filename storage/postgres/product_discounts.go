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

type ProductDiscountRepo struct {
	db *pgx.Conn
}

func NewProductDiscountuctRepo(db *pgx.Conn) *ProductDiscountRepo {
	return &ProductDiscountRepo{
		db: db,
	}
}

func (r *ProductDiscountRepo) CreateProductDiscount(ctx context.Context, req *product_service.CreateProductDiscountRequest) (*product_service.ProductDiscount, error) {
	if req.ProductDiscount.Id == "" {
		req.ProductDiscount.Id = uuid.NewString()
	}

	query := `
		INSERT INTO product_discounts (
			id,
			product_id,
			discount_id,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			$1, $2, $3, NOW(), NOW(), 0
		) RETURNING id, created_at, updated_at
	`

	productDiscountModel := makeProductDiscountModel(req.ProductDiscount)

	err := r.db.QueryRow(ctx, query,
		productDiscountModel.Id,
		productDiscountModel.ProductId,
		productDiscountModel.DiscountId,
	).Scan(&productDiscountModel.Id, &productDiscountModel.CreatedAt, &productDiscountModel.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return makeProductDiscountProto(productDiscountModel), nil
}

func (r *ProductDiscountRepo) GetProductDiscount(ctx context.Context, req *product_service.GetProductDiscountRequest) (*product_service.ProductDiscount, error) {
	var productDiscountModel models.ProductDiscount

	query := `
		SELECT 
			id,
			product_id,
			discount_id,
			created_at,
			updated_at,
			deleted_at
		FROM product_discounts
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&productDiscountModel.Id,
		&productDiscountModel.ProductId,
		&productDiscountModel.DiscountId,
		&productDiscountModel.CreatedAt,
		&productDiscountModel.UpdatedAt,
		&productDiscountModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeProductDiscountProto(productDiscountModel), nil
}

func (r *ProductDiscountRepo) UpdateProductDiscount(ctx context.Context, req *product_service.UpdateProductDiscountRequest) (*product_service.ProductDiscount, error) {
	query := `
		UPDATE product_discounts
		SET 
			product_id = $1,
			discount_id = $2,
			updated_at = NOW()
		WHERE id = $3 AND deleted_at = 0
		RETURNING id, product_id, discount_id, created_at, updated_at
	`

	productDiscountModel := makeProductDiscountModel(req.ProductDiscount)

	err := r.db.QueryRow(ctx, query,
		productDiscountModel.ProductId,
		productDiscountModel.DiscountId,
		productDiscountModel.Id,
	).Scan(
		&productDiscountModel.Id,
		&productDiscountModel.ProductId,
		&productDiscountModel.DiscountId,
		&productDiscountModel.CreatedAt,
		&productDiscountModel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeProductDiscountProto(productDiscountModel), nil
}

func (r *ProductDiscountRepo) DeleteProductDiscount(ctx context.Context, req *product_service.DeleteProductDiscountRequest) (*product_service.DeleteProductDiscountResponse, error) {
	query := `
		UPDATE product_discounts
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at = 0
	`

	_, err := r.db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		return nil, err
	}

	return &product_service.DeleteProductDiscountResponse{
		Message: "Product discount deleted successfully",
	}, nil
}
func (r *ProductDiscountRepo) ListProductDiscounts(ctx context.Context, req *product_service.ListProductDiscountsRequest) (*product_service.ListProductDiscountsResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			product_id,
			discount_id,
			created_at,
			updated_at,
			deleted_at
		FROM 
			product_discounts
		WHERE deleted_at = 0
	`

	filter := ""

	if req.ProductId != "" {
		filter += fmt.Sprintf(" AND product_id = $%d", count)
		args = append(args, req.ProductId)
		count++
	}

	if req.DiscountId != "" {
		filter += fmt.Sprintf(" AND discount_id = $%d", count)
		args = append(args, req.DiscountId)
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

	totalCountQuery := "SELECT count(*) FROM product_discounts WHERE 1=1 AND deleted_at = 0" + filter
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

	var productDiscountList []*product_service.ProductDiscount

	for rows.Next() {
		var productDiscountModel models.ProductDiscount
		err = rows.Scan(
			&productDiscountModel.Id,
			&productDiscountModel.ProductId,
			&productDiscountModel.DiscountId,
			&productDiscountModel.CreatedAt,
			&productDiscountModel.UpdatedAt,
			&productDiscountModel.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		productDiscountList = append(productDiscountList, makeProductDiscountProto(productDiscountModel))
	}

	return &product_service.ListProductDiscountsResponse{
		ProductDiscounts: productDiscountList,
		Total:            int32(totalCount),
	}, nil
}

// Convert db model to proto model
func makeProductDiscountProto(pd models.ProductDiscount) *product_service.ProductDiscount {
	return &product_service.ProductDiscount{
		Id:         pd.Id,
		ProductId:  pd.ProductId,
		DiscountId: pd.DiscountId,
		CreatedAt:  timestamppb.New(pd.CreatedAt),
		UpdatedAt:  timestamppb.New(pd.UpdatedAt),
	}
}

// Convert proto model to db model
func makeProductDiscountModel(pd *product_service.ProductDiscount) models.ProductDiscount {
	return models.ProductDiscount{
		Id:         pd.Id,
		ProductId:  pd.ProductId,
		DiscountId: pd.DiscountId,
		CreatedAt:  pd.CreatedAt.AsTime(),
		UpdatedAt:  pd.UpdatedAt.AsTime(),
	}
}
