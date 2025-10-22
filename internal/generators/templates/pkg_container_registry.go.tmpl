package container

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// ModuleRegistry é o registro central para componentes dos módulos
type ModuleRegistry struct {
	container *Container

	// Handlers HTTP
	httpHandlers map[string]HTTPHandler

	// Handlers gRPC
	grpcServices []GRPCServiceRegistrar

	// Repositórios
	repositories map[string]interface{}

	// Application Services
	appServices map[string]interface{}

	// Event Subscribers
	eventSubscribers []EventSubscriber

	mu sync.RWMutex
}

// HTTPHandler representa um handler HTTP que pode ser registrado
type HTTPHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
}

// GRPCServiceRegistrar representa um serviço gRPC que pode ser registrado
type GRPCServiceRegistrar interface {
	RegisterService(server *grpc.Server)
}

// EventSubscriber representa um subscriber de eventos
type EventSubscriber interface {
	Subscribe(eventBus interface{}) error
}

// NewModuleRegistry cria um novo registry de módulos
func NewModuleRegistry(container *Container) *ModuleRegistry {
	return &ModuleRegistry{
		container:        container,
		httpHandlers:     make(map[string]HTTPHandler),
		grpcServices:     make([]GRPCServiceRegistrar, 0),
		repositories:     make(map[string]interface{}),
		appServices:      make(map[string]interface{}),
		eventSubscribers: make([]EventSubscriber, 0),
	}
}

// Container retorna o container DI
func (r *ModuleRegistry) Container() *Container {
	return r.container
}

// RegisterHTTPHandler registra um handler HTTP
func (r *ModuleRegistry) RegisterHTTPHandler(name string, handler HTTPHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.httpHandlers[name] = handler
	r.container.Register(fmt.Sprintf("http.handler.%s", name), handler)
}

// RegisterGRPCService registra um serviço gRPC
func (r *ModuleRegistry) RegisterGRPCService(service GRPCServiceRegistrar) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.grpcServices = append(r.grpcServices, service)
}

// RegisterRepository registra um repositório
func (r *ModuleRegistry) RegisterRepository(name string, repo interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.repositories[name] = repo
	r.container.Register(fmt.Sprintf("repository.%s", name), repo)
}

// RegisterApplicationService registra um application service
func (r *ModuleRegistry) RegisterApplicationService(name string, service interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.appServices[name] = service
	r.container.Register(fmt.Sprintf("app.service.%s", name), service)
}

// RegisterEventSubscriber registra um subscriber de eventos
func (r *ModuleRegistry) RegisterEventSubscriber(subscriber EventSubscriber) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.eventSubscribers = append(r.eventSubscribers, subscriber)
}

// GetHTTPHandlers retorna todos os handlers HTTP registrados
func (r *ModuleRegistry) GetHTTPHandlers() map[string]HTTPHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Retorna uma cópia para segurança
	handlers := make(map[string]HTTPHandler, len(r.httpHandlers))
	for k, v := range r.httpHandlers {
		handlers[k] = v
	}
	return handlers
}

// GetGRPCServices retorna todos os serviços gRPC registrados
func (r *ModuleRegistry) GetGRPCServices() []GRPCServiceRegistrar {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Retorna uma cópia para segurança
	services := make([]GRPCServiceRegistrar, len(r.grpcServices))
	copy(services, r.grpcServices)
	return services
}

// GetRepository obtém um repositório por nome
func (r *ModuleRegistry) GetRepository(name string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repo, exists := r.repositories[name]
	if !exists {
		return nil, fmt.Errorf("repository '%s' not found", name)
	}
	return repo, nil
}

// GetApplicationService obtém um application service por nome
func (r *ModuleRegistry) GetApplicationService(name string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.appServices[name]
	if !exists {
		return nil, fmt.Errorf("application service '%s' not found", name)
	}
	return service, nil
}

// InitializeEventSubscribers inicializa todos os subscribers de eventos
func (r *ModuleRegistry) InitializeEventSubscribers(ctx context.Context, eventBus interface{}) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, subscriber := range r.eventSubscribers {
		if err := subscriber.Subscribe(eventBus); err != nil {
			return fmt.Errorf("failed to subscribe events: %w", err)
		}
	}
	return nil
}

// RegisterHTTPRoutes registra todas as rotas HTTP no router
func (r *ModuleRegistry) RegisterHTTPRoutes(router *gin.RouterGroup) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for name, handler := range r.httpHandlers {
		fmt.Printf("Registering HTTP routes for module: %s\n", name)
		handler.RegisterRoutes(router)
	}
}

// RegisterGRPCServices registra todos os serviços gRPC no servidor
func (r *ModuleRegistry) RegisterGRPCServices(server *grpc.Server) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, service := range r.grpcServices {
		service.RegisterService(server)
	}
}

// Stats retorna estatísticas do registry
func (r *ModuleRegistry) Stats() RegistryStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return RegistryStats{
		HTTPHandlers:     len(r.httpHandlers),
		GRPCServices:     len(r.grpcServices),
		Repositories:     len(r.repositories),
		AppServices:      len(r.appServices),
		EventSubscribers: len(r.eventSubscribers),
	}
}

// RegistryStats contém estatísticas do registry
type RegistryStats struct {
	HTTPHandlers     int
	GRPCServices     int
	Repositories     int
	AppServices      int
	EventSubscribers int
}

// String retorna uma representação em string das estatísticas
func (s RegistryStats) String() string {
	return fmt.Sprintf(
		"Registry Stats: HTTP Handlers=%d, gRPC Services=%d, Repositories=%d, App Services=%d, Event Subscribers=%d",
		s.HTTPHandlers, s.GRPCServices, s.Repositories, s.AppServices, s.EventSubscribers,
	)
}
