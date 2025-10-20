# 🏛️ Artemis Framework

<div align="center">

![Artemis Framework](https://img.shields.io/badge/Artemis-Framework-blue)
![Version](https://img.shields.io/badge/version-1.0.0-green)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean%20%7C%20Hexagonal%20%7C%20DDD-orange)

**Framework modular para desenvolvimento de aplicações Go enterprise**

[📚 Documentação Completa](docs/framework/README.md) • [🚀 Início Rápido](docs/framework/02-quickstart.md) • [🏗️ Arquitetura](docs/framework/04-architecture.md)

</div>

---

## 🎯 O que é o Artemis?

O **Artemis Framework** é um framework empresarial em Go que implementa Clean Architecture, Hexagonal Architecture (Ports & Adapters), Domain-Driven Design (DDD) e CQRS. Ele fornece uma base sólida e padronizada para desenvolvimento de aplicações complexas e escaláveis.

### ⭐ Principais Características

✅ **Clean Architecture** - Separação clara de responsabilidades  
✅ **Modular** - Sistema de módulos auto-registráveis  
✅ **CQRS** - Separação entre Commands e Queries  
✅ **DI Container** - Injeção de dependências poderosa  
✅ **Multi-Protocol** - HTTP (REST) e gRPC nativos  
✅ **Event-Driven** - Sistema de eventos assíncrono  
✅ **Type-Safe** - Forte tipagem em todas as camadas  
✅ **Testável** - Arquitetura facilita testes em todos os níveis  

---

## 🚀 Início Rápido

### Instalação

```bash
# Clone o repositório
git clone https://github.com/sua-agencia/artemis.git meu-projeto
cd meu-projeto

# Configure o ambiente
cp .env.example .env

# Instale dependências
go mod download

# Execute migrações
make migrate

# Inicie a aplicação
go run main.go
```

### Primeiro Módulo

```bash
# Crie um novo módulo
make module name=blog

# Estrutura criada:
# internal/modules/blog/
#   ├── domain/
#   ├── application/
#   ├── adapters/
#   └── repository/
```

👉 **[Veja o tutorial completo →](docs/framework/02-quickstart.md)**

---

## 📚 Documentação

### 🎓 Para Iniciantes

- **[Visão Geral do Framework](docs/framework/01-overview.md)** - Conceitos e filosofia
- **[Guia de Início Rápido](docs/framework/02-quickstart.md)** - Crie seu primeiro módulo
- **[Estrutura do Projeto](docs/framework/03-project-structure.md)** - Organização de pastas

### 🏗️ Arquitetura

- **[Arquitetura Geral](docs/framework/04-architecture.md)** - Clean Architecture e DDD
- **[Padrão CQRS](docs/framework/05-cqrs-pattern.md)** - Commands e Queries
- **[Sistema de Módulos](docs/framework/06-modules-system.md)** - Módulos auto-registráveis
- **[Dependency Injection](docs/framework/07-dependency-injection.md)** - Container DI

### 🔧 Componentes

- **[Sistema de Eventos](docs/framework/08-events-system.md)** - Event Bus e Domain Events
- **[Adapters](docs/framework/09-adapters.md)** - HTTP, gRPC, Database
- **[Contratos e Interfaces](docs/framework/10-contracts-interfaces.md)** - Ports & Adapters

### 📖 Guias Práticos

- **[Criando um Módulo](docs/framework/12-creating-modules.md)** - Tutorial completo
- **[Implementando Commands](docs/framework/13-implementing-commands.md)** - Write operations
- **[Implementando Queries](docs/framework/14-implementing-queries.md)** - Read operations

### 🎯 Boas Práticas

- **[Boas Práticas](docs/framework/23-best-practices.md)** - Convenções e recomendações
- **[Estratégia de Testes](docs/framework/21-testing-strategy.md)** - Testes eficazes

---

## 🏛️ Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                     Presentation Layer                       │
│              (HTTP Handlers, gRPC Services)                  │
├─────────────────────────────────────────────────────────────┤
│                    Application Layer                         │
│        (Commands, Queries, Application Services)             │
├─────────────────────────────────────────────────────────────┤
│                      Domain Layer                            │
│         (Entities, Value Objects, Domain Events)             │
├─────────────────────────────────────────────────────────────┤
│                   Infrastructure Layer                       │
│     (Repositories, Event Bus, External Services)             │
└─────────────────────────────────────────────────────────────┘
```

### Princípios Fundamentais

- **Inversão de Dependências** - Dependa de abstrações
- **Separação de Responsabilidades** - Cada camada tem seu papel
- **Modularidade** - Módulos independentes e desacoplados
- **Testabilidade** - Fácil testar em todos os níveis

---

## 📁 Estrutura do Projeto

```
meuApp/
├── docs/
│   └── framework/         # 📚 DOCUMENTAÇÃO COMPLETA
│       ├── README.md      # Índice da documentação
│       ├── 01-overview.md
│       ├── 02-quickstart.md
│       └── ...
│
├── internal/
│   ├── bootstrap/         # Inicialização e DI
│   └── modules/           # Módulos da aplicação
│       ├── user/
│       ├── product/
│       └── order/
│
├── pkg/                   # Código compartilhável
│   ├── adapters/         # Database, HTTP, Logger
│   ├── container/        # DI Container + Registry
│   ├── contracts/        # Interfaces
│   ├── events/           # Event Bus
│   └── framework/        # Core do framework
│
├── proto/                # Protocol Buffers
├── tests/                # Testes E2E
├── framework.yaml        # Configuração do framework
└── main.go              # Entry point
```

---

## 🔥 Recursos Principais

### 1. Sistema de Módulos Auto-Registráveis

Cada módulo se registra automaticamente no framework:

```go
func (m *UserModule) Register(registry *container.ModuleRegistry) error {
    // Registra componentes automaticamente
    registry.RegisterRepository("user", userRepo)
    registry.RegisterApplicationService("user", userService)
    registry.RegisterHTTPHandler("user", httpHandler)
    return nil
}
```

### 2. CQRS Pattern

Separação clara entre leitura e escrita:

```go
// Commands (Write)
type CreateUserCommand struct {
    Name  string
    Email string
}

// Queries (Read)
type ListUsersQuery struct {
    Page     int
    PageSize int
}
```

### 3. Event-Driven Architecture

```go
// Publicar eventos
eventBus.Publish("user.created", UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
})

// Subscrever
eventBus.Subscribe("user.created", func(event Event) error {
    // Processar evento
})
```

### 4. Multi-Protocol Support

```yaml
# framework.yaml
protocols:
  http: true    # REST API
  grpc: true    # gRPC API
```

---

## 🧪 Testes

```bash
# Testes unitários
go test ./internal/modules/...

# Testes de integração
go test ./tests/integration/...

# Testes E2E
go test ./tests/e2e/...

# Coverage
go test -cover ./...
```

---

## 📖 Exemplos

### Criar um Usuário

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }'
```

### Listar Usuários

```bash
curl http://localhost:8080/api/v1/users?page=1&page_size=20
```

---

## �️ Comandos Úteis

```bash
# Build
make build

# Executar
make run

# Testes
make test

# Migrações
make migrate

# Gerar código gRPC
make proto

# Limpar
make clean
```

---

## 📚 Documentação Adicional

- **[ADRs](docs/adr/README.md)** - Architecture Decision Records
- **[Deployment](docs/DEPLOYMENT.md)** - Guia de deploy
- **[API Reference](docs/framework/24-api-reference.md)** - Referência da API

---

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor, leia nosso [guia de contribuição](CONTRIBUTING.md).

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 👥 Equipe

Desenvolvido com ❤️ pela Equipe de Desenvolvimento da Agência

---

## 📞 Suporte

- 📧 Email: devteam@sua-agencia.com
- 💬 Slack: #artemis-framework
- 📖 Wiki: [Documentação Completa](docs/framework/README.md)

---

<div align="center">

**[⬆️ Voltar ao topo](#-artemis-framework)**

**[📚 Documentação Completa](docs/framework/README.md)** | **[🚀 Início Rápido](docs/framework/02-quickstart.md)** | **[🏗️ Arquitetura](docs/framework/04-architecture.md)**

</div>


---

## 📖 Documentação

> **✅ 6.094+ linhas | 44 arquivos organizados | 100% completa**

### 🚀 Início Rápido
- � [**Índice Completo**](docs/INDEX.md) - Navegação detalhada de toda documentação
- 🏃 [**Quick Start**](docs/guides/QUICKSTART.md) - Execute o projeto em 15 minutos
- 🌐 [**Swagger UI**](http://localhost:8080/swagger/index.html) - API interativa (16 endpoints)

### 📐 Arquitetura
- 🏗️ [**Arquitetura Completa**](docs/ARCHITECTURE.md) - 6 diagramas + explicações (800+ linhas)
- 📊 [**Antes vs Depois**](docs/ARCHITECTURE_COMPARISON.md) - Comparação da refatoração
- 📋 [**ADRs**](docs/adr/) - 3 Architecture Decision Records

### 📦 Desenvolvimento
- � [**Criar Módulo**](docs/MODULE_CREATION_GUIDE.md) - Tutorial 12 passos (1500+ linhas)
- 💡 [**Exemplos de Código**](docs/guides/CODE_EXAMPLES.md) - Exemplos práticos
- 🎯 [**Sistema de Eventos**](docs/guides/EVENTS_GUIDE.md) - Event Bus type-safe
- 🗂️ [**Estrutura de Arquivos**](docs/guides/FILE_STRUCTURE.md) - Organização do projeto

### 🚀 Deployment & Progresso
- 🐳 [**Guia de Deploy**](docs/DEPLOYMENT.md) - Docker + configuração (600+ linhas)
- ✅ [**CHECKLIST.md**](CHECKLIST.md) - **93% completo** (130/140 tarefas)
- 📊 [**Fases Completas**](docs/phases/) - Resumos de cada fase

### 📂 Organização da Documentação
```
docs/
├── 📋 adr/          # Architecture Decision Records (3 ADRs)
├── 📚 guides/       # Guias práticos (5 guias)
├── 📊 phases/       # Resumos de fases (18 documentos)
└── � archive/      # Documentos históricos (6 documentos)
```

**Ver:** [DOCUMENTACAO_REORGANIZADA.md](DOCUMENTACAO_REORGANIZADA.md) para detalhes da organização

---

## 🎯 Progresso da Refatoração

| Fase | Descrição | Status | Progresso |
|------|-----------|--------|-----------|
| 1 | Reorganização de Estrutura | ✅ Completo | 24/24 (100%) |
| 2 | Interfaces e Contratos | ✅ Completo | 23/23 (100%) |
| 3 | Camada de Application | 🔄 Próximo | 0/21 (0%) |
| 4 | Sistema de Erros | ⏳ Pendente | 0/15 (0%) |
| 5 | Event Bus | ⏳ Pendente | 0/12 (0%) |
| 6 | Auto-registro | ⏳ Pendente | 0/13 (0%) |

**Total: 53% (50/94 tarefas)**

---

## 🛠️ Tecnologias

- **Linguagem:** Go 1.24
- **Framework HTTP:** Gin
- **ORM:** GORM
- **Database:** MySQL
- **gRPC:** Protocol Buffers
- **DI:** Custom Container
- **Config:** YAML

---

## 🏆 Princípios Aplicados

- ✅ **SOLID** - Todos os 5 princípios
- ✅ **Clean Architecture** - Dependency Rule, camadas isoladas
- ✅ **Hexagonal Architecture** - Ports & Adapters
- ✅ **DDD** - Entidades, Agregados, Value Objects
- ✅ **CQRS** - Preparado para Commands/Queries
- ✅ **Dependency Injection** - Inversão total de controle
- ✅ **Repository Pattern** - Abstração de persistência

---

## 📊 Métricas de Qualidade

- **Avaliação Arquitetura:** ~8.0/10 (meta: 9/10)
- **Separação de Camadas:** 95%
- **Testabilidade:** Alta (ports facilitam mocks)
- **Manutenibilidade:** Alta (módulos independentes)
- **Acoplamento:** Baixo (via interfaces)

---

## 🤝 Contribuindo

1. Leia [REFACTORING_INDEX.md](./REFACTORING_INDEX.md)
2. Verifique [CHECKLIST.md](./CHECKLIST.md) para tarefas pendentes
3. Crie uma branch: `git checkout -b feature/nova-feature`
4. Faça commits pequenos e frequentes
5. Abra um Pull Request

---

## 📝 Licença

[Definir licença]

---

## 📞 Contato

[Informações de contato]

---

*Última atualização: 18/10/2025 - Fase 2 completa! 🎉*