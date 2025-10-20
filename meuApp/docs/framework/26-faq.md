# ❓ FAQ - Perguntas Frequentes

## 📋 Índice
- [Instalação e Setup](#instalação-e-setup)
- [Desenvolvimento](#desenvolvimento)
- [Arquitetura](#arquitetura)
- [Database e Migrations](#database-e-migrations)
- [Testes](#testes)
- [Deployment](#deployment)
- [Performance](#performance)
- [Troubleshooting](#troubleshooting)

---

## 🔧 Instalação e Setup

### Como instalar o Artemis Framework?

```bash
# 1. Clone o repositório
git clone <repo-url>
cd meuApp

# 2. Copie as variáveis de ambiente
cp .env.example .env

# 3. Configure o .env
nano .env

# 4. Inicie o MySQL
make docker-up

# 5. Execute as migrations
make migrate

# 6. Inicie a aplicação
go run main.go
```

**Documentação:** [02-quickstart.md](02-quickstart.md)

---

### Quais são os requisitos mínimos?

**Software Necessário:**
- Go 1.21 ou superior
- Docker & Docker Compose (para MySQL)
- Make (opcional, mas recomendado)

**Hardware Mínimo:**
- CPU: 2 cores
- RAM: 4 GB
- Disk: 10 GB

---

### Como configurar o banco de dados?

**1. Usando Docker (Recomendado):**

```bash
# Inicie o MySQL
make docker-up

# Ou manualmente
docker-compose up -d mysql
```

**2. MySQL Local:**

```env
# .env
DB_HOST=localhost
DB_PORT=3306
DB_NAME=artemis
DB_USER=root
DB_PASSWORD=secret
```

**3. Teste a conexão:**

```bash
# Via make
make db-shell

# Ou manualmente
mysql -h localhost -u root -p artemis
```

**Documentação:** [16-configuration-bootstrap.md](16-configuration-bootstrap.md)

---

### Erro: "database connection failed"

**Causa Comum:** MySQL não está rodando ou credenciais incorretas

**Soluções:**

```bash
# 1. Verifique se MySQL está rodando
docker ps | grep mysql

# 2. Verifique os logs
docker logs meuApp_mysql_1

# 3. Aguarde MySQL inicializar (pode levar 30s)
sleep 30

# 4. Teste conexão
make db-shell

# 5. Verifique as credenciais no .env
cat .env | grep DB_
```

---

## 💻 Desenvolvimento

### Como criar um novo módulo?

**Passo a passo completo:**

```bash
# 1. Crie a estrutura de diretórios
mkdir -p internal/modules/category/{domain,application/{commands,queries},repository,adapters/{http,grpc},ports,dto,tests}

# 2. Crie as entidades de domínio
# internal/modules/category/domain/category.go

# 3. Crie os ports (interfaces)
# internal/modules/category/ports/ports.go

# 4. Implemente Commands e Queries
# internal/modules/category/application/commands/
# internal/modules/category/application/queries/

# 5. Implemente Repository
# internal/modules/category/repository/

# 6. Crie HTTP Handlers
# internal/modules/category/adapters/http/

# 7. Registre o módulo
# internal/modules/category_module.go
```

**Documentação Completa:** [12-creating-modules.md](12-creating-modules.md)

---

### Como implementar um novo Command?

**Template básico:**

```go
// CreateCategoryCommand
type CreateCategoryCommand struct {
    Name        string
    Description string
}

// CreateCategoryHandler
type CreateCategoryHandler struct {
    repo     ports.CategoryRepository
    eventBus contracts.EventPublisher
    logger   contracts.Logger
}

func (h *CreateCategoryHandler) Handle(ctx context.Context, cmd CreateCategoryCommand) (*domain.Category, error) {
    // 1. Validar
    if cmd.Name == "" {
        return nil, errors.NewValidationError("name is required", nil)
    }

    // 2. Criar entidade
    category, err := domain.NewCategory(uuid.New().String(), cmd.Name, cmd.Description)
    if err != nil {
        return nil, err
    }

    // 3. Salvar
    if err := h.repo.Create(ctx, category); err != nil {
        return nil, errors.WrapError(err, "failed to create category")
    }

    // 4. Publicar evento
    event := contracts.Event{
        Type:    "category.created",
        Payload: category,
    }
    h.eventBus.Publish(ctx, event)

    return category, nil
}
```

**Documentação:** [13-implementing-commands.md](13-implementing-commands.md)

---

### Como implementar paginação?

**No Repository:**

```go
func (r *Repository) ListPaginated(ctx context.Context, filters Filters) (*PaginatedResult, error) {
    query := r.db.WithContext(ctx).Model(&Model{})

    // Aplicar filtros
    if filters.CategoryID != nil {
        query = query.Where("category_id = ?", *filters.CategoryID)
    }

    // Contar total
    var totalItems int64
    if err := query.Count(&totalItems).Error; err != nil {
        return nil, err
    }

    // Paginação
    page := filters.Page
    if page < 1 {
        page = 1
    }
    pageSize := filters.PageSize
    if pageSize < 1 {
        pageSize = 10
    }

    offset := (page - 1) * pageSize
    
    var items []Model
    if err := query.Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
        return nil, err
    }

    totalPages := int(totalItems) / pageSize
    if int(totalItems)%pageSize != 0 {
        totalPages++
    }

    return &PaginatedResult{
        Items:      items,
        TotalItems: totalItems,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    }, nil
}
```

**Documentação:** [14-implementing-queries.md](14-implementing-queries.md)

---

### Como adicionar validação customizada?

**1. Criar validator customizado:**

```go
// pkg/validation/custom_validators.go
func init() {
    validate.RegisterValidation("cpf", validateCPF)
}

func validateCPF(fl validator.FieldLevel) bool {
    cpf := fl.Field().String()
    // Lógica de validação
    return isValidCPF(cpf)
}
```

**2. Usar na struct:**

```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=3"`
    Email string `json:"email" validate:"required,email"`
    CPF   string `json:"cpf" validate:"required,cpf"`
}
```

**Documentação:** [18-validation.md](18-validation.md)

---

### Como usar eventos (Event Bus)?

**1. Publicar evento:**

```go
// No Command Handler
event := contracts.Event{
    Type:      "user.created",
    Payload:   user,
    Timestamp: time.Now(),
}
h.eventBus.Publish(ctx, event)
```

**2. Criar subscriber:**

```go
// internal/modules/notification/subscribers/user_created_subscriber.go
func HandleUserCreated(ctx context.Context, event contracts.Event) error {
    user := event.Payload.(*domain.User)
    
    // Enviar email de boas-vindas
    return emailService.SendWelcome(user.Email())
}
```

**3. Registrar subscriber:**

```go
// No módulo
eventBus.Subscribe("user.created", HandleUserCreated)
```

**Documentação:** [08-events-system.md](08-events-system.md), [15-working-with-events.md](15-working-with-events.md)

---

## 🏗️ Arquitetura

### Qual a diferença entre Command e Query?

**Command (Write):**
- ✅ **Modifica** estado da aplicação
- ✅ Retorna entidade criada/modificada ou erro
- ✅ Dispara eventos de domínio
- ✅ Validações complexas
- ❌ Não otimizado para leitura

**Query (Read):**
- ✅ **Lê** dados sem modificar
- ✅ Otimizado para leitura (joins, paginação)
- ✅ Pode retornar DTOs customizados
- ❌ Não deve modificar estado
- ❌ Não dispara eventos

**Exemplo:**

```go
// ✅ Command (Write)
CreateProductCommand    // Cria produto
UpdateStockCommand      // Atualiza estoque
DeleteProductCommand    // Deleta produto

// ✅ Query (Read)
GetProductByIDQuery     // Busca por ID
ListProductsQuery       // Lista com filtros
SearchProductsQuery     // Busca textual
```

**Documentação:** [05-cqrs-pattern.md](05-cqrs-pattern.md)

---

### O que são Ports e Adapters?

**Ports (Interfaces):**
- Definem **o que** o módulo precisa
- Ficam em `ports/ports.go`
- Exemplo: `ProductRepository`, `EmailService`

**Adapters (Implementações):**
- Definem **como** implementar
- Ficam em `adapters/`, `repository/`
- Exemplo: `MySQLProductRepository`, `SendGridEmailService`

**Benefícios:**
- ✅ Facilita testes (mock ports)
- ✅ Troca de implementação sem mudar lógica
- ✅ Dependências invertidas (Clean Architecture)

**Documentação:** [04-architecture.md](04-architecture.md), [10-contracts-interfaces.md](10-contracts-interfaces.md)

---

### Por que usar DDD (Domain-Driven Design)?

**Benefícios:**

1. **Linguagem Ubíqua**: Código reflete negócio
2. **Separação clara**: Domain isolado de infraestrutura
3. **Invariantes protegidas**: Validações no domain
4. **Eventos de domínio**: Comunicação entre módulos
5. **Testabilidade**: Domain puro, sem dependências

**Exemplo:**

```go
// ✅ Bom - Domain com invariantes
func (p *Product) DecreaseStock(quantity int) error {
    if quantity <= 0 {
        return errors.New("quantity must be positive")
    }
    if p.stock < quantity {
        return errors.New("insufficient stock")
    }
    p.stock -= quantity
    return nil
}

// ❌ Ruim - Repository com lógica de negócio
func (r *Repository) DecreaseStock(id string, quantity int) error {
    // Lógica de negócio no repositório!
}
```

**Documentação:** [04-architecture.md](04-architecture.md)

---

## 💾 Database e Migrations

### Como criar migrations?

**1. Criar arquivo de migration:**

```sql
-- migrations/20231020_create_categories_table.up.sql
CREATE TABLE categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_name ON categories(name);
```

**2. Criar rollback:**

```sql
-- migrations/20231020_create_categories_table.down.sql
DROP TABLE IF EXISTS categories;
```

**3. Executar migration:**

```bash
make migrate
```

**Documentação:** [16-configuration-bootstrap.md](16-configuration-bootstrap.md)

---

### Como fazer rollback de migration?

```bash
# Rollback da última migration
make migrate-rollback

# Ou manualmente
go run cmd/migrate/main.go down
```

---

### Erro: "table already exists"

**Causa:** Migration já foi executada

**Solução:**

```bash
# 1. Verifique migrations executadas
SELECT * FROM schema_migrations;

# 2. Se necessário, limpe o banco
make db-reset

# 3. Execute migrations novamente
make migrate
```

---

## 🧪 Testes

### Como rodar os testes?

```bash
# Todos os testes
make test

# Apenas unit tests
make test-unit

# Apenas integration tests
make test-integration

# E2E tests
make test-e2e

# Com coverage
make test-coverage
```

**Documentação:** [21-testing-strategy.md](21-testing-strategy.md)

---

### Como criar mocks?

**Opção 1: Manual (Simples)**

```go
type MockProductRepository struct {
    mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
    args := m.Called(ctx, product)
    return args.Error(0)
}
```

**Opção 2: Usando testify/mock**

```go
// No teste
mockRepo := new(MockProductRepository)
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)

// Executar teste
handler := NewCreateProductHandler(mockRepo, mockEventBus, mockLogger)
result, err := handler.Handle(ctx, cmd)

// Verificar
mockRepo.AssertExpectations(t)
```

**Documentação:** [21-testing-strategy.md](21-testing-strategy.md)

---

### Erro: "no such table" em testes de integração

**Causa:** SQLite in-memory não tem tabelas criadas

**Solução:**

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    require.NoError(t, err)

    // 🔧 Auto migrate TODAS as tabelas necessárias
    err = db.AutoMigrate(
        &ProductModel{},
        &CategoryModel{},
        &UserModel{},
    )
    require.NoError(t, err)

    return db
}
```

---

## 🚀 Deployment

### Como fazer deploy em produção?

**1. Build da aplicação:**

```bash
# Build otimizado
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go
```

**2. Usando Docker:**

```bash
# Build da imagem
docker build -t meuapp:latest .

# Push para registry
docker tag meuapp:latest registry.com/meuapp:latest
docker push registry.com/meuapp:latest

# Deploy
docker run -d \
  -p 8080:8080 \
  --env-file .env.production \
  registry.com/meuapp:latest
```

**3. Usando Docker Compose:**

```bash
# Production
docker-compose -f docker-compose.prod.yml up -d
```

**Documentação:** [DEPLOYMENT.md](../DEPLOYMENT.md)

---

### Como configurar variáveis de ambiente?

**Arquivo `.env.production`:**

```env
# App
APP_ENV=production
APP_DEBUG=false
APP_PORT=8080

# Database
DB_HOST=prod-db.example.com
DB_PORT=3306
DB_NAME=artemis_prod
DB_USER=artemis_user
DB_PASSWORD=<strong-password>

# JWT
JWT_SECRET=<generate-strong-secret-256-bits>
JWT_EXPIRATION=24h

# Redis
REDIS_HOST=prod-redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=<redis-password>

# Email
SENDGRID_API_KEY=<your-sendgrid-key>
```

**Gerar secrets seguros:**

```bash
# JWT Secret
openssl rand -base64 32

# Password
openssl rand -base64 16
```

**Documentação:** [16-configuration-bootstrap.md](16-configuration-bootstrap.md)

---

### Como configurar health checks?

**1. Endpoint já existe:**

```bash
curl http://localhost:8080/health
```

**Response:**

```json
{
  "status": "healthy",
  "app": "meuApp",
  "env": "production",
  "modules": 3
}
```

**2. Kubernetes:**

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

---

### Erro: "panic: runtime error" em produção

**Causa:** Panic não tratado

**Solução:**

```go
// 1. Adicione recovery middleware (já incluído)
router.Use(gin.Recovery())

// 2. Capture panics em goroutines
go func() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic recovered: %v", r)
        }
    }()
    // Código que pode dar panic
}()

// 3. Use contexto com timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

**Documentação:** [19-error-handling.md](19-error-handling.md)

---

## ⚡ Performance

### Como otimizar queries no database?

**1. Use índices:**

```sql
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_created ON products(created_at);
```

**2. Evite N+1 queries:**

```go
// ❌ Ruim - N+1 queries
products := repo.List(ctx)
for _, p := range products {
    category := categoryRepo.GetByID(ctx, p.CategoryID) // N queries!
}

// ✅ Bom - Preload
db.Preload("Category").Find(&products)
```

**3. Use paginação:**

```go
// Sempre pagine resultados grandes
result, err := repo.ListPaginated(ctx, ports.ProductFilters{
    Page:     1,
    PageSize: 20,
})
```

**4. Select apenas campos necessários:**

```go
// ❌ Ruim
db.Select("*").Find(&products)

// ✅ Bom
db.Select("id", "name", "price").Find(&products)
```

**Documentação:** [14-implementing-queries.md](14-implementing-queries.md)

---

### Como usar cache?

**Exemplo com Redis:**

```go
func (s *ProductService) GetByID(ctx context.Context, id string) (*domain.Product, error) {
    // 1. Tentar cache
    cacheKey := fmt.Sprintf("product:%s", id)
    cached, err := s.cache.Get(ctx, cacheKey)
    if err == nil && cached != nil {
        return cached.(*domain.Product), nil
    }

    // 2. Buscar no DB
    product, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. Salvar no cache (5 minutos)
    s.cache.Set(ctx, cacheKey, product, 5*time.Minute)

    return product, nil
}
```

**Invalidação de cache:**

```go
// Quando produto é atualizado
func (h *UpdateProductHandler) Handle(ctx context.Context, cmd UpdateProductCommand) error {
    // Atualizar DB
    if err := h.repo.Update(ctx, product); err != nil {
        return err
    }

    // Invalidar cache
    cacheKey := fmt.Sprintf("product:%s", product.ID())
    h.cache.Delete(ctx, cacheKey)

    return nil
}
```

---

### Aplicação está lenta, como debugar?

**1. Profile CPU:**

```bash
# Iniciar com profiling
go run main.go -cpuprofile=cpu.prof

# Analisar
go tool pprof cpu.prof
```

**2. Profile Memory:**

```bash
# Endpoint de profiling (adicionar)
import _ "net/http/pprof"

// Analisar
go tool pprof http://localhost:8080/debug/pprof/heap
```

**3. Logs estruturados:**

```go
// Medir tempo de execução
start := time.Now()
result, err := handler.Handle(ctx, cmd)
duration := time.Since(start)

logger.Info("Command executed", 
    contracts.Field{Key: "command", Value: "CreateProduct"},
    contracts.Field{Key: "duration_ms", Value: duration.Milliseconds()},
)
```

**4. Database query logs:**

```go
// Habilitar logs SQL
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // ou logger.Silent
})
```

---

## 🔧 Troubleshooting

### Erro: "cannot find package"

**Causa:** Dependências não instaladas

**Solução:**

```bash
# Baixar dependências
go mod download

# Sincronizar go.mod e go.sum
go mod tidy

# Verificar vendor (se usado)
go mod vendor
```

---

### Erro: "address already in use"

**Causa:** Porta 8080 já está em uso

**Solução:**

```bash
# 1. Encontrar processo usando a porta
lsof -i :8080
# ou
netstat -tulpn | grep 8080

# 2. Matar processo
kill -9 <PID>

# 3. Ou usar porta diferente
APP_PORT=8081 go run main.go
```

---

### Erro: "context deadline exceeded"

**Causa:** Operação levou muito tempo

**Solução:**

```go
// 1. Aumentar timeout
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()

// 2. Ou usar contexto sem timeout para operações longas
ctx := context.Background()

// 3. Verificar queries lentas
// - Adicionar índices
// - Otimizar queries
// - Usar paginação
```

---

### Erro: "too many open files"

**Causa:** Limite de file descriptors atingido

**Solução:**

```bash
# 1. Verificar limite atual
ulimit -n

# 2. Aumentar limite (temporário)
ulimit -n 4096

# 3. Permanente (Linux)
# Editar /etc/security/limits.conf
* soft nofile 4096
* hard nofile 8192

# 4. Verificar conexões de database
# Ajustar MaxOpenConns no GORM
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
```

---

### Como debugar "panic: runtime error: invalid memory address"?

**Causa:** Nil pointer dereference

**Solução:**

```go
// ❌ Ruim - Pode causar panic
user := userRepo.GetByID(ctx, id)
fmt.Println(user.Name()) // panic se user = nil

// ✅ Bom - Verificar nil
user, err := userRepo.GetByID(ctx, id)
if err != nil {
    return nil, err
}
if user == nil {
    return nil, errors.NewNotFoundError("User", id)
}
fmt.Println(user.Name())
```

**Usar linter:**

```bash
# Instalar
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Executar
golangci-lint run
```

---

### Testes passam localmente mas falham no CI

**Causas comuns:**

1. **Timezone diferente:**

```go
// ❌ Ruim
now := time.Now()

// ✅ Bom
now := time.Now().UTC()
```

2. **Race conditions:**

```bash
# Rodar com race detector
go test -race ./...
```

3. **Dependências de ordem:**

```go
// ❌ Ruim - Testes dependem de ordem
func TestA(t *testing.T) { /* cria user */ }
func TestB(t *testing.T) { /* usa user de TestA */ }

// ✅ Bom - Testes independentes
func TestA(t *testing.T) { 
    user := createTestUser() 
}
func TestB(t *testing.T) { 
    user := createTestUser() 
}
```

4. **Database state:**

```go
// Sempre limpar banco entre testes
func setupTestDB(t *testing.T) *gorm.DB {
    db := createDB()
    t.Cleanup(func() {
        cleanupDB(db)
    })
    return db
}
```

---

## 📚 Recursos Adicionais

### Onde encontrar mais exemplos?

- **[02-quickstart.md](02-quickstart.md)** - Tutorial completo com Blog module
- **[12-creating-modules.md](12-creating-modules.md)** - Criação de Category module
- **[15-working-with-events.md](15-working-with-events.md)** - 6 exemplos práticos de eventos
- **[21-testing-strategy.md](21-testing-strategy.md)** - Exemplos de testes

---

### Como contribuir com a documentação?

1. **Encontrou erro?** Abra uma issue
2. **Tem sugestão?** Abra uma PR
3. **Dúvida não respondida?** Adicione aqui

---

### Onde obter suporte?

- 📖 **Documentação:** [docs/framework/README.md](README.md)
- 🐛 **Issues:** GitHub Issues
- 💬 **Discussões:** GitHub Discussions
- 📧 **Email:** support@artemis.dev

---

**[⬅️ Complete Examples](25-complete-examples.md)** | **[Índice](README.md)** | **[Overview ➡️](01-overview.md)**
