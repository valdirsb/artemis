# 🎯 Visão Geral do Artemis Framework

## 📋 Índice
- [O que é o Artemis Framework?](#o-que-é-o-artemis-framework)
- [Filosofia e Princípios](#filosofia-e-princípios)
- [Principais Características](#principais-características)
- [Quando Usar o Artemis?](#quando-usar-o-artemis)
- [Comparação com Outros Frameworks](#comparação-com-outros-frameworks)

---

## 🎪 O que é o Artemis Framework?

O **Artemis Framework** é um framework empresarial em Go projetado para facilitar o desenvolvimento de aplicações complexas e escaláveis. Ele combina os melhores padrões de arquitetura de software em uma solução coesa e fácil de usar.

### 🎯 Objetivo

Fornecer uma base sólida e padronizada para o desenvolvimento de aplicações na Agência, permitindo que desenvolvedores foquem na lógica de negócio enquanto o framework cuida da infraestrutura.

### 🏗️ Fundamentos Arquiteturais

O Artemis é construído sobre quatro pilares arquiteturais:

1. **Clean Architecture (Arquitetura Limpa)**
   - Separação clara de responsabilidades em camadas
   - Regras de negócio independentes de frameworks e tecnologias

2. **Hexagonal Architecture (Ports & Adapters)**
   - Aplicação isolada de detalhes de implementação externos
   - Fácil substituição de componentes de infraestrutura

3. **Domain-Driven Design (DDD)**
   - Foco no domínio e lógica de negócio
   - Linguagem ubíqua entre desenvolvedores e stakeholders

4. **CQRS (Command Query Responsibility Segregation)**
   - Separação entre operações de leitura e escrita
   - Otimização independente de cada tipo de operação

---

## 💡 Filosofia e Princípios

### 🎨 Princípios SOLID

O framework é construído seguindo rigorosamente os princípios SOLID:

#### **S** - Single Responsibility Principle
Cada componente tem uma única responsabilidade bem definida.

```go
// ❌ Ruim: Classe faz múltiplas coisas
type UserService struct {
    // Lógica de negócio, persistência, envio de email...
}

// ✅ Bom: Responsabilidades separadas
type CreateUserHandler struct {
    repo      UserRepository
    hasher    PasswordHasher
    emailSvc  EmailService
}
```

#### **O** - Open/Closed Principle
Aberto para extensão, fechado para modificação.

```go
// Interfaces permitem extensão sem modificar código existente
type Repository interface {
    Save(entity interface{}) error
    Find(id string) (interface{}, error)
}

// Você pode criar novas implementações sem tocar no código original
type MySQLRepository struct { /* ... */ }
type MongoRepository struct { /* ... */ }
type RedisRepository struct { /* ... */ }
```

#### **L** - Liskov Substitution Principle
Subtipos devem ser substituíveis por seus tipos base.

```go
// Qualquer implementação de Logger pode ser usada
type Logger interface {
    Info(msg string)
    Error(msg string)
}

// Ambas podem ser usadas intercambiavelmente
type ConsoleLogger struct {}
type FileLogger struct {}
```

#### **I** - Interface Segregation Principle
Clientes não devem depender de interfaces que não usam.

```go
// ❌ Ruim: Interface muito grande
type Repository interface {
    Create(entity interface{}) error
    Read(id string) (interface{}, error)
    Update(entity interface{}) error
    Delete(id string) error
    BulkInsert(entities []interface{}) error
    Transaction(fn func() error) error
}

// ✅ Bom: Interfaces segregadas
type Reader interface {
    Find(id string) (interface{}, error)
}

type Writer interface {
    Save(entity interface{}) error
}
```

#### **D** - Dependency Inversion Principle
Dependa de abstrações, não de implementações concretas.

```go
// ✅ Handler depende de interface, não de implementação
type CreateUserHandler struct {
    repo UserRepository  // Interface, não MySQLRepository
}
```

### 🔄 Outros Princípios Fundamentais

#### 📦 Modularidade
- Aplicação dividida em módulos independentes
- Cada módulo é autocontido e desacoplado
- Facilita manutenção e evolução do código

#### 🧪 Testabilidade
- Código fácil de testar em todos os níveis
- Dependency Injection facilita uso de mocks
- Testes unitários, integração e E2E

#### 🔌 Extensibilidade
- Fácil adicionar novos módulos e funcionalidades
- Adapters permitem integração com diferentes tecnologias
- Sistema de plugins e hooks

#### 📈 Escalabilidade
- Arquitetura preparada para crescimento
- Módulos podem ser distribuídos em microserviços
- Suporte a padrões de escalabilidade (cache, filas, etc.)

---

## ⭐ Principais Características

### 1. 🎯 Sistema de Módulos Auto-Registráveis

Cada módulo se registra automaticamente no framework:

```go
// Módulo se auto-registra com todas suas dependências
func (m *UserModule) Register(registry *container.ModuleRegistry) error {
    registry.RegisterRepository("user", userRepo)
    registry.RegisterApplicationService("user", userService)
    registry.RegisterHTTPHandler("user", httpHandler)
    registry.RegisterGRPCService(grpcService)
    return nil
}
```

**Benefícios:**
- ✅ Menos código boilerplate
- ✅ Configuração centralizada
- ✅ Fácil adicionar/remover módulos
- ✅ Dependências explícitas

### 2. 💉 Dependency Injection Container

Container DI type-safe e poderoso:

```go
// Registro de serviços
container.Register("logger", logger)
container.RegisterSingleton("database", func() interface{} {
    return connectDB()
})

// Resolução com type safety
var logger Logger
container.GetAs("logger", &logger)
```

**Benefícios:**
- ✅ Gerenciamento automático de dependências
- ✅ Lazy loading com singletons
- ✅ Type safety
- ✅ Facilita testes (injeção de mocks)

### 3. 📡 Sistema de Eventos Assíncrono

Event Bus para comunicação entre módulos:

```go
// Publisher
eventBus.Publish("user.created", UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
})

// Subscriber
eventBus.Subscribe("user.created", func(event Event) error {
    // Enviar email de boas-vindas
    return emailService.SendWelcome(event.Email)
})
```

**Benefícios:**
- ✅ Desacoplamento entre módulos
- ✅ Fácil implementar side-effects
- ✅ Auditoria e rastreabilidade
- ✅ Suporte a Event Sourcing

### 4. 🔀 Padrão CQRS

Separação clara entre leitura e escrita:

```go
// Commands (Write) - Modificam estado
type CreateUserCommand struct {
    Name  string
    Email string
}

// Queries (Read) - Apenas leitura
type GetUserQuery struct {
    UserID string
}

type ListUsersQuery struct {
    Page     int
    PageSize int
}
```

**Benefícios:**
- ✅ Código mais organizado e limpo
- ✅ Otimização independente de reads/writes
- ✅ Facilita implementação de CQRS avançado
- ✅ Melhor performance

### 5. 🌐 Multi-Protocol Support

Suporte nativo a múltiplos protocolos:

```yaml
# framework.yaml
protocols:
  http: true     # REST API
  grpc: true     # gRPC API
  websockets: false
  graphql: false
```

**Protocolos Suportados:**
- ✅ HTTP/REST (Gin)
- ✅ gRPC
- 🔄 WebSockets (em desenvolvimento)
- 🔄 GraphQL (planejado)

### 6. 📄 Paginação Integrada

Sistema de paginação consistente:

```go
// HTTP - Cursor pagination
query := &queries.ListUsersQuery{
    PageSize: 20,
    Cursor:   "eyJpZCI6MTIzfQ==",
}

// gRPC - Page-based pagination
query := &queries.ListUsersQuery{
    Page:     1,
    PageSize: 20,
}
```

**Benefícios:**
- ✅ Paginação eficiente e escalável
- ✅ Suporte a cursor e offset pagination
- ✅ Metadados de paginação completos
- ✅ Consistente entre HTTP e gRPC

### 7. 🔧 Configuração Declarativa

Configuração simples via YAML:

```yaml
framework:
  name: "minha-app"
  version: "1.0.0"

core:
  http_server: true
  event_system: true
  logging: true

database:
  mysql: true

modules:
  user: true
  product: true
  order: true
```

### 8. 🧪 Testing-Friendly

Arquitetura facilita todos os tipos de teste:

```go
// Testes unitários - Isola handlers
func TestCreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    handler := NewCreateUserHandler(mockRepo, ...)
    
    // Testa apenas o handler
}

// Testes de integração - Usa DB real
func TestUserRepository(t *testing.T) {
    db := setupTestDB()
    repo := NewMySQLUserRepository(db)
    
    // Testa integração com DB
}
```

---

## 🎯 Quando Usar o Artemis?

### ✅ Casos de Uso Ideais

#### 1. **Aplicações Enterprise**
- Sistemas complexos com múltiplos módulos
- Requisitos de escalabilidade
- Múltiplas integrações

#### 2. **APIs RESTful e gRPC**
- Microserviços
- Backends para mobile/web
- Integrações B2B

#### 3. **Sistemas com Domínio Complexo**
- Regras de negócio sofisticadas
- Múltiplos bounded contexts
- DDD é essencial

#### 4. **Projetos com Múltiplos Desenvolvedores**
- Estrutura padronizada
- Convenções claras
- Fácil onboarding

#### 5. **Aplicações de Longa Duração**
- Manutenção facilitada
- Evolução controlada
- Código sustentável

### ❌ Quando NÃO Usar

#### 1. **Aplicações Simples**
- CRUD básico sem regras complexas
- Pode ser "overhead" desnecessário

#### 2. **Protótipos Rápidos**
- Quando velocidade > estrutura
- MVPs com prazo muito curto

#### 3. **Scripts e Ferramentas CLI**
- Ferramentas simples de linha de comando
- Processamento em batch simples

#### 4. **Projetos Pessoais Pequenos**
- Se a estrutura completa não faz sentido
- Framework pode ser pesado demais

---

## 🆚 Comparação com Outros Frameworks

### vs. Framework Padrão Go (net/http)

| Aspecto | Go Padrão | Artemis |
|---------|-----------|---------|
| **Estrutura** | Você define | Estrutura definida |
| **Boilerplate** | Muito | Mínimo |
| **DI Container** | Manual | Automático |
| **Modularidade** | Manual | Built-in |
| **Curva de Aprendizado** | Baixa | Média |
| **Manutenibilidade** | Média | Alta |
| **Performance** | Ótima | Muito Boa |

### vs. Gin Framework

| Aspecto | Gin | Artemis |
|---------|-----|---------|
| **Foco** | HTTP Routing | Arquitetura Completa |
| **CQRS** | Não | Sim |
| **DDD** | Não | Sim |
| **Módulos** | Não | Sim |
| **gRPC** | Adicionar manualmente | Built-in |
| **Eventos** | Não | Sim |

### vs. Go-kit

| Aspecto | Go-kit | Artemis |
|---------|--------|---------|
| **Foco** | Microserviços | Aplicações Enterprise |
| **Opinião** | Baixa | Alta |
| **Estrutura** | Flexível | Definida |
| **DDD/CQRS** | Você implementa | Built-in |
| **Curva** | Alta | Média |

---

## 🚀 Próximos Passos

Agora que você entende o framework, continue sua jornada:

1. **[Guia de Início Rápido](02-quickstart.md)** - Crie seu primeiro módulo
2. **[Estrutura do Projeto](03-project-structure.md)** - Entenda a organização
3. **[Arquitetura Geral](04-architecture.md)** - Mergulhe na arquitetura

---

## 📚 Recursos Adicionais

- [Documentação de Decisões Arquiteturais (ADRs)](../adr/README.md)
- [Exemplos Completos](25-complete-examples.md)
- [FAQ](26-faq.md)

---

**[⬅️ Voltar ao Índice](README.md)** | **[Próximo: Guia de Início Rápido ➡️](02-quickstart.md)**
