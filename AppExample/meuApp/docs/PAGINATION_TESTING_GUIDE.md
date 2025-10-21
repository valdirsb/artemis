# 🧪 Testando Paginação - Exemplos Práticos

## 📋 Cenários de Teste

Este guia fornece exemplos práticos para testar a paginação via HTTP REST e gRPC.

## 🌐 Testes HTTP REST

### Setup do Ambiente de Teste

```bash
# Iniciar o servidor
cd /home/junior/projetos/artemis/meuApp
go run main.go

# Em outro terminal, preparar comandos curl
export BASE_URL="http://localhost:8080/api/v1"
```

### 1. Testar Products

#### ✅ Teste 1: Paginação básica
```bash
# Primeira página (default)
curl -s "$BASE_URL/products" | jq '.'

# Primeira página explícita
curl -s "$BASE_URL/products?page=1&page_size=10" | jq '.'

# Resposta esperada:
{
  "items": [...],
  "pagination": {
    "total_items": 100,
    "page": 1,
    "page_size": 10,
    "total_pages": 10
  }
}
```

#### ✅ Teste 2: Diferentes tamanhos de página
```bash
# 5 items por página
curl -s "$BASE_URL/products?page=1&page_size=5" | jq '.pagination'

# 20 items por página
curl -s "$BASE_URL/products?page=1&page_size=20" | jq '.pagination'

# Máximo permitido (100)
curl -s "$BASE_URL/products?page=1&page_size=100" | jq '.pagination'

# Acima do máximo (deve limitar a 100)
curl -s "$BASE_URL/products?page=1&page_size=200" | jq '.pagination.page_size'
# Esperado: 100
```

#### ✅ Teste 3: Navegação entre páginas
```bash
# Página 1
curl -s "$BASE_URL/products?page=1&page_size=10" | jq '.pagination'

# Página 2
curl -s "$BASE_URL/products?page=2&page_size=10" | jq '.pagination'

# Última página
curl -s "$BASE_URL/products?page=10&page_size=10" | jq '.pagination'

# Página além do total (deve retornar vazio)
curl -s "$BASE_URL/products?page=999&page_size=10" | jq '.items | length'
# Esperado: 0
```

#### ✅ Teste 4: Valores inválidos
```bash
# Page negativo (deve usar default 1)
curl -s "$BASE_URL/products?page=-1&page_size=10" | jq '.pagination.page'
# Esperado: 1

# Page zero (deve usar default 1)
curl -s "$BASE_URL/products?page=0&page_size=10" | jq '.pagination.page'
# Esperado: 1

# Page_size negativo (deve usar default 10)
curl -s "$BASE_URL/products?page=1&page_size=-5" | jq '.pagination.page_size'
# Esperado: 10

# Parâmetros não numéricos (deve usar defaults)
curl -s "$BASE_URL/products?page=abc&page_size=xyz" | jq '.pagination'
# Esperado: page=1, page_size=10
```

### 2. Testar Users

```bash
# Listar usuários - página 1
curl -s "$BASE_URL/users?page=1&page_size=10" | jq '.'

# Verificar total de usuários
curl -s "$BASE_URL/users?page=1&page_size=1" | jq '.pagination.total_items'

# Navegar páginas
for page in {1..3}; do
  echo "=== Página $page ==="
  curl -s "$BASE_URL/users?page=$page&page_size=5" | jq '.items[].id'
done
```

### 3. Testar Orders

#### Listar todos os pedidos
```bash
# Primeira página
curl -s "$BASE_URL/orders?page=1&page_size=10" | jq '.'

# Contar total
curl -s "$BASE_URL/orders?page=1&page_size=1" | jq '.pagination.total_items'
```

#### Listar pedidos de um usuário específico
```bash
# Substitua pelo user_id real
USER_ID="0a326204-a40d-4f28-a7f9-8c87bf1f87e6"

# Primeira página dos pedidos do usuário
curl -s "$BASE_URL/orders/user/$USER_ID?page=1&page_size=10" | jq '.'

# Segunda página
curl -s "$BASE_URL/orders/user/$USER_ID?page=2&page_size=5" | jq '.'

# Verificar total de pedidos do usuário
curl -s "$BASE_URL/orders/user/$USER_ID?page=1&page_size=1" | \
  jq '.pagination.total_items'
```

## 🔌 Testes gRPC

### Setup do Ambiente gRPC

```bash
# Instalar grpcurl se ainda não tiver
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Verificar se o servidor gRPC está rodando
grpcurl -plaintext localhost:50051 list
```

### 1. Testar ProductService

