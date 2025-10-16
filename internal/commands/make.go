package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/valdirsb/artemis/internal/generators"
)

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Gera componentes do projeto",
	Long:  `Gera diferentes tipos de componentes como módulos, migrations, etc.`,
}

var makeModuleCmd = &cobra.Command{
	Use:   "module [name]",
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

var makeMigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Gera uma nova migration",
	Long: `Gera uma nova migration com timestamp.

Exemplo:
  artemis make:migration create_users_table
  artemis make:migration add_email_to_users`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		migrationName := args[0]

		fmt.Printf("📝 Gerando migration: %s\n", migrationName)

		generator := generators.NewMigrationGenerator()
		if err := generator.Generate(migrationName); err != nil {
			return fmt.Errorf("erro ao gerar migration: %w", err)
		}

		fmt.Printf("✅ Migration '%s' gerada com sucesso!\n", migrationName)

		return nil
	},
}

func init() {
	makeCmd.AddCommand(makeModuleCmd)
	makeCmd.AddCommand(makeMigrationCmd)
}
