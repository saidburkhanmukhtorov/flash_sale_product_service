package main

import (
	"fmt"
	"log"
	"net"

	"github.com/flash_sale/flash_sale_product_service/config"
	"github.com/flash_sale/flash_sale_product_service/genproto/product_service"
	"github.com/flash_sale/flash_sale_product_service/service"
	"github.com/flash_sale/flash_sale_product_service/storage/postgres"
	"github.com/flash_sale/flash_sale_product_service/storage/redis"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	// Initialize PostgreSQL storage
	pgStorage, err := postgres.NewStoragePg(cfg)
	if err != nil {
		log.Fatalf("failed to initialize PostgreSQL storage: %v", err)
	}

	// Initialize Redis client
	redisClient, err := redis.Connect(&cfg)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize gRPC server
	lis, err := net.Listen("tcp", cfg.ProductServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register gRPC services
	product_service.RegisterProductServiceServer(s, service.NewProductService(pgStorage))
	product_service.RegisterDiscountServiceServer(s, service.NewDiscountService(pgStorage))
	product_service.RegisterFlashSaleEventServiceServer(s, service.NewFlashSaleEventService(pgStorage, redisClient))
	product_service.RegisterProductDiscountServiceServer(s, service.NewProductDiscountService(pgStorage, redisClient))
	product_service.RegisterFlashSaleEventProductServiceServer(s, service.NewFlashSaleEventProductService(pgStorage))

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
