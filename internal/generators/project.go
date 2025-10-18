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
		"internal/bootstrap",
		"internal/modules",
		"internal/modules/user/adapters/",
		"internal/modules/user/domain/",
		"internal/modules/user/handler/",
		"internal/modules/user/ports/",
		"internal/modules/user/repository/",
		"internal/modules/user/service/",

		"internal/modules/product/adapters/",
		"internal/modules/product/domain/",
		"internal/modules/product/handler/",
		"internal/modules/product/ports/",
		"internal/modules/product/repository/",
		"internal/modules/product/service/",

		"internal/modules/order/adapters/",
		"internal/modules/order/domain/",
		"internal/modules/order/handler/",
		"internal/modules/order/ports/",
		"internal/modules/order/repository/",
		"internal/modules/order/service/",

		"internal/routes",
		"internal/shared/config",
		"internal/shared/database",
		"internal/shared/logger",
		"internal/shared/middleware",
		"pkg/container",
		"pkg/contracts",
		"pkg/events",
		"pkg/framework/interfaces",
		"pkg/framework/providers/cache",
		"pkg/framework/providers/grpc",
		"pkg/framework/providers/maps",
		"pkg/framework/providers/payment",
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
		"main.go":                              "main.go.tmpl",
		"README.md":                            "README.md.tmpl",
		"go.mod":                               "go.mod.tmpl",
		".env":                                 "env.tmpl",
		"Makefile":                             "Makefile.tmpl",
		"framework.yaml":                       "framework.yaml.tmpl",
		"internal/bootstrap/bootstrap.go":      "bootstrap.go.tmpl",
		"internal/bootstrap/mock.go":           "mock.go.tmpl",
		"internal/bootstrap/start_handlers.go": "start_handlers.go.tmpl",
		"internal/bootstrap/start_repositories.go":   "start_repositories.go.tmpl",
		"internal/bootstrap/start_services.go":       "start_services.go.tmpl",
		"internal/routes/routes.go":                  "routes.go.tmpl",
		"internal/shared/config/config.go":           "config.go.tmpl",
		"internal/shared/database/database.go":       "database.go.tmpl",
		"internal/shared/logger/logger.go":           "logger.go.tmpl",
		"internal/shared/middleware/middleware.go":   "middleware.go.tmpl",
		"pkg/container/container.go":                 "container.go.tmpl",
		"pkg/contracts/infrastructure.go":            "infrastructure.go.tmpl",
		"pkg/contracts/interfaces_order.go":          "interfaces_order.go.tmpl",
		"pkg/contracts/interfaces_product.go":        "interfaces_product.go.tmpl",
		"pkg/contracts/interfaces_user.go":           "interfaces_user.go.tmpl",
		"pkg/contracts/interfaces.go":                "interfaces.go.tmpl",
		"pkg/events/eventbus.go":                     "eventbus.go.tmpl",
		"pkg/framework/interfaces/provider.go":       "provider.go.tmpl",
		"pkg/framework/providers/cache/redis.go":     "redis.go.tmpl",
		"pkg/framework/providers/grpc/grpc.go":       "grpc.go.tmpl",
		"pkg/framework/providers/maps/maps.go":       "maps.go.tmpl",
		"pkg/framework/providers/payment/payment.go": "payment.go.tmpl",
		"pkg/framework/providers/provider.go":        "provider2.go.tmpl",
		"pkg/framework/config.go":                    "config-framework.go.tmpl",
		"pkg/framework/framework.go":                 "framework.go.tmpl",
		"proto/order.proto":                          "order.proto.tmpl",
		"proto/product.proto":                        "product.proto.tmpl",
		"proto/user.proto":                           "user.proto.tmpl",
		"internal/modules/README.md":                 "README-module.md.tmpl",

		//Module User
		"internal/modules/user/adapters/password_hasher.go":   "mod_user_password_hasher.go.tmpl",
		"internal/modules/user/domain/user.go":                "mod_user_domain_user.go.tmpl",
		"internal/modules/user/handler/user_handler.go":       "mod_user_handler_user_handler.go.tmpl",
		"internal/modules/user/handler/user_grpc_handler.go":  "mod_user_handler_user_grpc_handler.go.tmpl",
		"internal/modules/user/ports/ports.go":                "mod_user_ports.go.tmpl",
		"internal/modules/user/repository/user_repository.go": "mod_user_repository.go.tmpl",
		"internal/modules/user/service/user_service.go":       "mod_user_service.go.tmpl",
		"internal/modules/user/README.md":                     "README-module-user.md.tmpl",

		//Module Product
		"internal/modules/product/domain/product.go":                "mod_product_domain_product.go.tmpl",
		"internal/modules/product/handler/product_handler.go":       "mod_product_handler_product_handler.go.tmpl",
		"internal/modules/product/handler/product_grpc_handler.go":  "mod_product_handler_product_grpc_handler.go.tmpl",
		"internal/modules/product/repository/product_repository.go": "mod_product_repository.go.tmpl",
		"internal/modules/product/service/product_service.go":       "mod_product_service.go.tmpl",
		"internal/modules/product/README.md":                        "README-module-product.md.tmpl",

		//Module Order
		"internal/modules/order/domain/order.go":                "mod_order_domain_order.go.tmpl",
		"internal/modules/order/handler/order_handler.go":       "mod_order_handler_order_handler.go.tmpl",
		"internal/modules/order/handler/order_grpc_handler.go":  "mod_order_handler_order_grpc_handler.go.tmpl",
		"internal/modules/order/repository/order_repository.go": "mod_order_repository.go.tmpl",
		"internal/modules/order/service/order_service.go":       "mod_order_service.go.tmpl",
		"internal/modules/order/README.md":                      "README-module-order.md.tmpl",
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

	tmpl, err := template.ParseFiles(filepath.Join("pkg", "generator", "templates", templateName))
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
