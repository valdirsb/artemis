# 24. API Reference - Referência Completa das APIs

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [HTTP REST API](#http-rest-api)
3. [gRPC API](#grpc-api)
4. [Autenticação & Autorização](#autenticação--autorização)
5. [Paginação](#paginação)
6. [Filtros & Ordenação](#filtros--ordenação)
7. [Tratamento de Erros](#tratamento-de-erros)
8. [Rate Limiting](#rate-limiting)
9. [Versionamento](#versionamento)
10. [Swagger/OpenAPI](#swaggeropenapi)

---

## 📌 Visão Geral

O Artemis Framework expõe duas interfaces de API:

### 🌐 HTTP REST API
- **Base URL:** `http://localhost:8080/api/v1`
- **Formato:** JSON
- **Protocolo:** HTTP/1.1
- **Autenticação:** Bearer Token (JWT)

### 🔌 gRPC API
- **Endpoint:** `localhost:50051`
- **Formato:** Protocol Buffers
- **Protocolo:** HTTP/2
- **Autenticação:** Metadata Token

### 📊 Comparação

| Característica | HTTP REST | gRPC |
|----------------|-----------|------|
| Performance | Boa | Excelente |
| Browser Support | ✅ Sim | ❌ Limitado |
| Streaming | ❌ Não | ✅ Sim |
| Contratos | Swagger | Proto files |
| Debugging | Fácil | Moderado |

---

## 🌐 HTTP REST API

### Base URL

```
http://localhost:8080/api/v1
```

### Headers Comuns

```http
Content-Type: application/json
Accept: application/json
Authorization: Bearer {token}
X-Request-ID: {uuid}
```

---

## 🛍️ Products API

### 1. Create Product

Cria um novo produto.

**Endpoint:** `POST /api/v1/products`

**Request Body:**
```json
{
  "name": "Laptop Dell",
  "description": "Laptop Dell Inspiron 15",
  "category_id": "electronics-123",
  "price": 2999.99,
  "stock": 50
}
```

**Response:** `201 Created`
```json
{
  "id": "prod-abc123",
  "name": "Laptop Dell",
  "description": "Laptop Dell Inspiron 15",
  "category_id": "electronics-123",
  "price": 2999.99,
  "stock": 50,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**Validações:**
- `name`: obrigatório, max 255 caracteres
- `price`: obrigatório, > 0
- `stock`: obrigatório, >= 0
- `category_id`: obrigatório, UUID válido

**Erros:**
```json
// 400 Bad Request
{
  "error": "validation_error",
  "message": "Invalid product data",
  "details": [
    {
      "field": "price",
      "message": "must be greater than 0"
    }
  ]
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Laptop Dell",
    "description": "Laptop Dell Inspiron 15",
    "category_id": "electronics-123",
    "price": 2999.99,
    "stock": 50
  }'
```

---

### 2. Get Product by ID

Busca um produto pelo ID.

**Endpoint:** `GET /api/v1/products/:id`

**Parameters:**
- `id` (path, required): Product ID (UUID)

**Response:** `200 OK`
```json
{
  "id": "prod-abc123",
  "name": "Laptop Dell",
  "description": "Laptop Dell Inspiron 15",
  "category_id": "electronics-123",
  "price": 2999.99,
  "stock": 50,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**Erros:**
```json
// 404 Not Found
{
  "error": "not_found",
  "message": "Product not found",
  "details": {
    "product_id": "prod-abc123"
  }
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/products/prod-abc123 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 3. List Products

Lista produtos com paginação e filtros.

**Endpoint:** `GET /api/v1/products`

**Query Parameters:**
- `page` (int, optional): Página (default: 1)
- `page_size` (int, optional): Itens por página (default: 10, max: 100)
- `category_id` (string, optional): Filtrar por categoria
- `min_price` (float, optional): Preço mínimo
- `max_price` (float, optional): Preço máximo
- `in_stock` (bool, optional): Somente em estoque
- `sort_by` (string, optional): Campo para ordenar (name, price, created_at)
- `sort_order` (string, optional): Ordem (asc, desc)

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": "prod-abc123",
      "name": "Laptop Dell",
      "description": "Laptop Dell Inspiron 15",
      "category_id": "electronics-123",
      "price": 2999.99,
      "stock": 50,
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    },
    {
      "id": "prod-def456",
      "name": "Mouse Logitech",
      "description": "Mouse Wireless MX Master 3",
      "category_id": "electronics-123",
      "price": 399.99,
      "stock": 120,
      "created_at": "2025-01-14T15:20:00Z",
      "updated_at": "2025-01-14T15:20:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 45,
    "total_pages": 5,
    "has_next": true,
    "has_previous": false
  }
}
```

**cURL Example:**
```bash
# Lista todos os produtos
curl -X GET "http://localhost:8080/api/v1/products?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filtra por categoria e preço
curl -X GET "http://localhost:8080/api/v1/products?category_id=electronics-123&min_price=100&max_price=1000&in_stock=true" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Ordenar por preço descendente
curl -X GET "http://localhost:8080/api/v1/products?sort_by=price&sort_order=desc" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 4. Update Product

Atualiza um produto existente.

**Endpoint:** `PUT /api/v1/products/:id`

**Parameters:**
- `id` (path, required): Product ID

**Request Body:**
```json
{
  "name": "Laptop Dell XPS 15",
  "description": "Laptop Dell XPS 15 - Updated",
  "category_id": "electronics-123",
  "price": 3499.99,
  "stock": 45
}
```

**Response:** `200 OK`
```json
{
  "id": "prod-abc123",
  "name": "Laptop Dell XPS 15",
  "description": "Laptop Dell XPS 15 - Updated",
  "category_id": "electronics-123",
  "price": 3499.99,
  "stock": 45,
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T11:45:00Z"
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/products/prod-abc123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Laptop Dell XPS 15",
    "description": "Laptop Dell XPS 15 - Updated",
    "category_id": "electronics-123",
    "price": 3499.99,
    "stock": 45
  }'
```

---

### 5. Update Product Stock

Atualiza apenas o estoque do produto.

**Endpoint:** `PUT /api/v1/products/:id/stock`

**Request Body:**
```json
{
  "quantity": 100,
  "operation": "add"  // "add" ou "set"
}
```

**Response:** `200 OK`
```json
{
  "id": "prod-abc123",
  "name": "Laptop Dell XPS 15",
  "stock": 145,
  "updated_at": "2025-01-15T12:00:00Z"
}
```

**cURL Example:**
```bash
# Adicionar estoque
curl -X PUT http://localhost:8080/api/v1/products/prod-abc123/stock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"quantity": 100, "operation": "add"}'

# Definir estoque (substituir)
curl -X PUT http://localhost:8080/api/v1/products/prod-abc123/stock \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"quantity": 50, "operation": "set"}'
```

---

### 6. Delete Product

Remove um produto.

**Endpoint:** `DELETE /api/v1/products/:id`

**Response:** `204 No Content`

**Erros:**
```json
// 404 Not Found
{
  "error": "not_found",
  "message": "Product not found"
}

// 409 Conflict
{
  "error": "conflict",
  "message": "Cannot delete product with active orders"
}
```

**cURL Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/prod-abc123 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 👥 Users API

### 1. Create User

Cria um novo usuário.

**Endpoint:** `POST /api/v1/users`

**Request Body:**
```json
{
  "name": "João Silva",
  "email": "joao@example.com",
  "password": "SecurePass123!"
}
```

**Response:** `201 Created`
```json
{
  "id": "user-xyz789",
  "name": "João Silva",
  "email": "joao@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**Validações:**
- `name`: obrigatório, 3-100 caracteres
- `email`: obrigatório, formato válido, único
- `password`: obrigatório, mínimo 6 caracteres

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com",
    "password": "SecurePass123!"
  }'
```

---

### 2. Get User by ID

Busca um usuário pelo ID.

**Endpoint:** `GET /api/v1/users/:id`

**Response:** `200 OK`
```json
{
  "id": "user-xyz789",
  "name": "João Silva",
  "email": "joao@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/users/user-xyz789 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 3. Update User

Atualiza dados do usuário.

**Endpoint:** `PUT /api/v1/users/:id`

**Request Body:**
```json
{
  "name": "João Pedro Silva",
  "email": "joao.silva@example.com"
}
```

**Response:** `200 OK`
```json
{
  "id": "user-xyz789",
  "name": "João Pedro Silva",
  "email": "joao.silva@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T14:00:00Z"
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/user-xyz789 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "João Pedro Silva",
    "email": "joao.silva@example.com"
  }'
```

---

### 4. Delete User

Remove um usuário.

**Endpoint:** `DELETE /api/v1/users/:id`

**Response:** `204 No Content`

**cURL Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/user-xyz789 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 5. Login / Validate User

Valida credenciais e retorna token JWT.

**Endpoint:** `POST /api/v1/users/login`

**Request Body:**
```json
{
  "email": "joao@example.com",
  "password": "SecurePass123!"
}
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": "user-xyz789",
    "name": "João Silva",
    "email": "joao@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-01-16T10:30:00Z"
}
```

**Erros:**
```json
// 401 Unauthorized
{
  "error": "unauthorized",
  "message": "Invalid email or password"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "password": "SecurePass123!"
  }'
```

---

## 📦 Orders API

### 1. Create Order

Cria um novo pedido.

**Endpoint:** `POST /api/v1/orders`

**Request Body:**
```json
{
  "user_id": "user-xyz789",
  "items": [
    {
      "product_id": "prod-abc123",
      "quantity": 2
    },
    {
      "product_id": "prod-def456",
      "quantity": 1
    }
  ]
}
```

**Response:** `201 Created`
```json
{
  "id": "order-123456",
  "user_id": "user-xyz789",
  "items": [
    {
      "product_id": "prod-abc123",
      "product_name": "Laptop Dell",
      "quantity": 2,
      "unit_price": 2999.99,
      "subtotal": 5999.98
    },
    {
      "product_id": "prod-def456",
      "product_name": "Mouse Logitech",
      "quantity": 1,
      "unit_price": 399.99,
      "subtotal": 399.99
    }
  ],
  "total": 6399.97,
  "status": "pending",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**Validações:**
- `user_id`: obrigatório, usuário deve existir
- `items`: obrigatório, mínimo 1 item
- `items[].product_id`: produto deve existir e ter estoque
- `items[].quantity`: > 0

**Erros:**
```json
// 400 Bad Request - Estoque insuficiente
{
  "error": "insufficient_stock",
  "message": "Product 'Laptop Dell' has insufficient stock",
  "details": {
    "product_id": "prod-abc123",
    "requested": 10,
    "available": 5
  }
}

// 404 Not Found - Produto não existe
{
  "error": "not_found",
  "message": "Product not found",
  "details": {
    "product_id": "prod-invalid"
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "user_id": "user-xyz789",
    "items": [
      {"product_id": "prod-abc123", "quantity": 2},
      {"product_id": "prod-def456", "quantity": 1}
    ]
  }'
```

---

### 2. Get Order by ID

Busca um pedido pelo ID.

**Endpoint:** `GET /api/v1/orders/:id`

**Response:** `200 OK`
```json
{
  "id": "order-123456",
  "user_id": "user-xyz789",
  "user_name": "João Silva",
  "items": [
    {
      "product_id": "prod-abc123",
      "product_name": "Laptop Dell",
      "quantity": 2,
      "unit_price": 2999.99,
      "subtotal": 5999.98
    }
  ],
  "total": 6399.97,
  "status": "confirmed",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:35:00Z"
}
```

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/orders/order-123456 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 3. Update Order Status

Atualiza o status de um pedido.

**Endpoint:** `PUT /api/v1/orders/:id/status`

**Request Body:**
```json
{
  "status": "confirmed"
}
```

**Status Válidos:**
- `pending` → `confirmed`
- `confirmed` → `shipped`
- `shipped` → `delivered`
- Qualquer status → `cancelled`

**Response:** `200 OK`
```json
{
  "id": "order-123456",
  "status": "confirmed",
  "updated_at": "2025-01-15T10:35:00Z"
}
```

**Erros:**
```json
// 400 Bad Request - Transição inválida
{
  "error": "invalid_status_transition",
  "message": "Cannot change status from 'delivered' to 'confirmed'"
}
```

**cURL Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/orders/order-123456/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"status": "confirmed"}'
```

---

### 4. Cancel Order

Cancela um pedido.

**Endpoint:** `POST /api/v1/orders/:id/cancel`

**Response:** `200 OK`
```json
{
  "id": "order-123456",
  "status": "cancelled",
  "cancelled_at": "2025-01-15T11:00:00Z"
}
```

**Regras:**
- Só pode cancelar se status for `pending` ou `confirmed`
- Estoque dos produtos é restaurado

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/orders/order-123456/cancel \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 5. List Orders by User

Lista todos os pedidos de um usuário.

**Endpoint:** `GET /api/v1/orders/user/:user_id`

**Query Parameters:**
- `page` (int, optional): Página (default: 1)
- `page_size` (int, optional): Itens por página (default: 10)
- `status` (string, optional): Filtrar por status

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": "order-123456",
      "user_id": "user-xyz789",
      "total": 6399.97,
      "status": "confirmed",
      "created_at": "2025-01-15T10:30:00Z"
    },
    {
      "id": "order-789012",
      "user_id": "user-xyz789",
      "total": 1299.99,
      "status": "delivered",
      "created_at": "2025-01-10T14:20:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 15,
    "total_pages": 2,
    "has_next": true,
    "has_previous": false
  }
}
```

**cURL Example:**
```bash
# Todos os pedidos do usuário
curl -X GET "http://localhost:8080/api/v1/orders/user/user-xyz789?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filtrar por status
curl -X GET "http://localhost:8080/api/v1/orders/user/user-xyz789?status=confirmed" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🔌 gRPC API

### Configuração do Cliente

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
    conn, err := grpc.Dial(
        "localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatalf("failed to connect: %v", err)
    }
    defer conn.Close()

    // Criar cliente
    productClient := pb.NewProductServiceClient(conn)
    userClient := pb.NewUserServiceClient(conn)
    orderClient := pb.NewOrderServiceClient(conn)

    // Usar o cliente
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Exemplo: Criar produto
    resp, err := productClient.CreateProduct(ctx, &pb.CreateProductRequest{
        Name:        "Laptop Dell",
        Description: "Laptop Dell Inspiron 15",
        Price:       2999.99,
        Quantity:    50,
    })
    if err != nil {
        log.Fatalf("CreateProduct failed: %v", err)
    }

    log.Printf("Product created: %v", resp.Product)
}
```

---

## 🛍️ ProductService (gRPC)

### Proto Definition

```protobuf
syntax = "proto3";

package api;

option go_package = "go-modular-monolith/proto";

service ProductService {
  rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
  rpc GetProduct(GetProductRequest) returns (GetProductResponse);
  rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse);
  rpc DeleteProduct(DeleteProductRequest) returns (DeleteProductResponse);
  rpc ListProducts(ListProductsRequest) returns (ListProductsResponse);
  rpc SearchProducts(SearchProductsRequest) returns (ListProductsResponse);
}

message Product {
  string id = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  int32 quantity = 5;
  string created_at = 6;
  string updated_at = 7;
}

message CreateProductRequest {
  string name = 1;
  string description = 2;
  double price = 3;
  int32 quantity = 4;
}

message CreateProductResponse {
  Product product = 1;
  string message = 2;
}

message ListProductsRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message ListProductsResponse {
  repeated Product products = 1;
  int32 total_items = 2;
  int32 page = 3;
  int32 page_size = 4;
  int32 total_pages = 5;
  string message = 6;
}
```

### 1. CreateProduct

```go
ctx := context.Background()

req := &pb.CreateProductRequest{
    Name:        "Laptop Dell",
    Description: "Laptop Dell Inspiron 15",
    Price:       2999.99,
    Quantity:    50,
}

resp, err := productClient.CreateProduct(ctx, req)
if err != nil {
    log.Fatalf("CreateProduct failed: %v", err)
}

fmt.Printf("Product ID: %s\n", resp.Product.Id)
```

**Response:**
```go
&CreateProductResponse{
    Product: &Product{
        Id:          "prod-abc123",
        Name:        "Laptop Dell",
        Description: "Laptop Dell Inspiron 15",
        Price:       2999.99,
        Quantity:    50,
        CreatedAt:   "2025-01-15T10:30:00Z",
        UpdatedAt:   "2025-01-15T10:30:00Z",
    },
    Message: "Product created successfully",
}
```

---

### 2. GetProduct

```go
req := &pb.GetProductRequest{
    Id: "prod-abc123",
}

resp, err := productClient.GetProduct(ctx, req)
if err != nil {
    // Handle error (e.g., codes.NotFound)
    log.Fatalf("GetProduct failed: %v", err)
}

fmt.Printf("Product: %v\n", resp.Product)
```

**Erros gRPC:**
```go
import "google.golang.org/grpc/codes"
import "google.golang.org/grpc/status"

if err != nil {
    st, ok := status.FromError(err)
    if ok {
        switch st.Code() {
        case codes.NotFound:
            fmt.Println("Product not found")
        case codes.InvalidArgument:
            fmt.Println("Invalid product ID")
        default:
            fmt.Printf("Error: %v\n", st.Message())
        }
    }
}
```

---

### 3. ListProducts

```go
req := &pb.ListProductsRequest{
    Page:     1,
    PageSize: 10,
}

resp, err := productClient.ListProducts(ctx, req)
if err != nil {
    log.Fatalf("ListProducts failed: %v", err)
}

fmt.Printf("Total: %d products\n", resp.TotalItems)
for _, product := range resp.Products {
    fmt.Printf("- %s: $%.2f\n", product.Name, product.Price)
}
```

**Response:**
```go
&ListProductsResponse{
    Products: []*Product{
        {Id: "prod-abc123", Name: "Laptop Dell", Price: 2999.99},
        {Id: "prod-def456", Name: "Mouse Logitech", Price: 399.99},
    },
    TotalItems: 45,
    Page:       1,
    PageSize:   10,
    TotalPages: 5,
    Message:    "Products retrieved successfully",
}
```

---

### 4. UpdateProduct

```go
req := &pb.UpdateProductRequest{
    Id:          "prod-abc123",
    Name:        "Laptop Dell XPS 15",
    Description: "Laptop Dell XPS 15 - Updated",
    Price:       3499.99,
    Quantity:    45,
}

resp, err := productClient.UpdateProduct(ctx, req)
if err != nil {
    log.Fatalf("UpdateProduct failed: %v", err)
}

fmt.Printf("Product updated: %v\n", resp.Product)
```

---

### 5. DeleteProduct

```go
req := &pb.DeleteProductRequest{
    Id: "prod-abc123",
}

resp, err := productClient.DeleteProduct(ctx, req)
if err != nil {
    log.Fatalf("DeleteProduct failed: %v", err)
}

fmt.Printf("Message: %s\n", resp.Message)
```

---

## 👥 UserService (gRPC)

### Proto Definition

```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc GetUserByEmail(GetUserByEmailRequest) returns (GetUserResponse);
}

message User {
  string id = 1;
  string name = 2;
  string email = 3;
  string created_at = 4;
  string updated_at = 5;
}

message CreateUserRequest {
  string name = 1;
  string email = 2;
  string password = 3;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}
```

### Exemplos

```go
// CreateUser
userResp, err := userClient.CreateUser(ctx, &pb.CreateUserRequest{
    Name:     "João Silva",
    Email:    "joao@example.com",
    Password: "SecurePass123!",
})

// GetUser
user, err := userClient.GetUser(ctx, &pb.GetUserRequest{
    Id: "user-xyz789",
})

// GetUserByEmail
user, err := userClient.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{
    Email: "joao@example.com",
})

// ListUsers
users, err := userClient.ListUsers(ctx, &pb.ListUsersRequest{
    Page:     1,
    PageSize: 10,
})
```

---

## 📦 OrderService (gRPC)

### Proto Definition

```protobuf
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (UpdateOrderStatusResponse);
  rpc DeleteOrder(DeleteOrderRequest) returns (DeleteOrderResponse);
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
  rpc GetOrdersByUser(GetOrdersByUserRequest) returns (ListOrdersResponse);
}

message Order {
  string id = 1;
  string user_id = 2;
  repeated OrderItem items = 3;
  double total_price = 4;
  string status = 5;
  string created_at = 6;
  string updated_at = 7;
}

message OrderItem {
  string product_id = 1;
  int32 quantity = 2;
  double price = 3;
}

message CreateOrderRequest {
  string user_id = 1;
  repeated CreateOrderItem items = 2;
}

message CreateOrderItem {
  string product_id = 1;
  int32 quantity = 2;
}

message GetOrdersByUserRequest {
  string user_id = 1;
  int32 page = 2;
  int32 page_size = 3;
}
```

### Exemplos

```go
// CreateOrder
orderResp, err := orderClient.CreateOrder(ctx, &pb.CreateOrderRequest{
    UserId: "user-xyz789",
    Items: []*pb.CreateOrderItem{
        {ProductId: "prod-abc123", Quantity: 2},
        {ProductId: "prod-def456", Quantity: 1},
    },
})

// GetOrder
order, err := orderClient.GetOrder(ctx, &pb.GetOrderRequest{
    Id: "order-123456",
})

// UpdateOrderStatus
updated, err := orderClient.UpdateOrderStatus(ctx, &pb.UpdateOrderStatusRequest{
    Id:     "order-123456",
    Status: "confirmed",
})

// GetOrdersByUser (com paginação)
orders, err := orderClient.GetOrdersByUser(ctx, &pb.GetOrdersByUserRequest{
    UserId:   "user-xyz789",
    Page:     1,
    PageSize: 10,
})

fmt.Printf("User has %d orders\n", orders.Total)
for _, order := range orders.Orders {
    fmt.Printf("Order %s: $%.2f (%s)\n", order.Id, order.TotalPrice, order.Status)
}
```

---

## 🔐 Autenticação & Autorização

### JWT Token

O framework usa **JWT (JSON Web Tokens)** para autenticação.

#### 1. Obter Token

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "password": "SecurePass123!"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLXh5ejc4OSIsImVtYWlsIjoiam9hb0BleGFtcGxlLmNvbSIsImV4cCI6MTcwNTQwMTYwMH0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  "expires_at": "2025-01-16T10:30:00Z"
}
```

#### 2. Usar Token

**HTTP REST:**
```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**gRPC:**
```go
import "google.golang.org/grpc/metadata"

// Adicionar token aos metadados
md := metadata.Pairs("authorization", "Bearer "+token)
ctx := metadata.NewOutgoingContext(context.Background(), md)

// Fazer chamada com contexto
resp, err := productClient.GetProduct(ctx, req)
```

#### 3. Payload do Token

```json
{
  "sub": "user-xyz789",
  "email": "joao@example.com",
  "name": "João Silva",
  "roles": ["user"],
  "iat": 1705315200,
  "exp": 1705401600
}
```

### Middleware de Autenticação

```go
// pkg/adapters/http/middleware/auth.go
func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Missing authorization header"})
            c.Abort()
            return
        }

        // Extrair token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // Validar token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Extrair claims
        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["sub"])
        c.Set("email", claims["email"])

        c.Next()
    }
}
```

### Roles & Permissions

```go
// Verificar role
func RequireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRoles := c.GetStringSlice("roles")
        
        hasRole := false
        for _, r := range userRoles {
            if r == role {
                hasRole = true
                break
            }
        }

        if !hasRole {
            c.JSON(403, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// Uso
router.DELETE("/products/:id", 
    AuthMiddleware(secret),
    RequireRole("admin"),
    productHandler.DeleteProduct,
)
```

---

## 📄 Paginação

### Padrão de Paginação

Todos os endpoints de listagem suportam paginação:

**Query Parameters:**
- `page` (int): Número da página (começa em 1)
- `page_size` (int): Itens por página (default: 10, max: 100)

**Response:**
```json
{
  "items": [...],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total_items": 45,
    "total_pages": 5,
    "has_next": true,
    "has_previous": false
  }
}
```

### Cálculo de Páginas

```
total_pages = ceil(total_items / page_size)
has_next = page < total_pages
has_previous = page > 1
```

### Exemplos

```bash
# Primeira página (10 itens)
curl "http://localhost:8080/api/v1/products?page=1&page_size=10"

# Segunda página (20 itens)
curl "http://localhost:8080/api/v1/products?page=2&page_size=20"

# Última página
curl "http://localhost:8080/api/v1/products?page=5&page_size=10"
```

### Navegação

```javascript
// JavaScript example
async function fetchAllProducts() {
    let page = 1;
    let hasNext = true;
    const allProducts = [];

    while (hasNext) {
        const response = await fetch(
            `http://localhost:8080/api/v1/products?page=${page}&page_size=50`
        );
        const data = await response.json();
        
        allProducts.push(...data.items);
        hasNext = data.pagination.has_next;
        page++;
    }

    return allProducts;
}
```

---

## 🔍 Filtros & Ordenação

### Filtros Disponíveis

#### Products
- `category_id`: Filtrar por categoria
- `min_price`: Preço mínimo
- `max_price`: Preço máximo
- `in_stock`: Somente em estoque (true/false)
- `search`: Busca por nome ou descrição

#### Orders
- `status`: Filtrar por status (pending, confirmed, shipped, delivered, cancelled)
- `user_id`: Filtrar por usuário
- `start_date`: Data inicial (ISO 8601)
- `end_date`: Data final (ISO 8601)

### Ordenação

**Query Parameters:**
- `sort_by`: Campo para ordenar
- `sort_order`: Ordem (asc, desc)

**Campos ordenáveis:**

**Products:**
- `name`
- `price`
- `stock`
- `created_at`

**Orders:**
- `total`
- `status`
- `created_at`

**Users:**
- `name`
- `email`
- `created_at`

### Exemplos Combinados

```bash
# Produtos de eletrônicos, entre R$100 e R$500, ordenados por preço
curl "http://localhost:8080/api/v1/products?\
category_id=electronics-123&\
min_price=100&\
max_price=500&\
sort_by=price&\
sort_order=asc&\
page=1&\
page_size=20"

# Pedidos confirmados de um usuário, últimos primeiro
curl "http://localhost:8080/api/v1/orders/user/user-xyz789?\
status=confirmed&\
sort_by=created_at&\
sort_order=desc"

# Buscar produtos por nome
curl "http://localhost:8080/api/v1/products?\
search=laptop&\
in_stock=true"
```

---

## ❌ Tratamento de Erros

### HTTP Status Codes

| Status | Significado | Uso |
|--------|-------------|-----|
| 200 | OK | Sucesso (GET, PUT) |
| 201 | Created | Recurso criado (POST) |
| 204 | No Content | Sucesso sem retorno (DELETE) |
| 400 | Bad Request | Dados inválidos |
| 401 | Unauthorized | Sem autenticação |
| 403 | Forbidden | Sem permissão |
| 404 | Not Found | Recurso não encontrado |
| 409 | Conflict | Conflito (ex: email já existe) |
| 422 | Unprocessable Entity | Validação falhou |
| 429 | Too Many Requests | Rate limit excedido |
| 500 | Internal Server Error | Erro do servidor |
| 503 | Service Unavailable | Serviço indisponível |

### Formato de Erro HTTP

```json
{
  "error": "error_code",
  "message": "Human-readable error message",
  "details": {
    "field": "additional context"
  },
  "request_id": "req-abc123",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

### Exemplos de Erros

**Validation Error:**
```json
{
  "error": "validation_error",
  "message": "Invalid input data",
  "details": [
    {
      "field": "email",
      "message": "must be a valid email address"
    },
    {
      "field": "price",
      "message": "must be greater than 0"
    }
  ]
}
```

**Not Found:**
```json
{
  "error": "not_found",
  "message": "Product not found",
  "details": {
    "product_id": "prod-invalid"
  }
}
```

**Conflict:**
```json
{
  "error": "conflict",
  "message": "Email already exists",
  "details": {
    "email": "joao@example.com"
  }
}
```

**Unauthorized:**
```json
{
  "error": "unauthorized",
  "message": "Invalid or expired token"
}
```

### gRPC Status Codes

| Code | Significado | HTTP Equiv |
|------|-------------|-----------|
| OK | Sucesso | 200 |
| CANCELLED | Cancelado | 499 |
| INVALID_ARGUMENT | Argumento inválido | 400 |
| NOT_FOUND | Não encontrado | 404 |
| ALREADY_EXISTS | Já existe | 409 |
| PERMISSION_DENIED | Sem permissão | 403 |
| UNAUTHENTICATED | Não autenticado | 401 |
| RESOURCE_EXHAUSTED | Limite excedido | 429 |
| UNIMPLEMENTED | Não implementado | 501 |
| INTERNAL | Erro interno | 500 |
| UNAVAILABLE | Indisponível | 503 |
| DEADLINE_EXCEEDED | Timeout | 504 |

### Tratamento de Erro gRPC

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// Server side
func (s *ProductGRPCHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
    if req.Id == "" {
        return nil, status.Error(codes.InvalidArgument, "product ID is required")
    }

    product, err := s.productService.GetProductByID(ctx, req.Id)
    if err != nil {
        // Converter AppError para gRPC status
        if errors.Is(err, ErrNotFound) {
            return nil, status.Error(codes.NotFound, "product not found")
        }
        return nil, status.Error(codes.Internal, "internal error")
    }

    return &pb.GetProductResponse{Product: toProtoProduct(product)}, nil
}

// Client side
resp, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: "invalid"})
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        switch st.Code() {
        case codes.NotFound:
            fmt.Println("Product not found")
        case codes.InvalidArgument:
            fmt.Printf("Invalid argument: %s\n", st.Message())
        default:
            fmt.Printf("Error: %s\n", st.Message())
        }
    }
}
```

---

## 🚦 Rate Limiting

### Configuração

```yaml
# framework.yaml
rate_limiting:
  enabled: true
  requests_per_minute: 60
  burst: 10
