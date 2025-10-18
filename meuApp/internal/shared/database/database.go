package database

import (
	"meuApp/pkg/adapters/database/mysql"
	"os"

	"gorm.io/gorm"
)

// DatabaseConfig é um alias para manter compatibilidade
type DatabaseConfig = mysql.DatabaseConfig

// GetDefaultConfig retorna a configuração do banco
func GetDefaultConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "3306"),
		Username: getEnvOrDefault("DB_USERNAME", "root"),
		Password: getEnvOrDefault("DB_PASSWORD", "123456"),
		Database: getEnvOrDefault("DB_DATABASE", "app_db"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Connect estabelece conexão com o banco MySQL
func Connect(config *DatabaseConfig) (*gorm.DB, error) {
	return mysql.Connect(config)
}

// AutoMigrate executa as migrações necessárias
func AutoMigrate(db *gorm.DB) error {
	return mysql.AutoMigrate(db)
}

// SeedDatabase popula o banco com dados iniciais
func SeedDatabase(db *gorm.DB) error {
	return mysql.SeedDatabase(db)
}
