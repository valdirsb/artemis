# Resumo: Paginação nos Handlers gRPC

## 📋 Visão Geral

Implementação de paginação nos handlers gRPC para os módulos **User**, **Order** e **Product**, permitindo que clientes gRPC façam consultas paginadas de forma eficiente.

## 🔧 Arquivos Modificados

### 1. Proto Definitions

#### `proto/order.proto`
```protobuf
service OrderService {
  // Novo método adicionado
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
  
  // Método atualizado com paginação
  rpc GetOrdersByUser(GetOrdersByUserRequest) returns (ListOrdersResponse);
}

message ListOrdersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message GetOrdersByUserRequest {
  string user_id = 1;
  int32 page = 2;
  int32 page_size = 3;
}
```

### 2. User gRPC Handler

#### `internal/modules/user/adapters/grpc/user_grpc_handler.go`

**Antes:**
```go
func (s *UserGRPCHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    // Retornava resposta vazia
    return &pb.ListUsersResponse{
        Users:    []*pb.User{},
        Total:    0,
        Page:     req.Page,
        PageSize: req.PageSize,
        Message:  "Users retrieved successfully",
    }, nil
}
```

**Depois:**
```go
func (s *UserGRPCHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    // Parse pagination parameters
    page := 1
    if req.Page > 0 {
        page = int(req.Page)
    }

    pageSize := 10
    if req.PageSize > 0 {
        pageSize = int(req.PageSize)
    }
    if pageSize > 100 {
        pageSize = 100
    }

    // Call service with pagination
    usersPaginated, err := s.userService.ListUsers(ctx, page, pageSize)
    if err != nil {
        return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list users: %v", err))
    }

    // Convert to proto messages
    var protoUsers []*pb.User
    for _, user := range usersPaginated.Items {
        userResponse := dto.ToUserResponse(user)
        protoUsers = append(protoUsers, &pb.User{
            Id:        userResponse.ID,
            Name:      userResponse.Username,
            Email:     userResponse.Email,
            CreatedAt: userResponse.CreatedAt.Format(time.RFC3339),
            UpdatedAt: userResponse.UpdatedAt.Format(time.RFC3339),
        })
    }

    return &pb.ListUsersResponse{
        Users:    protoUsers,
        Total:    int32(usersPaginated.TotalItems),
        Page:     int32(usersPaginated.Page),
        PageSize: int32(usersPaginated.PageSize),
        Message:  "Users retrieved successfully",
    }, nil
}
```

### 3. Order gRPC Handler

#### `internal/modules/order/adapters/grpc/order_grpc_handler.go`

**Adicionado método ListOrders:**
```go
func (s *OrderGRPCHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
    // Parse pagination parameters
    page := 1
    if req.Page > 0 {
        page = int(req.Page)
    }

    pageSize := 10
    if req.PageSize > 0 {
        pageSize = int(req.PageSize)
    }
    if pageSize > 100 {
        pageSize = 100
    }

    // Call service with pagination
    ordersPaginated, err := s.orderService.ListOrders(ctx, page, pageSize)
    if err != nil {
        return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list orders: %v", err))
    }

    // Convert to proto messages
    var protoOrders []*pb.Order
    for _, order := range ordersPaginated.Items {
        orderResponse := dto.ToOrderResponse(order)
        protoOrder := &pb.Order{
            Id:         orderResponse.ID,
            UserId:     orderResponse.UserID,
            TotalPrice: orderResponse.Total,
            Status:     orderResponse.Status,
            CreatedAt:  orderResponse.CreatedAt.Format(time.RFC3339),
            UpdatedAt:  orderResponse.UpdatedAt.Format(time.RFC3339),
        }

        for _, item := range orderResponse.Items {
            protoOrder.Items = append(protoOrder.Items, &pb.OrderItem{
                ProductId: item.ProductID,
                Quantity:  int32(item.Quantity),
                Price:     item.Price,
            })
        }

        protoOrders = append(protoOrders, protoOrder)
    }

    return &pb.ListOrdersResponse{
        Orders:   protoOrders,
        Total:    int32(ordersPaginated.TotalItems),
        Page:     int32(ordersPaginated.Page),
        PageSize: int32(ordersPaginated.PageSize),
        Message:  "Orders retrieved successfully",
    }, nil
}
```

