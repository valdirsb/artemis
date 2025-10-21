package generators

import (
	"fmt"
	"os"
	"path/filepath"

	// "github.com/valdirsb/artemis/internal/templates"

	"text/template"
)

type Project struct {
	Name string
}

// ProjectGenerator gera a estrutura completa de um projeto Artemis
type ProjectGenerator struct{}

// NewProjectGenerator cria uma nova instância do gerador de projetos
func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{}
}

// Generate cria a estrutura completa do projeto
func (g *ProjectGenerator) Generate(projectPath, projectName string) error {

	project := Project{Name: projectName}

	// Criar diretório do projeto
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório do projeto: %w", err)
	}

	// Estrutura de diretórios
	dirs := []string{
		// Internal structure
		"internal/bootstrap",
		"internal/routes",
		"internal/modules",

		// User module
		"internal/modules/user/adapters/grpc",
		"internal/modules/user/adapters/http",
		"internal/modules/user/application/commands",
		"internal/modules/user/application/queries",
		"internal/modules/user/application/services",
		"internal/modules/user/domain",
		"internal/modules/user/dto",
		"internal/modules/user/ports",
		"internal/modules/user/repository",
		"internal/modules/user/tests",

		// Product module
		"internal/modules/product/adapters/grpc",
		"internal/modules/product/adapters/http",
		"internal/modules/product/application/commands",
		"internal/modules/product/application/queries",
		"internal/modules/product/application/services",
		"internal/modules/product/domain",
		"internal/modules/product/dto",
		"internal/modules/product/ports",
		"internal/modules/product/repository",
		"internal/modules/product/tests",

		// Order module
		"internal/modules/order/adapters/grpc",
		"internal/modules/order/adapters/http",
		"internal/modules/order/application/commands",
		"internal/modules/order/application/queries",
		"internal/modules/order/application/services",
		"internal/modules/order/domain",
		"internal/modules/order/dto",
		"internal/modules/order/ports",
		"internal/modules/order/repository",
		"internal/modules/order/tests",

		// Pkg structure
		"pkg/adapters/database/mysql",
		"pkg/adapters/email",
		"pkg/adapters/http/middleware",
		"pkg/adapters/logger",
		"pkg/config",
		"pkg/container",
		"pkg/contracts",
		"pkg/errors",
		"pkg/events",
		"pkg/framework/interfaces",
		"pkg/framework/providers/cache",
		"pkg/framework/providers/grpc",
		"pkg/framework/providers/maps",
		"pkg/framework/providers/payment",
		"pkg/proto",

		// Tests
		"tests/e2e",

		// Documentation
		"docs",

		//Proto
		"proto",
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
		// Arquivos raiz
		"main.go":        "main.go.tmpl",
		"README.md":      "README.md.tmpl",
		"go.mod":         "go.mod.tmpl",
		".env":           "env.tmpl",
		"Makefile":       "Makefile.tmpl",
		"framework.yaml": "framework.yaml.tmpl",

		// Internal/Bootstrap
		"internal/bootstrap/bootstrap_registry.go": "bootstrap_registry.go.tmpl",
		"internal/bootstrap/mock.go":               "mock.go.tmpl",

		// Pkg/Config
		"pkg/config/config.go": "config.go.tmpl",

		// Pkg/Container
		"pkg/container/container.go": "container.go.tmpl",

		// Pkg/Events
		"pkg/events/eventbus.go": "eventbus.go.tmpl",

		// Pkg/Framework
		"pkg/framework/framework.go": "framework.go.tmpl",
		"pkg/framework/config.go":    "framework_config.go.tmpl",

		// Pkg/Contracts
		"pkg/contracts/interfaces.go": "interfaces.go.tmpl",

		// Pkg/Adapters/Database/MySQL
		"pkg/adapters/database/mysql/connection.go": "mysql_connection.go.tmpl",
		"pkg/adapters/database/mysql/migrations.go": "mysql_migrations.go.tmpl",

		// Proto files
		"proto/user.proto":    "user.proto.tmpl",
		"proto/product.proto": "product.proto.tmpl",
		"proto/order.proto":   "order.proto.tmpl",

		// User Module
		"internal/modules/user_module.go":                                        "user_module.go.tmpl",
		"internal/modules/user/domain/user.go":                                   "modules_user_domain_user.go.tmpl",
		"internal/modules/user/repository/user_repository.go":                    "modules_user_repository_user_repository.go.tmpl",
		"internal/modules/user/repository/user_model.go":                         "modules_user_repository_user_model.go.tmpl",
		"internal/modules/user/ports/ports.go":                                   "modules_user_ports_ports.go.tmpl",
		"internal/modules/user/application/services/user_application_service.go": "modules_user_application_services_user_application_service.go.tmpl",

		// TODO: Adicionar mais templates dos módulos conforme necessário
		// Este é um exemplo básico com os templates essenciais criados
	} // Criar arquivos
	for filePath, content := range files {

		fullPath := filepath.Join(projectPath, filePath)

		if err := generateProjectFile(fullPath, content, project); err != nil {
			return err
		}
	}

	return nil
}

// generateFile creates a file from a template.
func generateProjectFile(filePath, templateName string, project Project) error {

	fmt.Println("Criando o arquivo:", filePath)

	tmpl, err := template.ParseFiles(filepath.Join("internal", "generators", "templates", templateName))
	if err != nil {

		fmt.Println("Erro 1")
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("path: ", filePath)
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, project)
}
