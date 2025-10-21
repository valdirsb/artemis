package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"meuApp/internal/modules/product/domain"
	"meuApp/internal/modules/product/ports"
	"meuApp/internal/modules/product/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB cria um banco de dados SQLite em memória para testes
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// Auto migrate
	err = db.AutoMigrate(&repository.ProductModel{})
	require.NoError(t, err)

	return db
}

// createTestProduct cria um produto de teste
func createTestProduct(name, categoryID string, price float64, stock int) *domain.Product {
	product, err := domain.NewProduct(
		generateID(),
		name,
		"Test description for "+name,
		categoryID,
		price,
		stock,
	)
	if err != nil {
		panic("Failed to create test product: " + err.Error())
	}
	return product
}

// generateID gera um ID simples para testes
func generateID() string {
	return "test-" + time.Now().Format("20060102150405.999999999")
}

func TestProductRepository_Create(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	product := createTestProduct("Laptop", "electronics", 999.99, 10)

	// Act
	err := repo.Create(ctx, product)

	// Assert
	require.NoError(t, err)

	// Verify product was created in database
	var count int64
	db.Model(&repository.ProductModel{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestProductRepository_Create_MultipleProducts(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	product1 := createTestProduct("Laptop", "electronics", 999.99, 10)
	product2 := createTestProduct("Mouse", "electronics", 29.99, 50)

	// Act
	err1 := repo.Create(ctx, product1)
	err2 := repo.Create(ctx, product2)

	// Assert
	require.NoError(t, err1)
	require.NoError(t, err2)

	var count int64
	db.Model(&repository.ProductModel{}).Count(&count)
	assert.Equal(t, int64(2), count)
}

func TestProductRepository_GetByID_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	product := createTestProduct("Keyboard", "electronics", 79.99, 25)
	err := repo.Create(ctx, product)
	require.NoError(t, err)

	// Act
	retrievedProduct, err := repo.GetByID(ctx, product.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, product.ID, retrievedProduct.ID)
	assert.Equal(t, product.Name, retrievedProduct.Name)
	assert.Equal(t, product.Price, retrievedProduct.Price)
	assert.Equal(t, product.Stock, retrievedProduct.Stock)
	assert.Equal(t, product.CategoryID, retrievedProduct.CategoryID)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Act
	retrievedProduct, err := repo.GetByID(ctx, "non-existent-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedProduct)
	assert.Contains(t, err.Error(), "not found")
}

func TestProductRepository_Update_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	product := createTestProduct("Monitor", "electronics", 299.99, 15)
	err := repo.Create(ctx, product)
	require.NoError(t, err)

	// Act - Update price and stock
	product.Price = 249.99
	product.Stock = 20
	err = repo.Update(ctx, product)

	// Assert
	require.NoError(t, err)

	// Verify update
	updatedProduct, err := repo.GetByID(ctx, product.ID)
	require.NoError(t, err)
	assert.Equal(t, 249.99, updatedProduct.Price)
	assert.Equal(t, 20, updatedProduct.Stock)
}

func TestProductRepository_Delete_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	product := createTestProduct("Headphones", "electronics", 149.99, 30)
	err := repo.Create(ctx, product)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, product.ID)

	// Assert
	require.NoError(t, err)

	// Verify deletion
	deletedProduct, err := repo.GetByID(ctx, product.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedProduct)
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Act
	err := repo.Delete(ctx, "non-existent-id")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestProductRepository_List_All(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Create multiple products
	products := []*domain.Product{
		createTestProduct("Product1", "category1", 10.0, 100),
		createTestProduct("Product2", "category2", 20.0, 200),
		createTestProduct("Product3", "category1", 30.0, 300),
		createTestProduct("Product4", "category3", 40.0, 400),
	}

	for _, product := range products {
		err := repo.Create(ctx, product)
		require.NoError(t, err)
		time.Sleep(1 * time.Millisecond)
	}

	// Act
	filters := ports.ProductFilters{}
	retrievedProducts, err := repo.List(ctx, filters)

	// Assert
	require.NoError(t, err)
	assert.Len(t, retrievedProducts, 4)
}

func TestProductRepository_List_ByCategory(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Create products in different categories
	product1 := createTestProduct("Laptop", "electronics", 999.99, 10)
	product2 := createTestProduct("Mouse", "electronics", 29.99, 50)
	product3 := createTestProduct("Desk", "furniture", 299.99, 5)

	repo.Create(ctx, product1)
	repo.Create(ctx, product2)
	repo.Create(ctx, product3)

	// Act - Filter by electronics category
	categoryID := "electronics"
	filters := ports.ProductFilters{
		CategoryID: &categoryID,
	}
	retrievedProducts, err := repo.List(ctx, filters)

	// Assert
	require.NoError(t, err)
	assert.Len(t, retrievedProducts, 2)
	for _, p := range retrievedProducts {
		assert.Equal(t, "electronics", p.CategoryID)
	}
}

func TestProductRepository_List_ByPriceRange(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Create products with different prices
	repo.Create(ctx, createTestProduct("Cheap", "cat1", 10.0, 100))
	repo.Create(ctx, createTestProduct("Medium", "cat1", 50.0, 100))
	repo.Create(ctx, createTestProduct("Expensive", "cat1", 100.0, 100))

	// Act - Filter by price range (30.0 - 80.0)
	minPrice := 30.0
	maxPrice := 80.0
	filters := ports.ProductFilters{
		MinPrice: &minPrice,
		MaxPrice: &maxPrice,
	}
	retrievedProducts, err := repo.List(ctx, filters)

	// Assert
	require.NoError(t, err)
	assert.Len(t, retrievedProducts, 1)
	assert.Equal(t, "Medium", retrievedProducts[0].Name)
	assert.Equal(t, 50.0, retrievedProducts[0].Price)
}

func TestProductRepository_List_EmptyDatabase(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Act
	filters := ports.ProductFilters{}
	products, err := repo.List(ctx, filters)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, products)
}

