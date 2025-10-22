# 🧪 Fase 7.2 - Testes: Plano de Implementação

> **Status:** 🚧 **EM ANDAMENTO**  
> **Data de Início:** 19 de Outubro de 2025  
> **Objetivo:** Implementar cobertura de testes abrangente (target: 80%+)

---

## 🎯 Objetivos da Fase

1. **Unit Tests:** Testar componentes isolados (registry, events, domain logic)
2. **Integration Tests:** Testar interação entre camadas (repository + service + handler)
3. **E2E Tests:** Testar fluxos completos via HTTP e gRPC
4. **Coverage Report:** Gerar relatórios e garantir 80%+ de cobertura
5. **CI/CD Ready:** Preparar testes para execução em pipeline

---

## 📊 Escopo de Testes

### 1. Unit Tests (0/15 tarefas)

#### 1.1 ModuleRegistry (`pkg/container/registry_test.go`)
- [ ] TestNewModuleRegistry
- [ ] TestRegisterHTTPHandler
- [ ] TestRegisterGRPCService
- [ ] TestRegisterRepository
- [ ] TestRegisterApplicationService
- [ ] TestGetStats

#### 1.2 Event Bus (`pkg/events/`)
- [ ] TestEventBus_Subscribe
- [ ] TestEventBus_Publish
- [ ] TestTypedEventPublisher
- [ ] TestEventHandlers

#### 1.3 Domain Entities
- [ ] TestUserDomain_Validation
- [ ] TestProductDomain_Validation
- [ ] TestOrderDomain_Validation
- [ ] TestOrderCalculateTotalAmount

#### 1.4 DTOs e Mappers
- [ ] TestUserMapper_ToDTO
- [ ] TestProductMapper_ToDTO

**Total:** 15 testes unitários

---

### 2. Integration Tests (0/9 tarefas)

#### 2.1 User Module (`internal/modules/user/tests/integration/`)
- [ ] TestUserRepository_CRUD
- [ ] TestCreateUserCommand
- [ ] TestGetUserQuery
- [ ] TestUserApplicationService_FullFlow

#### 2.2 Product Module (`internal/modules/product/tests/integration/`)
- [ ] TestProductRepository_CRUD
- [ ] TestCreateProductCommand
- [ ] TestGetProductQuery

#### 2.3 Order Module (`internal/modules/order/tests/integration/`)
- [ ] TestOrderRepository_CRUD
- [ ] TestCreateOrderCommand_WithDependencies

**Total:** 9 testes de integração

---

### 3. E2E Tests (0/8 tarefas)

#### 3.1 HTTP Endpoints (`tests/e2e/http/`)
- [ ] TestUserEndpoints_CreateGetListDelete
- [ ] TestProductEndpoints_CreateGetListUpdate
- [ ] TestOrderEndpoints_CreateGetUpdateStatus

#### 3.2 gRPC Services (`tests/e2e/grpc/`)
- [ ] TestUserGRPCService
- [ ] TestProductGRPCService
- [ ] TestOrderGRPCService

#### 3.3 Cross-Module Flows
- [ ] TestOrderCreation_WithUserAndProduct
- [ ] TestEventPropagation_AcrossModules

**Total:** 8 testes end-to-end

---

## 🛠️ Ferramentas e Bibliotecas

### Testing Framework
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
    "github.com/stretchr/testify/mock"
)
```

### Mocking
- **testify/mock:** Mocks para interfaces
- **gomock (opcional):** Alternative mocking framework
- **sqlmock:** Mock para database queries

### Database Testing
- **dockertest:** Containers efêmeros para testes
- **testcontainers-go:** Alternative para Docker
- **in-memory SQLite:** Testes rápidos sem Docker

### HTTP Testing
- **httptest:** stdlib para testar handlers HTTP
- **testify/suite:** Test suites para setup/teardown

### gRPC Testing
- **grpc.testing:** Test server/client
- **bufconn:** In-memory connection para gRPC

### Coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

---

## 📁 Estrutura de Diretórios

```
meuApp/
├── tests/                                    # Testes E2E
│   ├── e2e/
│   │   ├── http/
│   │   │   ├── user_test.go
│   │   │   ├── product_test.go
│   │   │   └── order_test.go
│   │   ├── grpc/
│   │   │   ├── user_test.go
│   │   │   ├── product_test.go
│   │   │   └── order_test.go
│   │   └── helpers/
│   │       ├── testserver.go
│   │       └── fixtures.go
│   └── README.md
│
├── pkg/
│   ├── container/
│   │   ├── registry.go
│   │   └── registry_test.go                 # Unit tests
│   ├── events/
│   │   ├── eventbus.go
│   │   ├── eventbus_test.go                 # Unit tests
│   │   ├── typed.go
│   │   └── typed_test.go                    # Unit tests
│   └── errors/
│       ├── errors.go
│       └── errors_test.go                   # Unit tests
│
└── internal/
    └── modules/
        ├── user/
        │   ├── tests/
        │   │   ├── unit/
        │   │   │   ├── domain_test.go
        │   │   │   └── dto_test.go
        │   │   └── integration/
        │   │       ├── repository_test.go
        │   │       ├── commands_test.go
        │   │       └── queries_test.go
        │   └── ...
        ├── product/
        │   ├── tests/
        │   │   ├── unit/
        │   │   └── integration/
        │   └── ...
        └── order/
            ├── tests/
            │   ├── unit/
            │   └── integration/
            └── ...
