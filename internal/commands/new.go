package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/valdirsb/artemis/internal/generators"
)

var newCmd = &cobra.Command{
	Use:   "new [project_name]",
	Short: "Cria um novo projeto Artemis",
	Long: `Cria um novo projeto Artemis com toda a estrutura base necessária.

Exemplo:
  artemis new myapp        # Cria projeto 'myapp' no diretório atual
  artemis new ./myapp      # Cria projeto 'myapp' no caminho especificado`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		// Extrair nome do projeto se for um caminho
		if strings.Contains(projectName, "/") {
			projectName = filepath.Base(projectName)
		}

		fmt.Printf("🏹 Criando novo projeto Artemis: %s\n", projectName)

		generator := generators.NewProjectGenerator()
		if err := generator.Generate(args[0], projectName); err != nil {
			return fmt.Errorf("erro ao criar projeto: %w", err)
		}

		fmt.Printf("✅ Projeto '%s' criado com sucesso!\n", projectName)
		fmt.Println()
		fmt.Println("Próximos passos:")
		fmt.Printf("  cd %s\n", args[0])
		fmt.Println("  go mod tidy")
		fmt.Println("  artemis make:module users")

		return nil
	},
}
