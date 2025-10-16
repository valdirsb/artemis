package generators

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/valdirsb/artemis/internal/templates"
)

// ProjectGenerator gera a estrutura completa de um projeto Artemis
type ProjectGenerator struct{}

// NewProjectGenerator cria uma nova instância do gerador de projetos
func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{}
}

// Generate cria a estrutura completa do projeto
func (g *ProjectGenerator) Generate(projectPath, projectName string) error {
	// Criar diretório do projeto
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %w", err)
	}

	// Estrutura de diretórios
	dirs := []string{
		"cmd/server",
		"internal/bootstrap",
		"internal/shared/config",
		"internal/shared/database",
		"internal/shared/logger",
		"internal/shared/middleware",
		"internal/modules",
		"pkg/contracts",
		"pkg/container",
		"pkg/events",
		"migrations",
		"docs",
	}

	// Criar diretórios
	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
		}
	}

	// Gerar arquivos base
	files := map[string]string{
		"go.mod":                           templates.GoModTemplate(projectName),
		"README.md":                        templates.ReadmeTemplate(projectName),
		"cmd/server/main.go":               templates.MainTemplate(projectName),
		"internal/bootstrap/bootstrap.go":  templates.BootstrapTemplate(projectName),
		"internal/shared/config/config.go": templates.ConfigTemplate(projectName),
		"pkg/contracts/interfaces.go":      templates.ContractsTemplate(projectName),
		"pkg/container/container.go":       templates.ContainerTemplate(projectName),
		".gitignore":                       templates.GitignoreTemplate(),
		"Dockerfile":                       templates.DockerfileTemplate(projectName),
		"docker-compose.yml":               templates.DockerComposeTemplate(projectName),
	}

	// Criar arquivos
	for filePath, content := range files {
		fullPath := filepath.Join(projectPath, filePath)

		// Criar diretório pai se necessário
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório pai de %s: %w", filePath, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo %s: %w", filePath, err)
		}
	}

	return nil
}
