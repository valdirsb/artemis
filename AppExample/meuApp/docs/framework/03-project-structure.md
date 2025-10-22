# 📂 Estrutura do Projeto

## 📋 Índice
- [Visão Geral](#visão-geral)
- [Estrutura de Diretórios](#estrutura-de-diretórios)
- [Descrição Detalhada](#descrição-detalhada)
- [Estrutura de um Módulo](#estrutura-de-um-módulo)
- [Convenções de Nomenclatura](#convenções-de-nomenclatura)

---

## 🎯 Visão Geral

O Artemis Framework segue uma estrutura bem definida que facilita navegação, manutenção e escalabilidade do projeto.

### Princípios da Organização

✅ **Separação por Camadas** - Clean Architecture em prática  
✅ **Modularidade** - Cada módulo é independente  
✅ **Coesão Alta** - Arquivos relacionados ficam próximos  
✅ **Acoplamento Baixo** - Módulos não dependem entre si  

---

## 🗂️ Estrutura de Diretórios

```
meuApp/
├── cmd/                        # Entry points da aplicação
│   └── server/
│       └── main.go
│
├── internal/                   # Código privado da aplicação
│   ├── bootstrap/             # Inicialização e configuração
│   │   ├── bootstrap_registry.go
│   │   └── mock.go
│   │
│   ├── modules/               # Módulos da aplicação
│   │   ├── user/             # Módulo de usuários
│   │   ├── product/          # Módulo de produtos
│   │   ├── order/            # Módulo de pedidos
│   │   ├── user_module.go
│   │   ├── product_module.go
│   │   └── order_module.go
│   │
│   └── routes/               # Configuração de rotas
│       └── routes.go
│
├── pkg/                       # Código reutilizável (pode ser compartilhado)
│   ├── adapters/             # Implementações de infraestrutura
│   │   ├── database/        # Adapters de banco de dados
│   │   ├── http/            # Adapters HTTP
│   │   └── logger/          # Adapters de logging
│   │
│   ├── config/              # Configurações da aplicação
│   │   └── config.go
│   │
│   ├── container/           # Dependency Injection Container
│   │   ├── container.go
│   │   └── registry.go
│   │
│   ├── contracts/           # Interfaces e contratos
│   │   ├── infrastructure.go
│   │   ├── interfaces.go
│   │   └── interfaces_*.go
│   │
│   ├── errors/              # Tratamento de erros
│   │   └── errors.go
│   │
│   ├── events/              # Sistema de eventos
│   │   └── event_bus.go
│   │
│   ├── framework/           # Core do framework
│   │   ├── framework.go
│   │   ├── config.go
│   │   ├── interfaces/
│   │   └── providers/
│   │
│   └── proto/              # Protocol Buffers gerados
│       └── *.pb.go
│
├── proto/                   # Definições Protocol Buffers
│   ├── user.proto
│   ├── product.proto
│   └── order.proto
│
├── docs/                    # Documentação
│   ├── framework/          # Documentação do framework (NOVO!)
│   ├── adr/               # Architecture Decision Records
│   ├── guides/            # Guias e tutoriais
│   └── archive/           # Documentos antigos
│
├── tests/                   # Testes
│   ├── e2e/               # Testes end-to-end
│   ├── integration/       # Testes de integração
│   └── fixtures/          # Dados de teste
│
├── migrations/             # Migrações de banco de dados
│   └── *.sql
│
├── scripts/                # Scripts utilitários
│   ├── migrate.sh
│   └── seed.sh
│
├── .env.example           # Exemplo de variáveis de ambiente
├── .gitignore
├── framework.yaml         # Configuração do framework
├── go.mod
├── go.sum
├── main.go               # Entry point principal
├── Makefile             # Comandos make
└── README.md
```

---

## 📖 Descrição Detalhada

### 📁 `/cmd` - Entry Points

Contém os pontos de entrada da aplicação. Cada subdiretório pode representar um binário diferente.

```
cmd/
└── server/
    └── main.go    # Entry point do servidor HTTP/gRPC
```

**Quando usar:**
- Múltiplos binários (server, worker, cli)
- Projetos maiores com vários executáveis

### 📁 `/internal` - Código Privado

Código que não pode ser importado por outros projetos (garantido pelo compilador Go).

#### `/internal/bootstrap`

**Responsabilidade:** Inicialização e configuração da aplicação

```go
// bootstrap_registry.go
func FrameworkBootstrapWithRegistry(configPath string) (*container.Container, *container.ModuleRegistry, *framework.Framework, error)
```

**Arquivos:**
- `bootstrap_registry.go` - Bootstrap com ModuleRegistry
- `mock.go` - Mocks para testes

#### `/internal/modules`

**Responsabilidade:** Módulos de domínio da aplicação

```
modules/
├── user/                  # Módulo completo
│   ├── domain/           # Entidades e regras
│   ├── application/      # Casos de uso
│   ├── adapters/         # HTTP, gRPC, etc
│   ├── repository/       # Persistência
│   └── dto/             # Data Transfer Objects
│
├── product/              # Outro módulo
└── user_module.go       # Registro do módulo
```

**Cada módulo é autocontido e independente!**

#### `/internal/routes`

**Responsabilidade:** Configuração centralizada de rotas

```go
// routes.go
func RegisterRoutes(router *gin.Engine, registry *container.ModuleRegistry) {
    api := router.Group("/api/v1")
    registry.RegisterHTTPRoutes(api)
}
```

### 📁 `/pkg` - Código Compartilhável

Código que pode ser importado por outros projetos. Use com moderação!

#### `/pkg/adapters`

**Implementações de infraestrutura:**

```
adapters/
├── database/
│   └── mysql/
│       ├── connection.go
│       └── migrations.go
│
├── http/
│   └── client.go
│
└── logger/
    ├── console_logger.go
    └── file_logger.go
```

#### `/pkg/config`

**Configurações da aplicação:**

```go
// config.go
type Config struct {
    DBHost     string `env:"DB_HOST"`
    DBPort     int    `env:"DB_PORT"`
    ServerPort int    `env:"SERVER_PORT"`
}

func LoadConfig() (*Config, error)
```

#### `/pkg/container`

**Sistema de Dependency Injection:**

```
container/
├── container.go       # DI Container
├── registry.go        # ModuleRegistry
└── registry_test.go   # Testes
```

**Componentes principais:**
- `Container` - Gerencia dependências
- `ModuleRegistry` - Registra módulos
- Interfaces para HTTP, gRPC, Repositories

#### `/pkg/contracts`

**Interfaces e contratos do sistema:**

```
contracts/
├── infrastructure.go      # Contratos de infraestrutura
├── interfaces.go         # Interfaces gerais
├── interfaces_user.go    # Interfaces do módulo User
├── interfaces_product.go # Interfaces do módulo Product
└── interfaces_order.go   # Interfaces do módulo Order
```

**Exemplos:**
```go
// infrastructure.go
type Logger interface {
    Info(msg string, args ...interface{})
    Error(msg string, args ...interface{})
}

type EmailService interface {
    Send(to, subject, body string) error
}

// interfaces_user.go
type UserRepository interface {
    Save(user *User) error
    FindByID(id string) (*User, error)
}
```

#### `/pkg/events`

**Sistema de eventos assíncrono:**

```go
// event_bus.go
type EventBus struct {
    subscribers map[string][]EventHandler
}

func (eb *EventBus) Publish(eventName string, data interface{})
func (eb *EventBus) Subscribe(eventName string, handler EventHandler)
```

#### `/pkg/framework`

**Core do framework:**

```
framework/
├── framework.go        # Framework principal
├── config.go          # Configuração (framework.yaml)
│
├── interfaces/
│   └── module.go      # Interface Module
│
└── providers/
    ├── http/          # HTTP Provider
    ├── grpc/          # gRPC Provider
    └── database/      # Database Provider
```

### 📁 `/proto` - Protocol Buffers

Definições `.proto` e código gerado:

```
proto/
├── user.proto         # Definição gRPC do User
├── user.pb.go        # Código gerado
├── user_grpc.pb.go   # Código gRPC gerado
├── product.proto
└── order.proto
```

**Gerar código:**
```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

### 📁 `/docs` - Documentação

```
docs/
├── framework/          # 📚 Documentação do Framework (PRINCIPAL)
│   ├── README.md
│   ├── 01-overview.md
│   ├── 02-quickstart.md
│   └── ...
│
├── adr/               # Architecture Decision Records
│   ├── 001-clean-architecture.md
│   ├── 002-cqrs-pattern.md
│   └── 003-module-auto-registration.md
│
├── guides/            # Guias específicos
│   ├── QUICKSTART.md
│   └── CODE_EXAMPLES.md
│
└── archive/           # Documentos antigos/refatoração
    └── ...
```

### 📁 `/tests` - Testes

```
tests/
├── e2e/              # Testes end-to-end
│   ├── user_test.go
│   └── order_test.go
│
├── integration/      # Testes de integração
│   └── repository_test.go
│
└── fixtures/         # Dados para testes
    └── users.json
```

### 📄 Arquivos Raiz

- **`main.go`** - Entry point principal
- **`framework.yaml`** - Configuração do framework
- **`.env.example`** - Template de variáveis de ambiente
- **`Makefile`** - Comandos úteis (build, test, migrate)
- **`go.mod/go.sum`** - Dependências Go

---

## 🏗️ Estrutura de um Módulo

Cada módulo segue a mesma estrutura padronizada:

```
internal/modules/exemplo/
│
├── domain/                          # 🎯 CAMADA DE DOMÍNIO
│   ├── entities/                   # Entidades do domínio
│   │   ├── exemplo.go
│   │   └── exemplo_test.go
│   │
│   ├── value_objects/              # Value Objects
│   │   └── status.go
│   │
│   ├── events/                     # Eventos de domínio
│   │   ├── exemplo_created.go
│   │   └── exemplo_updated.go
│   │
│   └── services/                   # Domain Services
│       └── exemplo_domain_service.go
│
├── application/                     # 💼 CAMADA DE APLICAÇÃO
│   ├── commands/                   # Commands (Write)
│   │   ├── create_exemplo.go
│   │   ├── update_exemplo.go
│   │   ├── delete_exemplo.go
│   │   └── *_test.go
│   │
│   ├── queries/                    # Queries (Read)
│   │   ├── get_exemplo.go
│   │   ├── list_exemplos.go
│   │   └── *_test.go
│   │
│   └── services/                   # Application Services
│       └── exemplo_service.go
│
├── adapters/                        # 🔌 CAMADA DE APRESENTAÇÃO
│   ├── http/                       # REST API
│   │   ├── handler.go
│   │   ├── middleware.go
│   │   └── handler_test.go
│   │
│   ├── grpc/                       # gRPC
│   │   ├── service.go
│   │   └── service_test.go
│   │
│   └── subscribers/                # Event Subscribers
│       └── exemplo_subscriber.go
│
├── repository/                      # 💾 CAMADA DE INFRAESTRUTURA
│   ├── repository.go               # Interface
│   ├── mysql_repository.go         # Implementação MySQL
│   ├── cache_repository.go         # Implementação com Cache
│   └── *_test.go
│
├── dto/                            # Data Transfer Objects
│   ├── request.go                  # DTOs de entrada
│   ├── response.go                 # DTOs de saída
│   └── mapper.go                   # Conversões
│
├── errors.go                       # Erros específicos do módulo
├── README.md                       # Documentação do módulo
└── exemplo_module.go              # ⚙️ REGISTRO DO MÓDULO
```

### Responsabilidades de Cada Camada

#### 🎯 Domain (Domínio)

**O quê:** Regras de negócio puras  
**Depende de:** Nada (totalmente independente)

```go
// domain/entities/user.go
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
```

#### 💼 Application (Aplicação)

**O quê:** Casos de uso e orquestração  
**Depende de:** Domain, Interfaces

```go
// application/commands/create_user.go
type CreateUserHandler struct {
    repo UserRepository  // Interface!
}

func (h *CreateUserHandler) Handle(cmd *CreateUserCommand) error {
    user := &User{...}
    return h.repo.Save(user)
}
```

#### 🔌 Adapters (Adaptadores)

**O quê:** Interface com mundo externo  
**Depende de:** Application

```go
// adapters/http/handler.go
type UserHTTPHandler struct {
    service UserApplicationService
}

func (h *UserHTTPHandler) CreateUser(c *gin.Context) {
    // Converte HTTP Request → Command
    // Chama service
    // Converte resultado → HTTP Response
}
```

#### 💾 Repository (Repositório)

**O quê:** Persistência de dados  
**Depende de:** Domain (entities)

```go
// repository/mysql_repository.go
type MySQLUserRepository struct {
    db *gorm.DB
}

func (r *MySQLUserRepository) Save(user *User) error {
    // Converte Entity → Model
    // Salva no banco
}
```

---

## 📝 Convenções de Nomenclatura

### Arquivos

```
✅ snake_case para arquivos: user_repository.go
✅ *_test.go para testes: user_repository_test.go
✅ Nomes descritivos: create_user_command.go
```

### Packages

```go
✅ Singular e minúsculo: package user (não users)
✅ Significativo: package repository (não repo)
✅ Evite generic: entities, não models
```

### Structs e Interfaces

```go
✅ PascalCase: type UserRepository interface {}
✅ Interfaces terminam em -er quando possível: Reader, Writer
✅ Implementations especificam tecnologia: MySQLUserRepository
```

### Métodos

```go
✅ PascalCase: func (h *Handler) CreateUser()
✅ Verbo + Substantivo: CreateUser, GetUser, ListUsers
✅ Bool methods começam com Is/Has: IsValid(), HasPermission()
```

### Constantes

```go
✅ PascalCase ou UPPER_CASE:
const MaxRetries = 3
const DEFAULT_PAGE_SIZE = 20
```

---

## 🎯 Boas Práticas

### 1. Um Arquivo por Responsabilidade

```
❌ Ruim:
user_handlers.go  (todos os handlers em um arquivo)

✅ Bom:
create_user.go
update_user.go
delete_user.go
```

### 2. Testes ao Lado do Código

```
user/
├── create_user.go
└── create_user_test.go  ← mesmo diretório
```

### 3. README em Módulos Complexos

```
user/
├── README.md  ← Explica o módulo
├── domain/
├── application/
└── ...
```

### 4. Interfaces Próximas ao Uso

```go
// Na application layer
type UserRepository interface {
    Save(user *User) error
}

// Na infrastructure layer
type MySQLUserRepository struct {
    // implementa UserRepository
}
```

---

## 📚 Próximos Passos

- **[Arquitetura Geral](04-architecture.md)** - Entenda a arquitetura em profundidade
- **[Sistema de Módulos](06-modules-system.md)** - Como módulos funcionam
- **[Criando um Módulo](12-creating-modules.md)** - Tutorial detalhado

---

**[⬅️ Guia de Início Rápido](02-quickstart.md)** | **[Índice](README.md)** | **[Arquitetura ➡️](04-architecture.md)**
