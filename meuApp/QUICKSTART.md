# 🚀 Guia de Início Rápido - Refatoração

> **Antes de começar:** Faça backup do código atual e crie uma branch para refatoração

## 📋 Preparação

### 1. Criar Branch de Trabalho
```bash
git checkout -b refactor/architecture-improvements
git push -u origin refactor/architecture-improvements
```

### 2. Fazer Backup
```bash
# Criar tag do estado atual
git tag -a v1.0.0-before-refactor -m "Estado antes da refatoração"
git push origin v1.0.0-before-refactor

# Ou criar branch de backup
git branch backup/pre-refactor
```

### 3. Configurar Ambiente
```bash
# Instalar ferramentas de análise
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Verificar estado atual
go build ./...
go test ./...
```

---

## 🎯 Começar: Fase 1.1 - Mover Database Models

### Passo 1: Criar arquivo para UserModel

```bash
# Criar arquivo
touch internal/modules/user/repository/user_model.go
```

### Passo 2: Copiar conteúdo

Abra `internal/shared/database/database.go` e copie:
- `UserModel` struct
- `TableName()` method
- `ToContract()` method
- `FromContract()` method

Cole em `internal/modules/user/repository/user_model.go`

### Passo 3: Ajustar package e imports

```go
package repository

import (
    "time"
    "meuApp/pkg/dto" // Quando DTOs forem movidos
)

type UserModel struct {
    ID        string    `gorm:"primaryKey;size:36"`
    Username  string    `gorm:"uniqueIndex;size:50;not null"`
    Email     string    `gorm:"uniqueIndex;size:100;not null"`
    Password  string    `gorm:"size:255;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string {
    return "users"
}

// ... ToContract e FromContract
```

### Passo 4: Atualizar user_repository.go

Trocar imports:
```go
// Antes
"meuApp/internal/shared/database"

// Depois
// Usar diretamente UserModel do mesmo package
```

Trocar referências:
```go
// Antes
var userModel database.UserModel

// Depois
var userModel UserModel
```

### Passo 5: Testar
```bash
cd internal/modules/user
go build ./...

# Volta para raiz
cd ../../..
go build ./...
```

### Passo 6: Repetir para Product e Order

Mesmo processo para:
- `internal/modules/product/repository/product_model.go`
- `internal/modules/order/repository/order_model.go`

### Passo 7: Atualizar database.go

Remover models de `internal/shared/database/database.go`, mantendo apenas:
- Configurações de conexão
- Função `Connect()`
- Função `AutoMigrate()` (atualizar para importar models)

```go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        // Importar de cada módulo
        &userRepo.UserModel{},
        &productRepo.ProductModel{},
        &orderRepo.OrderModel{},
        &orderRepo.OrderItemModel{},
    )
}
```

### Passo 8: Commit
```bash
git add .
git commit -m "refactor: move database models to respective repositories"
```

---

## 📊 Validação de Cada Fase

Após completar cada fase, execute:

### 1. Build Check
```bash
go build ./...
```

### 2. Tests
```bash
go test ./...
```

### 3. Lint
```bash
golangci-lint run
```

### 4. Complexity Check
```bash
gocyclo -over 15 .
```

### 5. Import Check
```bash
goimports -w .
```

---

## 🔄 Workflow Recomendado

### Para cada tarefa:

1. **Ler** a descrição no REFACTORING_PLAN.md
2. **Criar** branch específica (opcional)
   ```bash
   git checkout -b refactor/task-1.1.1
   ```
3. **Implementar** a mudança
4. **Testar** localmente
5. **Commit** com mensagem descritiva
   ```bash
   git commit -m "refactor(phase1): move UserModel to user repository"
   ```
6. **Marcar** no CHECKLIST.md
7. **Atualizar** progresso
8. **Merge** ou continuar

### Mensagens de Commit

Usar convenção:
```
refactor(phase1): descrição curta

Detalhe do que foi feito
Razão da mudança

Refs: #issue (se houver)
```

Exemplos:
```
refactor(phase1): move database models to repositories
refactor(phase2): remove interface duplication
refactor(phase3): add user use cases layer
fix(phase1): correct import paths after reorganization
docs: update architecture documentation
```

---

## 🐛 Troubleshooting

### Erro: Import Cycle

**Sintoma:** `import cycle not allowed`

**Solução:**
1. Verificar dependências circulares
2. Mover interfaces para `ports/`
3. Usar dependency injection

### Erro: Interface Type Assertion

**Sintoma:** `panic: interface conversion`

**Solução:**
1. Verificar registro no container
2. Verificar nome da chave
3. Usar type assertion segura:
```go
handler, ok := container.Get("userHandler")
if !ok {
    // handle error
}
```

### Erro: Database Migration

**Sintoma:** Tabelas não encontradas

**Solução:**
1. Verificar AutoMigrate() atualizado
2. Rodar migrações manualmente
3. Verificar connection string

### Erro: Tests Failing

**Sintoma:** Testes quebram após refactor

**Solução:**
1. Atualizar mocks
2. Atualizar imports
3. Revisar assertions

---

## 📈 Checklist de Preparação

Antes de começar cada fase:

- [ ] Branch criada
- [ ] Backup feito
- [ ] Testes passando
- [ ] Build funcionando
- [ ] Documentação lida
- [ ] Equipe avisada (se aplicável)

Após completar cada fase:

- [ ] Testes passando
- [ ] Build funcionando
- [ ] Lint limpo
- [ ] Documentação atualizada
- [ ] Checklist marcado
- [ ] Commit realizado
- [ ] Code review (se aplicável)

---

## 🎓 Recursos de Aprendizado

### Arquitetura
- [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [DDD Quickly - InfoQ](https://www.infoq.com/minibooks/domain-driven-design-quickly/)

### Go Específico
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Best Practices](https://go.dev/doc/effective_go)
- [SOLID in Go](https://dave.cheney.net/2016/08/20/solid-go-design)

### Padrões
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
- [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)

---

## 💡 Dicas

1. **Não tenha pressa** - Refatoração é um processo iterativo
2. **Teste constantemente** - Após cada mudança significativa
3. **Commite frequentemente** - Pequenos commits são mais fáceis de reverter
4. **Documente decisões** - Crie ADRs para decisões importantes
5. **Peça feedback** - Code review é essencial
6. **Mantenha funcionando** - Sempre mantenha o código compilável
7. **Use feature flags** - Para mudanças grandes, habilite progressivamente

---

## 📞 Precisa de Ajuda?

Se encontrar dificuldades:

1. Consulte o [REFACTORING_PLAN.md](./REFACTORING_PLAN.md) detalhado
2. Verifique a seção de Troubleshooting acima
3. Revise commits anteriores
4. Crie uma issue para discussão
5. Consulte a equipe

---

**Boa refatoração! 🚀**
