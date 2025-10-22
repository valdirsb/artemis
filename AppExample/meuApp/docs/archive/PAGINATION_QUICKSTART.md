# 🚀 Guia Rápido - API com Paginação

## Endpoints Disponíveis

### 📦 Products (Produtos)

```bash
# Listar produtos (página 1, 10 itens)
curl "http://localhost:8080/api/v1/products"

# Paginação customizada
curl "http://localhost:8080/api/v1/products?page=2&page_size=20"

# Com filtros + paginação
curl "http://localhost:8080/api/v1/products?category_id=cat-123&min_price=10.0&max_price=100.0&page=1&page_size=15"
```

**Resposta:**
```json
{
  "products": [
    {
      "id": "prod-123",
      "name": "Notebook",
      "description": "High performance laptop",
      "price": 2999.90,
      "stock": 10,
      "category_id": "cat-123",
      "created_at": "2025-10-20T10:00:00Z",
      "updated_at": "2025-10-20T10:00:00Z"
    }
  ],
  "total_items": 150,
  "page": 1,
  "page_size": 10,
  "total_pages": 15
}
```

---

### 👥 Users (Usuários)

```bash
# Listar usuários (página 1, 10 itens)
curl "http://localhost:8080/api/v1/users"

# Paginação customizada
curl "http://localhost:8080/api/v1/users?page=3&page_size=25"
```

**Resposta:**
```json
{
  "users": [
    {
      "id": "user-456",
      "username": "john_doe",
      "email": "john@example.com",
      "created_at": "2025-10-15T08:30:00Z",
      "updated_at": "2025-10-20T10:00:00Z"
    }
  ],
  "total_items": 275,
  "page": 3,
  "page_size": 25,
  "total_pages": 11
}
```

---

### 🛒 Orders (Pedidos)

```bash
# Listar pedidos (página 1, 10 itens)
curl "http://localhost:8080/api/v1/orders"

# Paginação customizada
curl "http://localhost:8080/api/v1/orders?page=2&page_size=15"
```

**Resposta:**
```json
{
  "orders": [
    {
      "id": "order-789",
      "user_id": "user-456",
      "items": [
        {
          "product_id": "prod-123",
          "quantity": 2,
          "price": 2999.90
        }
      ],
      "status": "confirmed",
      "total": 5999.80,
      "created_at": "2025-10-20T09:45:00Z",
      "updated_at": "2025-10-20T10:15:00Z"
    }
  ],
  "total_items": 85,
  "page": 2,
  "page_size": 15,
  "total_pages": 6
}
```

---

## ⚙️ Parâmetros

### Query Parameters

| Parâmetro | Tipo | Padrão | Limite | Descrição |
|-----------|------|--------|--------|-----------|
| `page` | int | 1 | - | Número da página |
| `page_size` | int | 10 | 100* | Itens por página |

*Limite de 100 aplicado apenas para User e Order

### Valores Especiais

- `page=0` ou negativo → usa `page=1`
- `page_size=0` ou negativo → usa `page_size=10`
- `page_size > 100` (User/Order) → usa `page_size=100`

---

## 📊 Estrutura de Resposta

Todos os endpoints paginados retornam:

```json
{
  "items": [...],        // Array de objetos (products, users ou orders)
  "total_items": 150,    // Total de registros no banco
  "page": 1,             // Página atual
  "page_size": 10,       // Itens por página
  "total_pages": 15      // Total de páginas calculado
}
```

---

## 🎯 Casos de Uso

### Navegação de Páginas

```javascript
// Primeira página
fetch('/api/v1/products?page=1&page_size=10')

// Próxima página
fetch('/api/v1/products?page=2&page_size=10')

// Última página (baseado em total_pages da resposta)
fetch('/api/v1/products?page=15&page_size=10')
```

### Ajustar Tamanho da Página

```javascript
// Mobile: menos itens
fetch('/api/v1/products?page_size=5')

// Desktop: mais itens
fetch('/api/v1/products?page_size=30')

// API/Backend: máximo permitido
fetch('/api/v1/products?page_size=100')
```

### Produtos com Filtros

```bash
# Categoria específica
curl "/api/v1/products?category_id=electronics&page=1&page_size=20"

# Faixa de preço
curl "/api/v1/products?min_price=100&max_price=500&page=1"

# Apenas em estoque
curl "/api/v1/products?in_stock=true&page=1&page_size=15"

# Combinação de filtros
curl "/api/v1/products?category_id=electronics&min_price=100&in_stock=true&page=2&page_size=20"
```

---

## 💻 Exemplos de Código

### JavaScript/TypeScript

```typescript
async function fetchPaginatedData(endpoint: string, page = 1, pageSize = 10) {
  const url = `${endpoint}?page=${page}&page_size=${pageSize}`;
  const response = await fetch(url);
  const data = await response.json();
  
  console.log(`Página ${data.page} de ${data.total_pages}`);
  console.log(`Mostrando ${data.items.length} de ${data.total_items} itens`);
  
  return data;
}

// Uso
const products = await fetchPaginatedData('/api/v1/products', 1, 20);
const users = await fetchPaginatedData('/api/v1/users', 2, 15);
const orders = await fetchPaginatedData('/api/v1/orders', 1, 10);
```

