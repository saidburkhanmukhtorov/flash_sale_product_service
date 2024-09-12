package models

import "time"

// Database Model
type Product struct {
	Id            string    `db:"id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	BasePrice     float32   `db:"base_price"`
	CurrentPrice  float32   `db:"current_price"`
	ImageUrl      string    `db:"image_url"`
	StockQuantity int32     `db:"stock_quantity"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	DeletedAt     int64     `db:"deleted_at"`
}

// Discount represents a discount model for the database.
type Discount struct {
	Id            string    `db:"id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	DiscountType  string    `db:"discount_type"` // Possible values: 'PERCENTAGE', 'FIXED_AMOUNT'
	DiscountValue float32   `db:"discount_value"`
	StartDate     time.Time `db:"start_date"`
	EndDate       time.Time `db:"end_date"`
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	DeletedAt     int64     `db:"deleted_at"`
}

// FlashSaleEvent represents a flash sale event model for the database.
type FlashSaleEvent struct {
	Id          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
	Status      string    `db:"status"`     // Possible values: 'UPCOMING', 'ACTIVE', 'ENDED'
	EventType   string    `db:"event_type"` // Possible values: 'FLASH_SALE', 'PROMOTION'
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	DeletedAt   int64     `db:"deleted_at"`
}

// ProductDiscount represents a product discount model for the database.
type ProductDiscount struct {
	Id         string    `db:"id"`
	ProductId  string    `db:"product_id"`
	DiscountId string    `db:"discount_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	DeletedAt  int64     `db:"deleted_at"`
}

// FlashSaleEventProduct represents a flash sale event product model for the database.
type FlashSaleEventProduct struct {
	Id                 string    `db:"id"`
	EventId            string    `db:"event_id"`
	ProductId          string    `db:"product_id"`
	DiscountPercentage float32   `db:"discount_percentage"`
	SalePrice          float32   `db:"sale_price"`
	AvailableQuantity  int32     `db:"available_quantity"`
	OriginalStock      int32     `db:"original_stock"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	DeletedAt          int64     `db:"deleted_at"`
}
