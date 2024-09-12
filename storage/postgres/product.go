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

type ProductRepo struct {
	db *pgx.Conn
}

func NewProductRepo(db *pgx.Conn) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (r *ProductRepo) CreateProduct(ctx context.Context, req *product_service.CreateProductRequest) (*product_service.Product, error) {
	if req.Product.Id == "" {
		req.Product.Id = uuid.NewString()
	}

	query := `
		INSERT INTO products (
			id,
			name,
			description,
			base_price,
			current_price,
			image_url,
			stock_quantity,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0
		) RETURNING id, created_at, updated_at, deleted_at
	`

	// Convert proto model to db model before inserting
	productModel := makeProductModel(req.Product)

	err := r.db.QueryRow(ctx, query,
		productModel.Id,
		productModel.Name,
		productModel.Description,
		productModel.BasePrice,
		productModel.CurrentPrice,
		productModel.ImageUrl,
		productModel.StockQuantity,
	).Scan(&productModel.Id, &productModel.CreatedAt, &productModel.UpdatedAt, &productModel.DeletedAt)

	if err != nil {
		return nil, err
	}

	// Convert db model back to proto model for response
	return makeProductProto(productModel), nil
}

func (r *ProductRepo) GetProduct(ctx context.Context, req *product_service.GetProductRequest) (*product_service.Product, error) {
	var (
		productModel models.Product // Use database model for querying
	)
	query := `
		SELECT 
			id,
			name,
			description,
			base_price,
			current_price,
			image_url,
			stock_quantity,
			created_at,
			updated_at,
			deleted_at
		FROM products
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&productModel.Id,
		&productModel.Name,
		&productModel.Description,
		&productModel.BasePrice,
		&productModel.CurrentPrice,
		&productModel.ImageUrl,
		&productModel.StockQuantity,
		&productModel.CreatedAt,
		&productModel.UpdatedAt,
		&productModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	// Convert db model to proto model for response
	return makeProductProto(productModel), nil
}

func (r *ProductRepo) ListProducts(ctx context.Context, req *product_service.ListProductsRequest) (*product_service.ListProductsResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			name,
			description,
			base_price,
			current_price,
			image_url,
			stock_quantity,
			created_at,
			updated_at,
			deleted_at
		FROM 
			products
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

	if req.MinPrice > 0 {
		filter += fmt.Sprintf(" AND current_price >= $%d", count)
		args = append(args, req.MinPrice)
		count++
	}

	if req.MaxPrice > 0 {
		filter += fmt.Sprintf(" AND current_price <= $%d", count)
		args = append(args, req.MaxPrice)
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

	totalCountQuery := "SELECT count(*) FROM products WHERE 1=1 AND deleted_at = 0" + filter
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

	var productList []*product_service.Product

	for rows.Next() {
		var productModel models.Product // Use database model for scanning
		err = rows.Scan(
			&productModel.Id,
			&productModel.Name,
			&productModel.Description,
			&productModel.BasePrice,
			&productModel.CurrentPrice,
			&productModel.ImageUrl,
			&productModel.StockQuantity,
			&productModel.CreatedAt,
			&productModel.UpdatedAt,
			&productModel.DeletedAt, // Scan deleted_at
		)
		if err != nil {
			return nil, err
		}

		// Convert db model to proto model before appending to list
		productList = append(productList, makeProductProto(productModel))
	}

	return &product_service.ListProductsResponse{
		Products: productList,
		Total:    int32(totalCount),
	}, nil
}

func (r *ProductRepo) UpdateProduct(ctx context.Context, req *product_service.UpdateProductRequest) (*product_service.Product, error) {
	query := `
		UPDATE products
		SET 
			name = $1,
			description = $2,
			base_price = $3,
			current_price = $4,
			image_url = $5,
			stock_quantity = $6,
			updated_at = NOW()
		WHERE id = $7 AND deleted_at = 0
		RETURNING id, name, description, base_price, current_price, image_url, stock_quantity, created_at, updated_at, deleted_at
	`

	// Convert proto model to db model before updating
	productModel := makeProductModel(req.Product)

	err := r.db.QueryRow(ctx, query,
		productModel.Name,
		productModel.Description,
		productModel.BasePrice,
		productModel.CurrentPrice,
		productModel.ImageUrl,
		productModel.StockQuantity,
		productModel.Id,
	).Scan(
		&productModel.Id,
		&productModel.Name,
		&productModel.Description,
		&productModel.BasePrice,
		&productModel.CurrentPrice,
		&productModel.ImageUrl,
		&productModel.StockQuantity,
		&productModel.CreatedAt,
		&productModel.UpdatedAt,
		&productModel.DeletedAt, // Scan deleted_at
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	// Convert db model back to proto model for response
	return makeProductProto(productModel), nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, req *product_service.DeleteProductRequest) (*product_service.DeleteProductResponse, error) {
	query := `
		UPDATE products
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at = 0
	`

	_, err := r.db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		return nil, err
	}

	return &product_service.DeleteProductResponse{
		Message: "Product deleted successfully",
	}, nil
}

// Convert db model to proto model
func makeProductProto(product models.Product) *product_service.Product {
	protoProduct := &product_service.Product{
		Id:            product.Id,
		Name:          product.Name,
		Description:   product.Description,
		BasePrice:     product.BasePrice,
		CurrentPrice:  product.CurrentPrice,
		ImageUrl:      product.ImageUrl,
		StockQuantity: product.StockQuantity,
		CreatedAt:     timestamppb.New(product.CreatedAt),
		UpdatedAt:     timestamppb.New(product.UpdatedAt),
	}

	return protoProduct
}

// Convert proto model to db model
func makeProductModel(product *product_service.Product) models.Product {
	modelProduct := models.Product{
		Id:            product.Id,
		Name:          product.Name,
		Description:   product.Description,
		BasePrice:     product.BasePrice,
		CurrentPrice:  product.CurrentPrice,
		ImageUrl:      product.ImageUrl,
		StockQuantity: product.StockQuantity,
		CreatedAt:     product.CreatedAt.AsTime(),
		UpdatedAt:     product.UpdatedAt.AsTime(),
	}

	return modelProduct
}
