# 🎨 Boas Práticas e Convenções

## 📋 Índice
- [Princípios Gerais](#princípios-gerais)
- [Estrutura de Código](#estrutura-de-código)
- [Nomenclatura](#nomenclatura)
- [Tratamento de Erros](#tratamento-de-erros)
- [Testes](#testes)
- [Performance](#performance)
- [Segurança](#segurança)

---

## 🎯 Princípios Gerais

### SOLID

Sempre siga os princípios SOLID:

```go
// ✅ Single Responsibility
type UserRepository struct {
    db *gorm.DB
}

// ✅ Open/Closed - Extensível via interfaces
type Logger interface {
    Log(msg string)
}

// ✅ Liskov Substitution - Substituível
var logger Logger = NewConsoleLogger() // ou NewFileLogger()

// ✅ Interface Segregation - Interfaces pequenas
type Reader interface { Read() }
type Writer interface { Write() }

// ✅ Dependency Inversion - Dependa de abstrações
type Handler struct {
    repo UserRepository  // Interface, não implementação
}
```

### DRY (Don't Repeat Yourself)

```go
// ❌ Ruim: Código duplicado
func CreateUser(...) { validate(); save(); log(); }
func UpdateUser(...) { validate(); save(); log(); }

// ✅ Bom: Reutilização
func (s *Service) save(user *User) error {
    if err := s.validate(user); err != nil { return err }
    if err := s.repo.Save(user); err != nil { return err }
    s.logger.Info("User saved", user.ID)
    return nil
}
```

### KISS (Keep It Simple, Stupid)

```go
// ❌ Over-engineering
type AbstractUserFactoryProvider interface {
    CreateFactoryInstance() UserFactory
}

// ✅ Simples e direto
func NewUser(name, email string) *User {
    return &User{Name: name, Email: email}
}
```

---

## 📐 Estrutura de Código

### Organização de Imports

```go
import (
    // 1. Standard library
    "context"
    "errors"
    "fmt"
    
    // 2. External packages
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    // 3. Internal packages
    "meuApp/internal/modules/user/domain"
    "meuApp/pkg/contracts"
)
```

### Ordem de Declarações em Structs

```go
type UserHandler struct {
    // 1. Dependencies (interfaces first)
    repo      UserRepository
    logger    Logger
    emailSvc  EmailService
    
    // 2. Configuration
    config    *Config
    
    // 3. State (avoid if possible)
    cache     map[string]*User
}
```

### Ordem de Métodos

```go
// 1. Constructor
func NewUserService() *UserService { }

// 2. Public methods (alfabética)
func (s *UserService) CreateUser() { }
func (s *UserService) DeleteUser() { }
func (s *UserService) GetUser() { }

// 3. Private methods (alfabética)
func (s *UserService) validate() { }
func (s *UserService) sanitize() { }
```

---

## 📝 Nomenclatura

### Packages

```go
✅ package user       // Singular
❌ package users      // Plural

✅ package repository // Completo
❌ package repo       // Abreviado
```

### Variáveis

```go
// Contexto
ctx context.Context    // Sempre "ctx"

// Errors
err error              // Sempre "err"

// Receivers
func (h *Handler)      // Primeira letra do tipo
func (s *Service)

// Loops
for i, user := range users {  // i para índice, nome significativo
    
}

// Booleans - Começam com is/has/can
isValid := true
hasPermission := false
canAccess := true
```

### Funções e Métodos

```go
// Construtores
func NewUser() *User
func NewUserRepository() UserRepository

// Getters (sem prefixo Get)
func (u *User) Email() string  // ✅
func (u *User) GetEmail()      // ❌

// Setters
func (u *User) SetEmail(email string)

// Conversões
func (u *User) ToDTO() *UserDTO
func (dto *UserDTO) ToEntity() *User

// Validações
func (u *User) Validate() error
func (u *User) IsValid() bool

// CRUD Operations
func (r *Repository) Create(user *User) error
func (r *Repository) FindByID(id string) (*User, error)
func (r *Repository) Update(user *User) error
func (r *Repository) Delete(id string) error
```

---

## ⚠️ Tratamento de Erros

### Sempre Trate Erros

```go
// ❌ Ignorar erros
user, _ := repo.FindByID(id)

// ✅ Sempre trate
user, err := repo.FindByID(id)
if err != nil {
    return nil, fmt.Errorf("failed to find user: %w", err)
}
```

### Use Erros Customizados

```go
// errors.go
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidEmail     = errors.New("invalid email")
    ErrDuplicateEmail   = errors.New("email already exists")
)

// Uso
func (r *Repository) FindByEmail(email string) (*User, error) {
    user, err := r.db.Where("email = ?", email).First()
    if err == gorm.ErrRecordNotFound {
        return nil, ErrUserNotFound  // Erro específico do domínio
    }
    return user, err
}
```

### Wrap Errors com Contexto

```go
// ❌ Erro genérico
return err

// ✅ Com contexto
return fmt.Errorf("failed to create user %s: %w", email, err)
```

### Error Handling Pattern

```go
func (h *Handler) CreateUser(cmd *CreateUserCommand) error {
    // 1. Validação de entrada
    if err := cmd.Validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // 2. Lógica de negócio
    user, err := h.buildUser(cmd)
    if err != nil {
        return fmt.Errorf("failed to build user: %w", err)
    }
    
    // 3. Persistência
    if err := h.repo.Save(user); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    // 4. Side effects (não falha a operação principal)
    if err := h.emailSvc.SendWelcome(user.Email); err != nil {
        h.logger.Warn("failed to send welcome email", err)
        // Não retorna erro - é um side effect
    }
    
    return nil
}
```

---

## 🧪 Testes

### Nomenclatura de Testes

```go
// Padrão: Test<Function>_<Scenario>_<ExpectedResult>

func TestCreateUser_ValidInput_Success(t *testing.T) {}
func TestCreateUser_InvalidEmail_ReturnsError(t *testing.T) {}
func TestCreateUser_DuplicateEmail_ReturnsError(t *testing.T) {}
```

### Estrutura AAA (Arrange, Act, Assert)

```go
func TestCreateUser_ValidInput_Success(t *testing.T) {
    // Arrange (Setup)
    mockRepo := new(MockUserRepository)
    handler := NewCreateUserHandler(mockRepo)
    cmd := &CreateUserCommand{
        Name:  "John",
        Email: "john@example.com",
    }
    mockRepo.On("Save", mock.Anything).Return(nil)
    
    // Act (Execute)
    err := handler.Handle(cmd)
    
    // Assert (Verify)
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Table-Driven Tests

```go
func TestUserEmail_Validate(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"empty", "", true},
        {"missing domain", "user@", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            email := Email{value: tt.email}
            err := email.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Use Mocks Apropriadamente

```go
// Para testes unitários
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Save(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

// Para testes de integração use banco real (container Docker)
func TestUserRepository_Integration(t *testing.T) {
    db := setupTestDB()  // Container Docker com MySQL
    defer db.Close()
    
    repo := NewMySQLUserRepository(db)
    // Testa com banco de verdade
}
```

---

## ⚡ Performance

### Database Queries

```go
// ❌ N+1 Problem
users, _ := repo.FindAll()
for _, user := range users {
    orders, _ := orderRepo.FindByUserID(user.ID)  // N queries!
}

// ✅ Eager Loading
users, _ := repo.FindAllWithOrders()  // 1 query com JOIN
```

### Use Paginação

```go
// ❌ Carrega tudo
func ListUsers() ([]*User, error) {
    return repo.FindAll()  // Pode retornar milhões
}

// ✅ Paginação
func ListUsers(page, pageSize int) ([]*User, int64, error) {
    return repo.FindAll(page, pageSize)
}
```

### Context com Timeout

```go
// ✅ Sempre use timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := service.LongOperation(ctx)
```

### Cache Estratégico

```go
type CachedUserRepository struct {
    repo  UserRepository
    cache Cache
}

func (r *CachedUserRepository) FindByID(id string) (*User, error) {
    // Tenta cache primeiro
    if cached, found := r.cache.Get(id); found {
        return cached.(*User), nil
    }
    
    // Busca do banco
    user, err := r.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // Armazena em cache
    r.cache.Set(id, user, 5*time.Minute)
    return user, nil
}
```

---

## 🔒 Segurança

### Validação de Entrada

```go
// ✅ Sempre valide input do usuário
type CreateUserCommand struct {
    Name     string `validate:"required,min=3,max=100"`
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8"`
}

func (cmd *CreateUserCommand) Validate() error {
    validate := validator.New()
    return validate.Struct(cmd)
}
```

### Sanitização

```go
import "html"

func (cmd *CreateUserCommand) Sanitize() {
    cmd.Name = html.EscapeString(strings.TrimSpace(cmd.Name))
    cmd.Email = strings.ToLower(strings.TrimSpace(cmd.Email))
}
```

### Senhas

```go
// ❌ NUNCA armazene senhas em plain text
user.Password = cmd.Password

// ✅ Sempre faça hash
hashedPassword, err := passwordHasher.Hash(cmd.Password)
if err != nil {
    return err
}
user.Password = hashedPassword
```

### SQL Injection

```go
// ❌ String concatenation
query := "SELECT * FROM users WHERE email = '" + email + "'"

// ✅ Prepared statements (GORM faz automaticamente)
db.Where("email = ?", email).First(&user)
```

### Autenticação e Autorização

```go
// Sempre verifique permissões
func (h *DeleteUserHandler) Handle(cmd *DeleteUserCommand) error {
    // 1. Verificar autenticação
    if cmd.ActorID == "" {
        return errors.New("unauthenticated")
    }
    
    // 2. Verificar autorização
    if !h.authService.HasPermission(cmd.ActorID, "user.delete") {
        return errors.New("unauthorized")
    }
    
    // 3. Executar ação
    return h.repo.Delete(cmd.UserID)
}
```

---

## 📚 Documentação

### Comentários em Código

```go
// ✅ Comente o "porquê", não o "o quê"

// ❌ Ruim: Óbvio
// Incrementa o contador
counter++

// ✅ Bom: Explica decisão
// Usamos sleep aqui devido a rate limiting da API externa
time.Sleep(1 * time.Second)
```

### Godoc

```go
// Package user provides user management functionality.
// It implements user creation, authentication, and profile management.
package user

// UserRepository defines the interface for user persistence.
// Implementations should handle connection pooling and transactions appropriately.
type UserRepository interface {
    // Save persists a user to the database.
    // Returns ErrDuplicateEmail if email already exists.
    Save(user *User) error
    
    // FindByID retrieves a user by ID.
    // Returns ErrUserNotFound if user doesn't exist.
    FindByID(id string) (*User, error)
}
```

---

## 🎯 Checklist de Code Review

Antes de fazer commit, verifique:

### Geral
- [ ] Código segue convenções do projeto
- [ ] Sem código comentado desnecessário
- [ ] Sem prints de debug
- [ ] Imports organizados

### Lógica
- [ ] Erros são tratados adequadamente
- [ ] Validações de entrada implementadas
- [ ] Regras de negócio no lugar correto (domain layer)
- [ ] Não há código duplicado

### Performance
- [ ] Queries otimizadas (sem N+1)
- [ ] Paginação implementada
- [ ] Context com timeout onde apropriado

### Segurança
- [ ] Input sanitizado
- [ ] Senhas com hash
- [ ] Autorização verificada
- [ ] Sem secrets no código

### Testes
- [ ] Testes unitários implementados
- [ ] Coverage > 80%
- [ ] Testes passando
- [ ] Edge cases cobertos

---

## 📚 Recursos Adicionais

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Clean Code em Go](https://github.com/Pungyeon/clean-go-article)

---

**[⬅️ Índice](README.md)** | **[API Reference ➡️](24-api-reference.md)**
