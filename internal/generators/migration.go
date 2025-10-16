package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/valdirsb/artemis/internal/templates"
)

// MigrationGenerator gera migrations
type MigrationGenerator struct{}

// NewMigrationGenerator cria uma nova instância do gerador de migrations
func NewMigrationGenerator() *MigrationGenerator {
	return &MigrationGenerator{}
}

// Generate cria uma nova migration
func (g *MigrationGenerator) Generate(migrationName string) error {
	// Verificar se estamos em um projeto Artemis
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("não foi encontrado go.mod - certifique-se de estar em um projeto Go válido")
	}

	// Criar diretório de migrations se não existir
	migrationsDir := "migrations"
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de migrations: %w", err)
	}

	// Gerar timestamp
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.sql", timestamp, migrationName)
	filePath := filepath.Join(migrationsDir, fileName)

	// Gerar conteúdo da migration
	content := templates.MigrationTemplate(migrationName, timestamp)

	// Criar arquivo
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("erro ao criar migration %s: %w", fileName, err)
	}

	return nil
}
