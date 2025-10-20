# Exemplos de Uso - Paginação de Produtos

## 1. Listagem Básica (Sem Paginação Explícita)

### Requisição
```bash
curl -X GET "http://localhost:8080/api/v1/products"
```

### Comportamento
- Página: 1 (padrão)
- Itens por página: 10 (padrão)

### Resposta Esperada
```json
{
  "products": [...],
  "total_items": 45,
  "page": 1,
  "page_size": 10,
  "total_pages": 5
}
```

---

## 2. Navegação entre Páginas

### Primeira Página
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&page_size=10"
```

### Segunda Página
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=2&page_size=10"
```

### Última Página
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=5&page_size=10"
```

---

## 3. Ajustando Tamanho da Página

### 5 itens por página
```bash
curl -X GET "http://localhost:8080/api/v1/products?page_size=5"
```

### 20 itens por página
```bash
curl -X GET "http://localhost:8080/api/v1/products?page_size=20"
```

### 50 itens por página
```bash
curl -X GET "http://localhost:8080/api/v1/products?page_size=50"
```

---

## 4. Paginação com Filtros

### Produtos de uma Categoria Específica
```bash
curl -X GET "http://localhost:8080/api/v1/products?category_id=cat-123&page=1&page_size=10"
```

### Produtos por Faixa de Preço
```bash
curl -X GET "http://localhost:8080/api/v1/products?min_price=50.0&max_price=200.0&page=1&page_size=15"
```

### Produtos em Estoque com Paginação
```bash
curl -X GET "http://localhost:8080/api/v1/products?in_stock=true&page=1&page_size=20"
```

### Combinação de Filtros
```bash
curl -X GET "http://localhost:8080/api/v1/products?category_id=cat-123&min_price=10.0&max_price=100.0&in_stock=true&page=2&page_size=15"
```

---

## 5. Casos Extremos

### Página Inexistente (retorna vazio)
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=999&page_size=10"
```

**Resposta:**
```json
{
  "products": [],
  "total_items": 45,
  "page": 999,
  "page_size": 10,
  "total_pages": 5
}
```

### Valores Inválidos (usa padrões)
```bash
# Page = 0 ou negativo → usa page = 1
curl -X GET "http://localhost:8080/api/v1/products?page=0&page_size=10"

# Page_size = 0 ou negativo → usa page_size = 10
curl -X GET "http://localhost:8080/api/v1/products?page=1&page_size=0"
```

---

## 6. Implementação Frontend

### JavaScript/TypeScript

```typescript
interface PaginatedResponse {
  products: Product[];
  total_items: number;
  page: number;
  page_size: number;
  total_pages: number;
}

async function fetchProducts(page: number = 1, pageSize: number = 10) {
  const response = await fetch(
    `http://localhost:8080/api/v1/products?page=${page}&page_size=${pageSize}`
  );
  const data: PaginatedResponse = await response.json();
  return data;
}

// Uso
const result = await fetchProducts(1, 20);
console.log(`Mostrando ${result.products.length} de ${result.total_items} produtos`);
console.log(`Página ${result.page} de ${result.total_pages}`);
```

### React Component Example

```tsx
import { useState, useEffect } from 'react';

function ProductList() {
  const [data, setData] = useState(null);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);

  useEffect(() => {
    fetch(`/api/v1/products?page=${page}&page_size=${pageSize}`)
      .then(res => res.json())
      .then(setData);
  }, [page, pageSize]);

  if (!data) return <div>Loading...</div>;

  return (
    <div>
      <div className="products">
        {data.products.map(product => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
      
      <div className="pagination">
        <button 
          onClick={() => setPage(p => Math.max(1, p - 1))}
          disabled={page === 1}
        >
          Anterior
        </button>
        
        <span>Página {page} de {data.total_pages}</span>
        
        <button 
          onClick={() => setPage(p => Math.min(data.total_pages, p + 1))}
          disabled={page === data.total_pages}
        >
          Próxima
        </button>
      </div>
      
      <div className="info">
        Mostrando {data.products.length} de {data.total_items} produtos
      </div>
    </div>
  );
}
```

---

## 7. Testes com cURL

### Script de Teste Completo

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1/products"

echo "=== Teste 1: Listagem padrão ==="
curl -s "$BASE_URL" | jq '.total_items, .page, .page_size'

echo -e "\n=== Teste 2: Segunda página ==="
curl -s "$BASE_URL?page=2&page_size=10" | jq '.page, .products | length'

echo -e "\n=== Teste 3: Com filtro de categoria ==="
curl -s "$BASE_URL?category_id=cat-123&page=1&page_size=5" | jq '.products | length'

echo -e "\n=== Teste 4: Faixa de preço ==="
curl -s "$BASE_URL?min_price=50.0&max_price=100.0" | jq '.total_items'

echo -e "\n=== Teste 5: Produtos em estoque ==="
curl -s "$BASE_URL?in_stock=true&page_size=20" | jq '.products | length, .total_items'
```

---

## 8. Métricas e Performance

### Comparação de Performance

**Sem Paginação (retornando 1000 produtos):**
- Tempo de resposta: ~2.5s
- Tamanho da resposta: ~500KB
- Uso de memória: Alto

**Com Paginação (20 produtos por página):**
- Tempo de resposta: ~150ms
- Tamanho da resposta: ~10KB
- Uso de memória: Baixo

### Recomendações

- **Mobile**: 10-15 itens por página
- **Desktop**: 20-30 itens por página
- **API/Backend**: 50-100 itens por página
- **Exports/Reports**: 500-1000 itens por página

---

## 9. Tratamento de Erros

### Servidor Retorna 500
```json
{
  "error": "failed to list products: database connection lost"
}
```

### Sem Resultados
```json
{
  "products": [],
  "total_items": 0,
  "page": 1,
  "page_size": 10,
  "total_pages": 0
}
```

---

## 10. Boas Práticas

✅ **Fazer:**
- Sempre validar `page` e `page_size` no frontend
- Implementar indicador de loading durante fetch
- Mostrar informações de paginação ao usuário
- Cachear resultados quando possível
- Implementar scroll infinito quando apropriado

❌ **Evitar:**
- Page_size muito grande (>1000)
- Não tratar casos de página vazia
- Não mostrar feedback de loading
- Requisições simultâneas desnecessárias