```

### Headers de Rate Limit

```http
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1705315260
```

### Resposta quando Limite Excedido

```json
{
  "error": "rate_limit_exceeded",
  "message": "Too many requests. Please try again later.",
  "retry_after": 45
}
```

HTTP Status: `429 Too Many Requests`

### Implementação

```go
import "github.com/ulule/limiter/v3"
import "github.com/ulule/limiter/v3/drivers/middleware/gin"

// Criar limiter
rate := limiter.Rate{
    Period: 1 * time.Minute,
    Limit:  60,
}

store := memory.NewStore()
instance := limiter.New(store, rate)

// Middleware
middleware := ginlimiter.NewMiddleware(instance)
router.Use(middleware)
```

---

## 🔄 Versionamento

### Strategy: URL Path Versioning

```
/api/v1/products
/api/v2/products
```

### Versões Disponíveis

- **v1**: Versão atual (estável)

### Breaking Changes

Quando introduzir breaking changes:
1. Criar nova versão (v2)
2. Manter v1 por pelo menos 6 meses
3. Documentar deprecation

### Deprecation Header

```http
X-API-Deprecation: true
X-API-Sunset-Date: 2025-07-15T00:00:00Z
```

---

## 📖 Swagger/OpenAPI

### Acessar Swagger UI

```
http://localhost:8080/swagger/index.html
```

### Gerar Documentação

```bash
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar docs
swag init

