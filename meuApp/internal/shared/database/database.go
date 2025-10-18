package database

import (
	"fmt"
	"log"
	"os"
	"time"

	orderRepo "meuApp/internal/modules/order/repository"
	productRepo "meuApp/internal/modules/product/repository"
	userRepo "meuApp/internal/modules/user/repository"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig contém as configurações de conexão do banco
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// GetDefaultConfig retorna a configuração do banco usando variáveis de ambiente ou padrões
func GetDefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "3306"),
		Username: getEnvOrDefault("DB_USERNAME", "root"),
		Password: getEnvOrDefault("DB_PASSWORD", "123456"),
		Database: getEnvOrDefault("DB_DATABASE", "app_db"),
	}
}

// getEnvOrDefault obtém valor de variável de ambiente ou retorna padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Connect estabelece conexão com o banco MySQL
func Connect(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configurar connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configurações do pool de conexões
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to MySQL database")
	return db, nil
}

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
		// Não falhar se o seed der erro, apenas avisar
	}

	return nil
}

// SeedDatabase popula o banco com dados iniciais
func SeedDatabase(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Verificar se já existem produtos para evitar duplicação
	var count int64
	if err := db.Model(&productRepo.ProductModel{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count products: %w", err)
	}

	if count > 0 {
		log.Printf("Database already has %d products, skipping seed", count)
		return nil
	}

	// Seeds de produtos
	products := []productRepo.ProductModel{
		{
			ID:          "prod-001",
			Name:        "iPhone 15 Pro Max",
			Description: "Apple iPhone 15 Pro Max 256GB - Titânio Natural com câmera profissional de 48MP",
			Price:       8999.99,
			Stock:       15,
			CategoryID:  "electronics",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-002",
			Name:        "MacBook Air M2",
			Description: "MacBook Air 13\" com chip M2, 8GB RAM, 256GB SSD - Cor Meia-noite",
			Price:       12999.99,
			Stock:       8,
			CategoryID:  "computers",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-003",
			Name:        "Samsung Galaxy S24 Ultra",
			Description: "Samsung Galaxy S24 Ultra 512GB - Preto com S Pen incluída e câmera de 200MP",
			Price:       7499.99,
			Stock:       12,
			CategoryID:  "electronics",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-004",
			Name:        "Dell XPS 13",
			Description: "Notebook Dell XPS 13 Intel Core i7, 16GB RAM, 512GB SSD, Tela InfinityEdge",
			Price:       9999.99,
			Stock:       6,
			CategoryID:  "computers",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-005",
			Name:        "AirPods Pro (3ª geração)",
			Description: "Apple AirPods Pro com cancelamento ativo de ruído e case de carregamento MagSafe",
			Price:       2499.99,
			Stock:       25,
			CategoryID:  "accessories",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-006",
			Name:        "Sony WH-1000XM5",
			Description: "Fone de ouvido Sony WH-1000XM5 com cancelamento de ruído premium e 30h de bateria",
			Price:       1899.99,
			Stock:       18,
			CategoryID:  "accessories",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-007",
			Name:        "iPad Air (5ª geração)",
			Description: "iPad Air 10.9\" com chip M1, 256GB, Wi-Fi + Cellular - Azul-céu",
			Price:       6499.99,
			Stock:       10,
			CategoryID:  "tablets",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-008",
			Name:        "Nintendo Switch OLED",
			Description: "Console Nintendo Switch modelo OLED com tela de 7\" e 64GB de armazenamento",
			Price:       2799.99,
			Stock:       20,
			CategoryID:  "gaming",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-009",
			Name:        "PlayStation 5",
			Description: "Console Sony PlayStation 5 com SSD ultrarrápido e controle DualSense",
			Price:       4999.99,
			Stock:       5,
			CategoryID:  "gaming",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-010",
			Name:        "Microsoft Surface Pro 9",
			Description: "Surface Pro 9 Intel i7, 16GB RAM, 512GB SSD com teclado Type Cover incluso",
			Price:       11999.99,
			Stock:       7,
			CategoryID:  "tablets",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-011",
			Name:        "LG OLED C3 55\"",
			Description: "Smart TV LG OLED C3 55\" 4K com webOS, Dolby Vision IQ e Gaming Hub",
			Price:       6999.99,
			Stock:       4,
			CategoryID:  "tv",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "prod-012",
			Name:        "Apple Watch Series 9",
			Description: "Apple Watch Series 9 GPS 45mm caixa de alumínio com pulseira esportiva",
			Price:       3999.99,
			Stock:       14,
			CategoryID:  "wearables",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Inserir produtos em lote
	if err := db.CreateInBatches(products, 100).Error; err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	log.Printf("Successfully seeded %d products", len(products))
	return nil
}