**Atualizado GetOrdersByUser:**
```go
func (s *OrderGRPCHandler) GetOrdersByUser(ctx context.Context, req *pb.GetOrdersByUserRequest) (*pb.ListOrdersResponse, error) {
    if req.UserId == "" {
        return nil, status.Error(codes.InvalidArgument, "user_id is required")
    }

    // Parse pagination parameters
    page := 1
    if req.Page > 0 {
        page = int(req.Page)
    }

    pageSize := 10
    if req.PageSize > 0 {
        pageSize = int(req.PageSize)
    }
    if pageSize > 100 {
        pageSize = 100
    }

    // Call service with pagination
    ordersPaginated, err := s.orderService.GetOrdersByUserIDPaginated(ctx, req.UserId, page, pageSize)
    if err != nil {
        return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get orders: %v", err))
    }

    // Convert to proto messages
    var protoOrders []*pb.Order
    for _, order := range ordersPaginated.Items {
        orderResponse := dto.ToOrderResponse(order)
        protoOrder := &pb.Order{
            Id:         orderResponse.ID,
            UserId:     orderResponse.UserID,
            TotalPrice: orderResponse.Total,
            Status:     orderResponse.Status,
            CreatedAt:  orderResponse.CreatedAt.Format(time.RFC3339),
            UpdatedAt:  orderResponse.UpdatedAt.Format(time.RFC3339),
        }

        for _, item := range orderResponse.Items {
            protoOrder.Items = append(protoOrder.Items, &pb.OrderItem{
                ProductId: item.ProductID,
                Quantity:  int32(item.Quantity),
                Price:     item.Price,
            })
        }

        protoOrders = append(protoOrders, protoOrder)
    }

    return &pb.ListOrdersResponse{
        Orders:   protoOrders,
        Total:    int32(ordersPaginated.TotalItems),
        Page:     int32(ordersPaginated.Page),
        PageSize: int32(ordersPaginated.PageSize),
        Message:  "Orders retrieved successfully",
    }, nil
}
```

### 4. Product gRPC Handler

O handler de Product já estava com paginação implementada ✅

## 🎯 Funcionalidades Implementadas

### User Service
- ✅ `ListUsers(page, page_size)` - Lista todos os usuários com paginação

### Order Service
- ✅ `ListOrders(page, page_size)` - Lista todos os pedidos com paginação
- ✅ `GetOrdersByUser(user_id, page, page_size)` - Lista pedidos de um usuário com paginação

### Product Service
- ✅ `ListProducts(page, page_size)` - Lista todos os produtos com paginação (já implementado)

## 📊 Parâmetros de Paginação

Todos os métodos gRPC seguem o mesmo padrão:

| Parâmetro | Tipo  | Descrição                    | Default | Máximo |
|-----------|-------|------------------------------|---------|--------|
| page      | int32 | Número da página             | 1       | -      |
| page_size | int32 | Quantidade de itens por página| 10      | 100    |

## 📦 Resposta Paginada

Estrutura de resposta para todas as listagens:

```protobuf
message ListXxxResponse {
  repeated Xxx items = 1;      // Lista de items da página atual
  int32 total = 2;             // Total de items no banco
  int32 page = 3;              // Página atual
  int32 page_size = 4;         // Tamanho da página
  string message = 5;          // Mensagem de status
}
```

## 🔄 Validações

Todos os handlers implementam as seguintes validações:

1. **Page**: Se não fornecido ou menor que 1, usa default (1)
2. **PageSize**: 
   - Se não fornecido ou menor que 1, usa default (10)
   - Se maior que 100, limita a 100 para evitar sobrecarga

## 🧪 Exemplo de Uso (Cliente gRPC)

### Go Client
```go
// Listar usuários
req := &pb.ListUsersRequest{
    Page:     1,
    PageSize: 20,
}
resp, err := client.ListUsers(ctx, req)

// Listar pedidos
req := &pb.ListOrdersRequest{
    Page:     1,
    PageSize: 10,
}
resp, err := client.ListOrders(ctx, req)

// Listar pedidos de um usuário
req := &pb.GetOrdersByUserRequest{
    UserId:   "user-123",
    Page:     1,
    PageSize: 15,
}
resp, err := client.GetOrdersByUser(ctx, req)
```

## ✅ Resultados

- ✅ Todos os handlers gRPC compilam sem erros
- ✅ Proto files regenerados com sucesso
- ✅ Padrão consistente em todos os módulos
- ✅ Validações implementadas
- ✅ Limites de segurança aplicados (max 100 items por página)

## 🎉 Conclusão

A paginação foi implementada com sucesso em todos os handlers gRPC, permitindo:

1. **Eficiência**: Clientes podem requisitar apenas os dados necessários
2. **Performance**: Redução de carga no servidor e tráfego de rede
3. **Escalabilidade**: Sistema pronto para grandes volumes de dados
4. **Consistência**: Mesmo padrão em HTTP e gRPC
5. **Segurança**: Limites aplicados para evitar abusos

Próximos passos recomendados:
- [ ] Criar testes de integração para os endpoints gRPC
- [ ] Adicionar métricas de performance
- [ ] Documentar exemplos de uso em outras linguagens
