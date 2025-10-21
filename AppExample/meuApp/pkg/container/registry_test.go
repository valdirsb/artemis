package container

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// Mock implementations for testing

type mockHTTPHandler struct {
	routesRegistered bool
}

func (m *mockHTTPHandler) RegisterRoutes(router *gin.RouterGroup) {
	m.routesRegistered = true
}

type mockGRPCService struct {
	serviceRegistered bool
}

func (m *mockGRPCService) RegisterService(server *grpc.Server) {
	m.serviceRegistered = true
}

type mockEventSubscriber struct {
	subscribed bool
}

func (m *mockEventSubscriber) Subscribe(eventBus interface{}) error {
	m.subscribed = true
	return nil
}

type mockRepository struct {
	name string
}

type mockApplicationService struct {
	name string
}

// Tests

func TestNewModuleRegistry(t *testing.T) {
	// Arrange
	container := NewContainer()

	// Act
	registry := NewModuleRegistry(container)

	// Assert
	assert.NotNil(t, registry)
	assert.NotNil(t, registry.container)
	assert.NotNil(t, registry.httpHandlers)
	assert.NotNil(t, registry.grpcServices)
	assert.NotNil(t, registry.repositories)
	assert.NotNil(t, registry.appServices)
	assert.NotNil(t, registry.eventSubscribers)

	// Verify stats are initialized correctly
	stats := registry.Stats()
	assert.Equal(t, 0, stats.HTTPHandlers)
	assert.Equal(t, 0, stats.GRPCServices)
	assert.Equal(t, 0, stats.Repositories)
	assert.Equal(t, 0, stats.AppServices)
	assert.Equal(t, 0, stats.EventSubscribers)
}

func TestModuleRegistry_Container(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)

	// Act
	result := registry.Container()

	// Assert
	assert.Equal(t, container, result)
}

func TestModuleRegistry_RegisterHTTPHandler(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	handler := &mockHTTPHandler{}

	// Act
	registry.RegisterHTTPHandler("test-handler", handler)

	// Assert
	stats := registry.Stats()
	assert.Equal(t, 1, stats.HTTPHandlers)

	// Verify handler can be retrieved
	handlers := registry.GetHTTPHandlers()
	assert.Len(t, handlers, 1)
	assert.Equal(t, handler, handlers["test-handler"])

	// Verify handler is registered in container
	resolved, err := container.Get("http.handler.test-handler")
	require.NoError(t, err)
	assert.Equal(t, handler, resolved)
}

func TestModuleRegistry_RegisterHTTPHandler_Multiple(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	handler1 := &mockHTTPHandler{}
	handler2 := &mockHTTPHandler{}

	// Act
	registry.RegisterHTTPHandler("handler1", handler1)
	registry.RegisterHTTPHandler("handler2", handler2)

	// Assert
	stats := registry.Stats()
	assert.Equal(t, 2, stats.HTTPHandlers)

	handlers := registry.GetHTTPHandlers()
	assert.Len(t, handlers, 2)
	assert.Equal(t, handler1, handlers["handler1"])
	assert.Equal(t, handler2, handlers["handler2"])
}

func TestModuleRegistry_RegisterGRPCService(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	service := &mockGRPCService{}

	// Act
	registry.RegisterGRPCService(service)

	// Assert
	stats := registry.Stats()
	assert.Equal(t, 1, stats.GRPCServices)

	services := registry.GetGRPCServices()
	assert.Len(t, services, 1)
	assert.Equal(t, service, services[0])
}

func TestModuleRegistry_RegisterGRPCService_Multiple(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	service1 := &mockGRPCService{}
	service2 := &mockGRPCService{}

	// Act
	registry.RegisterGRPCService(service1)
	registry.RegisterGRPCService(service2)

	// Assert
	stats := registry.Stats()
	assert.Equal(t, 2, stats.GRPCServices)

	services := registry.GetGRPCServices()
	assert.Len(t, services, 2)
}

func TestModuleRegistry_RegisterRepository(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	repo := &mockRepository{name: "test"}

	// Act
	registry.RegisterRepository("test-repo", repo)

	// Assert
	stats := registry.Stats()
	assert.Equal(t, 1, stats.Repositories)

	// Verify repository can be retrieved by name
	retrieved, err := registry.GetRepository("test-repo")
	require.NoError(t, err)
	assert.Equal(t, repo, retrieved)

	// Verify repository is registered in container
	resolved, err := container.Get("repository.test-repo")
	require.NoError(t, err)
	assert.Equal(t, repo, resolved)
}

func TestModuleRegistry_GetRepository_NotFound(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)

	// Act
	retrieved, err := registry.GetRepository("non-existent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Contains(t, err.Error(), "repository 'non-existent' not found")
}

func TestModuleRegistry_RegisterApplicationService(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	service := &mockApplicationService{name: "test"}

	// Act
	registry.RegisterApplicationService("test-service", service)

	// Assert
	stats := registry.Stats()
	assert.Equal(t, 1, stats.AppServices)

	// Verify service can be retrieved by name
	retrieved, err := registry.GetApplicationService("test-service")
	require.NoError(t, err)
	assert.Equal(t, service, retrieved)

	// Verify service is registered in container
	resolved, err := container.Get("app.service.test-service")
	require.NoError(t, err)
	assert.Equal(t, service, resolved)
}

