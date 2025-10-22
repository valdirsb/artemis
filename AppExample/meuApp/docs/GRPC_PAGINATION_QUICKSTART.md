# 🚀 Guia Rápido: Usando Paginação via gRPC

## 📋 Pré-requisitos

```bash
# Instalar grpcurl (ferramenta CLI para testar gRPC)
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Ou usando homebrew (macOS)
brew install grpcurl
```

## 🎯 Endpoints gRPC Disponíveis

| Serviço | Método | Descrição |
|---------|--------|-----------|
| ProductService | ListProducts | Lista produtos com paginação |
| UserService | ListUsers | Lista usuários com paginação |
| OrderService | ListOrders | Lista todos os pedidos com paginação |
| OrderService | GetOrdersByUser | Lista pedidos de um usuário com paginação |

## 🔧 Exemplos de Uso

### 1. Listar Produtos

#### Comando básico (usando defaults: page=1, page_size=10)
```bash
grpcurl -plaintext localhost:50051 api.ProductService/ListProducts
```

#### Com paginação customizada
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 20
}' localhost:50051 api.ProductService/ListProducts
```

#### Segunda página
```bash
grpcurl -plaintext -d '{
  "page": 2,
  "page_size": 10
}' localhost:50051 api.ProductService/ListProducts
```

#### Resposta esperada
```json
{
  "products": [
    {
      "id": "product-id-1",
      "name": "Product 1",
      "description": "Description",
      "price": 99.99,
      "quantity": 10,
      "createdAt": "2025-10-20T10:00:00Z",
      "updatedAt": "2025-10-20T10:00:00Z"
    }
  ],
  "totalItems": 100,
  "page": 1,
  "pageSize": 20,
  "totalPages": 5,
  "message": "Products retrieved successfully"
}
```

### 2. Listar Usuários

#### Comando básico
```bash
grpcurl -plaintext localhost:50051 api.UserService/ListUsers
```

#### Com paginação
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 15
}' localhost:50051 api.UserService/ListUsers
```

#### Resposta esperada
```json
{
  "users": [
    {
      "id": "user-id-1",
      "name": "John Doe",
      "email": "john@example.com",
      "createdAt": "2025-10-20T10:00:00Z",
      "updatedAt": "2025-10-20T10:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "pageSize": 15,
  "message": "Users retrieved successfully"
}
```

### 3. Listar Todos os Pedidos

#### Comando básico
```bash
grpcurl -plaintext localhost:50051 api.OrderService/ListOrders
```

#### Com paginação
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50051 api.OrderService/ListOrders
```

#### Resposta esperada
```json
{
  "orders": [
    {
      "id": "order-id-1",
      "userId": "user-id-1",
      "items": [
        {
          "productId": "product-id-1",
          "quantity": 2,
          "price": 99.99
        }
      ],
      "totalPrice": 199.98,
      "status": "pending",
      "createdAt": "2025-10-20T10:00:00Z",
      "updatedAt": "2025-10-20T10:00:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "pageSize": 10,
  "message": "Orders retrieved successfully"
}
```

### 4. Listar Pedidos de um Usuário Específico

#### Comando com user_id
```bash
grpcurl -plaintext -d '{
  "user_id": "0a326204-a40d-4f28-a7f9-8c87bf1f87e6",
  "page": 1,
  "page_size": 10
}' localhost:50051 api.OrderService/GetOrdersByUser
```

#### Apenas primeira página
```bash
grpcurl -plaintext -d '{
  "user_id": "0a326204-a40d-4f28-a7f9-8c87bf1f87e6"
}' localhost:50051 api.OrderService/GetOrdersByUser
```

## 📝 Usando em Código Go

### Setup do Cliente

```go
package main

