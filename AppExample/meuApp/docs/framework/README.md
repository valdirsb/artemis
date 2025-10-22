# 🏛️ Artemis Framework - Documentação Oficial

<div align="center">

![Artemis Framework](https://img.shields.io/badge/Artemis-Framework-blue)
![Version](https://img.shields.io/badge/version-1.0.0-green)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

**Framework modular baseado em Clean Architecture, Hexagonal Architecture, DDD e CQRS para desenvolvimento de aplicações Go enterprise**

[Início Rápido](#início-rápido) • [Arquitetura](#arquitetura) • [Guias](#guias) • [API](#api)

</div>

---

## 📚 Índice da Documentação

### 🚀 Primeiros Passos
- **[01. Visão Geral do Framework](01-overview.md)** - Introdução, conceitos e filosofia
- **[02. Guia de Início Rápido](02-quickstart.md)** - Tutorial passo a passo para criar seu primeiro módulo
- **[03. Estrutura do Projeto](03-project-structure.md)** - Organização de pastas e arquivos

### 🏗️ Arquitetura
- **[04. Arquitetura Geral](04-architecture.md)** - Clean Architecture, Hexagonal Architecture e DDD
- **[05. Padrão CQRS](05-cqrs-pattern.md)** - Command Query Responsibility Segregation
- **[06. Sistema de Módulos](06-modules-system.md)** - Auto-registro e estrutura de módulos
- **[07. Dependency Injection](07-dependency-injection.md)** - Container DI e ModuleRegistry

### 🔌 Componentes Core
- **[08. Sistema de Eventos](08-events-system.md)** - Event Bus, Publishers e Subscribers
- **[09. Adapters](09-adapters.md)** - HTTP, gRPC, Database e outros
- **[10. Contratos e Interfaces](10-contracts-interfaces.md)** - Ports & Adapters pattern
- **[11. Repositórios](11-repositories.md)** - Padrão Repository e persistência

### 📖 Guias Práticos
- **[12. Criando um Novo Módulo](12-creating-modules.md)** - Tutorial completo com exemplos
- **[13. Implementando Commands](13-implementing-commands.md)** - Write operations
- **[14. Implementando Queries](14-implementing-queries.md)** - Read operations com paginação
- **[15. Trabalhando com Eventos](15-working-with-events.md)** - Eventos de domínio e integração
- **[16. Configuração e Bootstrap](16-configuration-bootstrap.md)** - Framework.yaml e inicialização

### 🔧 Recursos Avançados
- **[17. Paginação](17-pagination.md)** - Sistema de paginação para HTTP e gRPC
- **[18. Validação](18-validation.md)** - Validação de dados e regras de negócio
- **[19. Tratamento de Erros](19-error-handling.md)** - Error handling patterns
- **[20. Logging](20-logging.md)** - Sistema de logs estruturados

### 🧪 Testes
- **[21. Estratégia de Testes](21-testing-strategy.md)** - Testes unitários, integração e E2E
- **[22. Mocks e Fixtures](22-mocks-fixtures.md)** - Criando mocks e dados de teste

### 🛠️ Referências
- **[23. Boas Práticas](23-best-practices.md)** - Convenções e recomendações
- **[24. API Reference](24-api-reference.md)** - Referência completa da API
- **[25. Exemplos Completos](25-complete-examples.md)** - Casos de uso reais
- **[26. FAQ](26-faq.md)** - Perguntas frequentes

---

## 🎯 Início Rápido

### Pré-requisitos
```bash
- Go 1.21+
- MySQL 8.0+
- Docker (opcional)
```

### Instalação Rápida

```bash
# 1. Clone o framework
git clone https://github.com/sua-agencia/artemis.git
cd artemis

# 2. Configure as variáveis de ambiente
cp .env.example .env

# 3. Instale as dependências
go mod download

# 4. Execute as migrações
make migrate

# 5. Inicie a aplicação
go run main.go
```

### Criando seu Primeiro Módulo

```bash
# Use o gerador de módulos
make module name=blog

# Isso criará a estrutura completa:
# internal/modules/blog/
#   ├── domain/
#   ├── application/
#   ├── adapters/
#   └── repository/
```

👉 **[Veja o tutorial completo no Guia de Início Rápido →](02-quickstart.md)**

---

## 🏛️ Arquitetura

O Artemis Framework implementa uma arquitetura limpa e modular:

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

✅ **Separação de Responsabilidades** - Cada camada tem uma responsabilidade clara  
✅ **Inversão de Dependências** - Dependa de abstrações, não de implementações  
✅ **Modularidade** - Módulos independentes e desacoplados  
✅ **Testabilidade** - Código fácil de testar em todos os níveis  
✅ **Escalabilidade** - Pronto para crescer horizontalmente  

👉 **[Entenda a arquitetura em detalhes →](04-architecture.md)**

---

## 📦 Estrutura de um Módulo

Cada módulo segue uma estrutura consistente:

```
internal/modules/exemplo/
├── domain/                 # Entidades e regras de negócio
│   ├── entities/
│   ├── value_objects/
│   └── events/
├── application/            # Casos de uso
│   ├── commands/          # Operações de escrita (CUD)
│   ├── queries/           # Operações de leitura (R)
│   └── services/          # Application Services
├── adapters/              # Implementações de interfaces
│   ├── http/             # REST API handlers
│   ├── grpc/             # gRPC handlers
│   └── subscribers/      # Event subscribers
├── repository/            # Persistência de dados
│   └── mysql_repository.go
└── dto/                   # Data Transfer Objects
    ├── requests.go
    └── responses.go
```

👉 **[Veja como criar módulos →](12-creating-modules.md)**

---

## 🔥 Recursos Principais

### 🎯 Auto-registro de Módulos
```go
// Cada módulo se auto-registra no framework
func (m *UserModule) Register(registry *container.ModuleRegistry) error {
    // Registra repositórios, serviços, handlers...
    registry.RegisterRepository("user", userRepo)
    registry.RegisterApplicationService("user", userService)
    registry.RegisterHTTPHandler("user", httpHandler)
    return nil
}
```

### 💉 Dependency Injection
```go
// Container DI poderoso e type-safe
container.Register("service", service)
service, err := container.Get("service")

// Ou com type assertion automático
container.GetAs("service", &myService)
```

### 📡 Sistema de Eventos
```go
// Publique eventos de domínio
eventBus.Publish("user.created", UserCreatedEvent{
    UserID: user.ID,
    Email:  user.Email,
})

// Subscribe para eventos
eventBus.Subscribe("user.created", func(event Event) error {
    // Processa o evento
})
```

### 🌐 Multi-Protocol Support
```yaml
# framework.yaml
protocols:
  http: true    # REST API
  grpc: true    # gRPC API
  websockets: false
  graphql: false
```

---

## 📚 Guias por Caso de Uso

### Para Desenvolvedores Iniciantes
1. [Visão Geral do Framework](01-overview.md)
2. [Guia de Início Rápido](02-quickstart.md)
3. [Estrutura do Projeto](03-project-structure.md)
4. [Criando um Novo Módulo](12-creating-modules.md)

### Para Arquitetos de Software
1. [Arquitetura Geral](04-architecture.md)
2. [Padrão CQRS](05-cqrs-pattern.md)
3. [Sistema de Módulos](06-modules-system.md)
4. [Dependency Injection](07-dependency-injection.md)

### Para DevOps
1. [Configuração e Bootstrap](16-configuration-bootstrap.md)
2. [Estratégia de Testes](21-testing-strategy.md)
3. [Deployment](../DEPLOYMENT.md)

---

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor, leia nosso guia de contribuição.

---

## 📄 Licença

MIT License - veja [LICENSE](../../LICENSE) para detalhes.

---

## 📞 Suporte

- 📧 Email: devteam@sua-agencia.com
- 💬 Slack: #artemis-framework
- 📖 Wiki: [wiki.sua-agencia.com/artemis](https://wiki.sua-agencia.com/artemis)

---

<div align="center">

**Desenvolvido com ❤️ pela Equipe de Desenvolvimento da Agência**

[⬆️ Voltar ao topo](#-artemis-framework---documentação-oficial)

</div>
