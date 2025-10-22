# 🧪 Testing Strategy

## 📋 Índice
- [Visão Geral](#visão-geral)
- [Tipos de Testes](#tipos-de-testes)
- [Testes Unitários](#testes-unitários)
- [Testes de Integração](#testes-de-integração)
- [Testes E2E](#testes-e2e)
- [Mocks e Stubs](#mocks-e-stubs)
- [Test Coverage](#test-coverage)
- [Boas Práticas](#boas-práticas)
- [Exemplos Práticos](#exemplos-práticos)

---

## 🎯 Visão Geral

O Artemis Framework adota uma **estratégia de testes em camadas**, garantindo qualidade e confiabilidade em todos os níveis da aplicação.

```
┌────────────────────────────────────────────────┐
│           E2E Tests (10%)                      │
│   • Testa fluxo completo                       │
│   • HTTP/gRPC endpoints                        │
│   • Database real                              │
└──────────────────┬─────────────────────────────┘
                   │
┌────────────────────────────────────────────────┐
│      Integration Tests (20%)                   │
│   • Testa componentes integrados               │
│   • Database (SQLite in-memory)                │
│   • Repositórios + Domain                      │
└──────────────────┬─────────────────────────────┘
                   │
┌────────────────────────────────────────────────┐
│         Unit Tests (70%)                       │
│   • Testa unidades isoladas                    │
│   • Commands, Queries, Domain                  │
│   • Mocks para dependências                    │
└────────────────────────────────────────────────┘
```

### Pirâmide de Testes

```
           /\
          /E2\      ← Poucos testes, lentos, alto valor
         /____\
        / INT  \    ← Testes médios, velocidade média
       /________\
      /   UNIT   \  ← Muitos testes, rápidos, feedback imediato
     /____________\
```

---

## 📝 Tipos de Testes

### 1. Unit Tests (70%)

**O que testar:**
- ✅ Domain entities e value objects
- ✅ Commands handlers (com mocks)
- ✅ Query handlers (com mocks)
- ✅ Business logic
- ✅ Validações

**Não testar:**
- ❌ Database queries
- ❌ External APIs
- ❌ Framework code

### 2. Integration Tests (20%)

**O que testar:**
- ✅ Repositories com database real (SQLite)
- ✅ Commands + Repositories
- ✅ Queries + Repositories
- ✅ Event handlers com EventBus

**Não testar:**
- ❌ HTTP handlers (usar E2E)
- ❌ gRPC services (usar E2E)

### 3. E2E Tests (10%)

**O que testar:**
- ✅ Fluxos completos de negócio
- ✅ API REST endpoints
- ✅ gRPC services
- ✅ Autenticação/Autorização

---

## 🔬 Testes Unitários

### Setup Básico

**`internal/modules/product/application/commands/create_product_test.go`**

```go
package commands

import (
    "context"
    "errors"
    "testing"

    "meuApp/internal/modules/product/domain"
    "meuApp/internal/modules/product/ports"
    "meuApp/pkg/contracts"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

// Mock do ProductRepository
type MockProductRepository struct {
    mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
    args := m.Called(ctx, product)
    return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Product), args.Error(1)
}

// Mock do EventPublisher
type MockEventPublisher struct {
    mock.Mock
}

func (m *MockEventPublisher) Publish(ctx context.Context, event contracts.Event) error {
    args := m.Called(ctx, event)
    return args.Error(0)
}

// Mock do Logger
type MockLogger struct {
    mock.Mock
}

func (m *MockLogger) Info(msg string, fields ...contracts.Field) {
    m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...contracts.Field) {
    m.Called(msg, fields)
}

func (m *MockLogger) Debug(msg string, fields ...contracts.Field) {
    m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...contracts.Field) {
    m.Called(msg, fields)
}

// Tests

func TestCreateProductHandler_Handle_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockProductRepository)
    mockEventBus := new(MockEventPublisher)
    mockLogger := new(MockLogger)

    handler := NewCreateProductHandler(mockRepo, mockEventBus, mockLogger)

    cmd := CreateProductCommand{
        Name:        "Laptop",
        Description: "High-end laptop",
        CategoryID:  "electronics",
        Price:       999.99,
        Stock:       10,
    }

    // Expectations
    mockLogger.On("Info", mock.Anything, mock.Anything).Return()
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)
    mockEventBus.On("Publish", mock.Anything, mock.AnythingOfType("contracts.Event")).Return(nil)

    // Act
    result, err := handler.Handle(context.Background(), cmd)

    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Laptop", result.Name())
    assert.Equal(t, 999.99, result.Price())
    assert.Equal(t, 10, result.Stock())

    // Verify mocks
    mockRepo.AssertExpectations(t)
    mockEventBus.AssertExpectations(t)
    mockLogger.AssertExpectations(t)
}

func TestCreateProductHandler_Handle_InvalidName(t *testing.T) {
    // Arrange
    mockRepo := new(MockProductRepository)
    mockEventBus := new(MockEventPublisher)
    mockLogger := new(MockLogger)

    handler := NewCreateProductHandler(mockRepo, mockEventBus, mockLogger)

    cmd := CreateProductCommand{
        Name:        "", // Nome vazio
        Description: "High-end laptop",
        CategoryID:  "electronics",
        Price:       999.99,
        Stock:       10,
    }

    // Act
    result, err := handler.Handle(context.Background(), cmd)

    // Assert
    require.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "name")
}

func TestCreateProductHandler_Handle_NegativePrice(t *testing.T) {
    // Arrange
    mockRepo := new(MockProductRepository)
    mockEventBus := new(MockEventPublisher)
    mockLogger := new(MockLogger)

    handler := NewCreateProductHandler(mockRepo, mockEventBus, mockLogger)

    cmd := CreateProductCommand{
        Name:        "Laptop",
        Description: "High-end laptop",
        CategoryID:  "electronics",
        Price:       -100.00, // Preço negativo
        Stock:       10,
    }

    // Act
    result, err := handler.Handle(context.Background(), cmd)

    // Assert
    require.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "price")
}

func TestCreateProductHandler_Handle_RepositoryError(t *testing.T) {
    // Arrange
    mockRepo := new(MockProductRepository)
    mockEventBus := new(MockEventPublisher)
    mockLogger := new(MockLogger)

    handler := NewCreateProductHandler(mockRepo, mockEventBus, mockLogger)

    cmd := CreateProductCommand{
        Name:        "Laptop",
        Description: "High-end laptop",
        CategoryID:  "electronics",
        Price:       999.99,
        Stock:       10,
    }

    // Expectations
    mockLogger.On("Info", mock.Anything, mock.Anything).Return()
    mockLogger.On("Error", mock.Anything, mock.Anything).Return()
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).
        Return(errors.New("database error"))

    // Act
    result, err := handler.Handle(context.Background(), cmd)

    // Assert
    require.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "failed to save product")

    mockRepo.AssertExpectations(t)
    mockLogger.AssertExpectations(t)
}
```

### Table-Driven Tests

**`internal/modules/product/domain/product_test.go`**

```go
package domain

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewProduct(t *testing.T) {
    tests := []struct {
        name        string
        productName string
        description string
        categoryID  string
        price       float64
        stock       int
        wantErr     bool
        errContains string
    }{
        {
            name:        "valid product",
            productName: "Laptop",
            description: "High-end laptop",
            categoryID:  "electronics",
            price:       999.99,
            stock:       10,
            wantErr:     false,
        },
        {
            name:        "empty name",
            productName: "",
            description: "High-end laptop",
            categoryID:  "electronics",
            price:       999.99,
            stock:       10,
            wantErr:     true,
            errContains: "name",
        },
        {
            name:        "negative price",
            productName: "Laptop",
            description: "High-end laptop",
            categoryID:  "electronics",
            price:       -100.00,
            stock:       10,
            wantErr:     true,
            errContains: "price",
        },
        {
            name:        "negative stock",
            productName: "Laptop",
            description: "High-end laptop",
            categoryID:  "electronics",
            price:       999.99,
            stock:       -5,
            wantErr:     true,
            errContains: "stock",
        },
        {
            name:        "empty category",
            productName: "Laptop",
            description: "High-end laptop",
            categoryID:  "",
            price:       999.99,
            stock:       10,
            wantErr:     true,
            errContains: "category",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act
            product, err := NewProduct(
                "test-id",
                tt.productName,
                tt.description,
                tt.categoryID,
                tt.price,
                tt.stock,
            )

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errContains)
                assert.Nil(t, product)
            } else {
                require.NoError(t, err)
                assert.NotNil(t, product)
                assert.Equal(t, tt.productName, product.Name())
                assert.Equal(t, tt.price, product.Price())
                assert.Equal(t, tt.stock, product.Stock())
            }
        })
    }
}

func TestProduct_DecreaseStock(t *testing.T) {
    tests := []struct {
        name          string
        initialStock  int
        decreaseBy    int
        wantErr       bool
        expectedStock int
        errContains   string
    }{
        {
            name:          "successful decrease",
            initialStock:  10,
            decreaseBy:    3,
            wantErr:       false,
            expectedStock: 7,
        },
        {
            name:          "decrease to zero",
            initialStock:  5,
            decreaseBy:    5,
            wantErr:       false,
            expectedStock: 0,
        },
        {
            name:         "insufficient stock",
            initialStock: 5,
            decreaseBy:   10,
            wantErr:      true,
            errContains:  "insufficient stock",
        },
        {
            name:         "negative quantity",
            initialStock: 10,
            decreaseBy:   -5,
            wantErr:      true,
            errContains:  "quantity must be positive",
        },
        {
            name:         "zero quantity",
            initialStock: 10,
            decreaseBy:   0,
            wantErr:      true,
            errContains:  "quantity must be positive",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            product, err := NewProduct(
                "test-id",
                "Test Product",
                "Description",
                "category-1",
                99.99,
                tt.initialStock,
            )
            require.NoError(t, err)

            // Act
            err = product.DecreaseStock(tt.decreaseBy)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errContains)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.expectedStock, product.Stock())
            }
        })
    }
}
```

---

## 🔗 Testes de Integração

### Repository Integration Tests

**`internal/modules/product/tests/integration/repository_test.go`**

```go
package integration

import (
    "context"
    "testing"
    "time"

    "meuApp/internal/modules/product/domain"
    "meuApp/internal/modules/product/ports"
    "meuApp/internal/modules/product/repository"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// setupTestDB cria um banco de dados SQLite em memória para testes
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    require.NoError(t, err)

    // Auto migrate
    err = db.AutoMigrate(&repository.ProductModel{})
    require.NoError(t, err)

    return db
}

// createTestProduct cria um produto de teste
func createTestProduct(name, categoryID string, price float64, stock int) *domain.Product {
    product, err := domain.NewProduct(
        generateID(),
        name,
        "Test description for "+name,
        categoryID,
        price,
        stock,
    )
    if err != nil {
        panic("Failed to create test product: " + err.Error())
    }
    return product
}

// generateID gera um ID simples para testes
func generateID() string {
    return "test-" + time.Now().Format("20060102150405.999999999")
}

func TestProductRepository_Create(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := repository.NewMySQLProductRepository(db)
    ctx := context.Background()

    product := createTestProduct("Laptop", "electronics", 999.99, 10)

    // Act
    err := repo.Create(ctx, product)

    // Assert
    require.NoError(t, err)

    // Verify product was created in database
    var count int64
    db.Model(&repository.ProductModel{}).Count(&count)
    assert.Equal(t, int64(1), count)
}

func TestProductRepository_GetByID_Success(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := repository.NewMySQLProductRepository(db)
    ctx := context.Background()

    product := createTestProduct("Laptop", "electronics", 999.99, 10)
    err := repo.Create(ctx, product)
    require.NoError(t, err)

    // Act
    found, err := repo.GetByID(ctx, product.ID())

    // Assert
    require.NoError(t, err)
    assert.NotNil(t, found)
    assert.Equal(t, product.ID(), found.ID())
    assert.Equal(t, product.Name(), found.Name())
    assert.Equal(t, product.Price(), found.Price())
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := repository.NewMySQLProductRepository(db)
    ctx := context.Background()

    // Act
    found, err := repo.GetByID(ctx, "non-existent-id")

    // Assert
    require.Error(t, err)
    assert.Nil(t, found)
}

func TestProductRepository_Update(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := repository.NewMySQLProductRepository(db)
    ctx := context.Background()

    product := createTestProduct("Laptop", "electronics", 999.99, 10)
    err := repo.Create(ctx, product)
    require.NoError(t, err)

    // Act - Update price and stock
    err = product.SetPrice(1099.99)
    require.NoError(t, err)
    err = product.DecreaseStock(3)
    require.NoError(t, err)

    err = repo.Update(ctx, product)
    require.NoError(t, err)

    // Assert - Verify updates
    found, err := repo.GetByID(ctx, product.ID())
    require.NoError(t, err)
    assert.Equal(t, 1099.99, found.Price())
    assert.Equal(t, 7, found.Stock())
}

func TestProductRepository_Delete(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := repository.NewMySQLProductRepository(db)
    ctx := context.Background()

    product := createTestProduct("Laptop", "electronics", 999.99, 10)
    err := repo.Create(ctx, product)
    require.NoError(t, err)

    // Act
    err = repo.Delete(ctx, product.ID())

    // Assert
    require.NoError(t, err)

    // Verify deletion
    found, err := repo.GetByID(ctx, product.ID())
    assert.Error(t, err)
    assert.Nil(t, found)
}

func TestProductRepository_List(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := repository.NewMySQLProductRepository(db)
    ctx := context.Background()

    // Create multiple products
    product1 := createTestProduct("Laptop", "electronics", 999.99, 10)
    product2 := createTestProduct("Mouse", "electronics", 29.99, 50)
    product3 := createTestProduct("Keyboard", "electronics", 79.99, 30)

    err := repo.Create(ctx, product1)
    require.NoError(t, err)
    err = repo.Create(ctx, product2)
    require.NoError(t, err)
    err = repo.Create(ctx, product3)
    require.NoError(t, err)

    // Act
    products, err := repo.List(ctx, ports.ProductFilters{})

    // Assert
    require.NoError(t, err)
    assert.Len(t, products, 3)
}

func TestProductRepository_List_WithFilters(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    repo := repository.NewMySQLProductRepository(db)
    ctx := context.Background()

    // Create products in different categories
    product1 := createTestProduct("Laptop", "electronics", 999.99, 10)
    product2 := createTestProduct("Mouse", "electronics", 29.99, 50)
    product3 := createTestProduct("Desk", "furniture", 299.99, 5)

    err := repo.Create(ctx, product1)
    require.NoError(t, err)
    err = repo.Create(ctx, product2)
    require.NoError(t, err)
    err = repo.Create(ctx, product3)
    require.NoError(t, err)

    t.Run("filter by category", func(t *testing.T) {
        categoryID := "electronics"
        filters := ports.ProductFilters{
            CategoryID: &categoryID,
        }

        products, err := repo.List(ctx, filters)
        require.NoError(t, err)
        assert.Len(t, products, 2)
    })

    t.Run("filter by min price", func(t *testing.T) {
        minPrice := 100.0
        filters := ports.ProductFilters{
            MinPrice: &minPrice,
        }

        products, err := repo.List(ctx, filters)
        require.NoError(t, err)
        assert.Len(t, products, 2) // Laptop and Desk
    })

    t.Run("filter by max price", func(t *testing.T) {
        maxPrice := 100.0
        filters := ports.ProductFilters{
            MaxPrice: &maxPrice,
        }

        products, err := repo.List(ctx, filters)
        require.NoError(t, err)
        assert.Len(t, products, 2) // Mouse and (79.99 if exists)
    })
}
```

### Pagination Integration Tests

**`internal/modules/product/tests/integration/pagination_test.go`**

```go
package integration

import (
    "context"
    "fmt"
    "testing"

    "meuApp/internal/modules/product/ports"
    "meuApp/internal/modules/product/repository"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestProductRepository_ListPaginated(t *testing.T) {
    db := setupTestDB(t)
    repo := repository.NewMySQLProductRepository(db)
    ctx := context.Background()

    // Criar 25 produtos de teste
    for i := 1; i <= 25; i++ {
        product := createTestProduct(
            fmt.Sprintf("Product %d", i),
            "cat-test",
            float64(i*10),
            i,
        )
        err := repo.Create(ctx, product)
        require.NoError(t, err)
    }

    t.Run("primeira página com 10 itens", func(t *testing.T) {
        filters := ports.ProductFilters{
            Page:     1,
            PageSize: 10,
        }

        result, err := repo.ListPaginated(ctx, filters)
        require.NoError(t, err)
        assert.NotNil(t, result)
        assert.Equal(t, 10, len(result.Items))
        assert.Equal(t, int64(25), result.TotalItems)
        assert.Equal(t, 1, result.Page)
        assert.Equal(t, 10, result.PageSize)
        assert.Equal(t, 3, result.TotalPages)
    })

    t.Run("última página com 5 itens", func(t *testing.T) {
        filters := ports.ProductFilters{
            Page:     3,
            PageSize: 10,
        }

        result, err := repo.ListPaginated(ctx, filters)
        require.NoError(t, err)
        assert.Equal(t, 5, len(result.Items))
        assert.Equal(t, int64(25), result.TotalItems)
        assert.Equal(t, 3, result.Page)
    })

    t.Run("página inexistente retorna vazio", func(t *testing.T) {
        filters := ports.ProductFilters{
            Page:     10,
            PageSize: 10,
        }

        result, err := repo.ListPaginated(ctx, filters)
        require.NoError(t, err)
        assert.Equal(t, 0, len(result.Items))
        assert.Equal(t, int64(25), result.TotalItems)
    })
}
```

---

## 🌐 Testes E2E

### HTTP API E2E Tests

**`tests/e2e/product_api_test.go`**

```go
package e2e

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "meuApp/internal/bootstrap"
    "meuApp/internal/modules/product/dto"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) *gin.Engine {
    // Bootstrap application
    _, registry, _, err := bootstrap.FrameworkBootstrapWithRegistry("../../framework.yaml")
    require.NoError(t, err)

    // Setup router
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    api := router.Group("/api/v1")
    registry.RegisterHTTPRoutes(api)

    return router
}

func TestProductAPI_CreateProduct_E2E(t *testing.T) {
    // Arrange
    router := setupTestApp(t)

    requestBody := dto.CreateProductRequest{
        Name:        "Laptop E2E Test",
        Description: "Integration test product",
        CategoryID:  "electronics",
        Price:       999.99,
        Stock:       10,
    }

    body, err := json.Marshal(requestBody)
    require.NoError(t, err)

    req, err := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
    require.NoError(t, err)
    req.Header.Set("Content-Type", "application/json")

    // Act
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)

    var response dto.ProductResponse
    err = json.Unmarshal(w.Body.Bytes(), &response)
    require.NoError(t, err)

    assert.NotEmpty(t, response.ID)
    assert.Equal(t, "Laptop E2E Test", response.Name)
    assert.Equal(t, 999.99, response.Price)
}

func TestProductAPI_GetProduct_E2E(t *testing.T) {
    // Arrange
    router := setupTestApp(t)

    // Create product first
    createReq := dto.CreateProductRequest{
        Name:        "Mouse E2E Test",
        Description: "Test mouse",
        CategoryID:  "electronics",
        Price:       29.99,
        Stock:       50,
    }

    body, _ := json.Marshal(createReq)
    req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    var createResponse dto.ProductResponse
    json.Unmarshal(w.Body.Bytes(), &createResponse)

    // Act - Get product
    req, err := http.NewRequest("GET", "/api/v1/products/"+createResponse.ID, nil)
    require.NoError(t, err)

    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)

    var getResponse dto.ProductResponse
    err = json.Unmarshal(w.Body.Bytes(), &getResponse)
    require.NoError(t, err)

    assert.Equal(t, createResponse.ID, getResponse.ID)
    assert.Equal(t, "Mouse E2E Test", getResponse.Name)
}

func TestProductAPI_ListProducts_E2E(t *testing.T) {
    // Arrange
    router := setupTestApp(t)

    // Create multiple products
    products := []dto.CreateProductRequest{
        {Name: "Product 1", Description: "Desc 1", CategoryID: "cat1", Price: 10.0, Stock: 5},
        {Name: "Product 2", Description: "Desc 2", CategoryID: "cat1", Price: 20.0, Stock: 10},
        {Name: "Product 3", Description: "Desc 3", CategoryID: "cat2", Price: 30.0, Stock: 15},
    }

    for _, p := range products {
        body, _ := json.Marshal(p)
        req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
    }

    // Act
    req, err := http.NewRequest("GET", "/api/v1/products", nil)
    require.NoError(t, err)

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)

    var response []dto.ProductResponse
    err = json.Unmarshal(w.Body.Bytes(), &response)
    require.NoError(t, err)

    assert.GreaterOrEqual(t, len(response), 3)
}

func TestProductAPI_UpdateProduct_E2E(t *testing.T) {
    // Arrange
    router := setupTestApp(t)

    // Create product
    createReq := dto.CreateProductRequest{
        Name:        "Original Name",
        Description: "Original desc",
        CategoryID:  "electronics",
        Price:       99.99,
        Stock:       5,
    }

    body, _ := json.Marshal(createReq)
    req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    var createResponse dto.ProductResponse
    json.Unmarshal(w.Body.Bytes(), &createResponse)

    // Act - Update product
    updateReq := dto.UpdateProductRequest{
        Name:        stringPtr("Updated Name"),
        Description: stringPtr("Updated desc"),
        Price:       float64Ptr(149.99),
    }

    body, _ = json.Marshal(updateReq)
    req, err := http.NewRequest("PUT", "/api/v1/products/"+createResponse.ID, bytes.NewBuffer(body))
    require.NoError(t, err)
    req.Header.Set("Content-Type", "application/json")

    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)

    var updateResponse dto.ProductResponse
    err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
    require.NoError(t, err)

    assert.Equal(t, "Updated Name", updateResponse.Name)
    assert.Equal(t, 149.99, updateResponse.Price)
}

func TestProductAPI_DeleteProduct_E2E(t *testing.T) {
    // Arrange
    router := setupTestApp(t)

    // Create product
    createReq := dto.CreateProductRequest{
        Name:        "To Delete",
        Description: "Will be deleted",
        CategoryID:  "electronics",
        Price:       99.99,
        Stock:       5,
    }

    body, _ := json.Marshal(createReq)
    req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    var createResponse dto.ProductResponse
    json.Unmarshal(w.Body.Bytes(), &createResponse)

    // Act - Delete product
    req, err := http.NewRequest("DELETE", "/api/v1/products/"+createResponse.ID, nil)
    require.NoError(t, err)

    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusNoContent, w.Code)

    // Verify deletion
    req, _ = http.NewRequest("GET", "/api/v1/products/"+createResponse.ID, nil)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusNotFound, w.Code)
}

// Helper functions
func stringPtr(s string) *string {
    return &s
}

func float64Ptr(f float64) *float64 {
    return &f
}
```

---

## 🎭 Mocks e Stubs

### Usando testify/mock

**Instalação:**

```bash
go get github.com/stretchr/testify/mock
```

### Mock Generator

**Interface:**

```go
// internal/modules/product/ports/ports.go
type ProductRepository interface {
    Create(ctx context.Context, product *domain.Product) error
    GetByID(ctx context.Context, id string) (*domain.Product, error)
    Update(ctx context.Context, product *domain.Product) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filters ProductFilters) ([]*domain.Product, error)
}
```

**Mock gerado:**

```go
// internal/modules/product/ports/mocks/product_repository_mock.go
package mocks

import (
    "context"

    "meuApp/internal/modules/product/domain"
    "meuApp/internal/modules/product/ports"

    "github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
    mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
    args := m.Called(ctx, product)
    return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
    args := m.Called(ctx, product)
    return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id string) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockProductRepository) List(ctx context.Context, filters ports.ProductFilters) ([]*domain.Product, error) {
    args := m.Called(ctx, filters)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*domain.Product), args.Error(1)
}
```

---

## 📊 Test Coverage

### Makefile Commands

**`Makefile`**

```makefile
.PHONY: test test-unit test-integration test-e2e test-coverage

## test: Executa todos os testes
test:
	@echo "🧪 Running all tests..."
	@go test -v -race ./...

## test-unit: Executa testes unitários
test-unit:
	@echo "🔬 Running unit tests..."
	@go test -v -short ./...

## test-integration: Executa testes de integração
test-integration:
	@echo "🔗 Running integration tests..."
	@go test -v -run Integration ./...

## test-e2e: Executa testes E2E
test-e2e:
	@echo "🌐 Running E2E tests..."
	@go test -v ./tests/e2e/...

## test-coverage: Gera relatório de cobertura
test-coverage:
	@echo "📊 Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

## test-coverage-unit: Gera relatório de cobertura unitária
test-coverage-unit:
	@echo "📊 Generating unit test coverage..."
	@go test -short -coverprofile=coverage-unit.out ./...
	@go tool cover -html=coverage-unit.out -o coverage-unit.html
	@echo "✅ Unit coverage report: coverage-unit.html"

## test-watch: Executa testes em modo watch
test-watch:
	@which gotestsum > /dev/null || go install gotest.tools/gotestsum@latest
	@gotestsum --watch
```

### Comandos:

```bash
# Todos os testes
make test

# Testes unitários apenas
make test-unit

# Testes de integração
make test-integration

# Testes E2E
make test-e2e

# Coverage report
make test-coverage

# Watch mode
make test-watch
```

### Coverage Goals

```
Target Coverage:
- Unit Tests: 80%+
- Integration Tests: 70%+
- Overall: 75%+
```

---

## ✨ Boas Práticas

### 1. Arrange-Act-Assert (AAA)

```go
func TestExample(t *testing.T) {
    // Arrange - Setup
    user := createTestUser()
    
    // Act - Execute
    result := user.GetFullName()
    
    // Assert - Verify
    assert.Equal(t, "John Doe", result)
}
```

### 2. Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid", "test@example.com", false},
        {"invalid", "not-an-email", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 3. Test Names

```go
// ✅ Bom
func TestCreateProductHandler_Handle_Success(t *testing.T) {}
func TestCreateProductHandler_Handle_InvalidPrice(t *testing.T) {}

// ❌ Ruim
func TestProduct(t *testing.T) {}
func TestHandler1(t *testing.T) {}
```

### 4. Use Subtests

```go
func TestUserService(t *testing.T) {
    t.Run("create user success", func(t *testing.T) {
        // Test logic
    })
    
    t.Run("create user duplicate email", func(t *testing.T) {
        // Test logic
    })
}
```

### 5. Test Fixtures

```go
// tests/fixtures/users.go
package fixtures

func CreateTestUser(name, email string) *domain.User {
    user, _ := domain.NewUser("test-id", name, email)
    return user
}

func CreateTestProducts(count int) []*domain.Product {
    products := make([]*domain.Product, count)
    for i := 0; i < count; i++ {
        products[i] = CreateTestProduct(fmt.Sprintf("Product %d", i+1))
    }
    return products
}
```

---

## 📚 Próximos Passos

- **[Mocks & Fixtures](22-mocks-fixtures.md)** - Mocks e fixtures detalhados
- **[Best Practices](23-best-practices.md)** - Boas práticas gerais
- **[API Reference](24-api-reference.md)** - Documentação de APIs

---

**[⬅️ Error Handling](19-error-handling.md)** | **[Índice](README.md)** | **[API Reference ➡️](24-api-reference.md)**