### Python

```python
import requests

def fetch_paginated(endpoint, page=1, page_size=10):
    url = f"{endpoint}?page={page}&page_size={page_size}"
    response = requests.get(url)
    data = response.json()
    
    print(f"Página {data['page']} de {data['total_pages']}")
    print(f"Mostrando {len(data['items'])} de {data['total_items']} itens")
    
    return data

# Uso
products = fetch_paginated('http://localhost:8080/api/v1/products', page=1, page_size=20)
users = fetch_paginated('http://localhost:8080/api/v1/users', page=2, page_size=15)
orders = fetch_paginated('http://localhost:8080/api/v1/orders', page=1, page_size=10)
```

### React Component

```tsx
import { useState, useEffect } from 'react';

function PaginatedList({ endpoint, itemComponent: ItemComponent }) {
  const [data, setData] = useState(null);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    fetch(`${endpoint}?page=${page}&page_size=${pageSize}`)
      .then(res => res.json())
      .then(data => {
        setData(data);
        setLoading(false);
      });
  }, [endpoint, page, pageSize]);

  if (loading) return <div>Carregando...</div>;
  if (!data) return null;

  const items = data.products || data.users || data.orders;

  return (
    <div>
      <div className="items">
        {items.map(item => (
          <ItemComponent key={item.id} item={item} />
        ))}
      </div>
      
      <div className="pagination">
        <button 
          onClick={() => setPage(p => Math.max(1, p - 1))}
          disabled={page === 1}
        >
          ← Anterior
        </button>
        
        <span>
          Página {data.page} de {data.total_pages}
        </span>
        
        <button 
          onClick={() => setPage(p => Math.min(data.total_pages, p + 1))}
          disabled={page === data.total_pages}
        >
          Próxima →
        </button>
      </div>
      
      <div className="info">
        Mostrando {items.length} de {data.total_items} itens
      </div>
    </div>
  );
}

// Uso
<PaginatedList endpoint="/api/v1/products" itemComponent={ProductCard} />
<PaginatedList endpoint="/api/v1/users" itemComponent={UserCard} />
<PaginatedList endpoint="/api/v1/orders" itemComponent={OrderCard} />
```

---

## 🧪 Testando com cURL

### Script de Teste Completo

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "=== Testando Products ==="
echo "Página 1:"
curl -s "$BASE_URL/products?page=1&page_size=5" | jq '.total_items, .page, .page_size'

echo -e "\nPágina 2:"
curl -s "$BASE_URL/products?page=2&page_size=5" | jq '.page, .products | length'

echo -e "\n=== Testando Users ==="
curl -s "$BASE_URL/users?page=1&page_size=10" | jq '.total_items, .total_pages'

echo -e "\n=== Testando Orders ==="
curl -s "$BASE_URL/orders?page=1&page_size=10" | jq '.total_items, .orders | length'
```

---

## 📱 Recomendações por Plataforma

| Plataforma | page_size recomendado |
|------------|----------------------|
| **Mobile** | 5-10 |
| **Tablet** | 15-20 |
| **Desktop** | 20-30 |
| **API/Backend** | 50-100 |
| **Exports/Reports** | 100+ |

---

## ⚠️ Tratamento de Erros

### Página Inexistente

**Request:**
```bash
curl "/api/v1/products?page=999&page_size=10"
```

**Response:**
```json
{
  "products": [],
  "total_items": 150,
  "page": 999,
  "page_size": 10,
  "total_pages": 15
}
```

### Erro do Servidor

**Response:**
```json
{
  "error": "failed to list products: database connection lost"
}
```

---

## ✅ Checklist de Implementação Frontend

- [ ] Implementar controles de navegação (anterior/próximo)
- [ ] Mostrar informação de página atual
- [ ] Exibir total de itens e páginas
- [ ] Adicionar indicador de loading
- [ ] Tratar casos de lista vazia
- [ ] Implementar seletor de page_size
- [ ] Adicionar deep linking (URL com parâmetros)
- [ ] Cachear resultados quando apropriado
- [ ] Implementar scroll infinito (opcional)
- [ ] Adicionar keyboard shortcuts (opcional)

---

## 🎉 Pronto!

Agora você pode usar paginação em todos os módulos do sistema de forma consistente e eficiente!

**Documentação Adicional:**
- `PAGINATION_SUMMARY.md` - Resumo da implementação do módulo Product
- `PAGINATION_USER_ORDER_SUMMARY.md` - Resumo dos módulos User e Order
- `docs/PAGINATION_PRODUCT_MODULE.md` - Documentação técnica detalhada
- `docs/examples/PAGINATION_EXAMPLES.md` - Mais exemplos de uso
