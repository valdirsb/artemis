# Paginação no Módulo de Produtos

## Resumo
Implementação de paginação na listagem de produtos do módulo Product, seguindo os padrões da Clean Architecture e CQRS.

## Alterações Realizadas

### 1. Ports (`internal/modules/product/ports/ports.go`)
- **Adicionados campos de paginação em `ProductFilters`:**
  - `Page int`: Número da página (começa em 1)
  - `PageSize int`: Quantidade de itens por página

- **Criada struct `PaginatedResult`:**
  ```go
  type PaginatedResult struct {
      Items      []*domain.Product
      TotalItems int64
      Page       int
      PageSize   int
      TotalPages int
  }
  ```

- **Adicionado método em `ProductService`:**
  - `ListProductsPaginated(ctx context.Context, filters ProductFilters) (*PaginatedResult, error)`

- **Adicionado método em `ProductRepository`:**
  - `ListPaginated(ctx context.Context, filters ProductFilters) (*PaginatedResult, error)`

### 2. DTOs (`internal/modules/product/dto/`)

#### `responses.go`
- **Criada `PaginatedProductResponse`:**
  ```go
  type PaginatedProductResponse struct {
      Products   []ProductResponse `json:"products"`
      TotalItems int64             `json:"total_items"`
      Page       int               `json:"page"`
      PageSize   int               `json:"page_size"`
      TotalPages int               `json:"total_pages"`
  }
  ```

#### `mapper.go`
- **Adicionado `ToPaginatedProductResponse`:**
  - Converte `ports.PaginatedResult` para `PaginatedProductResponse`

### 3. Repository (`internal/modules/product/repository/product_repository.go`)
- **Implementado método `ListPaginated`:**
  - Conta o total de registros (`COUNT`)
  - Aplica filtros (categoria, preço, estoque)
  - Aplica paginação (`LIMIT` e `OFFSET`)
  - Calcula total de páginas
  - Retorna `PaginatedResult` com metadados

- **Valores padrão:**
  - `Page`: 1 (se não informado ou inválido)
  - `PageSize`: 10 (se não informado ou inválido)

### 4. Application Layer (`internal/modules/product/application/`)

#### `queries/list_products.go`
- **Atualizada `ListProductsQuery`:**
  - Adicionados campos `Page` e `PageSize`

- **Implementado método `HandlePaginated`:**
  - Processa query de listagem com paginação
  - Utiliza `productRepo.ListPaginated`
  - Retorna `PaginatedResult`

#### `services/product_application_service.go`
- **Implementado método `ListProductsPaginated`:**
  - Orquestra a query paginada
  - Converte filtros para query

### 5. HTTP Handler (`internal/modules/product/adapters/http/product_handler.go`)
- **Atualizado método `GetProducts`:**
  - Parse de query parameters:
    - `page`: número da página (padrão: 1)
    - `page_size`: itens por página (padrão: 10)
  - Utiliza `ListProductsPaginated` ao invés de `ListProducts`
  - Retorna `PaginatedProductResponse`

- **Atualizada documentação Swagger:**
  - Adicionados parâmetros `page` e `page_size`
  - Tipo de resposta alterado para `dto.PaginatedProductResponse`

## Como Usar

### Requisição HTTP

```bash
# Listar produtos com paginação padrão (página 1, 10 itens)
GET /api/v1/products

# Listar produtos - página 2, 20 itens por página
GET /api/v1/products?page=2&page_size=20

# Listar produtos com filtros e paginação
GET /api/v1/products?category_id=123&min_price=10.0&page=1&page_size=15
```

### Resposta JSON

```json
{
  "products": [
    {
      "id": "123",
      "name": "Produto 1",
      "description": "Descrição",
      "price": 99.90,
      "stock": 10,
      "category_id": "cat-123",
      "created_at": "2025-01-01T10:00:00Z",
      "updated_at": "2025-01-01T10:00:00Z"
    }
  ],
  "total_items": 45,
  "page": 1,
  "page_size": 10,
  "total_pages": 5
}
```

## Benefícios

1. **Performance**: Reduz a quantidade de dados trafegados na rede
2. **Escalabilidade**: Evita sobrecarga do banco de dados com queries muito grandes
3. **UX**: Melhor experiência do usuário com carregamento mais rápido
4. **Flexibilidade**: Cliente pode controlar a quantidade de dados por requisição
5. **Metadados**: Informações úteis para implementar UI de paginação

## Compatibilidade

O método `ListProducts` original foi mantido para compatibilidade com código existente. Novos endpoints devem utilizar a versão paginada.

## Testes

Para validar a implementação:

1. Teste sem parâmetros (deve usar valores padrão)
2. Teste com diferentes valores de `page` e `page_size`
3. Teste com filtros combinados com paginação
4. Teste valores inválidos (negativos, zero)
5. Teste navegação entre páginas

## Próximos Passos

- [ ] Adicionar testes unitários para `ListPaginated`
- [ ] Adicionar testes de integração para endpoint paginado
- [ ] Implementar paginação nos demais módulos (User, Order)
- [ ] Adicionar validação de limites máximos para `page_size`
- [ ] Considerar adicionar ordenação (sort) na paginação
