package events

import (
	"context"
	"fmt"
	"log"
)

// SimpleLogger é uma interface simples para logging
type SimpleLogger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

// defaultSimpleLogger implementação padrão
type defaultSimpleLogger struct{}

func (l *defaultSimpleLogger) Info(msg string) {
	log.Println("[INFO]", msg)
}

func (l *defaultSimpleLogger) Warn(msg string) {
	log.Println("[WARN]", msg)
}

func (l *defaultSimpleLogger) Error(msg string) {
	log.Println("[ERROR]", msg)
}

// UserCreatedHandlerFunc processa eventos de criação de usuário
func UserCreatedHandlerFunc(logger SimpleLogger) func(ctx context.Context, data UserCreatedEvent) error {
	if logger == nil {
		logger = &defaultSimpleLogger{}
	}
	return func(ctx context.Context, data UserCreatedEvent) error {
		logger.Info(fmt.Sprintf("Processing UserCreatedEvent for user: %s (%s)",
			data.Username, data.Email))

		// Aqui você pode:
		// - Enviar email de boas-vindas
		// - Criar perfil em outros sistemas
		// - Enviar notificação
		// - Registrar em analytics

		return nil
	}
}

// LowStockHandlerFunc processa eventos de estoque baixo
func LowStockHandlerFunc(logger SimpleLogger) func(ctx context.Context, data LowStockEvent) error {
	if logger == nil {
		logger = &defaultSimpleLogger{}
	}
	return func(ctx context.Context, data LowStockEvent) error {
		logger.Warn(fmt.Sprintf("⚠️  LOW STOCK ALERT: Product %s has only %d units left!",
			data.ProductName, data.CurrentStock))

		// Aqui você pode:
		// - Enviar notificação para admin
		// - Criar ordem de compra automática
		// - Alertar fornecedores
		// - Atualizar dashboard

		return nil
	}
}

// OrderCreatedHandlerFunc processa eventos de criação de pedido
func OrderCreatedHandlerFunc(logger SimpleLogger) func(ctx context.Context, data OrderCreatedEvent) error {
	if logger == nil {
		logger = &defaultSimpleLogger{}
	}
	return func(ctx context.Context, data OrderCreatedEvent) error {
		logger.Info(fmt.Sprintf("📦 New order created: %s (Total: $%.2f, Items: %d)",
			data.OrderID, data.TotalAmount, data.ItemCount))

		// Aqui você pode:
		// - Enviar email de confirmação
		// - Notificar sistema de estoque
		// - Atualizar métricas
		// - Iniciar processo de pagamento

		return nil
	}
}

// AuditLogHandlerFunc registra eventos para auditoria
func AuditLogHandlerFunc[T any]() func(ctx context.Context, data T) error {
	return func(ctx context.Context, data T) error {
		log.Printf("[AUDIT] Event data: %+v", data)

		// Aqui você pode:
		// - Salvar em banco de dados de auditoria
		// - Enviar para sistema de compliance
		// - Arquivar em storage

		return nil
	}
}
