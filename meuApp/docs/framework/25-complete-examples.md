# 25. Complete Examples - Exemplos Completos de Aplicações

## 📋 Índice

1. [E-Commerce Completo](#e-commerce-completo)
2. [Blog Platform](#blog-platform)
3. [Sistema de Autenticação](#sistema-de-autenticação)
4. [Sistema de Pagamentos](#sistema-de-pagamentos)
5. [Workflow com Saga Pattern](#workflow-com-saga-pattern)
6. [Sistema de Notificações](#sistema-de-notificações)

---

## 🛒 E-Commerce Completo

Exemplo de um e-commerce completo com produtos, pedidos, estoque e pagamentos.

### Estrutura de Módulos

```
internal/modules/
├── catalog/              # Gerenciamento de produtos
├── inventory/            # Controle de estoque
├── order/                # Processamento de pedidos
├── payment/              # Integração de pagamentos
├── shipping/             # Cálculo e rastreamento
└── customer/             # Perfil e histórico
```

---

### 1. Módulo de Catálogo (Catalog)

#### Domain Entity

```go
// internal/modules/catalog/domain/entities/product.go
package entities

import (
    "errors"
    "time"
    
    "github.com/google/uuid"
)

type Product struct {
    ID          string
    SKU         string
    Name        string
    Description string
    Category    *Category
    Price       Money
    Images      []ProductImage
    Attributes  map[string]string
    Status      ProductStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type ProductStatus string

const (
    ProductStatusActive   ProductStatus = "active"
    ProductStatusInactive ProductStatus = "inactive"
    ProductStatusDraft    ProductStatus = "draft"
)

type Money struct {
    Amount   float64
    Currency string
}

type ProductImage struct {
    URL       string
    IsPrimary bool
    Order     int
}

type Category struct {
    ID       string
    Name     string
    ParentID *string
}

// NewProduct - Factory method com validações
func NewProduct(sku, name string, price Money, category *Category) (*Product, error) {
    if sku == "" {
        return nil, errors.New("SKU is required")
    }
    if name == "" {
        return nil, errors.New("name is required")
    }
    if price.Amount <= 0 {
        return nil, errors.New("price must be greater than zero")
    }
    if category == nil {
        return nil, errors.New("category is required")
    }
    
    return &Product{
        ID:         uuid.New().String(),
        SKU:        sku,
        Name:       name,
        Price:      price,
        Category:   category,
        Status:     ProductStatusDraft,
        Attributes: make(map[string]string),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}

// Domain methods
func (p *Product) Activate() error {
    if p.Status == ProductStatusActive {
        return errors.New("product is already active")
    }
    p.Status = ProductStatusActive
    p.UpdatedAt = time.Now()
    return nil
}

func (p *Product) Deactivate() error {
    if p.Status == ProductStatusInactive {
        return errors.New("product is already inactive")
    }
    p.Status = ProductStatusInactive
    p.UpdatedAt = time.Now()
    return nil
}

func (p *Product) UpdatePrice(newPrice Money) error {
    if newPrice.Amount <= 0 {
        return errors.New("price must be greater than zero")
    }
    if newPrice.Currency != p.Price.Currency {
        return errors.New("currency mismatch")
    }
    p.Price = newPrice
    p.UpdatedAt = time.Now()
    return nil
}

func (p *Product) AddImage(url string, isPrimary bool) {
    order := len(p.Images)
    p.Images = append(p.Images, ProductImage{
        URL:       url,
        IsPrimary: isPrimary,
        Order:     order,
    })
}
```

#### Command: Create Product

```go
// internal/modules/catalog/application/commands/create_product.go
package commands

import (
    "context"
    "fmt"
    
    "meuApp/internal/modules/catalog/domain/entities"
    "meuApp/internal/modules/catalog/ports"
    "meuApp/pkg/contracts"
    "meuApp/pkg/errors"
)

type CreateProductCommand struct {
    SKU         string
    Name        string
    Description string
    CategoryID  string
    Price       float64
    Currency    string
    Images      []CreateProductImage
    Attributes  map[string]string
}

type CreateProductImage struct {
    URL       string
    IsPrimary bool
}

type CreateProductHandler struct {
    productRepo  ports.ProductRepository
    categoryRepo ports.CategoryRepository
    eventBus     contracts.EventBus
    logger       contracts.Logger
}

func NewCreateProductHandler(
    productRepo ports.ProductRepository,
    categoryRepo ports.CategoryRepository,
    eventBus contracts.EventBus,
    logger contracts.Logger,
) *CreateProductHandler {
    return &CreateProductHandler{
        productRepo:  productRepo,
        categoryRepo: categoryRepo,
        eventBus:     eventBus,
        logger:       logger,
    }
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd *CreateProductCommand) (*entities.Product, error) {
    // 1. Validar comando
    if err := h.validate(cmd); err != nil {
        return nil, err
    }
    
    // 2. Verificar se SKU já existe
    existing, err := h.productRepo.FindBySKU(ctx, cmd.SKU)
    if err == nil && existing != nil {
        return nil, errors.NewConflictError("Product with SKU already exists", map[string]interface{}{
            "sku": cmd.SKU,
        })
    }
    
    // 3. Buscar categoria
    category, err := h.categoryRepo.FindByID(ctx, cmd.CategoryID)
    if err != nil {
        return nil, errors.NewNotFoundError("Category not found", map[string]interface{}{
            "category_id": cmd.CategoryID,
        })
    }
    
    // 4. Criar produto
    product, err := entities.NewProduct(
        cmd.SKU,
        cmd.Name,
        entities.Money{Amount: cmd.Price, Currency: cmd.Currency},
        category,
    )
    if err != nil {
        return nil, errors.NewValidationError("Invalid product data", map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    product.Description = cmd.Description
    product.Attributes = cmd.Attributes
    
    // 5. Adicionar imagens
    for _, img := range cmd.Images {
        product.AddImage(img.URL, img.IsPrimary)
    }
    
    // 6. Salvar no repositório
    if err := h.productRepo.Create(ctx, product); err != nil {
        h.logger.Error("Failed to create product", map[string]interface{}{
            "error": err.Error(),
            "sku":   cmd.SKU,
        })
        return nil, errors.NewInternalError("Failed to create product", nil)
    }
    
    // 7. Publicar evento
    h.eventBus.Publish(ctx, &ProductCreatedEvent{
        ProductID:   product.ID,
        SKU:         product.SKU,
        Name:        product.Name,
        Price:       product.Price.Amount,
        CategoryID:  product.Category.ID,
        OccurredAt:  product.CreatedAt,
    })
    
    h.logger.Info("Product created successfully", map[string]interface{}{
        "product_id": product.ID,
        "sku":        product.SKU,
    })
    
    return product, nil
}

func (h *CreateProductHandler) validate(cmd *CreateProductCommand) error {
    if cmd.SKU == "" {
        return errors.NewValidationError("SKU is required", nil)
    }
    if cmd.Name == "" {
        return errors.NewValidationError("Name is required", nil)
    }
    if cmd.Price <= 0 {
        return errors.NewValidationError("Price must be greater than zero", nil)
    }
    if cmd.CategoryID == "" {
        return errors.NewValidationError("Category ID is required", nil)
    }
    return nil
}
```

---

### 2. Módulo de Pedidos (Order)

#### Domain Aggregate

```go
// internal/modules/order/domain/aggregates/order.go
package aggregates

import (
    "errors"
    "time"
    
    "github.com/google/uuid"
)

type Order struct {
    ID              string
    CustomerID      string
    Items           []OrderItem
    ShippingAddress Address
    BillingAddress  Address
    PaymentMethod   PaymentMethod
    Status          OrderStatus
    Subtotal        float64
    ShippingCost    float64
    Tax             float64
    Discount        float64
    Total           float64
    Notes           string
    CreatedAt       time.Time
    UpdatedAt       time.Time
    
    // Domain events (não persistidos)
    events []interface{}
}

type OrderItem struct {
    ProductID   string
    ProductSKU  string
    ProductName string
    Quantity    int
    UnitPrice   float64
    Subtotal    float64
}

type Address struct {
    Street     string
    Number     string
    Complement string
    City       string
    State      string
    ZipCode    string
    Country    string
}

type PaymentMethod struct {
    Type        string // credit_card, debit_card, pix, boleto
    CardLast4   string
    CardBrand   string
}

type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusConfirmed OrderStatus = "confirmed"
    OrderStatusPaid      OrderStatus = "paid"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusDelivered OrderStatus = "delivered"
    OrderStatusCancelled OrderStatus = "cancelled"
)

// Factory
func NewOrder(customerID string, shippingAddress Address) (*Order, error) {
    if customerID == "" {
        return nil, errors.New("customer ID is required")
    }
    
    order := &Order{
        ID:              uuid.New().String(),
        CustomerID:      customerID,
        Items:           make([]OrderItem, 0),
        ShippingAddress: shippingAddress,
        Status:          OrderStatusPending,
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
        events:          make([]interface{}, 0),
    }
    
    return order, nil
}

// Domain methods
func (o *Order) AddItem(productID, sku, name string, quantity int, unitPrice float64) error {
    if quantity <= 0 {
        return errors.New("quantity must be greater than zero")
    }
    if unitPrice <= 0 {
        return errors.New("unit price must be greater than zero")
    }
    
    // Verificar se produto já está no pedido
    for i, item := range o.Items {
        if item.ProductID == productID {
            // Atualizar quantidade
            o.Items[i].Quantity += quantity
            o.Items[i].Subtotal = float64(o.Items[i].Quantity) * o.Items[i].UnitPrice
            o.recalculateTotal()
            return nil
        }
    }
    
    // Adicionar novo item
    item := OrderItem{
        ProductID:   productID,
        ProductSKU:  sku,
        ProductName: name,
        Quantity:    quantity,
        UnitPrice:   unitPrice,
        Subtotal:    float64(quantity) * unitPrice,
    }
    
    o.Items = append(o.Items, item)
    o.recalculateTotal()
    o.UpdatedAt = time.Now()
    
    return nil
}

func (o *Order) RemoveItem(productID string) error {
    for i, item := range o.Items {
        if item.ProductID == productID {
            o.Items = append(o.Items[:i], o.Items[i+1:]...)
            o.recalculateTotal()
            o.UpdatedAt = time.Now()
            return nil
        }
    }
    return errors.New("item not found in order")
}

func (o *Order) SetShippingCost(cost float64) error {
    if cost < 0 {
        return errors.New("shipping cost cannot be negative")
    }
    o.ShippingCost = cost
    o.recalculateTotal()
    o.UpdatedAt = time.Now()
    return nil
}

func (o *Order) ApplyDiscount(discount float64) error {
    if discount < 0 {
        return errors.New("discount cannot be negative")
    }
    if discount > o.Subtotal {
        return errors.New("discount cannot be greater than subtotal")
    }
    o.Discount = discount
    o.recalculateTotal()
    o.UpdatedAt = time.Now()
    return nil
}

func (o *Order) Confirm() error {
    if o.Status != OrderStatusPending {
        return errors.New("only pending orders can be confirmed")
    }
    if len(o.Items) == 0 {
        return errors.New("order must have at least one item")
    }
    
    o.Status = OrderStatusConfirmed
    o.UpdatedAt = time.Now()
    
    // Adicionar evento de domínio
    o.addEvent(OrderConfirmedEvent{
        OrderID:    o.ID,
        CustomerID: o.CustomerID,
        Total:      o.Total,
        Items:      o.Items,
        OccurredAt: time.Now(),
    })
    
    return nil
}

func (o *Order) MarkAsPaid(paymentMethod PaymentMethod) error {
    if o.Status != OrderStatusConfirmed {
        return errors.New("only confirmed orders can be marked as paid")
    }
    
    o.Status = OrderStatusPaid
    o.PaymentMethod = paymentMethod
    o.UpdatedAt = time.Now()
    
    o.addEvent(OrderPaidEvent{
        OrderID:       o.ID,
        CustomerID:    o.CustomerID,
        Total:         o.Total,
        PaymentMethod: paymentMethod.Type,
        OccurredAt:    time.Now(),
    })
    
    return nil
}

func (o *Order) Ship(trackingCode string) error {
    if o.Status != OrderStatusPaid {
        return errors.New("only paid orders can be shipped")
    }
    
    o.Status = OrderStatusShipped
    o.UpdatedAt = time.Now()
    
    o.addEvent(OrderShippedEvent{
        OrderID:      o.ID,
        CustomerID:   o.CustomerID,
        TrackingCode: trackingCode,
        OccurredAt:   time.Now(),
    })
    
    return nil
}

func (o *Order) Deliver() error {
    if o.Status != OrderStatusShipped {
        return errors.New("only shipped orders can be delivered")
    }
    
    o.Status = OrderStatusDelivered
    o.UpdatedAt = time.Now()
    
    o.addEvent(OrderDeliveredEvent{
        OrderID:    o.ID,
        CustomerID: o.CustomerID,
        OccurredAt: time.Now(),
    })
    
    return nil
}

func (o *Order) Cancel(reason string) error {
    if o.Status == OrderStatusShipped || o.Status == OrderStatusDelivered {
        return errors.New("shipped or delivered orders cannot be cancelled")
    }
    if o.Status == OrderStatusCancelled {
        return errors.New("order is already cancelled")
    }
    
    o.Status = OrderStatusCancelled
    o.Notes = "Cancelled: " + reason
    o.UpdatedAt = time.Now()
    
    o.addEvent(OrderCancelledEvent{
        OrderID:    o.ID,
        CustomerID: o.CustomerID,
        Reason:     reason,
        OccurredAt: time.Now(),
    })
    
    return nil
}

// Private helpers
func (o *Order) recalculateTotal() {
    o.Subtotal = 0
    for _, item := range o.Items {
        o.Subtotal += item.Subtotal
    }
    
    // Total = Subtotal + Shipping + Tax - Discount
    o.Total = o.Subtotal + o.ShippingCost + o.Tax - o.Discount
}

func (o *Order) addEvent(event interface{}) {
    o.events = append(o.events, event)
}

func (o *Order) GetEvents() []interface{} {
    return o.events
}

func (o *Order) ClearEvents() {
    o.events = make([]interface{}, 0)
}
```

#### Command: Create Order

```go
// internal/modules/order/application/commands/create_order.go
package commands

import (
    "context"
    
    "meuApp/internal/modules/order/domain/aggregates"
    "meuApp/internal/modules/order/ports"
    "meuApp/pkg/contracts"
    "meuApp/pkg/errors"
)

type CreateOrderCommand struct {
    CustomerID      string
    Items           []CreateOrderItem
    ShippingAddress CreateOrderAddress
}

type CreateOrderItem struct {
    ProductID string
    Quantity  int
}

type CreateOrderAddress struct {
    Street     string
    Number     string
    Complement string
    City       string
    State      string
    ZipCode    string
    Country    string
}

type CreateOrderHandler struct {
    orderRepo    ports.OrderRepository
    productRepo  ports.ProductRepository
    inventoryRepo ports.InventoryRepository
    eventBus     contracts.EventBus
    logger       contracts.Logger
}

func NewCreateOrderHandler(
    orderRepo ports.OrderRepository,
    productRepo ports.ProductRepository,
    inventoryRepo ports.InventoryRepository,
    eventBus contracts.EventBus,
    logger contracts.Logger,
) *CreateOrderHandler {
    return &CreateOrderHandler{
        orderRepo:     orderRepo,
        productRepo:   productRepo,
        inventoryRepo: inventoryRepo,
        eventBus:      eventBus,
        logger:        logger,
    }
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd *CreateOrderCommand) (*aggregates.Order, error) {
    // 1. Validar comando
    if err := h.validate(cmd); err != nil {
        return nil, err
    }
    
    // 2. Criar aggregate
    address := aggregates.Address{
        Street:     cmd.ShippingAddress.Street,
        Number:     cmd.ShippingAddress.Number,
        Complement: cmd.ShippingAddress.Complement,
        City:       cmd.ShippingAddress.City,
        State:      cmd.ShippingAddress.State,
        ZipCode:    cmd.ShippingAddress.ZipCode,
        Country:    cmd.ShippingAddress.Country,
    }
    
    order, err := aggregates.NewOrder(cmd.CustomerID, address)
    if err != nil {
        return nil, errors.NewValidationError("Invalid order data", map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    // 3. Adicionar itens
    for _, item := range cmd.Items {
        // Buscar produto
        product, err := h.productRepo.FindByID(ctx, item.ProductID)
        if err != nil {
            return nil, errors.NewNotFoundError("Product not found", map[string]interface{}{
                "product_id": item.ProductID,
            })
        }
        
        // Verificar estoque
        hasStock, err := h.inventoryRepo.CheckAvailability(ctx, item.ProductID, item.Quantity)
        if err != nil {
            return nil, errors.NewInternalError("Failed to check inventory", nil)
        }
        if !hasStock {
            return nil, errors.NewConflictError("Insufficient stock", map[string]interface{}{
                "product_id": item.ProductID,
                "product":    product.Name,
                "requested":  item.Quantity,
            })
        }
        
        // Adicionar item ao pedido
        err = order.AddItem(
            product.ID,
            product.SKU,
            product.Name,
            item.Quantity,
            product.Price.Amount,
        )
        if err != nil {
            return nil, errors.NewValidationError("Failed to add item", map[string]interface{}{
                "error":      err.Error(),
                "product_id": item.ProductID,
            })
        }
    }
    
    // 4. Calcular frete (exemplo simplificado)
    shippingCost := h.calculateShipping(order)
    order.SetShippingCost(shippingCost)
    
    // 5. Salvar pedido
    if err := h.orderRepo.Create(ctx, order); err != nil {
        h.logger.Error("Failed to create order", map[string]interface{}{
            "error":       err.Error(),
            "customer_id": cmd.CustomerID,
        })
        return nil, errors.NewInternalError("Failed to create order", nil)
    }
    
    // 6. Publicar evento
    h.eventBus.Publish(ctx, &OrderCreatedEvent{
        OrderID:    order.ID,
        CustomerID: order.CustomerID,
        Total:      order.Total,
        Items:      order.Items,
        OccurredAt: order.CreatedAt,
    })
    
    h.logger.Info("Order created successfully", map[string]interface{}{
        "order_id":    order.ID,
        "customer_id": cmd.CustomerID,
        "total":       order.Total,
    })
    
    return order, nil
}

func (h *CreateOrderHandler) validate(cmd *CreateOrderCommand) error {
    if cmd.CustomerID == "" {
        return errors.NewValidationError("Customer ID is required", nil)
    }
    if len(cmd.Items) == 0 {
        return errors.NewValidationError("Order must have at least one item", nil)
    }
    if cmd.ShippingAddress.ZipCode == "" {
        return errors.NewValidationError("Shipping address is required", nil)
    }
    return nil
}

func (h *CreateOrderHandler) calculateShipping(order *aggregates.Order) float64 {
    // Lógica simplificada - na prática, integraria com API de correios
    baseRate := 10.0
    itemCount := len(order.Items)
    return baseRate + float64(itemCount)*2.0
}
```

---

### 3. Saga Pattern: Checkout Completo

```go
// internal/modules/checkout/sagas/checkout_saga.go
package sagas

import (
    "context"
    "fmt"
    
    "meuApp/pkg/contracts"
)

type CheckoutSaga struct {
    orderRepo     OrderRepository
    inventoryRepo InventoryRepository
    paymentSvc    PaymentService
    emailSvc      EmailService
    eventBus      contracts.EventBus
    logger        contracts.Logger
}

func NewCheckoutSaga(
    orderRepo OrderRepository,
    inventoryRepo InventoryRepository,
    paymentSvc PaymentService,
    emailSvc EmailService,
    eventBus contracts.EventBus,
    logger contracts.Logger,
) *CheckoutSaga {
    return &CheckoutSaga{
        orderRepo:     orderRepo,
        inventoryRepo: inventoryRepo,
        paymentSvc:    paymentSvc,
        emailSvc:      emailSvc,
        eventBus:      eventBus,
        logger:        logger,
    }
}

// Subscribe registra os handlers de eventos
func (s *CheckoutSaga) Subscribe() {
    s.eventBus.Subscribe("order.created", s.handleOrderCreated)
    s.eventBus.Subscribe("payment.completed", s.handlePaymentCompleted)
    s.eventBus.Subscribe("payment.failed", s.handlePaymentFailed)
    s.eventBus.Subscribe("inventory.reserved", s.handleInventoryReserved)
    s.eventBus.Subscribe("inventory.failed", s.handleInventoryFailed)
}

// Step 1: Order Created → Reserve Inventory
func (s *CheckoutSaga) handleOrderCreated(ctx context.Context, event interface{}) error {
    e := event.(*OrderCreatedEvent)
    
    s.logger.Info("Starting checkout saga", map[string]interface{}{
        "order_id": e.OrderID,
    })
    
    // Reservar estoque para cada item
    for _, item := range e.Items {
        err := s.inventoryRepo.Reserve(ctx, item.ProductID, item.Quantity, e.OrderID)
        if err != nil {
            // Falhou - publicar evento de falha
            s.eventBus.Publish(ctx, &InventoryReservationFailedEvent{
                OrderID:   e.OrderID,
                ProductID: item.ProductID,
                Reason:    err.Error(),
            })
            
            // Compensar - cancelar pedido
            return s.compensateOrderCreation(ctx, e.OrderID, "inventory reservation failed")
        }
    }
    
    // Sucesso - publicar evento
    s.eventBus.Publish(ctx, &InventoryReservedEvent{
        OrderID: e.OrderID,
        Items:   e.Items,
    })
    
    return nil
}

// Step 2: Inventory Reserved → Process Payment
func (s *CheckoutSaga) handleInventoryReserved(ctx context.Context, event interface{}) error {
    e := event.(*InventoryReservedEvent)
    
    // Buscar pedido
    order, err := s.orderRepo.FindByID(ctx, e.OrderID)
    if err != nil {
        return err
    }
    
    // Processar pagamento
    paymentResult, err := s.paymentSvc.ProcessPayment(ctx, PaymentRequest{
        OrderID:       order.ID,
        Amount:        order.Total,
        Currency:      "BRL",
        PaymentMethod: order.PaymentMethod,
    })
    
    if err != nil {
        // Falhou - publicar evento de falha
        s.eventBus.Publish(ctx, &PaymentFailedEvent{
            OrderID: e.OrderID,
            Reason:  err.Error(),
        })
        
        // Compensar - liberar estoque
        return s.compensateInventoryReservation(ctx, e.OrderID)
    }
    
    // Sucesso - publicar evento
    s.eventBus.Publish(ctx, &PaymentCompletedEvent{
        OrderID:       e.OrderID,
        TransactionID: paymentResult.TransactionID,
        Amount:        paymentResult.Amount,
    })
    
    return nil
}

// Step 3: Payment Completed → Update Order & Send Email
func (s *CheckoutSaga) handlePaymentCompleted(ctx context.Context, event interface{}) error {
    e := event.(*PaymentCompletedEvent)
    
    // Buscar pedido
    order, err := s.orderRepo.FindByID(ctx, e.OrderID)
    if err != nil {
        return err
    }
    
    // Atualizar status do pedido
    err = order.MarkAsPaid(order.PaymentMethod)
    if err != nil {
        return err
    }
    
    // Salvar
    if err := s.orderRepo.Update(ctx, order); err != nil {
        return err
    }
    
    // Enviar email de confirmação (não bloqueia se falhar)
    go func() {
        err := s.emailSvc.SendOrderConfirmation(context.Background(), order.CustomerID, order.ID)
        if err != nil {
            s.logger.Error("Failed to send confirmation email", map[string]interface{}{
                "order_id": order.ID,
                "error":    err.Error(),
            })
        }
    }()
    
    s.logger.Info("Checkout completed successfully", map[string]interface{}{
        "order_id":       order.ID,
        "transaction_id": e.TransactionID,
    })
    
    return nil
}

// Compensations (Rollback)

func (s *CheckoutSaga) handlePaymentFailed(ctx context.Context, event interface{}) error {
    e := event.(*PaymentFailedEvent)
    
    s.logger.Warning("Payment failed, starting compensation", map[string]interface{}{
        "order_id": e.OrderID,
        "reason":   e.Reason,
    })
    
    // Liberar estoque
    if err := s.compensateInventoryReservation(ctx, e.OrderID); err != nil {
        s.logger.Error("Failed to compensate inventory", map[string]interface{}{
            "order_id": e.OrderID,
            "error":    err.Error(),
        })
    }
    
    // Cancelar pedido
    return s.compensateOrderCreation(ctx, e.OrderID, "payment failed: "+e.Reason)
}

func (s *CheckoutSaga) handleInventoryFailed(ctx context.Context, event interface{}) error {
    e := event.(*InventoryReservationFailedEvent)
    
    s.logger.Warning("Inventory reservation failed", map[string]interface{}{
        "order_id":   e.OrderID,
        "product_id": e.ProductID,
        "reason":     e.Reason,
    })
    
    // Cancelar pedido
    return s.compensateOrderCreation(ctx, e.OrderID, "insufficient stock")
}

func (s *CheckoutSaga) compensateInventoryReservation(ctx context.Context, orderID string) error {
    order, err := s.orderRepo.FindByID(ctx, orderID)
    if err != nil {
        return err
    }
    
    // Liberar estoque de todos os itens
    for _, item := range order.Items {
        err := s.inventoryRepo.Release(ctx, item.ProductID, item.Quantity, orderID)
        if err != nil {
            s.logger.Error("Failed to release inventory", map[string]interface{}{
                "order_id":   orderID,
                "product_id": item.ProductID,
                "error":      err.Error(),
            })
        }
    }
    
    return nil
}

func (s *CheckoutSaga) compensateOrderCreation(ctx context.Context, orderID, reason string) error {
    order, err := s.orderRepo.FindByID(ctx, orderID)
    if err != nil {
        return err
    }
    
    // Cancelar pedido
    if err := order.Cancel(reason); err != nil {
        return err
    }
    
    // Salvar
    if err := s.orderRepo.Update(ctx, order); err != nil {
        return err
    }
    
    // Publicar evento
    s.eventBus.Publish(ctx, &OrderCancelledEvent{
        OrderID: orderID,
        Reason:  reason,
    })
    
    s.logger.Info("Order cancelled due to saga compensation", map[string]interface{}{
        "order_id": orderID,
        "reason":   reason,
    })
    
    return nil
}
```

#### Diagrama de Fluxo da Saga

```
┌─────────────┐
│Order Created│
└──────┬──────┘
       │
       ▼
┌──────────────────┐      ❌ Fail
│Reserve Inventory │──────────┐
└──────┬───────────┘          │
       │ ✅ Success           │
       ▼                      │
┌──────────────────┐      ❌ Fail
│Process Payment   │──────────┤
└──────┬───────────┘          │
       │ ✅ Success           │
       ▼                      │
┌──────────────────┐          │
│Update Order      │          │
│Status → Paid     │          │
└──────┬───────────┘          │
       │                      │
       ▼                      ▼
┌──────────────────┐   ┌───────────────┐
│Send Confirmation │   │  Compensation │
│Email             │   │  (Rollback)   │
└──────────────────┘   └───────────────┘
                              │
                              ▼
                       ┌──────────────┐
                       │Release Stock │
                       └──────┬───────┘
                              │
                              ▼
                       ┌──────────────┐
                       │Cancel Order  │
                       └──────────────┘
```

---

## 📝 Blog Platform

Sistema completo de blog com posts, comentários e tags.

### Domain Entity: Blog Post

```go
// internal/modules/blog/domain/entities/post.go
package entities

import (
    "errors"
    "strings"
    "time"
    
    "github.com/google/uuid"
)

type Post struct {
    ID            string
    Title         string
    Slug          string
    Content       string
    Excerpt       string
    AuthorID      string
    CategoryID    string
    Tags          []string
    Status        PostStatus
    ViewCount     int
    PublishedAt   *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
    
    // Relationships (não persistido diretamente)
    Author   *Author
    Category *Category
    Comments []Comment
}

type PostStatus string

const (
    PostStatusDraft     PostStatus = "draft"
    PostStatusPublished PostStatus = "published"
    PostStatusArchived  PostStatus = "archived"
)

type Author struct {
    ID     string
    Name   string
    Email  string
    Bio    string
    Avatar string
}

type Category struct {
    ID   string
    Name string
    Slug string
}

type Comment struct {
    ID        string
    PostID    string
    AuthorID  string
    Content   string
    IsApproved bool
    CreatedAt time.Time
}

// Factory
func NewPost(title, content, authorID, categoryID string) (*Post, error) {
    if title == "" {
        return nil, errors.New("title is required")
    }
    if content == "" {
        return nil, errors.New("content is required")
    }
    if authorID == "" {
        return nil, errors.New("author ID is required")
    }
    
    slug := generateSlug(title)
    excerpt := generateExcerpt(content, 200)
    
    return &Post{
        ID:         uuid.New().String(),
        Title:      title,
        Slug:       slug,
        Content:    content,
        Excerpt:    excerpt,
        AuthorID:   authorID,
        CategoryID: categoryID,
        Tags:       make([]string, 0),
        Status:     PostStatusDraft,
        ViewCount:  0,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}

// Domain methods
func (p *Post) Publish() error {
    if p.Status == PostStatusPublished {
        return errors.New("post is already published")
    }
    
    now := time.Now()
    p.Status = PostStatusPublished
    p.PublishedAt = &now
    p.UpdatedAt = now
    
    return nil
}

func (p *Post) Archive() error {
    if p.Status == PostStatusArchived {
        return errors.New("post is already archived")
    }
    
    p.Status = PostStatusArchived
    p.UpdatedAt = time.Now()
    
    return nil
}

func (p *Post) AddTag(tag string) {
    tag = strings.ToLower(strings.TrimSpace(tag))
    if tag == "" {
        return
    }
    
    // Verificar se já existe
    for _, t := range p.Tags {
        if t == tag {
            return
        }
    }
    
    p.Tags = append(p.Tags, tag)
    p.UpdatedAt = time.Now()
}

func (p *Post) RemoveTag(tag string) {
    for i, t := range p.Tags {
        if t == tag {
            p.Tags = append(p.Tags[:i], p.Tags[i+1:]...)
            p.UpdatedAt = time.Now()
            return
        }
    }
}

func (p *Post) IncrementView() {
    p.ViewCount++
}

func (p *Post) UpdateContent(title, content string) error {
    if title == "" {
        return errors.New("title is required")
    }
    if content == "" {
        return errors.New("content is required")
    }
    
    p.Title = title
    p.Content = content
    p.Excerpt = generateExcerpt(content, 200)
    p.Slug = generateSlug(title)
    p.UpdatedAt = time.Now()
    
    return nil
}

// Helpers
func generateSlug(title string) string {
    slug := strings.ToLower(title)
    slug = strings.ReplaceAll(slug, " ", "-")
    // Remove caracteres especiais (simplificado)
    return slug
}

func generateExcerpt(content string, maxLength int) string {
    if len(content) <= maxLength {
        return content
    }
    return content[:maxLength] + "..."
}
```

### Query: List Posts with Filters

```go
// internal/modules/blog/application/queries/list_posts.go
package queries

import (
    "context"
    
    "meuApp/internal/modules/blog/domain/entities"
    "meuApp/internal/modules/blog/ports"
)

type ListPostsQuery struct {
    Status     *entities.PostStatus
    CategoryID *string
    AuthorID   *string
    Tags       []string
    Search     *string
    SortBy     string // "created_at", "published_at", "view_count"
    SortOrder  string // "asc", "desc"
    Page       int
    PageSize   int
}

type ListPostsResult struct {
    Posts      []*entities.Post
    TotalItems int
    Page       int
    PageSize   int
    TotalPages int
}

type ListPostsHandler struct {
    postRepo ports.PostRepository
    logger   contracts.Logger
}

func NewListPostsHandler(postRepo ports.PostRepository, logger contracts.Logger) *ListPostsHandler {
    return &ListPostsHandler{
        postRepo: postRepo,
        logger:   logger,
    }
}

func (h *ListPostsHandler) Handle(ctx context.Context, query *ListPostsQuery) (*ListPostsResult, error) {
    // Aplicar defaults
    if query.Page <= 0 {
        query.Page = 1
    }
    if query.PageSize <= 0 {
        query.PageSize = 10
    }
    if query.SortBy == "" {
        query.SortBy = "created_at"
    }
    if query.SortOrder == "" {
        query.SortOrder = "desc"
    }
    
    // Buscar posts
    posts, total, err := h.postRepo.List(ctx, ports.PostFilters{
        Status:     query.Status,
        CategoryID: query.CategoryID,
        AuthorID:   query.AuthorID,
        Tags:       query.Tags,
        Search:     query.Search,
        SortBy:     query.SortBy,
        SortOrder:  query.SortOrder,
        Offset:     (query.Page - 1) * query.PageSize,
        Limit:      query.PageSize,
    })
    
    if err != nil {
        h.logger.Error("Failed to list posts", map[string]interface{}{
            "error": err.Error(),
        })
        return nil, err
    }
    
    totalPages := (total + query.PageSize - 1) / query.PageSize
    
    return &ListPostsResult{
        Posts:      posts,
        TotalItems: total,
        Page:       query.Page,
        PageSize:   query.PageSize,
        TotalPages: totalPages,
    }, nil
}
```

---

## 🔐 Sistema de Autenticação

### JWT Token Service

```go
// internal/modules/auth/services/jwt_service.go
package services

import (
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
    secretKey       string
    accessTokenTTL  time.Duration
    refreshTokenTTL time.Duration
}

type TokenPair struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int64
}

type Claims struct {
    UserID string   `json:"user_id"`
    Email  string   `json:"email"`
    Roles  []string `json:"roles"`
    jwt.RegisteredClaims
}

func NewJWTService(secretKey string) *JWTService {
    return &JWTService{
        secretKey:       secretKey,
        accessTokenTTL:  15 * time.Minute,
        refreshTokenTTL: 7 * 24 * time.Hour,
    }
}

func (s *JWTService) GenerateTokenPair(userID, email string, roles []string) (*TokenPair, error) {
    // Generate access token
    accessToken, err := s.generateAccessToken(userID, email, roles)
    if err != nil {
        return nil, err
    }
    
    // Generate refresh token
    refreshToken, err := s.generateRefreshToken(userID)
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
    }, nil
}

func (s *JWTService) generateAccessToken(userID, email string, roles []string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Roles:  roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) generateRefreshToken(userID string) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   userID,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenTTL)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.secretKey), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    return claims, nil
}

func (s *JWTService) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.secretKey), nil
    })
    
    if err != nil || !token.Valid {
        return nil, errors.New("invalid refresh token")
    }
    
    claims := token.Claims.(jwt.MapClaims)
    userID := claims["sub"].(string)
    
    // Buscar usuário e gerar novo token
    // (aqui você buscaria o usuário do banco e suas roles)
    
    return s.GenerateTokenPair(userID, "", []string{"user"})
}
```

### HTTP Middleware

```go
// pkg/adapters/http/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Missing authorization header",
            })
            c.Abort()
            return
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization header format",
            })
            c.Abort()
            return
        }
        
        token := parts[1]
        claims, err := jwtService.ValidateAccessToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }
        
        // Adicionar claims ao contexto
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("roles", claims.Roles)
        
        c.Next()
    }
}

func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRoles, exists := c.Get("roles")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "No roles found",
            })
            c.Abort()
            return
        }
        
        userRolesList := userRoles.([]string)
        hasRole := false
        
        for _, requiredRole := range roles {
            for _, userRole := range userRolesList {
                if userRole == requiredRole {
                    hasRole = true
                    break
                }
            }
        }
        
        if !hasRole {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Insufficient permissions",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

## 💳 Sistema de Pagamentos

### Payment Gateway Integration

```go
// internal/modules/payment/services/payment_gateway.go
package services

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "time"
)

type PaymentGateway interface {
    ProcessPayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error)
    RefundPayment(ctx context.Context, transactionID string, amount float64) error
    GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatus, error)
}

type PaymentRequest struct {
    OrderID       string
    Amount        float64
    Currency      string
    PaymentMethod PaymentMethodDetails
    Customer      CustomerDetails
}

type PaymentMethodDetails struct {
    Type      string // credit_card, debit_card, pix, boleto
    CardToken string
    CardLast4 string
    CardBrand string
}

type CustomerDetails struct {
    ID    string
    Name  string
    Email string
    CPF   string
}

type PaymentResult struct {
    TransactionID string
    Status        string
    Amount        float64
    Currency      string
    ProcessedAt   time.Time
    GatewayData   map[string]interface{}
}

type PaymentStatus struct {
    Status      string
    ProcessedAt time.Time
    RefundedAt  *time.Time
}

// Stripe Implementation
type StripeGateway struct {
    apiKey     string
    apiURL     string
    httpClient *http.Client
}

func NewStripeGateway(apiKey string) *StripeGateway {
    return &StripeGateway{
        apiKey: apiKey,
        apiURL: "https://api.stripe.com/v1",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (g *StripeGateway) ProcessPayment(ctx context.Context, req PaymentRequest) (*PaymentResult, error) {
    // Implementação simplificada
    // Na prática, faria chamada real à API do Stripe
    
    payload := map[string]interface{}{
        "amount":               int(req.Amount * 100), // Centavos
        "currency":             req.Currency,
        "source":               req.PaymentMethod.CardToken,
        "description":          "Order " + req.OrderID,
        "metadata": map[string]string{
            "order_id": req.OrderID,
        },
    }
    
    // Simular chamada HTTP
    result, err := g.makeAPICall(ctx, "POST", "/charges", payload)
    if err != nil {
        return nil, err
    }
    
    return &PaymentResult{
        TransactionID: result["id"].(string),
        Status:        result["status"].(string),
        Amount:        req.Amount,
        Currency:      req.Currency,
        ProcessedAt:   time.Now(),
        GatewayData:   result,
    }, nil
}

func (g *StripeGateway) RefundPayment(ctx context.Context, transactionID string, amount float64) error {
    payload := map[string]interface{}{
        "charge": transactionID,
        "amount": int(amount * 100),
    }
    
    _, err := g.makeAPICall(ctx, "POST", "/refunds", payload)
    return err
}

func (g *StripeGateway) GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatus, error) {
    result, err := g.makeAPICall(ctx, "GET", "/charges/"+transactionID, nil)
    if err != nil {
        return nil, err
    }
    
    status := &PaymentStatus{
        Status:      result["status"].(string),
        ProcessedAt: time.Now(), // Parse from result
    }
    
    return status, nil
}

func (g *StripeGateway) makeAPICall(ctx context.Context, method, path string, payload interface{}) (map[string]interface{}, error) {
    // Implementação simplificada
    // Na prática, construiria request HTTP real
    return map[string]interface{}{
        "id":     "ch_" + generateID(),
        "status": "succeeded",
    }, nil
}

func generateID() string {
    return "mock_id_123"
}
```

---

## 🔔 Sistema de Notificações

### Multi-Channel Notification System

```go
// internal/modules/notification/services/notification_service.go
package services

import (
    "context"
    "fmt"
)

type NotificationChannel string

const (
    ChannelEmail    NotificationChannel = "email"
    ChannelSMS      NotificationChannel = "sms"
    ChannelPush     NotificationChannel = "push"
    ChannelWebhook  NotificationChannel = "webhook"
)

type Notification struct {
    ID          string
    RecipientID string
    Channel     NotificationChannel
    Template    string
    Data        map[string]interface{}
    Priority    Priority
}

type Priority string

const (
    PriorityLow      Priority = "low"
    PriorityMedium   Priority = "medium"
    PriorityHigh     Priority = "high"
    PriorityCritical Priority = "critical"
)

type NotificationService struct {
    emailSender   EmailSender
    smsSender     SMSSender
    pushSender    PushSender
    webhookSender WebhookSender
    queue         NotificationQueue
}

func NewNotificationService(
    emailSender EmailSender,
    smsSender SMSSender,
    pushSender PushSender,
    webhookSender WebhookSender,
    queue NotificationQueue,
) *NotificationService {
    return &NotificationService{
        emailSender:   emailSender,
        smsSender:     smsSender,
        pushSender:    pushSender,
        webhookSender: webhookSender,
        queue:         queue,
    }
}

func (s *NotificationService) Send(ctx context.Context, notification *Notification) error {
    // Se for prioridade crítica, envia imediatamente
    if notification.Priority == PriorityCritical {
        return s.sendImmediately(ctx, notification)
    }
    
    // Caso contrário, adiciona na fila
    return s.queue.Enqueue(ctx, notification)
}

func (s *NotificationService) sendImmediately(ctx context.Context, notification *Notification) error {
    switch notification.Channel {
    case ChannelEmail:
        return s.sendEmail(ctx, notification)
    case ChannelSMS:
        return s.sendSMS(ctx, notification)
    case ChannelPush:
        return s.sendPush(ctx, notification)
    case ChannelWebhook:
        return s.sendWebhook(ctx, notification)
    default:
        return fmt.Errorf("unsupported channel: %s", notification.Channel)
    }
}

func (s *NotificationService) sendEmail(ctx context.Context, notification *Notification) error {
    // Renderizar template
    subject := s.renderTemplate(notification.Template+".subject", notification.Data)
    body := s.renderTemplate(notification.Template+".body", notification.Data)
    
    return s.emailSender.Send(ctx, EmailMessage{
        To:      notification.Data["email"].(string),
        Subject: subject,
        Body:    body,
        IsHTML:  true,
    })
}

func (s *NotificationService) sendSMS(ctx context.Context, notification *Notification) error {
    message := s.renderTemplate(notification.Template, notification.Data)
    
    return s.smsSender.Send(ctx, SMSMessage{
        To:      notification.Data["phone"].(string),
        Message: message,
    })
}

func (s *NotificationService) sendPush(ctx context.Context, notification *Notification) error {
    title := s.renderTemplate(notification.Template+".title", notification.Data)
    body := s.renderTemplate(notification.Template+".body", notification.Data)
    
    return s.pushSender.Send(ctx, PushMessage{
        DeviceToken: notification.Data["device_token"].(string),
        Title:       title,
        Body:        body,
        Data:        notification.Data,
    })
}

func (s *NotificationService) sendWebhook(ctx context.Context, notification *Notification) error {
    return s.webhookSender.Send(ctx, WebhookMessage{
        URL:     notification.Data["webhook_url"].(string),
        Payload: notification.Data,
    })
}

func (s *NotificationService) renderTemplate(template string, data map[string]interface{}) string {
    // Implementação simplificada
    // Na prática, usaria um template engine real
    return fmt.Sprintf("Template: %s", template)
}
```

---

## 📚 Próximos Passos

1. **[24-api-reference.md](24-api-reference.md)** - Referência completa de APIs
2. **[26-faq.md](26-faq.md)** - Perguntas frequentes
3. **[21-testing-strategy.md](21-testing-strategy.md)** - Como testar estes exemplos

---

## 💡 Dicas para Usar Estes Exemplos

### 1. Adapte para Seu Contexto

Estes exemplos são **templates** - adapte para suas necessidades:

```go
// ❌ Não copie direto
type Product struct {
    Price Money
}

// ✅ Adapte para seu caso
type YourProduct struct {
    BasePrice    float64
    SalePrice    float64
    PriceHistory []PriceChange
}
```

### 2. Comece Simples

Não implemente tudo de uma vez:

**Fase 1:** Entidades básicas + CRUD  
**Fase 2:** Validações de domínio  
**Fase 3:** Eventos  
**Fase 4:** Saga Pattern  

### 3. Teste Cada Camada

```go
// Teste domain logic
func TestProduct_UpdatePrice(t *testing.T) { ... }

// Teste command handler
func TestCreateProductHandler_Handle(t *testing.T) { ... }

// Teste HTTP handler
func TestProductHandler_CreateProduct(t *testing.T) { ... }
```

### 4. Use Mocks

```go
type mockProductRepo struct {
    mock.Mock
}

func (m *mockProductRepo) Create(ctx context.Context, product *Product) error {
    args := m.Called(ctx, product)
    return args.Error(0)
}
```

---

<div align="center">

**[⬆️ Voltar ao Topo](#25-complete-examples---exemplos-completos-de-aplicações)**

**Documentação criada com ❤️ pela Equipe Artemis**

</div>
