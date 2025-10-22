package mysql
import (
"fmt"
"log"
"os"
"time"

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
