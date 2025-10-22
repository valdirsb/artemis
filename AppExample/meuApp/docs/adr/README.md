# Architecture Decision Records (ADR)

Este diretório contém os registros de decisões arquiteturais importantes do projeto Artemis.

## O que é um ADR?

Um Architecture Decision Record (ADR) é um documento que captura uma decisão arquitetural importante juntamente com seu contexto e consequências.

## Formato

Cada ADR segue este formato:

- **Status**: Proposto | Aceito | Rejeitado | Substituído | Obsoleto
- **Data**: Data da decisão
- **Contexto**: O problema ou situação que levou à decisão
- **Decisão**: O que foi decidido
- **Consequências**: Impactos positivos e negativos da decisão
- **Alternativas Consideradas**: Outras opções que foram avaliadas

## Índice de ADRs

| # | Título | Status | Data |
|---|--------|--------|------|
| [001](./001-clean-architecture.md) | Adoção de Clean Architecture | Aceito | 2025-10-01 |
| [002](./002-cqrs-pattern.md) | Implementação do Padrão CQRS | Aceito | 2025-10-03 |
| [003](./003-module-auto-registration.md) | Sistema de Auto-Registro de Módulos | Aceito | 2025-10-15 |
| [004](./004-event-bus-type-safe.md) | Event Bus Type-Safe | Aceito | 2025-10-12 |
| [005](./005-hexagonal-architecture.md) | Hexagonal Architecture (Ports & Adapters) | Aceito | 2025-10-02 |
| [006](./006-multi-protocol-support.md) | Suporte Multi-Protocolo (HTTP + gRPC) | Aceito | 2025-10-05 |

## Como Criar um Novo ADR

1. Copie o template `000-template.md`
2. Renomeie para `XXX-titulo-descritivo.md`
3. Preencha todas as seções
4. Adicione ao índice acima
5. Commit e PR

## Referências

- [Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)
- [ADR GitHub Organization](https://adr.github.io/)