```bash
# Teste 1: Paginação básica
grpcurl -plaintext -d '{"page": 1, "page_size": 10}' \
  localhost:50051 api.ProductService/ListProducts | jq '.'

# Teste 2: Primeira página implícita (defaults)
grpcurl -plaintext -d '{}' \
  localhost:50051 api.ProductService/ListProducts | jq '.totalItems'

# Teste 3: Página grande
grpcurl -plaintext -d '{"page": 1, "page_size": 50}' \
  localhost:50051 api.ProductService/ListProducts | jq '.pageSize'

# Teste 4: Segunda página
grpcurl -plaintext -d '{"page": 2, "page_size": 10}' \
  localhost:50051 api.ProductService/ListProducts | jq '.page'
```

### 2. Testar UserService

```bash
# Teste 1: Listar usuários
grpcurl -plaintext -d '{"page": 1, "page_size": 10}' \
  localhost:50051 api.UserService/ListUsers | jq '.'

# Teste 2: Contar usuários
grpcurl -plaintext -d '{"page": 1, "page_size": 1}' \
  localhost:50051 api.UserService/ListUsers | jq '.total'

# Teste 3: Navegação
for page in {1..3}; do
  echo "=== Página $page ==="
  grpcurl -plaintext -d "{\"page\": $page, \"page_size\": 5}" \
    localhost:50051 api.UserService/ListUsers | jq '.users[].id'
done
```

### 3. Testar OrderService

#### ListOrders
```bash
# Teste 1: Listar todos os pedidos
grpcurl -plaintext -d '{"page": 1, "page_size": 10}' \
  localhost:50051 api.OrderService/ListOrders | jq '.'

# Teste 2: Contar total
grpcurl -plaintext -d '{"page": 1, "page_size": 1}' \
  localhost:50051 api.OrderService/ListOrders | jq '.total'

# Teste 3: Página específica
grpcurl -plaintext -d '{"page": 2, "page_size": 5}' \
  localhost:50051 api.OrderService/ListOrders | jq '.orders[].id'
```

#### GetOrdersByUser
```bash
# Substitua pelo user_id real
USER_ID="0a326204-a40d-4f28-a7f9-8c87bf1f87e6"

# Teste 1: Pedidos do usuário - primeira página
grpcurl -plaintext -d "{\"user_id\": \"$USER_ID\", \"page\": 1, \"page_size\": 10}" \
  localhost:50051 api.OrderService/GetOrdersByUser | jq '.'

# Teste 2: Sem especificar paginação (usa defaults)
grpcurl -plaintext -d "{\"user_id\": \"$USER_ID\"}" \
  localhost:50051 api.OrderService/GetOrdersByUser | jq '.total'

# Teste 3: Segunda página
grpcurl -plaintext -d "{\"user_id\": \"$USER_ID\", \"page\": 2, \"page_size\": 5}" \
  localhost:50051 api.OrderService/GetOrdersByUser | jq '.page'
```

## 🎯 Script de Teste Automatizado

### Bash Script para HTTP

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Testando Paginação HTTP REST"
echo "================================"

# Teste Products
echo -e "\n📦 Products:"
echo -n "  Total items: "
curl -s "$BASE_URL/products?page=1&page_size=1" | jq -r '.pagination.total_items'
echo -n "  Primeira página: "
curl -s "$BASE_URL/products?page=1&page_size=10" | jq -r '.pagination.page'
echo -n "  Items na página: "
curl -s "$BASE_URL/products?page=1&page_size=10" | jq -r '.items | length'

# Teste Users
echo -e "\n👥 Users:"
echo -n "  Total items: "
curl -s "$BASE_URL/users?page=1&page_size=1" | jq -r '.pagination.total_items'
echo -n "  Primeira página: "
curl -s "$BASE_URL/users?page=1&page_size=10" | jq -r '.pagination.page'
echo -n "  Items na página: "
curl -s "$BASE_URL/users?page=1&page_size=10" | jq -r '.items | length'

# Teste Orders
echo -e "\n📋 Orders:"
echo -n "  Total items: "
curl -s "$BASE_URL/orders?page=1&page_size=1" | jq -r '.pagination.total_items'
echo -n "  Primeira página: "
curl -s "$BASE_URL/orders?page=1&page_size=10" | jq -r '.pagination.page'
echo -n "  Items na página: "
curl -s "$BASE_URL/orders?page=1&page_size=10" | jq -r '.items | length'

echo -e "\n✅ Testes concluídos!"
```

### Go Test para gRPC

```go
package main