# Arquivos gerados:
# - docs/docs.go
# - docs/swagger.json
# - docs/swagger.yaml
```

### Anotações Swagger

```go
// @Summary Create a new product
// @Description Create a new product with name, description, price and stock
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product data"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/products [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // ...
}
```

### Download da Especificação

```bash
# JSON
curl http://localhost:8080/swagger/doc.json > api-spec.json

# YAML
curl http://localhost:8080/swagger/doc.yaml > api-spec.yaml
```

---

## 🔍 Health Checks

### Endpoints

**Health Check Simples:**
```
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

**Health Check Detalhado:**
```
GET /framework/info
```

**Response:**
```json
{
  "framework": "Artemis",
  "version": "1.0.0",
  "environment": "production",
  "uptime_seconds": 3600,
  "database": {
    "status": "connected",
    "ping_ms": 2
  },
  "modules": ["user", "product", "order"],
  "timestamp": "2025-01-15T10:30:00Z"
}
```

---

## 📊 Métricas & Monitoramento

### Prometheus Metrics

```
GET /metrics
```

**Métricas Disponíveis:**
- `http_requests_total{method, endpoint, status}`
- `http_request_duration_seconds{method, endpoint}`
- `grpc_requests_total{service, method, status}`
- `grpc_request_duration_seconds{service, method}`
- `database_query_duration_seconds{operation}`
- `active_connections`

