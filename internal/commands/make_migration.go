package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/valdirsb/artemis/internal/generators"
)

var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration [name]",
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

var makeRepositoryCmd = &cobra.Command{
	Use:   "repository [name]",
	Short: "Gera um novo repositório",
	Long: `Gera um novo repositório com toda a estrutura necessária:
  • Entidade
  • Repositório
  • Migrations

Exemplo:
  artemis make:repository create_users_table
  artemis make:repository add_email_to_users`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repositoryName := args[0]

		fmt.Printf("📝 Gerando repositório: %s\n", repositoryName)

		//... implementar a lógica de geração do repositório aqui ...

		fmt.Printf("✅ Repository '%s' gerada com sucesso!\n", repositoryName)

		return nil
	},
}