import (
    "context"
    "testing"
    "time"

    pb "meuApp/pkg/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func TestProductsPagination(t *testing.T) {
    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        t.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewProductServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    // Test 1: Default pagination
    resp, err := client.ListProducts(ctx, &pb.ListProductsRequest{})
    if err != nil {
        t.Fatalf("ListProducts failed: %v", err)
    }
    if resp.Page != 1 {
        t.Errorf("Expected page 1, got %d", resp.Page)
    }
    if resp.PageSize != 10 {
        t.Errorf("Expected pageSize 10, got %d", resp.PageSize)
    }

    // Test 2: Custom pagination
    resp, err = client.ListProducts(ctx, &pb.ListProductsRequest{
        Page:     2,
        PageSize: 5,
    })
    if err != nil {
        t.Fatalf("ListProducts failed: %v", err)
    }
    if resp.Page != 2 {
        t.Errorf("Expected page 2, got %d", resp.Page)
    }
    if resp.PageSize != 5 {
        t.Errorf("Expected pageSize 5, got %d", resp.PageSize)
    }

    // Test 3: Max page size limit
    resp, err = client.ListProducts(ctx, &pb.ListProductsRequest{
        Page:     1,
        PageSize: 200, // Should be limited to 100
    })
    if err != nil {
        t.Fatalf("ListProducts failed: %v", err)
    }
    if resp.PageSize > 100 {
        t.Errorf("PageSize should be limited to 100, got %d", resp.PageSize)
    }
}

func TestUsersPagination(t *testing.T) {
    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        t.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewUserServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    resp, err := client.ListUsers(ctx, &pb.ListUsersRequest{
        Page:     1,
        PageSize: 10,
    })
    if err != nil {
        t.Fatalf("ListUsers failed: %v", err)
    }

    if resp.Total < 0 {
        t.Error("Total should be >= 0")
    }
    if len(resp.Users) > int(resp.PageSize) {
        t.Errorf("Returned more users than page_size")
    }
}

func TestOrdersPagination(t *testing.T) {
    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        t.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewOrderServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    // Test ListOrders
    resp, err := client.ListOrders(ctx, &pb.ListOrdersRequest{
        Page:     1,
        PageSize: 10,
    })
    if err != nil {
        t.Fatalf("ListOrders failed: %v", err)
    }
    if resp.Page != 1 {
        t.Errorf("Expected page 1, got %d", resp.Page)
    }

    // Test GetOrdersByUser
    userResp, err := client.GetOrdersByUser(ctx, &pb.GetOrdersByUserRequest{
        UserId:   "test-user-id",
        Page:     1,
        PageSize: 5,
    })
    if err == nil {
        if userResp.PageSize != 5 {
            t.Errorf("Expected pageSize 5, got %d", userResp.PageSize)
        }
    }
}
```

## 📊 Performance Testing

### Apache Bench (HTTP)

```bash
# Teste de carga na listagem de produtos
ab -n 1000 -c 10 "http://localhost:8080/api/v1/products?page=1&page_size=10"

# Teste com diferentes tamanhos de página
ab -n 500 -c 5 "http://localhost:8080/api/v1/products?page=1&page_size=50"
```

### ghz (gRPC)

```bash
# Instalar ghz
go install github.com/bojand/ghz/cmd/ghz@latest

# Teste de carga gRPC
ghz --insecure \
  --proto proto/product.proto \
  --call api.ProductService/ListProducts \
  -d '{"page": 1, "page_size": 10}' \
  -n 1000 \
  -c 10 \
  localhost:50051
```

## ✅ Checklist de Testes

### HTTP REST
- [ ] ✅ Paginação básica (defaults)
- [ ] ✅ Página específica
- [ ] ✅ Page_size customizado
- [ ] ✅ Limite máximo de page_size
- [ ] ✅ Valores inválidos (negativos, zero)
- [ ] ✅ Parâmetros não numéricos
- [ ] ✅ Navegação entre páginas
- [ ] ✅ Última página
- [ ] ✅ Página além do limite
- [ ] ✅ Filtros com paginação (Products)
- [ ] ✅ Pedidos por usuário com paginação

### gRPC
- [ ] ✅ Paginação básica
- [ ] ✅ Defaults quando não especificado
- [ ] ✅ Page_size customizado
- [ ] ✅ Limite máximo
- [ ] ✅ ListProducts
- [ ] ✅ ListUsers
- [ ] ✅ ListOrders
- [ ] ✅ GetOrdersByUser

### Performance
- [ ] ⏳ Teste de carga HTTP (1000 requests)
- [ ] ⏳ Teste de carga gRPC (1000 requests)
- [ ] ⏳ Tempo de resposta < 100ms (página pequena)
- [ ] ⏳ Tempo de resposta < 500ms (página grande)

## 🎉 Resultado Esperado

Todos os testes devem passar com:
- ✅ Respostas corretas
- ✅ Paginação funcionando
- ✅ Validações aplicadas
- ✅ Performance adequada
- ✅ Erros tratados corretamente

## 📝 Notas

- Os testes assumem que o servidor está rodando em `localhost:8080` (HTTP) e `localhost:50051` (gRPC)
- Ajuste os IDs de usuário conforme necessário
- Use `jq` para formatar JSON
- Salve os scripts em arquivos executáveis para reutilização
