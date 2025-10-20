# 💾 Padrão Repository

## 📋 Índice
- [O que é Repository?](#o-que-é-repository)
- [Repository no Artemis](#repository-no-artemis)
- [Implementações](#implementações)
- [Padrões Avançados](#padrões-avançados)
- [Boas Práticas](#boas-práticas)

---

## 🎯 O que é Repository?

O **Repository Pattern** é um padrão que encapsula a lógica de acesso a dados, fornecendo uma interface de coleção para trabalhar com agregados de domínio.

### Conceito

```
┌──────────────────────────────────────────────────┐
│         APPLICATION LAYER                        │
│                                                  │
│  userService.CreateUser(user)                   │
│        │                                         │
│        ▼                                         │
│  userRepo.Save(user)  ← Interface               │
└────────┬─────────────────────────────────────────┘
         │
         │ Dependency Inversion
         │
┌────────▼─────────────────────────────────────────┐
│         INFRASTRUCTURE LAYER                     │
│                                                  │
│  MySQLUserRepository implements UserRepository  │
│        │                                         │
│        ▼                                         │
│     DATABASE                                     │
└──────────────────────────────────────────────────┘
```

### Sem Repository (❌ Ruim)

```go
// Application layer acessa DB diretamente
func (s *UserService) CreateUser(name, email string) error {
    // ❌ SQL direto no service
    _, err := s.db.Exec(
        "INSERT INTO users (name, email) VALUES (?, ?)",
        name, email,
    )
    return err
}
```

**Problemas:**
- ❌ Lógica de persistência espalhada
- ❌ Difícil de testar
- ❌ Acoplamento com tecnologia específica
- ❌ Difícil de trocar banco de dados

### Com Repository (✅ Bom)

```go
// Interface
type UserRepository interface {
    Save(user *User) error
    FindByID(id string) (*User, error)
}

// Application layer usa interface
func (s *UserService) CreateUser(name, email string) error {
    user := &User{Name: name, Email: email}
    return s.repo.Save(user) // ✅ Abstração
}

// Implementação MySQL
type MySQLUserRepository struct {
    db *gorm.DB
}

func (r *MySQLUserRepository) Save(user *User) error {
    return r.db.Create(user).Error
}
```

**Benefícios:**
- ✅ Lógica de persistência centralizada
- ✅ Fácil de testar (mock repository)
- ✅ Desacoplado de tecnologia
- ✅ Fácil trocar implementação

---

## 🏗️ Repository no Artemis

### Estrutura de um Repository

```
internal/modules/user/repository/
├── repository.go           # Interface
├── mysql_repository.go     # Implementação MySQL
├── cache_repository.go     # Decorator com cache
└── repository_test.go      # Testes
```

### 1. Interface do Repository

**`repository/repository.go`**

```go
package repository

import (
    "context"
    "meuApp/internal/modules/user/domain/entities"
)

// UserRepository define as operações de persistência
type UserRepository interface {
    // Create
    Save(ctx context.Context, user *entities.User) error
    
    // Read
    FindByID(ctx context.Context, id string) (*entities.User, error)
    FindByEmail(ctx context.Context, email string) (*entities.User, error)
    FindAll(ctx context.Context, page, pageSize int) ([]*entities.User, int64, error)
    
    // Update
    Update(ctx context.Context, user *entities.User) error
    
    // Delete
    Delete(ctx context.Context, id string) error
    
    // Queries específicas
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    CountActive(ctx context.Context) (int64, error)
}
```

### 2. Implementação MySQL

**`repository/mysql_repository.go`**

```go
package repository

import (
    "context"
    "meuApp/internal/modules/user/domain/entities"
    "gorm.io/gorm"
)

type MySQLUserRepository struct {
    db *gorm.DB
}

func NewMySQLUserRepository(db *gorm.DB) *MySQLUserRepository {
    return &MySQLUserRepository{db: db}
}

// UserModel é o modelo do banco de dados
type UserModel struct {
    ID        string `gorm:"primaryKey"`
    Name      string
    Email     string `gorm:"uniqueIndex"`
    Password  string
    Active    bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (UserModel) TableName() string {
    return "users"
}

// Save persiste um usuário
func (r *MySQLUserRepository) Save(ctx context.Context, user *entities.User) error {
    model := r.toModel(user)
    return r.db.WithContext(ctx).Create(model).Error
}

// FindByID busca usuário por ID
func (r *MySQLUserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
    var model UserModel
    
    err := r.db.WithContext(ctx).
        Where("id = ?", id).
        First(&model).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return r.toEntity(&model), nil
}

// FindByEmail busca usuário por email
func (r *MySQLUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
    var model UserModel
    
    err := r.db.WithContext(ctx).
        Where("email = ?", email).
        First(&model).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return r.toEntity(&model), nil
}

// FindAll lista usuários com paginação
func (r *MySQLUserRepository) FindAll(
    ctx context.Context,
    page, pageSize int,
) ([]*entities.User, int64, error) {
    var models []UserModel
    var total int64
    
    // Count total
    if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Query com paginação
    offset := (page - 1) * pageSize
    err := r.db.WithContext(ctx).
        Offset(offset).
        Limit(pageSize).
        Order("created_at DESC").
        Find(&models).Error
    
    if err != nil {
        return nil, 0, err
    }
    
    // Convert to entities
    users := make([]*entities.User, len(models))
    for i, model := range models {
        users[i] = r.toEntity(&model)
    }
    
    return users, total, nil
}

// Update atualiza um usuário
func (r *MySQLUserRepository) Update(ctx context.Context, user *entities.User) error {
    model := r.toModel(user)
    return r.db.WithContext(ctx).Save(model).Error
}

// Delete remove um usuário
func (r *MySQLUserRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&UserModel{}).Error
}

// ExistsByEmail verifica se email já existe
func (r *MySQLUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&UserModel{}).
        Where("email = ?", email).
        Count(&count).Error
    
    return count > 0, err
}

// CountActive conta usuários ativos
func (r *MySQLUserRepository) CountActive(ctx context.Context) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&UserModel{}).
        Where("active = ?", true).
        Count(&count).Error
    
    return count, err
}

// Mappers (Entity ↔ Model)
func (r *MySQLUserRepository) toEntity(model *UserModel) *entities.User {
    return &entities.User{
        ID:        model.ID,
        Name:      model.Name,
        Email:     model.Email,
        Password:  model.Password,
        Active:    model.Active,
        CreatedAt: model.CreatedAt,
        UpdatedAt: model.UpdatedAt,
    }
}

func (r *MySQLUserRepository) toModel(user *entities.User) *UserModel {
    return &UserModel{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Password:  user.Password,
        Active:    user.Active,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}
```

### 3. Erros Customizados

**`repository/errors.go`**

```go
package repository

import "errors"

var (
    ErrUserNotFound     = errors.New("user not found")
    ErrDuplicateEmail   = errors.New("email already exists")
    ErrInvalidUser      = errors.New("invalid user data")
)
```

---

## 🎨 Padrões Avançados

### 1. Repository com Cache (Decorator Pattern)

```go
type CachedUserRepository struct {
    repo  UserRepository
    cache Cache
    ttl   time.Duration
}

func NewCachedUserRepository(repo UserRepository, cache Cache) *CachedUserRepository {
    return &CachedUserRepository{
        repo:  repo,
        cache: cache,
        ttl:   5 * time.Minute,
    }
}

func (r *CachedUserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
    cacheKey := fmt.Sprintf("user:%s", id)
    
    // Tenta cache primeiro
    if cached, found := r.cache.Get(cacheKey); found {
        return cached.(*entities.User), nil
    }
    
    // Busca do repositório
    user, err := r.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Armazena em cache
    r.cache.Set(cacheKey, user, r.ttl)
    
    return user, nil
}

func (r *CachedUserRepository) Save(ctx context.Context, user *entities.User) error {
    // Salva no repositório
    if err := r.repo.Save(ctx, user); err != nil {
        return err
    }
    
    // Invalida cache
    cacheKey := fmt.Sprintf("user:%s", user.ID)
    r.cache.Delete(cacheKey)
    
    return nil
}
```

### 2. Specification Pattern

```go
// Specification define um critério de busca
type Specification interface {
    ToSQL() (string, []interface{})
}

// Exemplo: Usuários ativos com email específico
type ActiveUserWithEmailSpec struct {
    email string
}

func (s *ActiveUserWithEmailSpec) ToSQL() (string, []interface{}) {
    return "active = ? AND email = ?", []interface{}{true, s.email}
}

// Repository aceita specifications
func (r *MySQLUserRepository) FindBySpec(
    ctx context.Context,
    spec Specification,
) ([]*entities.User, error) {
    var models []UserModel
    
    sql, args := spec.ToSQL()
    err := r.db.WithContext(ctx).
        Where(sql, args...).
        Find(&models).Error
    
    if err != nil {
        return nil, err
    }
    
    users := make([]*entities.User, len(models))
    for i, model := range models {
        users[i] = r.toEntity(&model)
    }
    
    return users, nil
}

// Uso
spec := &ActiveUserWithEmailSpec{email: "john@example.com"}
users, err := repo.FindBySpec(ctx, spec)
```

### 3. Unit of Work Pattern

```go
type UnitOfWork struct {
    db *gorm.DB
    tx *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
    return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) Begin() error {
    uow.tx = uow.db.Begin()
    return uow.tx.Error
}

func (uow *UnitOfWork) Commit() error {
    return uow.tx.Commit().Error
}

func (uow *UnitOfWork) Rollback() error {
    return uow.tx.Rollback().Error
}

func (uow *UnitOfWork) UserRepository() UserRepository {
    return NewMySQLUserRepository(uow.tx)
}

func (uow *UnitOfWork) OrderRepository() OrderRepository {
    return NewMySQLOrderRepository(uow.tx)
}

// Uso
func (s *OrderService) CreateOrderWithUser(
    userName, userEmail string,
    items []OrderItem,
) error {
    uow := NewUnitOfWork(s.db)
    
    if err := uow.Begin(); err != nil {
        return err
    }
    defer func() {
        if r := recover(); r != nil {
            uow.Rollback()
            panic(r)
        }
    }()
    
    // Criar usuário
    user := &User{Name: userName, Email: userEmail}
    if err := uow.UserRepository().Save(context.Background(), user); err != nil {
        uow.Rollback()
        return err
    }
    
    // Criar pedido
    order := &Order{UserID: user.ID, Items: items}
    if err := uow.OrderRepository().Save(context.Background(), order); err != nil {
        uow.Rollback()
        return err
    }
    
    return uow.Commit()
}
```

### 4. Generic Repository (Go 1.18+)

```go
type Repository[T any] interface {
    Save(ctx context.Context, entity T) error
    FindByID(ctx context.Context, id string) (T, error)
    FindAll(ctx context.Context) ([]T, error)
    Delete(ctx context.Context, id string) error
}

type GenericRepository[T any] struct {
    db        *gorm.DB
    tableName string
}

func NewGenericRepository[T any](db *gorm.DB, tableName string) *GenericRepository[T] {
    return &GenericRepository[T]{
        db:        db,
        tableName: tableName,
    }
}

func (r *GenericRepository[T]) Save(ctx context.Context, entity T) error {
    return r.db.WithContext(ctx).Create(&entity).Error
}

func (r *GenericRepository[T]) FindByID(ctx context.Context, id string) (T, error) {
    var entity T
    err := r.db.WithContext(ctx).
        Where("id = ?", id).
        First(&entity).Error
    return entity, err
}

// Uso
userRepo := NewGenericRepository[User](db, "users")
productRepo := NewGenericRepository[Product](db, "products")
```

---

## 🎯 Boas Práticas

### 1. Separe Model de Entity

```go
// ✅ BOM: Model separado de Entity
// Domain Entity (regras de negócio)
type User struct {
    ID    string
    Email string
    Name  string
}

func (u *User) ChangeEmail(newEmail string) error {
    if !isValidEmail(newEmail) {
        return errors.New("invalid email")
    }
    u.Email = newEmail
    return nil
}

// Database Model (persistência)
type UserModel struct {
    ID        string `gorm:"primaryKey"`
    Email     string `gorm:"uniqueIndex"`
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Mapper
func toEntity(model *UserModel) *User {
    return &User{
        ID:    model.ID,
        Email: model.Email,
        Name:  model.Name,
    }
}
```

### 2. Use Context

```go
// ✅ Sempre aceite context
func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
    return r.db.WithContext(ctx).First(&user, id).Error
}

// ❌ Sem context
func (r *Repository) FindByID(id string) (*User, error) {
    return r.db.First(&user, id).Error
}
```

### 3. Retorne Erros Específicos

```go
// ✅ Erros específicos do domínio
var ErrUserNotFound = errors.New("user not found")

func (r *Repository) FindByID(id string) (*User, error) {
    var user User
    err := r.db.First(&user, id).Error
    
    if err == gorm.ErrRecordNotFound {
        return nil, ErrUserNotFound  // Erro do domínio
    }
    
    return &user, err
}

// Handler pode tratar especificamente
user, err := repo.FindByID(id)
if err == repository.ErrUserNotFound {
    return http.StatusNotFound
}
```

### 4. Não Exponha GORM Diretamente

```go
// ❌ Ruim: Vaza detalhes de implementação
func (r *Repository) GetDB() *gorm.DB {
    return r.db
}

// ✅ Bom: Interface abstrata
type UserRepository interface {
    FindByID(id string) (*User, error)
    Save(user *User) error
}
```

### 5. Use Transactions Adequadamente

```go
// ✅ Transaction para operações múltiplas
func (r *Repository) CreateUserWithProfile(user *User, profile *Profile) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        
        profile.UserID = user.ID
        if err := tx.Create(profile).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### 6. Paginação Eficiente

```go
// ✅ Retorne total junto com resultados
func (r *Repository) FindAll(page, pageSize int) ([]*User, int64, error) {
    var users []*User
    var total int64
    
    // Count
    r.db.Model(&User{}).Count(&total)
    
    // Query paginada
    offset := (page - 1) * pageSize
    err := r.db.
        Offset(offset).
        Limit(pageSize).
        Find(&users).Error
    
    return users, total, err
}
```

### 7. Índices e Performance

```go
type UserModel struct {
    ID        string `gorm:"primaryKey"`
    Email     string `gorm:"uniqueIndex"`           // Índice único
    Name      string `gorm:"index"`                 // Índice simples
    Status    string `gorm:"index:idx_status_date"` // Índice composto
    CreatedAt time.Time `gorm:"index:idx_status_date"`
}
```

---

## 🧪 Testes de Repository

### Teste Unitário (com Mock)

```go
func TestCreateUserHandler_Save(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    handler := NewCreateUserHandler(mockRepo)
    
    user := &User{Name: "John", Email: "john@example.com"}
    mockRepo.On("Save", mock.Anything, user).Return(nil)
    
    // Act
    err := handler.Handle(context.Background(), &CreateUserCommand{
        Name:  "John",
        Email: "john@example.com",
    })
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Teste de Integração (com DB Real)

```go
func TestMySQLUserRepository_Integration(t *testing.T) {
    // Setup: Container Docker com MySQL
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewMySQLUserRepository(db)
    ctx := context.Background()
    
    // Test Save
    t.Run("Save should persist user", func(t *testing.T) {
        user := &User{
            ID:    "user-1",
            Name:  "John",
            Email: "john@example.com",
        }
        
        err := repo.Save(ctx, user)
        assert.NoError(t, err)
    })
    
    // Test FindByID
    t.Run("FindByID should return saved user", func(t *testing.T) {
        user, err := repo.FindByID(ctx, "user-1")
        
        assert.NoError(t, err)
        assert.Equal(t, "John", user.Name)
        assert.Equal(t, "john@example.com", user.Email)
    })
}
```

---

## 📚 Próximos Passos

- **[Sistema de Módulos](06-modules-system.md)** - Como registrar repositories
- **[Adapters](09-adapters.md)** - Outros tipos de adapters
- **[Boas Práticas](23-best-practices.md)** - Padrões recomendados

---

**[⬅️ Contratos](10-contracts-interfaces.md)** | **[Índice](README.md)** | **[Criando Módulos ➡️](12-creating-modules.md)**
