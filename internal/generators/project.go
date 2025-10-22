package generators

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	// "github.com/valdirsb/artemis/internal/templates"

	"text/template"

	"github.com/valdirsb/artemis/internal/utils"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

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
		if err := utils.CreateDir(dirPath); err != nil {
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

		//Docs
		"docs/docs.go":      "docs_docs.go.tmpl",
		"docs/swagger.json": "docs_swagger.json.tmpl",
		"docs/swagger.yaml": "docs_swagger.yaml.tmpl",

		// Internal/Bootstrap
		"internal/bootstrap/bootstrap_registry.go": "bootstrap_registry.go.tmpl",
		"internal/bootstrap/mock.go":               "mock.go.tmpl",

		// Pkg/Config
		"pkg/config/config.go": "config.go.tmpl",

		// Pkg/Container
		"pkg/container/container.go": "pkg_container_container.go.tmpl",
		"pkg/container/registry.go":  "pkg_container_registry.go.tmpl",

		// Pkg/Events
		"pkg/events/eventbus.go": "pkg_events_eventbus.go.tmpl",
		"pkg/events/handlers.go": "pkg_events_handlers.go.tmpl",
		"pkg/events/typed.go":    "pkg_events_typed.go.tmpl",
		"pkg/events/types.go":    "pkg_events_types.go.tmpl",

		// Pkg/Framework
		"pkg/framework/framework.go":                 "pkg_framework_framework.go.tmpl",
		"pkg/framework/config.go":                    "pkg_framework_config.go.tmpl",
		"pkg/framework/providers/grpc/grpc.go":       "pkg_framework_providers_grpc_grpc.go.tmpl",
		"pkg/framework/interfaces/module.go":         "pkg_framework_interfaces_module.go.tmpl",
		"pkg/framework/interfaces/provider.go":       "pkg_framework_interfaces_provider.go.tmpl",
		"pkg/framework/providers/provider.go":        "pkg_framework_providers_provider.go.tmpl",
		"pkg/framework/providers/cache/redis.go":     "pkg_framework_providers_cache_redis.go.tmpl",
		"pkg/framework/providers/maps/maps.go":       "pkg_framework_providers_maps_maps.go.tmpl",
		"pkg/framework/providers/payment/payment.go": "pkg_framework_providers_payment_payment.go.tmpl",

		// Pkg/Contracts
		"pkg/contracts/interfaces.go":      "pkg_contracts_interfaces.go.tmpl",
		"pkg/contracts/infrastructure.go":  "pkg_contracts_infrastructure.go.tmpl",
		"pkg/contracts/interfaces_user.go": "pkg_contracts_interfaces_user.go.tmpl",

		// Pkg/Adapters/Database/MySQL
		"pkg/adapters/database/mysql/connection.go": "adapters_database_mysql_connection.go.tmpl",
		// "pkg/adapters/database/mysql/migrations.go": "adapters_database_mysql_migrations.go.tmpl",
		"pkg/adapters/database/mysql/migrations.go": "mysql_migrations.go.tmpl",

		// Pkg/Adapters/Logger
		"pkg/adapters/logger/logger.go": "adapters_logger_logger.go.tmpl",

		// Pkg/Adapters/HTTP/Middleware
		"pkg/adapters/http/middleware/middleware.go":    "adapters_http_middleware_middleware.go.tmpl",
		"pkg/adapters/http/middleware/error_handler.go": "adapters_http_middleware_error_handler.go.tmpl",

		// Pkg/Errors
		"pkg/errors/errors.go": "pkg_errors_errors.go.tmpl",

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
		"internal/modules/user/application/commands/create_user.go":              "modules_user_application_commands_create_user.go.tmpl",
		"internal/modules/user/application/commands/update_user.go":              "modules_user_application_commands_update_user.go.tmpl",
		"internal/modules/user/application/commands/delete_user.go":              "modules_user_application_commands_delete_user.go.tmpl",
		"internal/modules/user/application/commands/validate_credentials.go":     "modules_user_application_commands_validate_credentials.go.tmpl",
		"internal/modules/user/application/queries/get_user.go":                  "modules_user_application_queries_get_user.go.tmpl",
		"internal/modules/user/application/queries/list_users.go":                "modules_user_application_queries_list_users.go.tmpl",
		"internal/modules/user/application/queries/get_user_by_email.go":         "modules_user_application_queries_get_user_by_email.go.tmpl",
		"internal/modules/user/errors.go":                                        "modules_user_errors.go.tmpl",
		"internal/modules/user/dto/requests.go":                                  "modules_user_dto_requests.go.tmpl",
		"internal/modules/user/dto/responses.go":                                 "modules_user_dto_responses.go.tmpl",
		"internal/modules/user/dto/mapper.go":                                    "modules_user_dto_mapper.go.tmpl",
		"internal/modules/user/adapters/http/user_http_handler.go":               "modules_user_adapters_http_user_http_handler.go.tmpl",
		"internal/modules/user/adapters/grpc/user_grpc_handler.go":               "modules_user_adapters_grpc_user_grpc_handler.go.tmpl",
		"internal/modules/user/adapters/email_service.go":                        "modules_user_adapters_email_service.go.tmpl",
		"internal/modules/user/adapters/password_hasher.go":                      "modules_user_adapters_password_hasher.go.tmpl",
		"internal/modules/user/adapters/logger.go":                               "modules_user_adapters_logger.go.tmpl",
	}

	// Criar arquivos
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

	var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.tmpl"))

	// tmpl, err := template.ParseFiles(filepath.Join("internal", "generators", "templates", templateName))
	// if err != nil {

	// 	fmt.Println("Erro 1")
	// 	return err
	// }

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("path: ", filePath)
		return err
	}
	defer file.Close()

	return tmpl.ExecuteTemplate(file, templateName, project)
}
