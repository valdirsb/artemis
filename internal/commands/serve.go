package commands

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia o servidor de desenvolvimento",
	Long: `Inicia o servidor de desenvolvimento na porta 8080.

Exemplo:
  artemis serve           # Inicia na porta 8080
  artemis serve --port 3000  # Inicia na porta 3000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetString("port")

		fmt.Printf("🚀 Iniciando servidor Artemis na porta %s...\n", port)

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"message": "🏹 Artemis Framework", "version": "1.0.0", "status": "running"}`)
		})

		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status": "healthy", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
		})

		server := &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		}

		// Graceful shutdown
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c

			fmt.Println("\n🛑 Parando servidor...")
			server.Close()
		}()

		fmt.Printf("📡 Servidor rodando em http://localhost:%s\n", port)
		fmt.Println("📋 Rotas disponíveis:")
		fmt.Println("   GET /        - Status da aplicação")
		fmt.Println("   GET /health  - Health check")
		fmt.Println("\n💡 Pressione Ctrl+C para parar")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("erro ao iniciar servidor: %w", err)
		}

		return nil
	},
}

func init() {
	serveCmd.Flags().StringP("port", "p", "8080", "Porta do servidor")
}
