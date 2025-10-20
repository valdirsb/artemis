# 🔍 Implementando Queries - Read Operations

## 📋 Índice
- [O que são Queries?](#o-que-são-queries)
- [Anatomia de uma Query](#anatomia-de-uma-query)
- [Query Handlers](#query-handlers)
- [Paginação](#paginação)
- [Filtros e Ordenação](#filtros-e-ordenação)
- [Performance e Otimização](#performance-e-otimização)
- [Cache](#cache)
- [Exemplos Práticos](#exemplos-práticos)
- [Boas Práticas](#boas-práticas)

---

## 🎯 O que são Queries?

**Queries** são objetos que representam **requisições de leitura** de dados (read operations: SELECT, GET, LIST).

### Query vs Command

```
┌─────────────────────────────────────────────────┐
│              QUERIES (Read Side)                │
├─────────────────────────────────────────────────┤
│ • Apenas leem dados                             │
│ • Retornam DTOs                                 │
│ • Podem usar cache                              │
│ • Podem ter múltiplas representações            │
│ • Otimizadas para performance                   │
│ • Sem efeitos colaterais                        │
│ • Podem usar banco otimizado (read replica)     │
└─────────────────────────────────────────────────┘
```

### Fluxo de uma Query

```
HTTP Request (GET)
        │
        ▼
┌───────────────┐
│  HTTP Adapter │ Parse params
└───────┬───────┘
        │
        ▼
┌───────────────┐
│     Query     │ DTO com filtros
└───────┬───────┘
        │
        ▼
┌───────────────┐
│  Query Bus    │ Despacha query
└───────┬───────┘
        │
        ▼
┌───────────────┐
│ Query Handler │ Busca dados
│               │ 1. Valida params
│               │ 2. Busca do repo
│               │ 3. Converte para DTO
│               │ 4. Retorna resultado
└───────┬───────┘
        │
        ▼
┌───────────────┐
│  Repository   │ SELECT do DB
└───────┬───────┘
        │
        ▼
┌───────────────┐
│      DTO      │ Retorna ao cliente
└───────────────┘
```

---

## 🏗️ Anatomia de uma Query

### Estrutura Básica

```go
package queries

// Query simples
type GetProductByIDQuery struct {
    ProductID string
}

func (q *GetProductByIDQuery) QueryName() string {
    return "GetProductByIDQuery"
}
```

### Query com Parâmetros

```go
type ListProductsQuery struct {
    // Paginação
    Page     int
    PageSize int
    
    // Filtros
    CategoryID *string
    MinPrice   *float64
    MaxPrice   *float64
    InStock    *bool
    SearchTerm *string
    
    // Ordenação
    SortBy    string // "price", "name", "created_at"
    SortOrder string // "asc", "desc"
}

func (q *ListProductsQuery) QueryName() string {
    return "ListProductsQuery"
}

// Validação
func (q *ListProductsQuery) Validate() error {
    if q.Page < 1 {
        return errors.New("page must be greater than 0")
    }
    if q.PageSize < 1 || q.PageSize > 100 {
        return errors.New("page_size must be between 1 and 100")
    }
    
    validSortFields := []string{"price", "name", "created_at", "stock"}
    if q.SortBy != "" && !contains(validSortFields, q.SortBy) {
        return errors.New("invalid sort_by field")
    }
    
    return nil
}
```

### DTOs de Resposta

```go
package dto

import "time"

// DTO individual
type ProductDTO struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    SKU         string    `json:"sku"`
    Stock       int       `json:"stock"`
    CategoryID  string    `json:"category_id"`
    Active      bool      `json:"active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// DTO paginado
type PaginatedProductsDTO struct {
    Products   []*ProductDTO `json:"products"`
    Total      int64         `json:"total"`
    Page       int           `json:"page"`
    PageSize   int           `json:"page_size"`
    TotalPages int           `json:"total_pages"`
}
```

---

## 🔧 Query Handlers

### Handler Simples

```go
package queries

import (
    "context"
    "meuApp/internal/modules/product/dto"
    "meuApp/internal/modules/product/repository"
    "meuApp/pkg/contracts"
)

type GetProductByIDHandler struct {
    repo repository.ProductRepository
}

func NewGetProductByIDHandler(repo repository.ProductRepository) *GetProductByIDHandler {
    return &GetProductByIDHandler{repo: repo}
}

func (h *GetProductByIDHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    query := qry.(*GetProductByIDQuery)
    
    // Buscar produto
    product, err := h.repo.FindByID(ctx, query.ProductID)
    if err != nil {
        return nil, err
    }
    
    // Converter para DTO
    return &dto.ProductDTO{
        ID:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        SKU:         product.SKU,
        Stock:       product.Stock,
        CategoryID:  product.CategoryID,
        Active:      product.Active,
        CreatedAt:   product.CreatedAt,
        UpdatedAt:   product.UpdatedAt,
    }, nil
}
```

### Handler com Paginação

```go
type ListProductsHandler struct {
    repo repository.ProductRepository
}

func NewListProductsHandler(repo repository.ProductRepository) *ListProductsHandler {
    return &ListProductsHandler{repo: repo}
}

func (h *ListProductsHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    query := qry.(*ListProductsQuery)
    
    // Validar query
    if err := query.Validate(); err != nil {
        return nil, err
    }
    
    // Buscar produtos com filtros
    products, total, err := h.repo.FindAll(
        ctx,
        query.Page,
        query.PageSize,
        query.CategoryID,
        query.MinPrice,
        query.MaxPrice,
        query.InStock,
        query.SearchTerm,
        query.SortBy,
        query.SortOrder,
    )
    if err != nil {
        return nil, err
    }
    
    // Converter para DTOs
    productDTOs := make([]*dto.ProductDTO, len(products))
    for i, product := range products {
        productDTOs[i] = &dto.ProductDTO{
            ID:          product.ID,
            Name:        product.Name,
            Description: product.Description,
            Price:       product.Price,
            SKU:         product.SKU,
            Stock:       product.Stock,
            CategoryID:  product.CategoryID,
            Active:      product.Active,
            CreatedAt:   product.CreatedAt,
            UpdatedAt:   product.UpdatedAt,
        }
    }
    
    // Calcular total de páginas
    totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
    
    return &dto.PaginatedProductsDTO{
        Products:   productDTOs,
        Total:      total,
        Page:       query.Page,
        PageSize:   query.PageSize,
        TotalPages: totalPages,
    }, nil
}
```

### Handler com Cache

```go
type GetProductByIDHandler struct {
    repo  repository.ProductRepository
    cache contracts.Cache
}

func NewGetProductByIDHandler(
    repo repository.ProductRepository,
    cache contracts.Cache,
) *GetProductByIDHandler {
    return &GetProductByIDHandler{
        repo:  repo,
        cache: cache,
    }
}

func (h *GetProductByIDHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    query := qry.(*GetProductByIDQuery)
    
    // Tentar cache primeiro
    cacheKey := fmt.Sprintf("product:%s", query.ProductID)
    if cached, found := h.cache.Get(cacheKey); found {
        return cached.(*dto.ProductDTO), nil
    }
    
    // Buscar do banco
    product, err := h.repo.FindByID(ctx, query.ProductID)
    if err != nil {
        return nil, err
    }
    
    // Converter para DTO
    productDTO := &dto.ProductDTO{
        ID:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        SKU:         product.SKU,
        Stock:       product.Stock,
        CategoryID:  product.CategoryID,
        Active:      product.Active,
        CreatedAt:   product.CreatedAt,
        UpdatedAt:   product.UpdatedAt,
    }
    
    // Armazenar em cache (5 minutos)
    h.cache.Set(cacheKey, productDTO, 5*time.Minute)
    
    return productDTO, nil
}
```

---

## 📄 Paginação

### Paginação Básica (Offset-Based)

```go
type ListProductsQuery struct {
    Page     int // 1, 2, 3...
    PageSize int // 10, 20, 50...
}

// No Repository
func (r *MySQLProductRepository) FindAll(
    ctx context.Context,
    page, pageSize int,
) ([]*entities.Product, int64, error) {
    var products []ProductModel
    var total int64
    
    // Count total
    if err := r.db.WithContext(ctx).
        Model(&ProductModel{}).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Query paginada
    offset := (page - 1) * pageSize
    err := r.db.WithContext(ctx).
        Offset(offset).
        Limit(pageSize).
        Order("created_at DESC").
        Find(&products).Error
    
    if err != nil {
        return nil, 0, err
    }
    
    // Converter para entities
    entities := make([]*entities.Product, len(products))
    for i, model := range products {
        entities[i] = r.toEntity(&model)
    }
    
    return entities, total, nil
}
```

### Paginação Cursor-Based

```go
type ListProductsQuery struct {
    Cursor   *string // ID do último item
    PageSize int
}

// No Repository
func (r *MySQLProductRepository) FindAllCursor(
    ctx context.Context,
    cursor *string,
    pageSize int,
) ([]*entities.Product, *string, error) {
    var products []ProductModel
    
    query := r.db.WithContext(ctx)
    
    // Aplicar cursor
    if cursor != nil {
        query = query.Where("id > ?", *cursor)
    }
    
    // Query
    err := query.
        Order("id ASC").
        Limit(pageSize + 1). // +1 para saber se tem mais
        Find(&products).Error
    
    if err != nil {
        return nil, nil, err
    }
    
    // Verificar se tem próxima página
    var nextCursor *string
    hasMore := len(products) > pageSize
    if hasMore {
        products = products[:pageSize]
        lastID := products[len(products)-1].ID
        nextCursor = &lastID
    }
    
    // Converter
    entities := make([]*entities.Product, len(products))
    for i, model := range products {
        entities[i] = r.toEntity(&model)
    }
    
    return entities, nextCursor, nil
}
```

### DTO com Cursor

```go
type CursorPaginatedProductsDTO struct {
    Products   []*ProductDTO `json:"products"`
    NextCursor *string       `json:"next_cursor"`
    HasMore    bool          `json:"has_more"`
}
```

---

## 🔎 Filtros e Ordenação

### Query com Múltiplos Filtros

```go
type SearchProductsQuery struct {
    // Paginação
    Page     int
    PageSize int
    
    // Filtros
    SearchTerm *string  // Busca em name, description
    CategoryID *string
    MinPrice   *float64
    MaxPrice   *float64
    InStock    *bool
    Active     *bool
    
    // Ordenação
    SortBy    string // "price", "name", "created_at"
    SortOrder string // "asc", "desc"
}
```

### Implementação no Repository

```go
func (r *MySQLProductRepository) Search(
    ctx context.Context,
    page, pageSize int,
    searchTerm *string,
    categoryID *string,
    minPrice, maxPrice *float64,
    inStock, active *bool,
    sortBy, sortOrder string,
) ([]*entities.Product, int64, error) {
    var products []ProductModel
    var total int64
    
    // Base query
    query := r.db.WithContext(ctx).Model(&ProductModel{})
    
    // Aplicar filtros
    if searchTerm != nil && *searchTerm != "" {
        searchPattern := "%" + *searchTerm + "%"
        query = query.Where(
            "name LIKE ? OR description LIKE ?",
            searchPattern, searchPattern,
        )
    }
    
    if categoryID != nil {
        query = query.Where("category_id = ?", *categoryID)
    }
    
    if minPrice != nil {
        query = query.Where("price >= ?", *minPrice)
    }
    
    if maxPrice != nil {
        query = query.Where("price <= ?", *maxPrice)
    }
    
    if inStock != nil && *inStock {
        query = query.Where("stock > 0")
    }
    
    if active != nil {
        query = query.Where("active = ?", *active)
    }
    
    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Aplicar ordenação
    orderClause := "created_at DESC" // default
    if sortBy != "" {
        order := "ASC"
        if sortOrder == "desc" {
            order = "DESC"
        }
        orderClause = fmt.Sprintf("%s %s", sortBy, order)
    }
    query = query.Order(orderClause)
    
    // Paginação
    offset := (page - 1) * pageSize
    err := query.
        Offset(offset).
        Limit(pageSize).
        Find(&products).Error
    
    if err != nil {
        return nil, 0, err
    }
    
    // Converter
    entities := make([]*entities.Product, len(products))
    for i, model := range products {
        entities[i] = r.toEntity(&model)
    }
    
    return entities, total, nil
}
```

### Query Builder Pattern

```go
type ProductQueryBuilder struct {
    db *gorm.DB
}

func NewProductQueryBuilder(db *gorm.DB) *ProductQueryBuilder {
    return &ProductQueryBuilder{db: db}
}

func (b *ProductQueryBuilder) WithSearchTerm(term string) *ProductQueryBuilder {
    if term != "" {
        pattern := "%" + term + "%"
        b.db = b.db.Where("name LIKE ? OR description LIKE ?", pattern, pattern)
    }
    return b
}

func (b *ProductQueryBuilder) WithCategory(categoryID string) *ProductQueryBuilder {
    if categoryID != "" {
        b.db = b.db.Where("category_id = ?", categoryID)
    }
    return b
}

func (b *ProductQueryBuilder) WithPriceRange(min, max *float64) *ProductQueryBuilder {
    if min != nil {
        b.db = b.db.Where("price >= ?", *min)
    }
    if max != nil {
        b.db = b.db.Where("price <= ?", *max)
    }
    return b
}

func (b *ProductQueryBuilder) InStock() *ProductQueryBuilder {
    b.db = b.db.Where("stock > 0")
    return b
}

func (b *ProductQueryBuilder) OrderBy(field, order string) *ProductQueryBuilder {
    b.db = b.db.Order(fmt.Sprintf("%s %s", field, order))
    return b
}

func (b *ProductQueryBuilder) Paginate(page, pageSize int) *ProductQueryBuilder {
    offset := (page - 1) * pageSize
    b.db = b.db.Offset(offset).Limit(pageSize)
    return b
}

func (b *ProductQueryBuilder) Execute() ([]ProductModel, int64, error) {
    var products []ProductModel
    var total int64
    
    // Count
    if err := b.db.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Query
    if err := b.db.Find(&products).Error; err != nil {
        return nil, 0, err
    }
    
    return products, total, nil
}

// Uso
products, total, err := NewProductQueryBuilder(db).
    WithSearchTerm("laptop").
    WithCategory("electronics").
    WithPriceRange(&minPrice, &maxPrice).
    InStock().
    OrderBy("price", "ASC").
    Paginate(1, 20).
    Execute()
```

---

## ⚡ Performance e Otimização

### 1. Eager Loading (Evitar N+1)

```go
// ❌ RUIM: N+1 queries
func (r *MySQLProductRepository) FindAll(ctx context.Context) ([]*entities.Product, error) {
    var products []ProductModel
    r.db.Find(&products)
    
    for _, product := range products {
        // N queries adicionais!
        r.db.Where("id = ?", product.CategoryID).First(&product.Category)
    }
    
    return products, nil
}

// ✅ BOM: 1 query com JOIN
func (r *MySQLProductRepository) FindAll(ctx context.Context) ([]*entities.Product, error) {
    var products []ProductModel
    
    err := r.db.
        Preload("Category").  // Eager load
        Find(&products).Error
    
    return products, err
}
```

### 2. Select Apenas Campos Necessários

```go
// ❌ RUIM: SELECT *
func (r *MySQLProductRepository) FindAllNames(ctx context.Context) ([]string, error) {
    var products []ProductModel
    r.db.Find(&products) // Traz tudo
    
    names := make([]string, len(products))
    for i, p := range products {
        names[i] = p.Name
    }
    return names, nil
}

// ✅ BOM: SELECT apenas campos necessários
func (r *MySQLProductRepository) FindAllNames(ctx context.Context) ([]string, error) {
    var names []string
    err := r.db.
        Model(&ProductModel{}).
        Pluck("name", &names).Error
    
    return names, err
}
```

### 3. Índices no Banco

```go
type ProductModel struct {
    ID         string  `gorm:"primaryKey"`
    Name       string  `gorm:"index"`                    // Índice simples
    SKU        string  `gorm:"uniqueIndex"`              // Índice único
    CategoryID string  `gorm:"index:idx_category_price"` // Índice composto
    Price      float64 `gorm:"index:idx_category_price"`
    Active     bool    `gorm:"index"`
    Stock      int     `gorm:"index"`
    CreatedAt  time.Time `gorm:"index"`
}
```

### 4. Query com Raw SQL (quando necessário)

```go
func (r *MySQLProductRepository) FindTopSellingProducts(
    ctx context.Context,
    limit int,
) ([]*ProductSalesDTO, error) {
    var results []*ProductSalesDTO
    
    err := r.db.Raw(`
        SELECT 
            p.id,
            p.name,
            p.price,
            COUNT(oi.id) as total_sales,
            SUM(oi.quantity) as total_quantity
        FROM products p
        INNER JOIN order_items oi ON oi.product_id = p.id
        GROUP BY p.id
        ORDER BY total_sales DESC
        LIMIT ?
    `, limit).Scan(&results).Error
    
    return results, err
}
```

### 5. Prepared Statements

```go
// GORM usa prepared statements automaticamente
// Mas você pode forçar:
func (r *MySQLProductRepository) FindByCategory(
    ctx context.Context,
    categoryID string,
) ([]*entities.Product, error) {
    var products []ProductModel
    
    // Prepared statement
    stmt := r.db.Session(&gorm.Session{PrepareStmt: true})
    err := stmt.Where("category_id = ?", categoryID).Find(&products).Error
    
    return products, err
}
```

---

## 💾 Cache

### Cache de Query Completa

```go
type ListProductsHandler struct {
    repo  repository.ProductRepository
    cache contracts.Cache
}

func (h *ListProductsHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    query := qry.(*ListProductsQuery)
    
    // Gerar chave de cache
    cacheKey := h.generateCacheKey(query)
    
    // Tentar cache
    if cached, found := h.cache.Get(cacheKey); found {
        return cached.(*dto.PaginatedProductsDTO), nil
    }
    
    // Buscar do banco
    products, total, err := h.repo.FindAll(ctx, ...)
    if err != nil {
        return nil, err
    }
    
    // Converter para DTO
    result := h.toDTO(products, total, query)
    
    // Armazenar em cache (1 minuto)
    h.cache.Set(cacheKey, result, 1*time.Minute)
    
    return result, nil
}

func (h *ListProductsHandler) generateCacheKey(query *ListProductsQuery) string {
    return fmt.Sprintf(
        "products:list:page=%d:size=%d:cat=%v:sort=%s",
        query.Page,
        query.PageSize,
        query.CategoryID,
        query.SortBy,
    )
}
```

### Cache com Invalidação

```go
// Subscriber que invalida cache quando produto é atualizado
type CacheInvalidationSubscriber struct {
    cache contracts.Cache
}

func (s *CacheInvalidationSubscriber) Handle(
    ctx context.Context,
    event contracts.Event,
) error {
    switch e := event.(type) {
    case *events.ProductCreatedEvent:
        s.invalidateProductLists()
    case *events.ProductUpdatedEvent:
        s.invalidateProduct(e.ProductID)
        s.invalidateProductLists()
    case *events.ProductDeletedEvent:
        s.invalidateProduct(e.ProductID)
        s.invalidateProductLists()
    }
    return nil
}

func (s *CacheInvalidationSubscriber) invalidateProduct(productID string) {
    key := fmt.Sprintf("product:%s", productID)
    s.cache.Delete(key)
}

func (s *CacheInvalidationSubscriber) invalidateProductLists() {
    // Invalidar todas as listas (pattern matching)
    s.cache.DeletePattern("products:list:*")
}
```

---

## 📝 Exemplos Práticos

### Exemplo 1: Query Complexa com Joins

```go
type GetOrderWithDetailsQuery struct {
    OrderID string
}

type GetOrderWithDetailsHandler struct {
    repo repository.OrderRepository
}

func (h *GetOrderWithDetailsHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    query := qry.(*GetOrderWithDetailsQuery)
    
    // Buscar pedido com todos relacionamentos
    order, err := h.repo.FindByIDWithDetails(ctx, query.OrderID)
    if err != nil {
        return nil, err
    }
    
    // Converter para DTO detalhado
    return &dto.OrderDetailDTO{
        ID:         order.ID,
        OrderNumber: order.Number,
        Status:     order.Status,
        Total:      order.Total,
        User: &dto.UserDTO{
            ID:    order.User.ID,
            Name:  order.User.Name,
            Email: order.User.Email,
        },
        Items: h.toItemDTOs(order.Items),
        ShippingAddress: &dto.AddressDTO{
            Street:  order.ShippingAddress.Street,
            City:    order.ShippingAddress.City,
            ZipCode: order.ShippingAddress.ZipCode,
        },
        CreatedAt: order.CreatedAt,
    }, nil
}

// No Repository
func (r *MySQLOrderRepository) FindByIDWithDetails(
    ctx context.Context,
    id string,
) (*entities.Order, error) {
    var model OrderModel
    
    err := r.db.WithContext(ctx).
        Preload("User").
        Preload("Items.Product").
        Preload("ShippingAddress").
        Where("id = ?", id).
        First(&model).Error
    
    if err != nil {
        return nil, err
    }
    
    return r.toEntity(&model), nil
}
```

### Exemplo 2: Agregações

```go
type GetProductStatisticsQuery struct {
    ProductID string
}

type ProductStatisticsDTO struct {
    ProductID     string  `json:"product_id"`
    TotalSales    int64   `json:"total_sales"`
    TotalRevenue  float64 `json:"total_revenue"`
    AverageRating float64 `json:"average_rating"`
    ReviewCount   int64   `json:"review_count"`
}

type GetProductStatisticsHandler struct {
    db *gorm.DB
}

func (h *GetProductStatisticsHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    query := qry.(*GetProductStatisticsQuery)
    
    var stats ProductStatisticsDTO
    
    err := h.db.Raw(`
        SELECT 
            p.id as product_id,
            COALESCE(SUM(oi.quantity), 0) as total_sales,
            COALESCE(SUM(oi.quantity * oi.price), 0) as total_revenue,
            COALESCE(AVG(r.rating), 0) as average_rating,
            COUNT(DISTINCT r.id) as review_count
        FROM products p
        LEFT JOIN order_items oi ON oi.product_id = p.id
        LEFT JOIN reviews r ON r.product_id = p.id
        WHERE p.id = ?
        GROUP BY p.id
    `, query.ProductID).Scan(&stats).Error
    
    if err != nil {
        return nil, err
    }
    
    return &stats, nil
}
```

### Exemplo 3: Busca Full-Text

```go
type SearchProductsQuery struct {
    SearchTerm string
    Page       int
    PageSize   int
}

// No Repository (MySQL Full-Text Search)
func (r *MySQLProductRepository) FullTextSearch(
    ctx context.Context,
    searchTerm string,
    page, pageSize int,
) ([]*entities.Product, int64, error) {
    var products []ProductModel
    var total int64
    
    // Criar índice full-text primeiro:
    // CREATE FULLTEXT INDEX idx_product_search ON products(name, description)
    
    query := r.db.WithContext(ctx).
        Where("MATCH(name, description) AGAINST(? IN NATURAL LANGUAGE MODE)", searchTerm)
    
    // Count
    query.Model(&ProductModel{}).Count(&total)
    
    // Query
    offset := (page - 1) * pageSize
    err := query.
        Offset(offset).
        Limit(pageSize).
        Find(&products).Error
    
    // Converter...
    return entities, total, err
}
```

---

## 🎯 Boas Práticas

### 1. Retorne Sempre DTOs, Nunca Entities

```go
// ✅ BOM: Retorna DTO
func (h *GetProductHandler) Handle(...) (interface{}, error) {
    product, _ := h.repo.FindByID(...)
    
    return &dto.ProductDTO{
        ID:    product.ID,
        Name:  product.Name,
        Price: product.Price,
    }, nil
}

// ❌ RUIM: Retorna entity
func (h *GetProductHandler) Handle(...) (interface{}, error) {
    product, _ := h.repo.FindByID(...)
    return product, nil // ❌ Vaza detalhes de domínio
}
```

### 2. Valide Parâmetros de Paginação

```go
func (q *ListProductsQuery) Validate() error {
    if q.Page < 1 {
        q.Page = 1
    }
    if q.PageSize < 1 {
        q.PageSize = 10
    }
    if q.PageSize > 100 {
        q.PageSize = 100 // Limite máximo
    }
    return nil
}
```

### 3. Use Projection para Listas

```go
// DTO completo para detalhes
type ProductDetailDTO struct {
    ID          string
    Name        string
    Description string // Campo grande
    Price       float64
    SKU         string
    Stock       int
    // ... muitos campos
}

// DTO simplificado para listas
type ProductListItemDTO struct {
    ID    string
    Name  string
    Price float64
    Stock int
}
```

### 4. Implemente Timeout

```go
func (h *ListProductsHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    // Timeout de 5 segundos
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    products, total, err := h.repo.FindAll(ctx, ...)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, errors.New("query timeout")
        }
        return nil, err
    }
    
    return products, nil
}
```

### 5. Log Queries Lentas

```go
func (h *ListProductsHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    start := time.Now()
    
    products, total, err := h.repo.FindAll(ctx, ...)
    
    duration := time.Since(start)
    if duration > 1*time.Second {
        h.logger.Warn("Slow query detected", map[string]interface{}{
            "query":    qry.QueryName(),
            "duration": duration.Milliseconds(),
        })
    }
    
    return products, err
}
```

### 6. Use Read Replicas para Queries

```go
type ListProductsHandler struct {
    readDB  *gorm.DB  // Read replica
    writeDB *gorm.DB  // Master
}

func (h *ListProductsHandler) Handle(
    ctx context.Context,
    qry contracts.Query,
) (interface{}, error) {
    // Usar read replica para queries
    repo := repository.NewMySQLProductRepository(h.readDB)
    
    products, total, err := repo.FindAll(ctx, ...)
    return products, err
}
```

---

## 📚 Próximos Passos

- **[Commands](13-implementing-commands.md)** - Write operations
- **[Eventos](15-working-with-events.md)** - Event-driven queries
- **[Boas Práticas](23-best-practices.md)** - Performance tips

---

**[⬅️ Commands](13-implementing-commands.md)** | **[Índice](README.md)** | **[Eventos ➡️](15-working-with-events.md)**