func TestModuleRegistry_GetApplicationService_NotFound(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)

	// Act
	retrieved, err := registry.GetApplicationService("non-existent")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrieved)
	assert.Contains(t, err.Error(), "application service 'non-existent' not found")
}

func TestModuleRegistry_RegisterEventSubscriber(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	subscriber := &mockEventSubscriber{}

	// Act
	registry.RegisterEventSubscriber(subscriber)

	// Assert
	stats := registry.Stats()
	assert.Equal(t, 1, stats.EventSubscribers)
}

func TestModuleRegistry_InitializeEventSubscribers(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	subscriber1 := &mockEventSubscriber{}
	subscriber2 := &mockEventSubscriber{}

	registry.RegisterEventSubscriber(subscriber1)
	registry.RegisterEventSubscriber(subscriber2)

	eventBus := "mock-event-bus"
	ctx := context.Background()

	// Act
	err := registry.InitializeEventSubscribers(ctx, eventBus)

	// Assert
	require.NoError(t, err)
	assert.True(t, subscriber1.subscribed)
	assert.True(t, subscriber2.subscribed)
}

func TestModuleRegistry_RegisterHTTPRoutes(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	handler1 := &mockHTTPHandler{}
	handler2 := &mockHTTPHandler{}

	registry.RegisterHTTPHandler("handler1", handler1)
	registry.RegisterHTTPHandler("handler2", handler2)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	routerGroup := router.Group("/api")

	// Act
	registry.RegisterHTTPRoutes(routerGroup)

	// Assert
	assert.True(t, handler1.routesRegistered)
	assert.True(t, handler2.routesRegistered)
}

func TestModuleRegistry_RegisterGRPCServices(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	service1 := &mockGRPCService{}
	service2 := &mockGRPCService{}

	registry.RegisterGRPCService(service1)
	registry.RegisterGRPCService(service2)

	grpcServer := grpc.NewServer()

	// Act
	registry.RegisterGRPCServices(grpcServer)

	// Assert
	assert.True(t, service1.serviceRegistered)
	assert.True(t, service2.serviceRegistered)
}

func TestModuleRegistry_Stats(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)

	// Register various components
	registry.RegisterHTTPHandler("handler1", &mockHTTPHandler{})
	registry.RegisterHTTPHandler("handler2", &mockHTTPHandler{})
	registry.RegisterGRPCService(&mockGRPCService{})
	registry.RegisterRepository("repo1", &mockRepository{})
	registry.RegisterApplicationService("service1", &mockApplicationService{})
	registry.RegisterEventSubscriber(&mockEventSubscriber{})

	// Act
	stats := registry.Stats()

	// Assert
	assert.Equal(t, 2, stats.HTTPHandlers)
	assert.Equal(t, 1, stats.GRPCServices)
	assert.Equal(t, 1, stats.Repositories)
	assert.Equal(t, 1, stats.AppServices)
	assert.Equal(t, 1, stats.EventSubscribers)
}

func TestRegistryStats_String(t *testing.T) {
	// Arrange
	stats := RegistryStats{
		HTTPHandlers:     3,
		GRPCServices:     2,
		Repositories:     5,
		AppServices:      4,
		EventSubscribers: 1,
	}

	// Act
	result := stats.String()

	// Assert
	assert.Contains(t, result, "HTTP Handlers=3")
	assert.Contains(t, result, "gRPC Services=2")
	assert.Contains(t, result, "Repositories=5")
	assert.Contains(t, result, "App Services=4")
	assert.Contains(t, result, "Event Subscribers=1")
}

func TestModuleRegistry_GetHTTPHandlers_ReturnsCopy(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	handler := &mockHTTPHandler{}
	registry.RegisterHTTPHandler("test", handler)

	// Act
	handlers1 := registry.GetHTTPHandlers()
	handlers2 := registry.GetHTTPHandlers()

	// Assert - Verify we get different map instances (copies)
	assert.Equal(t, handlers1, handlers2)

	// Modify one copy and verify the other is not affected
	handlers1["new"] = &mockHTTPHandler{}
	assert.NotEqual(t, len(handlers1), len(handlers2))
}

func TestModuleRegistry_GetGRPCServices_ReturnsCopy(t *testing.T) {
	// Arrange
	container := NewContainer()
	registry := NewModuleRegistry(container)
	service := &mockGRPCService{}
	registry.RegisterGRPCService(service)

	// Act
	services1 := registry.GetGRPCServices()
	services2 := registry.GetGRPCServices()

	// Assert - Verify we get different slice instances (copies)
	assert.Equal(t, services1, services2)
	assert.Len(t, services1, 1)
	assert.Len(t, services2, 1)
}

// Benchmark tests

func BenchmarkModuleRegistry_RegisterHTTPHandler(b *testing.B) {
	container := NewContainer()
	registry := NewModuleRegistry(container)
	handler := &mockHTTPHandler{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		registry.RegisterHTTPHandler("test", handler)
	}
}

func BenchmarkModuleRegistry_GetHTTPHandlers(b *testing.B) {
	container := NewContainer()
	registry := NewModuleRegistry(container)
	registry.RegisterHTTPHandler("test", &mockHTTPHandler{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = registry.GetHTTPHandlers()
	}
}

func BenchmarkModuleRegistry_Stats(b *testing.B) {
	container := NewContainer()
	registry := NewModuleRegistry(container)
	registry.RegisterHTTPHandler("handler", &mockHTTPHandler{})
	registry.RegisterGRPCService(&mockGRPCService{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = registry.Stats()
	}
}