import (
    "context"
    "log"
    "time"

    pb "meuApp/pkg/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Conectar ao servidor gRPC
    conn, err := grpc.Dial("localhost:50051", 
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    // Criar clientes
    productClient := pb.NewProductServiceClient(conn)
    userClient := pb.NewUserServiceClient(conn)
    orderClient := pb.NewOrderServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // Usar os clientes...
}
```

### Listar Produtos

```go
func listProducts(client pb.ProductServiceClient, ctx context.Context) {
    req := &pb.ListProductsRequest{
        Page:     1,
        PageSize: 20,
    }

    resp, err := client.ListProducts(ctx, req)
    if err != nil {
        log.Fatalf("could not list products: %v", err)
    }

    log.Printf("Total products: %d", resp.TotalItems)
    log.Printf("Total pages: %d", resp.TotalPages)
    log.Printf("Current page: %d", resp.Page)
    
    for _, product := range resp.Products {
        log.Printf("Product: %s - %s", product.Id, product.Name)
    }
}
```

### Listar Usuários

```go
func listUsers(client pb.UserServiceClient, ctx context.Context) {
    req := &pb.ListUsersRequest{
        Page:     1,
        PageSize: 10,
    }

    resp, err := client.ListUsers(ctx, req)
    if err != nil {
        log.Fatalf("could not list users: %v", err)
    }

    log.Printf("Total users: %d", resp.Total)
    
    for _, user := range resp.Users {
        log.Printf("User: %s - %s", user.Id, user.Name)
    }
}
```

### Listar Pedidos

```go
func listOrders(client pb.OrderServiceClient, ctx context.Context) {
    req := &pb.ListOrdersRequest{
        Page:     1,
        PageSize: 10,
    }

    resp, err := client.ListOrders(ctx, req)
    if err != nil {
        log.Fatalf("could not list orders: %v", err)
    }

    log.Printf("Total orders: %d", resp.Total)
    
    for _, order := range resp.Orders {
        log.Printf("Order: %s - User: %s - Total: $%.2f", 
            order.Id, order.UserId, order.TotalPrice)
    }
}
```

### Listar Pedidos de um Usuário

```go
func getOrdersByUser(client pb.OrderServiceClient, ctx context.Context, userID string) {
    req := &pb.GetOrdersByUserRequest{
        UserId:   userID,
        Page:     1,
        PageSize: 5,
    }

    resp, err := client.GetOrdersByUser(ctx, req)
    if err != nil {
        log.Fatalf("could not get user orders: %v", err)
    }

    log.Printf("User has %d orders", resp.Total)
    
    for _, order := range resp.Orders {
        log.Printf("Order: %s - Status: %s - Total: $%.2f", 
            order.Id, order.Status, order.TotalPrice)
    }
}
```

### Implementação Completa com Paginação Automática

```go
func getAllProducts(client pb.ProductServiceClient, ctx context.Context) []*pb.Product {
    var allProducts []*pb.Product
    page := int32(1)
    pageSize := int32(50)

    for {
        req := &pb.ListProductsRequest{
            Page:     page,
            PageSize: pageSize,
        }

        resp, err := client.ListProducts(ctx, req)
        if err != nil {
            log.Fatalf("could not list products: %v", err)
        }

        allProducts = append(allProducts, resp.Products...)

        log.Printf("Fetched page %d of %d", page, resp.TotalPages)

        // Se chegou na última página, para
        if page >= resp.TotalPages {
            break
        }

        page++
    }

    log.Printf("Total products fetched: %d", len(allProducts))
    return allProducts
}
```

## 🐍 Exemplo em Python

```python
import grpc
import proto.product_pb2 as product_pb2
import proto.product_pb2_grpc as product_pb2_grpc

def list_products():
    # Conectar ao servidor
    channel = grpc.insecure_channel('localhost:50051')
    stub = product_pb2_grpc.ProductServiceStub(channel)

    # Criar requisição
    request = product_pb2.ListProductsRequest(
        page=1,
        page_size=20
    )

    # Fazer chamada
    response = stub.ListProducts(request)

    # Processar resposta
    print(f"Total products: {response.total_items}")
    print(f"Total pages: {response.total_pages}")
    
    for product in response.products:
        print(f"Product: {product.id} - {product.name}")

if __name__ == '__main__':
    list_products()
```

## 📊 Testando Performance

### Script para Benchmark

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    pb "meuApp/pkg/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func benchmark() {
    conn, err := grpc.Dial("localhost:50051", 
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewProductServiceClient(conn)
    ctx := context.Background()

    // Testar diferentes tamanhos de página
    pageSizes := []int32{10, 25, 50, 100}

    for _, pageSize := range pageSizes {
        start := time.Now()

        req := &pb.ListProductsRequest{
            Page:     1,
            PageSize: pageSize,
        }

        resp, err := client.ListProducts(ctx, req)
        if err != nil {
            log.Printf("Error with page_size %d: %v", pageSize, err)
            continue
        }

        duration := time.Since(start)

        fmt.Printf("PageSize: %d | Items: %d | Duration: %v\n", 
            pageSize, len(resp.Products), duration)
    }
}
```

## 🔍 Descobrindo Serviços Disponíveis

```bash
# Listar todos os serviços
grpcurl -plaintext localhost:50051 list

# Listar métodos de um serviço
grpcurl -plaintext localhost:50051 list api.ProductService

# Ver definição de um método
grpcurl -plaintext localhost:50051 describe api.ProductService.ListProducts
```

## 🛠️ Troubleshooting

### Erro: "connection refused"
```bash
# Verificar se o servidor está rodando
lsof -i :50051

# Ou
netstat -an | grep 50051
```

### Erro: "method not found"
```bash
# Verificar métodos disponíveis
grpcurl -plaintext localhost:50051 list api.OrderService
```

### Testar conectividade
```bash
# Health check (se implementado)
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

## 📚 Recursos Adicionais

- [gRPC Go Tutorial](https://grpc.io/docs/languages/go/quickstart/)
- [grpcurl Documentation](https://github.com/fullstorydev/grpcurl)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers/docs/proto3)

## 🎉 Conclusão

Todos os serviços gRPC agora suportam paginação de forma consistente e eficiente. Use os exemplos acima como ponto de partida para integrar com seus clientes!
