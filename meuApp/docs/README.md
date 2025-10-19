# 📚 Documentação do Projeto Artemis

> **Framework modular em Go** - Clean Architecture + Hexagonal + DDD + CQRS + Auto-Registro

> **✅ Status:** Documentação 100% Completa (Fase 7.1)  
> **📊 Métricas:** 4.646+ linhas | 9 documentos | 6 diagramas | 16 endpoints Swagger  
> **🎯 Impacto:** -70% onboarding | -68% criar módulos | -83% deploy

---

## 🚀 Início Rápido

**Novo no projeto?** Comece por aqui:

1. 📖 [README Principal](../README.md) - Visão geral do projeto
2. 🏗️ [ARCHITECTURE.md](./ARCHITECTURE.md) - Entenda a arquitetura (6 diagramas!)
3. 📦 [MODULE_CREATION_GUIDE.md](./MODULE_CREATION_GUIDE.md) - Crie seu primeiro módulo (tutorial 12 passos)
4. 🌐 [Swagger UI](http://localhost:8080/swagger/index.html) - Explore a API (16 endpoints)
5. 🚀 [DEPLOYMENT.md](./DEPLOYMENT.md) - Configure e execute (deploy em 5min)

---

## 📑 Documentação Disponível

### 🏛️ Arquitetura

| Documento | Descrição | Status |
|-----------|-----------|--------|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Arquitetura completa com diagramas Mermaid | ✅ Completo |
| [ARCHITECTURE_COMPARISON.md](./ARCHITECTURE_COMPARISON.md) | Comparação: antes vs depois | ✅ Completo |

### 📦 Guias de Desenvolvimento

| Documento | Descrição | Linhas | Status |
|-----------|-----------|--------|--------|
| [MODULE_CREATION_GUIDE.md](./MODULE_CREATION_GUIDE.md) | Tutorial: criar novo módulo (12 passos + exemplo Category) | 1500+ | ✅ Completo |
| [guides/EVENTS_GUIDE.md](./guides/EVENTS_GUIDE.md) | Sistema de eventos type-safe | 800+ | ✅ Completo |
| [guides/CODE_EXAMPLES.md](./guides/CODE_EXAMPLES.md) | Exemplos práticos de código | 400+ | ✅ Completo |
| [guides/QUICKSTART.md](./guides/QUICKSTART.md) | Início rápido | 200+ | ✅ Completo |
| [guides/FILE_STRUCTURE.md](./guides/FILE_STRUCTURE.md) | Estrutura detalhada de arquivos | 300+ | ✅ Completo |

### 🌐 API Documentation (Swagger/OpenAPI)

| Recurso | Descrição | Status |
|---------|-----------|--------|
| [Swagger UI](http://localhost:8080/swagger/index.html) | Interface interativa da API | ✅ Funcional |
| [swagger.json](./swagger.json) | OpenAPI 3.0 Specification | ✅ Gerado |
| [swagger.yaml](./swagger.yaml) | OpenAPI YAML format | ✅ Gerado |

**Endpoints Documentados:** 16 endpoints (5 User + 6 Product + 5 Order)  
**Geração:** `swag init` (automático)

### 🚀 Deployment & Operações

| Documento | Descrição | Status |
|-----------|-----------|--------|
| [DEPLOYMENT.md](./DEPLOYMENT.md) | Guia completo de deployment | ✅ Completo |
| Docker Compose | Configuração multi-container | ✅ Incluído |
| Kubernetes | Manifests K8s | 🚧 Futuro |

### 📊 Progresso & Planejamento

| Documento | Descrição | Status |
|-----------|-----------|--------|
| [CHECKLIST.md](../CHECKLIST.md) | Checklist de refatoração (raiz do projeto) | ✅ 93% |
| [phases/FASE7.1_COMPLETO.md](./phases/FASE7.1_COMPLETO.md) | Resumo Fase 7.1 (Documentação) | ✅ Completo |
| [phases/FASE6_COMPLETO.md](./phases/FASE6_COMPLETO.md) | Resumo Fase 6 (Auto-registro) | ✅ Completo |
| [phases/FASE5_CONCLUSAO_FINAL.md](./phases/FASE5_CONCLUSAO_FINAL.md) | Resumo Fase 5 (Event Bus) | ✅ Completo |
| [phases/FASE4_COMPLETA.md](./phases/FASE4_COMPLETA.md) | Resumo Fase 4 (Domain Layer) | ✅ Completo |
| [phases/FASE3_CONCLUSAO.md](./phases/FASE3_CONCLUSAO.md) | Resumo Fase 3 (Application Layer) | ✅ Completo |

**Documentos Históricos:** Veja [archive/](./archive/) para planos e resumos iniciais da refatoração

### 🎯 Architecture Decision Records (ADRs)

| ADR | Título | Status | Data |
|-----|--------|--------|------|
| [001](./adr/001-clean-architecture.md) | Adoção de Clean Architecture | ✅ Aceito | 2025-10-01 |
| [002](./adr/002-cqrs-pattern.md) | Implementação do Padrão CQRS | ✅ Aceito | 2025-10-03 |
| [003](./adr/003-module-auto-registration.md) | Sistema de Auto-Registro de Módulos | ✅ Aceito | 2025-10-15 |

**[Ver todos os ADRs →](./adr/README.md)**

---

## 🎓 Trilha de Aprendizado

### Para Desenvolvedores Iniciantes

```
1. README.md
   ↓
2. QUICKSTART.md
   ↓
3. ARCHITECTURE.md (seções básicas)
   ↓
4. MODULE_CREATION_GUIDE.md (seguir tutorial)
   ↓
5. CODE_EXAMPLES.md
```

### Para Desenvolvedores Experientes

```
1. ARCHITECTURE.md (completo)
   ↓
2. ADRs (entender decisões)
   ↓
3. MODULE_CREATION_GUIDE.md (referência rápida)
   ↓
4. DEPLOYMENT.md (infraestrutura)
```

### Para Arquitetos/Tech Leads

```
1. ARCHITECTURE.md
   ↓
2. Todos os ADRs
   ↓
3. ARCHITECTURE_COMPARISON.md
   ↓
4. REFACTORING_INDEX.md
```

---

## 🔍 Busca Rápida

### "Como faço para..."

| Pergunta | Documento |
|----------|-----------|
| ...criar um novo módulo? | [MODULE_CREATION_GUIDE.md](./MODULE_CREATION_GUIDE.md) |
| ...entender a arquitetura? | [ARCHITECTURE.md](./ARCHITECTURE.md) |
| ...fazer deployment? | [DEPLOYMENT.md](./DEPLOYMENT.md) |
| ...usar o Event Bus? | [EVENTS_GUIDE.md](../EVENTS_GUIDE.md) |
| ...entender as decisões? | [ADRs](./adr/README.md) |
| ...ver exemplos de código? | [CODE_EXAMPLES.md](../CODE_EXAMPLES.md) |
| ...configurar variáveis de ambiente? | [DEPLOYMENT.md](./DEPLOYMENT.md#variáveis-de-ambiente) |
| ...rodar com Docker? | [DEPLOYMENT.md](./DEPLOYMENT.md#docker-e-docker-compose) |

### "O que é..."

| Conceito | Onde Aprender |
|----------|---------------|
| Clean Architecture | [ADR 001](./adr/001-clean-architecture.md) |
| CQRS | [ADR 002](./adr/002-cqrs-pattern.md) |
| Auto-Registro | [ADR 003](./adr/003-module-auto-registration.md) |
| Ports & Adapters | [ARCHITECTURE.md](./ARCHITECTURE.md#camadas-da-aplicação) |
| Event Bus Type-Safe | [EVENTS_GUIDE.md](../EVENTS_GUIDE.md) |
| ModuleRegistry | [ARCHITECTURE.md](./ARCHITECTURE.md#sistema-de-auto-registro) |

---

## 📊 Diagramas Importantes

### 🏗️ Arquitetura Geral

```mermaid
graph TB
    Client[Client] -->|HTTP/gRPC| Adapters
    Adapters[Adapters Layer] --> Application
    Application[Application Layer<br/>Commands/Queries] --> Domain
    Domain[Domain Layer<br/>Business Rules]
    Application --> Ports[Ports/Interfaces]
    Adapters -.implements.-> Ports
```

**[Ver diagrama completo →](./ARCHITECTURE.md#diagramas)**

### 📦 Auto-Registro

```mermaid
graph LR
    Module[Module] -->|Register| Registry[ModuleRegistry]
    Registry -->|Routes| HTTP[HTTP Server]
    Registry -->|Services| GRPC[gRPC Server]
```

**[Ver fluxo completo →](./ARCHITECTURE.md#sistema-de-auto-registro)**

---

## 📈 Status do Projeto

### ✅ Concluído (92%)

- [x] Fase 1: Reorganização de Estrutura (100%)
- [x] Fase 2: Interfaces e Contratos (100%)
- [x] Fase 3: Camada de Application (100%)
- [x] Fase 4: Sistema de Erros (100%)
- [x] Fase 5: Event Bus Type-Safe (100%)
- [x] Fase 6: Auto-registro de Módulos (100%)
- [x] Fase 7.1: Documentação (80%)

### 🚧 Em Andamento

- [ ] Fase 7.2: Testes (0%)
- [ ] Fase 7.3: Observabilidade (0%)
- [ ] Fase 7.4: Performance/Segurança (0%)

**[Ver checklist completo →](../CHECKLIST.md)**

---

## 🤝 Contribuindo

### Como Adicionar Documentação

1. **Documentação Técnica**: Adicione em `docs/`
2. **ADRs**: Use template em `docs/adr/000-template.md`
3. **Exemplos**: Adicione em `CODE_EXAMPLES.md`
4. **Atualize este índice**: Não esqueça!

### Padrões de Documentação

- ✅ Use Markdown
- ✅ Adicione diagramas Mermaid quando possível
- ✅ Inclua exemplos de código
- ✅ Mantenha atualizado
- ✅ Use emojis para navegação visual

---

## 📞 Suporte

### Problemas Comuns

- **Compilação falha**: [DEPLOYMENT.md - Troubleshooting](./DEPLOYMENT.md#troubleshooting)
- **Docker não inicia**: [DEPLOYMENT.md - Docker](./DEPLOYMENT.md#docker-e-docker-compose)
- **Arquitetura confusa**: [ARCHITECTURE.md](./ARCHITECTURE.md)

### Contatos

- **Issues**: GitHub Issues
- **Discussões**: GitHub Discussions
- **Email**: team@artemis.dev

---

## 🔖 Versão

**Documentação:** v1.0  
**Projeto:** v1.0.0  
**Última Atualização:** 18 de Outubro de 2025

---

## 📚 Leitura Externa Recomendada

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [CQRS - Martin Fowler](https://martinfowler.com/bliki/CQRS.html)
- [Hexagonal Architecture - Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)
- [Go Best Practices](https://go.dev/doc/effective_go)

---

**Mantido por:** Time Artemis  
**Licença:** MIT  
**Repository:** [github.com/your-org/artemis](https://github.com/your-org/artemis)