// Note: UpdateStock is a service-level operation, not a repository method
// The repository only has Create, GetByID, Update, Delete, and List
// To update stock, we use the Update method with the full product

func TestProductRepository_CRUD_FullFlow(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// 1. Create
	product := createTestProduct("TestProduct", "testcat", 99.99, 50)
	err := repo.Create(ctx, product)
	require.NoError(t, err)

	// 2. Read
	retrievedProduct, err := repo.GetByID(ctx, product.ID)
	require.NoError(t, err)
	assert.Equal(t, product.Name, retrievedProduct.Name)

	// 3. Update
	retrievedProduct.Price = 79.99
	retrievedProduct.Stock = 75
	err = repo.Update(ctx, retrievedProduct)
	require.NoError(t, err)

	// 4. Verify Update
	updatedProduct, err := repo.GetByID(ctx, product.ID)
	require.NoError(t, err)
	assert.Equal(t, 79.99, updatedProduct.Price)
	assert.Equal(t, 75, updatedProduct.Stock)

	// 5. Delete
	err = repo.Delete(ctx, product.ID)
	require.NoError(t, err)

	// 6. Verify Deletion
	deletedProduct, err := repo.GetByID(ctx, product.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedProduct)
}

func TestProductRepository_List_Pagination(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Note: The Product repository List method doesn't support pagination yet
	// This test demonstrates listing all products with filters
	// For pagination support, consider adding Limit/Offset to ProductFilters

	// Create multiple products
	for i := 1; i <= 5; i++ {
		product := createTestProduct(fmt.Sprintf("Product %d", i), "cat-1", float64(i*10), i)
		err := repo.Create(ctx, product)
		require.NoError(t, err)
	}

	// List all products (no pagination)
	products, err := repo.List(ctx, ports.ProductFilters{})
	require.NoError(t, err)
	assert.Equal(t, 5, len(products))
}

func TestProductRepository_List_ComplexFilters(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Create diverse products
	repo.Create(ctx, createTestProduct("ElectronicA", "electronics", 100.0, 10))
	repo.Create(ctx, createTestProduct("ElectronicB", "electronics", 200.0, 20))
	repo.Create(ctx, createTestProduct("FurnitureA", "furniture", 150.0, 5))
	repo.Create(ctx, createTestProduct("ElectronicC", "electronics", 50.0, 30))

	// Act - Complex filter: electronics category + price between 75-175
	categoryID := "electronics"
	minPrice := 75.0
	maxPrice := 175.0
	filters := ports.ProductFilters{
		CategoryID: &categoryID,
		MinPrice:   &minPrice,
		MaxPrice:   &maxPrice,
	}
	retrievedProducts, err := repo.List(ctx, filters)

	// Assert
	require.NoError(t, err)
	assert.Len(t, retrievedProducts, 1) // Only ElectronicA matches (100.0)
	assert.Equal(t, "ElectronicA", retrievedProducts[0].Name)
}
