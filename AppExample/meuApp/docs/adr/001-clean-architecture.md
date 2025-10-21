# ADR 001: Adoção de Clean Architecture

**Status:** ✅ Aceito  
**Data:** 01 de Outubro de 2025  
**Autores:** Time Artemis  

---

## Contexto

Precisávamos de uma arquitetura que garantisse:
- **Testabilidade**: Fácil criar testes unitários e de integração
- **Manutenibilidade**: Código fácil de entender e modificar
- **Independência de Frameworks**: Não ficar preso a tecnologias específicas
- **Escalabilidade**: Suportar crescimento do projeto

### Problema

Aplicações monolíticas tradicionais sofrem de:
- Acoplamento alto entre camadas
- Dificuldade para testar
- Dependências diretas de frameworks e bibliotecas
- Lógica de negócio misturada com infraestrutura

---

## Decisão

Adotar **Clean Architecture** (Arquitetura Limpa) conforme proposta por Robert C. Martin (Uncle Bob).

### Princípios Aplicados

1. **Dependency Rule** (Regra de Dependência)
   - Dependências apontam sempre para dentro (inward)
   - Camadas externas dependem de internas, nunca o contrário

2. **Separation of Concerns**
   - Cada camada tem responsabilidade específica
   - Domain, Application, Adapters, Framework

3. **Domain-Centric**
   - Domínio no centro, sem dependências externas
   - Regras de negócio puras e isoladas

### Estrutura de Camadas

```
┌─────────────────────────────────────┐
│  Framework & Drivers (External)    │  ← Gin, GORM, gRPC
├─────────────────────────────────────┤
│  Interface Adapters (HTTP/gRPC)    │  ← Handlers, Presenters
├─────────────────────────────────────┤
│  Application Business Rules         │  ← Use Cases, Services
├─────────────────────────────────────┤
│  Enterprise Business Rules (Domain) │  ← Entities, Value Objects
└─────────────────────────────────────┘
```

---

## Alternativas Consideradas

### 1. MVC Tradicional

**Prós:**
- Simples e bem conhecido
- Rápido para começar

**Contras:**
- ❌ Alto acoplamento
- ❌ Difícil de testar
- ❌ Lógica de negócio espalhada

**Decisão:** Rejeitado - não escala bem

### 2. Layered Architecture (3-tier)

**Prós:**
- Separação de responsabilidades
- Organização clara

**Contras:**
- ❌ Dependências entre camadas não controladas
- ❌ Domínio ainda depende de infraestrutura

**Decisão:** Rejeitado - não garante independência

### 3. Microservices (puro)

**Prós:**
- Escalabilidade horizontal
- Deployment independente

**Contras:**
- ❌ Complexidade operacional
- ❌ Overhead de rede
- ❌ Overkill para projeto atual

**Decisão:** Rejeitado - complexidade desnecessária neste momento

---

## Consequências

### ✅ Positivas

1. **Testabilidade Máxima**
   - Domain pode ser testado sem banco de dados
   - Mocks fáceis via interfaces
   - Testes rápidos

2. **Independência de Framework**
   - Fácil trocar Gin por outro router
   - Fácil trocar GORM por outro ORM
   - Domínio não conhece frameworks

3. **Manutenibilidade**
   - Código bem organizado
   - Fácil localizar bugs
   - Mudanças localizadas

4. **Escalabilidade**
   - Arquitetura preparada para crescimento
   - Fácil adicionar novos módulos
   - Base para microservices futuros

5. **Qualidade de Código**
   - SOLID principles aplicados
   - Baixo acoplamento
   - Alta coesão

### ⚠️ Negativas

1. **Curva de Aprendizado**
   - Time precisa entender os conceitos
   - Mais complexo que MVC simples

2. **Mais Código Inicial**
   - Mais arquivos e interfaces
   - Setup inicial mais longo

3. **Over-engineering para Features Simples**
   - CRUD simples pode parecer verboso
   - Trade-off: consistência vs simplicidade

### 🔧 Mitigações

- Documentação completa (ARCHITECTURE.md)
- Exemplos práticos
- Code review rigoroso
- Templates para novos módulos

---

## Validação

### Métricas de Sucesso

- [x] Zero dependências do domain em frameworks
- [x] Testes de domain sem dependências externas
- [x] Módulos isolados e independentes
- [x] Fácil trocar adaptadores (comprovado com refactoring)

### Testes Práticos

- Trocamos logger sem tocar no domain ✅
- Trocamos database adapter sem problemas ✅
- Adicionamos gRPC sem modificar use cases ✅

---

## Referências

- [The Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Clean Architecture Book](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164)
- [Hexagonal Architecture - Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)

---

**Status Final:** ✅ **ACEITO E BEM-SUCEDIDO**

Clean Architecture provou ser a escolha correta, entregando todos os benefícios esperados.
