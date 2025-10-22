package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/valdirsb/artemis/internal/generators"
)

var makeModuleCmd = &cobra.Command{
	Use:   "make:module [name]",
	Short: "Gera um novo módulo",
	Long: `Gera um novo módulo com toda a estrutura necessária:
  • Domain (entidades e repositório)
  • Service (casos de uso)
  • Repository (implementação)
  • Handler (controladores HTTP)

Exemplo:
  artemis make:module users     # Gera módulo 'users'
  artemis make:module products  # Gera módulo 'products'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleName := strings.ToLower(args[0])

		fmt.Printf("🏗️  Gerando módulo: %s\n", moduleName)

		generator := generators.NewModuleGenerator()
		if err := generator.Generate(moduleName); err != nil {
			return fmt.Errorf("erro ao gerar módulo: %w", err)
		}

		fmt.Printf("✅ Módulo '%s' gerado com sucesso!\n", moduleName)

		return nil
	},
}

var (
	crud bool
)

func init() {

	makeModuleCmd.Flags().BoolVarP(&crud, "crud", "c", false, "Generate CRUD endpoints")
}
