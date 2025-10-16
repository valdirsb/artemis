package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/valdirsb/artemis/internal/templates"
)

// ModuleGenerator gera módulos do projeto
type ModuleGenerator struct{}

// NewModuleGenerator cria uma nova instância do gerador de módulos
func NewModuleGenerator() *ModuleGenerator {
	return &ModuleGenerator{}
}

// Generate cria a estrutura completa do módulo
func (g *ModuleGenerator) Generate(moduleName string) error {
	// Verificar se estamos em um projeto Artemis
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("não foi encontrado go.mod - certifique-se de estar em um projeto Go válido")
	}

	// Nome do módulo em diferentes formatos
	moduleNameLower := strings.ToLower(moduleName)
	moduleNameTitle := strings.Title(moduleNameLower)

	// Diretório base do módulo
	moduleDir := filepath.Join("internal", "modules", moduleNameLower)

	// Estrutura de diretórios do módulo
	dirs := []string{
		filepath.Join(moduleDir, "domain"),
		filepath.Join(moduleDir, "ports"),
		filepath.Join(moduleDir, "service"),
		filepath.Join(moduleDir, "adapters"),
		filepath.Join(moduleDir, "repository"),
		filepath.Join(moduleDir, "handler"),
	}

	// Criar diretórios
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
		}
	}

	// Gerar arquivos do módulo
	files := map[string]string{
		filepath.Join(moduleDir, "domain", moduleNameLower+".go"):                templates.EntityTemplate(moduleNameTitle, moduleNameLower),
		filepath.Join(moduleDir, "domain", "repository.go"):                      templates.DomainRepositoryTemplate(moduleNameTitle, moduleNameLower),
		filepath.Join(moduleDir, "ports", "ports.go"):                            templates.PortsTemplate(moduleNameTitle, moduleNameLower),
		filepath.Join(moduleDir, "service", moduleNameLower+"_service.go"):       templates.ServiceTemplate(moduleNameTitle, moduleNameLower),
		filepath.Join(moduleDir, "repository", moduleNameLower+"_repository.go"): templates.RepositoryTemplate(moduleNameTitle, moduleNameLower),
		filepath.Join(moduleDir, "handler", moduleNameLower+"_handler.go"):       templates.HandlerTemplate(moduleNameTitle, moduleNameLower),
	}

	// Criar arquivos
	for filePath, content := range files {
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("erro ao criar arquivo %s: %w", filePath, err)
		}
	}

	return nil
}