### Exemplo

```prometheus
# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",endpoint="/api/v1/products",status="200"} 1523

# HELP http_request_duration_seconds HTTP request duration
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",endpoint="/api/v1/products",le="0.1"} 1200
http_request_duration_seconds_bucket{method="GET",endpoint="/api/v1/products",le="0.5"} 1500
http_request_duration_seconds_bucket{method="GET",endpoint="/api/v1/products",le="1"} 1523
```

---

## 🧪 Testando as APIs

### Postman Collection

Importe a collection:
```bash
curl http://localhost:8080/api/postman-collection.json > artemis-api.postman_collection.json
```

### cURL Scripts

```bash
#!/bin/bash

# Setup
BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

# Login
login() {
    RESPONSE=$(curl -s -X POST "$BASE_URL/users/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"admin@example.com","password":"admin123"}')
    
    TOKEN=$(echo $RESPONSE | jq -r '.token')
    echo "Token: $TOKEN"
}

# Create Product
create_product() {
    curl -X POST "$BASE_URL/products" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "name": "Test Product",
            "description": "Test Description",
            "category_id": "cat-123",
            "price": 99.99,
            "stock": 100
        }'
}

# Run
login
create_product
```

### grpcurl

```bash
# Instalar grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Listar serviços
grpcurl -plaintext localhost:50051 list

# Listar métodos
grpcurl -plaintext localhost:50051 list api.ProductService

# Chamar método
grpcurl -plaintext -d '{"id": "prod-abc123"}' \
    localhost:50051 api.ProductService/GetProduct

# Criar produto
grpcurl -plaintext -d '{
    "name": "Test Product",
    "description": "Test Description",
    "price": 99.99,
    "quantity": 100
}' localhost:50051 api.ProductService/CreateProduct
```

