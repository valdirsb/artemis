package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "artemis",
	Short: "🏹 Artemis - A modular Go framework",
	Long: `
🏹 Artemis Go Framework

Build monoliths, modularly.

Artemis é um framework modular para Go que traz uma estrutura inspirada no Laravel,
com foco em produtividade e organização de código.

Recursos:
  • Estrutura modular de alto desempenho
  • CLI poderoso para geração de código
  • Suporte a múltiplas tecnologias (REST, gRPC, GraphQL, etc)
  • Facilidade de configuração e extensibilidade`,
	Version: "1.0.0",
}

// Execute executa o comando root
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Adiciona comandos
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(makeModuleCmd)
	// rootCmd.AddCommand(makeMigrationCmd)
}
