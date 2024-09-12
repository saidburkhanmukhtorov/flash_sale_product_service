package postgres

import (
	"context"
	"fmt"

	"github.com/flash_sale/flash_sale_product_service/config"
	"github.com/flash_sale/flash_sale_product_service/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Storage implements the storage.StorageI interface for PostgreSQL.
type StoragePg struct {
	db                   *pgx.Conn
	productRepo          storage.ProductI
	discountRepo         storage.DiscountI
	flashSaleEventRepo   storage.FlashSaleEventI
	productDiscountRepo  storage.ProductDiscountI
	flashSaleProductRepo storage.FlashSaleEventProductI
	notificationRepo     storage.NotificationI
}

// NewStoragePg creates a new PostgreSQL storage instance.
func NewStoragePg(cfg config.Config) (storage.StorageI, error) {
	dbCon := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	db, err := pgx.Connect(context.Background(), dbCon)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	return &StoragePg{
		db:                   db,
		productRepo:          NewProductRepo(db),
		discountRepo:         NewDiscountuctRepo(db),
		flashSaleEventRepo:   NewFlashSaleRepo(db),
		productDiscountRepo:  NewProductDiscountuctRepo(db),
		flashSaleProductRepo: NewFlashSaleEventProductRepo(db),
		notificationRepo:     NewNotificationRepo(db),
	}, nil
}

// Close closes the PostgreSQL connection.
func (s *StoragePg) Close() {
	if err := s.db.Close(context.Background()); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			fmt.Printf("Error closing database connection: %s (Code: %s)\n", pgErr.Message, pgErr.Code)
		} else {
			fmt.Printf("Error closing database connection: %s\n", err.Error())
		}
	}
}

// Product returns the ProductI implementation for PostgreSQL.
func (s *StoragePg) Product() storage.ProductI {
	return s.productRepo
}

// Discount returns the DiscountI implementation for PostgreSQL.
func (s *StoragePg) Discount() storage.DiscountI {
	return s.discountRepo
}

// FlashSaleEvent returns the FlashSaleEventI implementation for PostgreSQL.
func (s *StoragePg) FlashSaleEvent() storage.FlashSaleEventI {
	return s.flashSaleEventRepo
}

// ProductDiscount returns the ProductDiscountI implementation for PostgreSQL.
func (s *StoragePg) ProductDiscount() storage.ProductDiscountI {
	return s.productDiscountRepo
}

// FlashSaleEventProduct returns the FlashSaleEventProductI implementation for PostgreSQL.
func (s *StoragePg) FlashSaleEventProduct() storage.FlashSaleEventProductI {
	return s.flashSaleProductRepo
}

// FlashSaleEventProduct returns the FlashSaleEventProductI implementation for PostgreSQL.
func (s *StoragePg) SendNotification() storage.NotificationI {
	return s.notificationRepo
}
