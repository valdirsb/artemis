# 🧪 Mocks e Fixtures

<div align="center">

**Guia Completo de Mocking e Test Data para o Artemis Framework**

[Conceitos](#-conceitos) • [Mocks](#-mocks) • [Fixtures](#-fixtures) • [Builders](#-builders)

</div>

---

## 📚 Índice

1. [Visão Geral](#-visão-geral)
2. [Mocks vs Stubs vs Fakes](#-mocks-vs-stubs-vs-fakes)
3. [Mocking Manual](#-mocking-manual)
4. [Mocking com Testify](#-mocking-com-testify)
5. [Fixtures e Test Data](#-fixtures-e-test-data)
6. [Test Data Builders](#-test-data-builders)
7. [Factory Pattern](#-factory-pattern)
8. [Database Fixtures](#-database-fixtures)
9. [HTTP Mocking](#-http-mocking)
10. [gRPC Mocking](#-grpc-mocking)
11. [Melhores Práticas](#-melhores-práticas)

---

## 🎯 Visão Geral

Mocks e fixtures são essenciais para testes isolados e reproduzíveis:

✅ **Mocks** - Objetos simulados que verificam interações  
✅ **Stubs** - Objetos com respostas pré-programadas  
✅ **Fakes** - Implementações simplificadas funcionais  
✅ **Fixtures** - Dados de teste pré-configurados  
✅ **Builders** - Padrão fluente para criar test data  
✅ **Factories** - Criação centralizada de objetos de teste  

### Por Que Usar?

**Problemas sem Mocks/Fixtures:**
- ❌ Testes lentos (dependências reais)
- ❌ Testes frágeis (dependem de estado externo)
- ❌ Difícil testar edge cases
- ❌ Setup complexo e repetitivo

**Benefícios:**
- ✅ Testes rápidos e isolados
- ✅ Controle total sobre comportamento
- ✅ Fácil simular erros e edge cases
- ✅ Testes determinísticos

---

## 🔍 Mocks vs Stubs vs Fakes

### Mock

**Definição:** Objeto que **verifica interações** - quantas vezes um método foi chamado, com quais parâmetros.

**Quando usar:** Quando você quer verificar **comportamento** (se algo foi chamado).

**Exemplo:**
```go
// Mock verifica se Save foi chamado exatamente 1 vez
mockRepo := new(MockUserRepository)
mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

handler.Handle(cmd)

mockRepo.AssertExpectations(t) // ✅ Verifica que Save foi chamado
```

### Stub

**Definição:** Objeto que **retorna valores pré-programados** - não verifica interações.

**Quando usar:** Quando você só precisa de **dados de retorno**.

**Exemplo:**
```go
// Stub só retorna valor fixo
stubRepo := &StubUserRepository{
    FindByIDFunc: func(ctx context.Context, id string) (*User, error) {
        return &User{ID: id, Name: "Test"}, nil
    },
}

user, _ := stubRepo.FindByID(ctx, "123")
// Não verifica quantas vezes foi chamado
```

### Fake

**Definição:** Implementação **simplificada mas funcional** - usa estruturas in-memory em vez de DB real.

**Quando usar:** Quando você precisa de **comportamento real simplificado**.

**Exemplo:**
```go
// Fake usa map em memória
type FakeUserRepository struct {
    users map[string]*User
}

func (r *FakeUserRepository) Save(ctx context.Context, user *User) error {
    r.users[user.ID] = user
    return nil
}

func (r *FakeUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    if user, ok := r.users[id]; ok {
        return user, nil
    }
    return nil, ErrNotFound
}
```

### Comparação

| Tipo | Verifica Interações | Estado Mantido | Complexidade |
|------|---------------------|----------------|--------------|
| **Mock** | ✅ Sim | ❌ Não | Média |
| **Stub** | ❌ Não | ❌ Não | Baixa |
| **Fake** | ❌ Não | ✅ Sim | Alta |

---

## 🔧 Mocking Manual

### 1. Mock Simples com Funcs

**Vantagem:** Zero dependências externas

```go
package repositories_test

import (
    "context"
    "testing"
    "meuApp/internal/modules/user/domain"
)

type MockUserRepository struct {
    SaveFunc    func(ctx context.Context, user *domain.User) error
    FindByIDFunc func(ctx context.Context, id string) (*domain.User, error)
    FindByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
    DeleteFunc  func(ctx context.Context, id string) error
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
    if m.SaveFunc != nil {
        return m.SaveFunc(ctx, user)
    }
    return nil // Comportamento padrão
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    if m.FindByIDFunc != nil {
        return m.FindByIDFunc(ctx, id)
    }
    return nil, ErrUserNotFound
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    if m.FindByEmailFunc != nil {
        return m.FindByEmailFunc(ctx, email)
    }
    return nil, ErrUserNotFound
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
    if m.DeleteFunc != nil {
        return m.DeleteFunc(ctx, id)
    }
    return nil
}
```

**Uso em testes:**
```go
func TestCreateUserHandler_SaveError(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{
        FindByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
            return nil, ErrUserNotFound // Não existe
        },
        SaveFunc: func(ctx context.Context, user *domain.User) error {
            return errors.New("database connection lost") // Simula erro
        },
    }
    
    handler := NewCreateUserHandler(mockRepo)
    cmd := &CreateUserCommand{
        Name:  "John",
        Email: "john@example.com",
    }
    
    // Act
    err := handler.Handle(context.Background(), cmd)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "database connection lost")
}
```

### 2. Mock com Estado

**Quando usar:** Você precisa verificar chamadas ou comportamento stateful

```go
type MockUserRepositoryWithState struct {
    SaveCalled      int
    FindByIDCalled  int
    LastSavedUser   *domain.User
    LastQueriedID   string
    users           map[string]*domain.User
}

func NewMockUserRepository() *MockUserRepositoryWithState {
    return &MockUserRepositoryWithState{
        users: make(map[string]*domain.User),
    }
}

func (m *MockUserRepositoryWithState) Save(ctx context.Context, user *domain.User) error {
    m.SaveCalled++
    m.LastSavedUser = user
    m.users[user.ID] = user
    return nil
}

func (m *MockUserRepositoryWithState) FindByID(ctx context.Context, id string) (*domain.User, error) {
    m.FindByIDCalled++
    m.LastQueriedID = id
    
    if user, ok := m.users[id]; ok {
        return user, nil
    }
    return nil, ErrUserNotFound
}

// Helpers de verificação
func (m *MockUserRepositoryWithState) AssertSaveCalledTimes(t *testing.T, expected int) {
    if m.SaveCalled != expected {
        t.Errorf("Expected Save to be called %d times, but was called %d times", expected, m.SaveCalled)
    }
}
```

**Uso:**
```go
func TestCreateUserHandler_CallsRepository(t *testing.T) {
    // Arrange
    mockRepo := NewMockUserRepository()
    handler := NewCreateUserHandler(mockRepo)
    
    // Act
    handler.Handle(ctx, &CreateUserCommand{Name: "John", Email: "john@example.com"})
    
    // Assert
    mockRepo.AssertSaveCalledTimes(t, 1)
    assert.Equal(t, "John", mockRepo.LastSavedUser.Name)
}
```

---

## 🎭 Mocking com Testify

### Instalação

```bash
go get github.com/stretchr/testify
```

### Mock Básico

```go
package repositories_test

import (
    "context"
    "testing"
    
    "meuApp/internal/modules/user/domain"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/assert"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    args := m.Called(ctx, id)
    
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    args := m.Called(ctx, email)
    
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}
```

### Configurando Expectations

```go
func TestCreateUserHandler_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    handler := NewCreateUserHandler(mockRepo)
    
    cmd := &CreateUserCommand{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Expectativa 1: Verificar se email já existe (não existe)
    mockRepo.On("FindByEmail", mock.Anything, "john@example.com").
        Return(nil, ErrUserNotFound)
    
    // Expectativa 2: Salvar usuário
    mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
        return u.Name == "John Doe" && u.Email == "john@example.com"
    })).Return(nil)
    
    // Act
    err := handler.Handle(context.Background(), cmd)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t) // Verifica que tudo foi chamado
}
```

### Matchers Avançados

```go
// 1. Qualquer valor
mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

// 2. Valor específico
mockRepo.On("FindByID", mock.Anything, "user-123").Return(user, nil)

// 3. Matcher customizado
mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
    return len(u.Name) > 3 && strings.Contains(u.Email, "@")
})).Return(nil)

// 4. Tipo específico
mockRepo.On("Save", mock.AnythingOfType("*context.emptyCtx"), 
    mock.AnythingOfType("*domain.User")).Return(nil)
```

### Verificações

```go
// 1. Método chamado exatamente N vezes
mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()
mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Twice()
mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Times(5)

// 2. Método pode ser chamado 0 ou mais vezes
mockRepo.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil).Maybe()

// 3. Verificar ordem de chamadas
mockRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(nil, ErrNotFound).Once()
mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

// 4. Verificar depois
mockRepo.AssertCalled(t, "Save", mock.Anything, mock.Anything)
mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
mockRepo.AssertNumberOfCalls(t, "Save", 3)
```

### Retornos Dinâmicos

```go
// 1. Retorno baseado em argumento
mockRepo.On("FindByID", mock.Anything, mock.Anything).
    Return(func(ctx context.Context, id string) (*domain.User, error) {
        if id == "user-123" {
            return &domain.User{ID: id, Name: "John"}, nil
        }
        return nil, ErrUserNotFound
    })

// 2. Side effects
callCount := 0
mockRepo.On("Save", mock.Anything, mock.Anything).
    Run(func(args mock.Arguments) {
        callCount++
        user := args.Get(1).(*domain.User)
        fmt.Printf("Saving user: %s\n", user.Name)
    }).
    Return(nil)
```

---

## 🗂️ Fixtures e Test Data

### 1. Fixture Simples

**`testdata/users.go`**

```go
package testdata

import "meuApp/internal/modules/user/domain"

var (
    ValidUser1 = &domain.User{
        ID:    "user-001",
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    ValidUser2 = &domain.User{
        ID:    "user-002",
        Name:  "Jane Smith",
        Email: "jane@example.com",
    }
    
    AdminUser = &domain.User{
        ID:    "admin-001",
        Name:  "Admin User",
        Email: "admin@example.com",
        Role:  "admin",
    }
)
```

**Uso:**
```go
func TestSomeFeature(t *testing.T) {
    user := testdata.ValidUser1
    // Usar em teste
}
```

### 2. Fixture Functions

**Quando usar:** Você precisa de dados únicos por teste

```go
package testdata

import (
    "fmt"
    "time"
    "meuApp/internal/modules/user/domain"
)

func NewTestUser(name, email string) *domain.User {
    return &domain.User{
        ID:        GenerateTestID(),
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }
}

func GenerateTestID() string {
    return fmt.Sprintf("test-%d", time.Now().UnixNano())
}

func NewRandomUser() *domain.User {
    id := GenerateTestID()
    return &domain.User{
        ID:    id,
        Name:  "User " + id,
        Email: id + "@test.com",
    }
}
```

**Uso:**
```go
func TestUserCreation(t *testing.T) {
    user1 := testdata.NewRandomUser()
    user2 := testdata.NewRandomUser()
    
    // user1 e user2 são únicos
}
```

### 3. Fixture com Variações

```go
package testdata

func NewTestProduct(overrides ...func(*domain.Product)) *domain.Product {
    product := &domain.Product{
        ID:    GenerateTestID(),
        Name:  "Test Product",
        Price: 99.99,
        Stock: 100,
    }
    
    for _, override := range overrides {
        override(product)
    }
    
    return product
}

// Variações
func WithName(name string) func(*domain.Product) {
    return func(p *domain.Product) {
        p.Name = name
    }
}

func WithPrice(price float64) func(*domain.Product) {
    return func(p *domain.Product) {
        p.Price = price
    }
}

func OutOfStock() func(*domain.Product) {
    return func(p *domain.Product) {
        p.Stock = 0
    }
}
```

**Uso:**
```go
func TestProductVariations(t *testing.T) {
    // Produto padrão
    p1 := testdata.NewTestProduct()
    
    // Produto com nome customizado
    p2 := testdata.NewTestProduct(testdata.WithName("Custom Product"))
    
    // Produto fora de estoque
    p3 := testdata.NewTestProduct(testdata.OutOfStock())
    
    // Múltiplas variações
    p4 := testdata.NewTestProduct(
        testdata.WithName("Expensive Product"),
        testdata.WithPrice(999.99),
    )
}
```

---

## 🏗️ Test Data Builders

### Builder Pattern

**Vantagem:** API fluente, legível, reutilizável

```go
package testdata

import (
    "time"
    "meuApp/internal/modules/user/domain"
)

type UserBuilder struct {
    id        string
    name      string
    email     string
    password  string
    role      string
    active    bool
    createdAt time.Time
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        id:        GenerateTestID(),
        name:      "Test User",
        email:     "test@example.com",
        password:  "hashed_password",
        role:      "user",
        active:    true,
        createdAt: time.Now(),
    }
}

func (b *UserBuilder) WithID(id string) *UserBuilder {
    b.id = id
    return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
    b.name = name
    return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
    b.email = email
    return b
}

func (b *UserBuilder) WithRole(role string) *UserBuilder {
    b.role = role
    return b
}

func (b *UserBuilder) Inactive() *UserBuilder {
    b.active = false
    return b
}

func (b *UserBuilder) CreatedAt(t time.Time) *UserBuilder {
    b.createdAt = t
    return b
}

func (b *UserBuilder) Build() *domain.User {
    return &domain.User{
        ID:        b.id,
        Name:      b.name,
        Email:     b.email,
        Password:  b.password,
        Role:      b.role,
        Active:    b.active,
        CreatedAt: b.createdAt,
    }
}

// Helper para admin
func (b *UserBuilder) AsAdmin() *UserBuilder {
    return b.WithRole("admin").WithEmail("admin@example.com")
}
```

**Uso:**
```go
func TestUserScenarios(t *testing.T) {
    // Usuário padrão
    user1 := testdata.NewUserBuilder().Build()
    
    // Admin
    admin := testdata.NewUserBuilder().AsAdmin().Build()
    
    // Usuário inativo
    inactive := testdata.NewUserBuilder().Inactive().Build()
    
    // Usuário customizado
    custom := testdata.NewUserBuilder().
        WithName("John Doe").
        WithEmail("john@example.com").
        WithRole("moderator").
        Build()
}
```

### Builder para Entidades Complexas

```go
type OrderBuilder struct {
    id        string
    userID    string
    items     []domain.OrderItem
    status    domain.OrderStatus
    total     float64
    createdAt time.Time
}

func NewOrderBuilder() *OrderBuilder {
    return &OrderBuilder{
        id:        GenerateTestID(),
        userID:    "user-001",
        items:     []domain.OrderItem{},
        status:    domain.OrderStatusPending,
        total:     0,
        createdAt: time.Now(),
    }
}

func (b *OrderBuilder) ForUser(userID string) *OrderBuilder {
    b.userID = userID
    return b
}

func (b *OrderBuilder) AddItem(productID string, quantity int, price float64) *OrderBuilder {
    b.items = append(b.items, domain.OrderItem{
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
        Total:     float64(quantity) * price,
    })
    b.total += float64(quantity) * price
    return b
}

func (b *OrderBuilder) WithStatus(status domain.OrderStatus) *OrderBuilder {
    b.status = status
    return b
}

func (b *OrderBuilder) Confirmed() *OrderBuilder {
    return b.WithStatus(domain.OrderStatusConfirmed)
}

func (b *OrderBuilder) Paid() *OrderBuilder {
    return b.WithStatus(domain.OrderStatusPaid)
}

func (b *OrderBuilder) Build() *domain.Order {
    return &domain.Order{
        ID:        b.id,
        UserID:    b.userID,
        Items:     b.items,
        Status:    b.status,
        Total:     b.total,
        CreatedAt: b.createdAt,
    }
}
```

**Uso:**
```go
func TestOrderProcessing(t *testing.T) {
    // Pedido simples
    order := testdata.NewOrderBuilder().
        AddItem("prod-1", 2, 50.00).
        AddItem("prod-2", 1, 100.00).
        Build()
    
    assert.Equal(t, 200.00, order.Total)
    assert.Equal(t, 2, len(order.Items))
    
    // Pedido pago
    paidOrder := testdata.NewOrderBuilder().
        AddItem("prod-1", 1, 99.99).
        Paid().
        Build()
    
    assert.Equal(t, domain.OrderStatusPaid, paidOrder.Status)
}
```

---

## 🏭 Factory Pattern

### Factory para Múltiplas Entidades

```go
package testdata

type TestFactory struct {
    userCounter    int
    productCounter int
    orderCounter   int
}

func NewTestFactory() *TestFactory {
    return &TestFactory{}
}

func (f *TestFactory) CreateUser(opts ...func(*domain.User)) *domain.User {
    f.userCounter++
    
    user := &domain.User{
        ID:    fmt.Sprintf("user-%03d", f.userCounter),
        Name:  fmt.Sprintf("User %d", f.userCounter),
        Email: fmt.Sprintf("user%d@example.com", f.userCounter),
    }
    
    for _, opt := range opts {
        opt(user)
    }
    
    return user
}

func (f *TestFactory) CreateProduct(opts ...func(*domain.Product)) *domain.Product {
    f.productCounter++
    
    product := &domain.Product{
        ID:    fmt.Sprintf("prod-%03d", f.productCounter),
        Name:  fmt.Sprintf("Product %d", f.productCounter),
        Price: 99.99,
        Stock: 100,
    }
    
    for _, opt := range opts {
        opt(product)
    }
    
    return product
}

func (f *TestFactory) CreateOrder(userID string, opts ...func(*domain.Order)) *domain.Order {
    f.orderCounter++
    
    order := &domain.Order{
        ID:     fmt.Sprintf("ord-%03d", f.orderCounter),
        UserID: userID,
        Status: domain.OrderStatusPending,
        Items:  []domain.OrderItem{},
    }
    
    for _, opt := range opts {
        opt(order)
    }
    
    return order
}

// Scenarios complexos
func (f *TestFactory) CreateUserWithOrders(numOrders int) (*domain.User, []*domain.Order) {
    user := f.CreateUser()
    
    orders := make([]*domain.Order, numOrders)
    for i := 0; i < numOrders; i++ {
        orders[i] = f.CreateOrder(user.ID)
    }
    
    return user, orders
}
```

**Uso:**
```go
func TestUserOrders(t *testing.T) {
    factory := testdata.NewTestFactory()
    
    // Criar usuário com 3 pedidos
    user, orders := factory.CreateUserWithOrders(3)
    
    assert.Equal(t, 3, len(orders))
    for _, order := range orders {
        assert.Equal(t, user.ID, order.UserID)
    }
}
```

---

## 💾 Database Fixtures

### SQLite In-Memory

```go
package testdata

import (
    "testing"
    
    "github.com/stretchr/testify/require"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    require.NoError(t, err)
    
    // Auto-migrate models
    err = db.AutoMigrate(
        &UserModel{},
        &ProductModel{},
        &OrderModel{},
        &OrderItemModel{},
    )
    require.NoError(t, err)
    
    return db
}
```

### Seed Data

```go
func SeedTestData(t *testing.T, db *gorm.DB) {
    // Users
    users := []UserModel{
        {ID: "user-001", Name: "John Doe", Email: "john@example.com"},
        {ID: "user-002", Name: "Jane Smith", Email: "jane@example.com"},
        {ID: "admin-001", Name: "Admin", Email: "admin@example.com", Role: "admin"},
    }
    
    for _, user := range users {
        err := db.Create(&user).Error
        require.NoError(t, err)
    }
    
    // Products
    products := []ProductModel{
        {ID: "prod-001", Name: "Product 1", Price: 50.00, Stock: 100},
        {ID: "prod-002", Name: "Product 2", Price: 75.00, Stock: 50},
        {ID: "prod-003", Name: "Product 3", Price: 100.00, Stock: 0}, // Out of stock
    }
    
    for _, product := range products {
        err := db.Create(&product).Error
        require.NoError(t, err)
    }
    
    // Orders
    order := OrderModel{
        ID:     "ord-001",
        UserID: "user-001",
        Status: "confirmed",
        Total:  150.00,
    }
    err := db.Create(&order).Error
    require.NoError(t, err)
}
```

**Uso:**
```go
func TestUserRepository_Integration(t *testing.T) {
    db := testdata.SetupTestDB(t)
    testdata.SeedTestData(t, db)
    
    repo := repository.NewMySQLUserRepository(db)
    
    user, err := repo.FindByID(context.Background(), "user-001")
    require.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
}
```

### Database Cleaner

```go
type DBCleaner struct {
    db *gorm.DB
}

func NewDBCleaner(db *gorm.DB) *DBCleaner {
    return &DBCleaner{db: db}
}

func (c *DBCleaner) Clean() error {
    tables := []string{"order_items", "orders", "products", "users"}
    
    for _, table := range tables {
        if err := c.db.Exec("DELETE FROM " + table).Error; err != nil {
            return err
        }
    }
    
    return nil
}

// Uso com defer
func TestWithCleanup(t *testing.T) {
    db := testdata.SetupTestDB(t)
    cleaner := testdata.NewDBCleaner(db)
    defer cleaner.Clean()
    
    // Teste aqui
}
```

---

## 🌐 HTTP Mocking

### httptest para Handlers

```go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestHTTPHandler_CreateUser(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    
    mockRepo := new(MockUserRepository)
    mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
    
    handler := NewHTTPHandler(mockRepo)
    router.POST("/users", handler.CreateUser)
    
    // Request payload
    payload := map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    }
    body, _ := json.Marshal(payload)
    
    // Make request
    req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    // Record response
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.NotEmpty(t, response["id"])
    
    mockRepo.AssertExpectations(t)
}
```

### Mock HTTP Client

```go
type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    if m.DoFunc != nil {
        return m.DoFunc(req)
    }
    return &http.Response{StatusCode: 200}, nil
}

// Uso
func TestExternalAPICall(t *testing.T) {
    mockClient := &MockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            // Verificar request
            assert.Equal(t, "POST", req.Method)
            assert.Equal(t, "https://api.example.com/users", req.URL.String())
            
            // Retornar mock response
            return &http.Response{
                StatusCode: 201,
                Body:       io.NopCloser(strings.NewReader(`{"id": "123"}`)),
            }, nil
        },
    }
    
    service := NewUserService(mockClient)
    result, err := service.CreateUser("John", "john@example.com")
    
    assert.NoError(t, err)
    assert.Equal(t, "123", result.ID)
}
```

---

## 🔌 gRPC Mocking

### bufconn para Testes

```go
package grpc_test

import (
    "context"
    "net"
    "testing"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/test/bufconn"
    pb "meuApp/proto"
)

func setupTestServer(t *testing.T) (*grpc.Server, *bufconn.Listener, pb.UserServiceClient) {
    buffer := 1024 * 1024
    listener := bufconn.Listen(buffer)
    
    server := grpc.NewServer()
    
    // Register your service
    mockService := &MockUserService{}
    pb.RegisterUserServiceServer(server, mockService)
    
    go func() {
        if err := server.Serve(listener); err != nil {
            t.Errorf("Server error: %v", err)
        }
    }()
    
    // Create client
    conn, _ := grpc.DialContext(context.Background(), "",
        grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
            return listener.Dial()
        }),
        grpc.WithInsecure(),
    )
    
    client := pb.NewUserServiceClient(conn)
    
    return server, listener, client
}

func TestGRPCCreateUser(t *testing.T) {
    server, listener, client := setupTestServer(t)
    defer server.Stop()
    defer listener.Close()
    
    req := &pb.CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    resp, err := client.CreateUser(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.Id)
}
```

---

## ✅ Melhores Práticas

### 1. Isole Testes

```go
// ✅ BOM: Cada teste é independente
func TestUserCreation(t *testing.T) {
    factory := testdata.NewTestFactory()
    user := factory.CreateUser()
    // Test aqui
}

func TestUserDeletion(t *testing.T) {
    factory := testdata.NewTestFactory() // Nova instância
    user := factory.CreateUser()
    // Test aqui
}

// ❌ RUIM: Testes compartilham estado
var sharedFactory = testdata.NewTestFactory()

func TestUserCreation(t *testing.T) {
    user := sharedFactory.CreateUser() // Afeta outros testes
}
```

### 2. Use Table-Driven Tests

```go
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        user    *domain.User
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid user",
            user:    testdata.NewUserBuilder().Build(),
            wantErr: false,
        },
        {
            name:    "empty name",
            user:    testdata.NewUserBuilder().WithName("").Build(),
            wantErr: true,
            errMsg:  "name is required",
        },
        {
            name:    "invalid email",
            user:    testdata.NewUserBuilder().WithEmail("invalid").Build(),
            wantErr: true,
            errMsg:  "invalid email format",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 3. Nomeie Fixtures Claramente

```go
// ✅ BOM: Nome descritivo
validUser := testdata.NewUserBuilder().Build()
inactiveAdmin := testdata.NewUserBuilder().AsAdmin().Inactive().Build()
outOfStockProduct := testdata.NewProductBuilder().OutOfStock().Build()

// ❌ RUIM: Nomes genéricos
user1 := testdata.NewUserBuilder().Build()
user2 := testdata.NewUserBuilder().AsAdmin().Inactive().Build()
product := testdata.NewProductBuilder().OutOfStock().Build()
```

### 4. Mantenha Fixtures Simples

```go
// ✅ BOM: Fixture simples, lógica no teste
user := testdata.NewUserBuilder().Build()
user.ChangeEmail("new@example.com") // Lógica no teste

// ❌ RUIM: Fixture com muita lógica
user := testdata.CreateUserWithEmailChangeHistory(3) // Muita complexidade
```

### 5. Documente Fixtures Complexas

```go
// CreateUserWithOrders cria um usuário de teste com N pedidos.
// Os pedidos são criados com status "pending" e items aleatórios.
// Use este helper quando precisar testar cenários de múltiplos pedidos.
func (f *TestFactory) CreateUserWithOrders(numOrders int) (*domain.User, []*domain.Order) {
    // ...
}
```

---

## 📋 Checklist de Mocking

### ✅ Para Cada Teste

- [ ] Mock apenas dependências externas (DB, HTTP, gRPC)
- [ ] Use fixtures para dados repetitivos
- [ ] Verifique expectations ao final (`AssertExpectations`)
- [ ] Não compartilhe estado entre testes
- [ ] Nomeie claramente cenários de teste

### ✅ Para Mocks

- [ ] Implemente apenas métodos necessários
- [ ] Use matchers apropriados (`mock.Anything`, `mock.MatchedBy`)
- [ ] Configure retornos para todos os cenários (sucesso, erro)
- [ ] Verifique chamadas quando relevante

### ✅ Para Fixtures

- [ ] Use builders para entidades complexas
- [ ] Crie factories para múltiplas entidades
- [ ] Mantenha fixtures em package separado (`testdata`)
- [ ] Documente fixtures não-triviais

---

## 🎓 Resumo

### O Que Aprendemos

✅ **Mocks** verificam comportamento (chamadas, parâmetros)  
✅ **Stubs** retornam valores pré-programados  
✅ **Fakes** são implementações in-memory funcionais  
✅ **Fixtures** fornecem dados de teste consistentes  
✅ **Builders** criam objetos com API fluente  
✅ **Factories** centralizam criação de test data  

### Ferramentas Principais

- **testify/mock** - Mocking poderoso
- **testify/assert** - Assertions legíveis
- **testify/require** - Assertions que param teste
- **httptest** - Testes de HTTP handlers
- **bufconn** - Testes de gRPC in-memory
- **SQLite** - Database in-memory para integração

### Padrão Recomendado

```go
// 1. Setup (Arrange)
mockRepo := new(MockUserRepository)
mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

factory := testdata.NewTestFactory()
user := factory.CreateUser()

// 2. Ação (Act)
err := handler.Handle(ctx, &CreateUserCommand{...})

// 3. Verificação (Assert)
assert.NoError(t, err)
mockRepo.AssertExpectations(t)
```

---

## 📚 Próximos Passos

- **[21. Estratégia de Testes ←](21-testing-strategy.md)** - Pirâmide de testes
- **[20. Logging ←](20-logging.md)** - Logging em testes
- **[23. Boas Práticas ←](23-best-practices.md)** - Code quality

---

<div align="center">

**[⬆️ Voltar ao Topo](#-mocks-e-fixtures)**

</div>