```

---

## 🎯 Estratégia de Implementação

### Fase 1: Unit Tests (Dias 1-2)
**Prioridade:** ALTA  
**Objetivo:** Testar lógica isolada sem dependências externas

1. Criar `pkg/container/registry_test.go`
2. Criar `pkg/events/eventbus_test.go`
3. Criar testes de domínio para cada módulo
4. Criar testes de DTO/Mappers

**Estimativa:** 4-6 horas

---

### Fase 2: Integration Tests (Dias 3-4)
**Prioridade:** ALTA  
**Objetivo:** Testar interação entre camadas

1. Setup de banco de dados para testes (dockertest ou SQLite)
2. User Module integration tests
3. Product Module integration tests
4. Order Module integration tests

**Estimativa:** 8-10 horas

---

### Fase 3: E2E Tests (Dias 5-6)
**Prioridade:** MÉDIA  
**Objetivo:** Testar fluxos completos

1. Setup do test server (HTTP + gRPC)
2. HTTP endpoint tests
3. gRPC service tests
4. Cross-module flow tests

**Estimativa:** 6-8 horas

---

### Fase 4: Coverage & CI (Dia 7)
**Prioridade:** MÉDIA  
**Objetivo:** Métricas e automação

1. Configurar coverage report
2. Adicionar comandos no Makefile
3. Criar GitHub Actions workflow (opcional)
4. Documentação de testes

**Estimativa:** 2-4 horas

---

## 📋 Checklist de Progresso

### 🔵 Unit Tests (0/15)
- [ ] ModuleRegistry (6 tests)
- [ ] Event Bus (4 tests)
- [ ] Domain Entities (3 tests)
- [ ] DTOs e Mappers (2 tests)

### 🟢 Integration Tests (0/9)
- [ ] User Module (4 tests)
- [ ] Product Module (3 tests)
- [ ] Order Module (2 tests)

### 🟡 E2E Tests (0/8)
- [ ] HTTP Endpoints (3 tests)
- [ ] gRPC Services (3 tests)
- [ ] Cross-Module (2 tests)

### 🟣 Infrastructure (0/5)
- [ ] Setup de test database
- [ ] Test helpers e fixtures
- [ ] Coverage configuration
- [ ] Makefile targets
- [ ] Documentação

**Total:** 0/37 tarefas

---

## 🔧 Comandos Úteis

### Executar Todos os Testes
```bash
make test
# ou
go test ./...
```

### Executar com Coverage
```bash
make test-coverage
# ou
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Executar Apenas Unit Tests
```bash
go test ./pkg/...
```

### Executar Apenas Integration Tests
```bash
go test ./internal/.../tests/integration/...
```

### Executar Apenas E2E Tests
```bash
go test ./tests/e2e/...
```

### Verbose Output
```bash
go test -v ./...
```

### Run Specific Test
```bash
go test -v -run TestUserRepository_Create ./internal/modules/user/tests/integration/
```

---

## 📈 Métricas de Sucesso

### Cobertura de Código
- **Mínimo aceitável:** 70%
- **Target:** 80%
- **Ideal:** 85%+

### Cobertura por Camada
- **Domain Layer:** 90%+ (crítico)
- **Application Layer:** 85%+ (alta prioridade)
- **Adapters Layer:** 70%+ (médio)
- **Infrastructure:** 60%+ (baixo)

### Performance
- **Unit Tests:** < 1s total
- **Integration Tests:** < 10s total
- **E2E Tests:** < 30s total

---

## 🎓 Boas Práticas

### 1. Nomenclatura
```go
// Unit Test
func TestFunctionName_Scenario_ExpectedBehavior(t *testing.T)

// Example
func TestCreateUser_ValidInput_ReturnsUserAndNoError(t *testing.T)
```

### 2. AAA Pattern (Arrange, Act, Assert)
```go
func TestExample(t *testing.T) {
    // Arrange
    input := "test"
    expected := "TEST"
    
    // Act
    result := ToUpper(input)
    
    // Assert
    assert.Equal(t, expected, result)
}
```

### 3. Table-Driven Tests
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 4. Mocking Interfaces
```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func TestWithMock(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("Create", mock.Anything).Return(nil)
    
    // Use mockRepo in test
    
    mockRepo.AssertExpectations(t)
}
```

### 5. Test Fixtures
```go
// helpers/fixtures.go
func CreateTestUser() *domain.User {
    return &domain.User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
}
```

---

## 🚀 Próximos Passos

1. ✅ Criar este documento de planejamento
2. 🎯 **AGORA:** Implementar Unit Tests (ModuleRegistry)
3. ⏭️ Implementar Integration Tests (User Module)
4. ⏭️ Implementar E2E Tests (HTTP endpoints)
5. ⏭️ Configurar coverage report

---

## 📚 Referências

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Go Test Coverage](https://go.dev/blog/cover)
- [Testing Best Practices](https://github.com/golang/go/wiki/TestComments)

---

**Última Atualização:** 19 de Outubro de 2025  
**Responsável:** Arquitetura Artemis  
**Próxima Revisão:** Após conclusão de cada fase
