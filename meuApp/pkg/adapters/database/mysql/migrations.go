package mysql

import (
	"fmt"
	"log"
	"time"

	orderRepo "meuApp/internal/modules/order/repository"
	productRepo "meuApp/internal/modules/product/repository"
	userRepo "meuApp/internal/modules/user/repository"

	"gorm.io/gorm"
)

// AutoMigrate executa as migrações necessárias
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&userRepo.UserModel{},
		&productRepo.ProductModel{},
		&orderRepo.OrderModel{},
		&orderRepo.OrderItemModel{},
	)
	if err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	log.Println("Database migration completed successfully")

	// Executar seeds
	if err := SeedDatabase(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	return nil
}

// SeedDatabase popula o banco com dados iniciais
func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	var count int64
	if err := db.Model(&productRepo.ProductModel{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count products: %w", err)
	}

	if count > 0 {
		log.Printf("Database already has %d products, skipping seed", count)
		return nil
	}

	products := []productRepo.ProductModel{
		{
			ID:          "prod-001",
			Name:        "iPhone 15 Pro Max",
			Description: "Apple iPhone 15 Pro Max 256GB",
			Price:       8999.99,
			Stock:       15,
			CategoryID:  "electronics",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.CreateInBatches(products, 100).Error; err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	log.Printf("Successfully seeded %d products", len(products))
	return nil
}