---

## 📚 Próximos Passos

1. **[25-complete-examples.md](25-complete-examples.md)** - Exemplos completos de aplicações
2. **[26-faq.md](26-faq.md)** - Perguntas frequentes
3. **[21-testing-strategy.md](21-testing-strategy.md)** - Estratégia de testes
4. **[19-error-handling.md](19-error-handling.md)** - Tratamento de erros

---

## 💡 Dicas Importantes

### 🎯 Performance

1. **Use paginação sempre** para evitar carregamento de muitos dados
2. **Configure timeouts** apropriados (5s para queries simples, 30s para operações complexas)
3. **Use gRPC para comunicação interna** entre serviços (mais rápido que HTTP)
4. **Cache responses** quando apropriado (Redis)

### 🔒 Segurança

1. **Sempre valide input** no backend (nunca confie no cliente)
2. **Use HTTPS** em produção
3. **Rate limiting** para prevenir abuso
4. **Sanitize erros** antes de enviar ao cliente (não expor detalhes internos)

### 🧪 Testes

1. **Teste ambas APIs** (HTTP e gRPC)
2. **Use mocks** para testes unitários
3. **Testes de integração** com banco real
4. **Contract testing** para validar proto/swagger

---

<div align="center">

**[⬆️ Voltar ao Topo](#24-api-reference---referência-completa-das-apis)**

**Documentação criada com ❤️ pela Equipe Artemis**

</div>
