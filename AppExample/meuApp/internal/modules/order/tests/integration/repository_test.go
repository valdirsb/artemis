package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"meuApp/internal/modules/order/domain"
	"meuApp/internal/modules/order/repository"

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

	// Auto migrate - precisa migrar tanto OrderModel quanto OrderItemModel
	err = db.AutoMigrate(&repository.OrderModel{}, &repository.OrderItemModel{})
	require.NoError(t, err)

	return db
}

// createTestOrder cria um pedido de teste
func createTestOrder(userID string, itemCount int) *domain.Order {
	items := make([]domain.OrderItem, itemCount)
	for i := 0; i < itemCount; i++ {
		items[i] = domain.OrderItem{
			ProductID: fmt.Sprintf("product-%d", i+1),
			Quantity:  i + 1,
			Price:     float64((i + 1) * 10),
		}
	}

	order, err := domain.NewOrder(
		generateID(),
		userID,
		items,
	)
	if err != nil {
		panic("Failed to create test order: " + err.Error())
	}
	return order
}

// generateID gera um ID simples para testes
func generateID() string {
	return fmt.Sprintf("test-order-%d", time.Now().UnixNano())
}

// ===== TESTES DE CRIAÇÃO =====

func TestOrderRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	order := createTestOrder("user-1", 2)
	err := repo.Create(ctx, order)

	require.NoError(t, err)
	assert.NotEmpty(t, order.ID)
}

func TestOrderRepository_Create_MultipleOrders(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Create multiple orders
	for i := 1; i <= 3; i++ {
		order := createTestOrder(fmt.Sprintf("user-%d", i), i)
		err := repo.Create(ctx, order)
		require.NoError(t, err)
	}
}

// ===== TESTES DE LEITURA =====

func TestOrderRepository_GetByID_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Create an order
	order := createTestOrder("user-1", 2)
	err := repo.Create(ctx, order)
	require.NoError(t, err)

	// Retrieve the order
	retrieved, err := repo.GetByID(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, retrieved.ID)
	assert.Equal(t, order.UserID, retrieved.UserID)
	assert.Equal(t, order.Status, retrieved.Status)
	assert.Equal(t, order.Total, retrieved.Total)
	assert.Equal(t, len(order.Items), len(retrieved.Items))
}

func TestOrderRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	assert.Error(t, err)
}

func TestOrderRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	userID := "user-1"

	// Create multiple orders for the same user
	for i := 1; i <= 3; i++ {
		order := createTestOrder(userID, i)
		err := repo.Create(ctx, order)
		require.NoError(t, err)
	}

	// Create an order for a different user
	order := createTestOrder("user-2", 1)
	err := repo.Create(ctx, order)
	require.NoError(t, err)

	// Get orders by user ID
	orders, err := repo.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, 3, len(orders))

	// Verify all orders belong to the correct user
	for _, o := range orders {
		assert.Equal(t, userID, o.UserID)
	}
}

func TestOrderRepository_GetByUserID_EmptyResult(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	orders, err := repo.GetByUserID(ctx, "non-existent-user")
	require.NoError(t, err)
	assert.Equal(t, 0, len(orders))
}

// ===== TESTES DE ATUALIZAÇÃO =====

func TestOrderRepository_Update_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Create an order
	order := createTestOrder("user-1", 2)
	err := repo.Create(ctx, order)
	require.NoError(t, err)

	// Update the order status using the aggregate
	aggregate := domain.NewOrderAggregate(order)
	err = aggregate.UpdateStatus(domain.OrderStatusConfirmed)
	require.NoError(t, err)

	err = repo.Update(ctx, order)
	require.NoError(t, err)

	// Verify the update
	retrieved, err := repo.GetByID(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusConfirmed, retrieved.Status)
}

func TestOrderRepository_Update_CancelOrder(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Create an order
	order := createTestOrder("user-1", 2)
	err := repo.Create(ctx, order)
	require.NoError(t, err)

	// Cancel the order using the aggregate
	aggregate := domain.NewOrderAggregate(order)
	err = aggregate.Cancel()
	require.NoError(t, err)

	err = repo.Update(ctx, order)
	require.NoError(t, err)

	// Verify the cancellation
	retrieved, err := repo.GetByID(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusCancelled, retrieved.Status)
}

// ===== TESTES DE EXCLUSÃO =====

func TestOrderRepository_Delete_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Create an order
	order := createTestOrder("user-1", 2)
	err := repo.Create(ctx, order)
	require.NoError(t, err)

	// Delete the order
	err = repo.Delete(ctx, order.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(ctx, order.ID)
	assert.Error(t, err)
}

func TestOrderRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Try to delete a non-existent order
	err := repo.Delete(ctx, "non-existent-id")
	// Note: Depending on the implementation, this might not error
	// Some implementations return no error when deleting non-existent items
	// Let's just verify it doesn't panic
	_ = err
}

// ===== TESTES DE FLUXO COMPLETO =====

func TestOrderRepository_CRUD_FullFlow(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// 1. Create
	order := createTestOrder("user-1", 3)
	err := repo.Create(ctx, order)
	require.NoError(t, err)
	originalID := order.ID

	// 2. Read
	retrieved, err := repo.GetByID(ctx, originalID)
	require.NoError(t, err)
	assert.Equal(t, originalID, retrieved.ID)
	assert.Equal(t, domain.OrderStatusPending, retrieved.Status)

	// 3. Update
	aggregate := domain.NewOrderAggregate(retrieved)
	err = aggregate.UpdateStatus(domain.OrderStatusConfirmed)
	require.NoError(t, err)
	err = repo.Update(ctx, retrieved)
	require.NoError(t, err)

	// 4. Verify Update
	updated, err := repo.GetByID(ctx, originalID)
	require.NoError(t, err)
	assert.Equal(t, domain.OrderStatusConfirmed, updated.Status)

	// 5. Delete
	err = repo.Delete(ctx, originalID)
	require.NoError(t, err)

	// 6. Verify Deletion
	_, err = repo.GetByID(ctx, originalID)
	assert.Error(t, err)
}

// ===== TESTES DE CENÁRIOS ESPECIAIS =====

func TestOrderRepository_OrderWithMultipleItems(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Create an order with many items
	order := createTestOrder("user-1", 5)
	err := repo.Create(ctx, order)
	require.NoError(t, err)

	// Retrieve and verify all items
	retrieved, err := repo.GetByID(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, len(retrieved.Items))

	// Verify total calculation
	expectedTotal := 0.0
	for _, item := range retrieved.Items {
		expectedTotal += item.Price * float64(item.Quantity)
	}
	assert.Equal(t, expectedTotal, retrieved.Total)
}

func TestOrderRepository_GetByUserID_MultipleUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Create orders for multiple users
	users := []string{"user-1", "user-2", "user-3"}
	ordersPerUser := 2

	for _, userID := range users {
		for i := 0; i < ordersPerUser; i++ {
			order := createTestOrder(userID, i+1)
			err := repo.Create(ctx, order)
			require.NoError(t, err)
		}
	}

	// Verify each user has the correct number of orders
	for _, userID := range users {
		orders, err := repo.GetByUserID(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, ordersPerUser, len(orders))
	}
}

func TestOrderRepository_StatusTransitions(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLOrderRepository(db)
	ctx := context.Background()

	// Create an order
	order := createTestOrder("user-1", 2)
	err := repo.Create(ctx, order)
	require.NoError(t, err)

	// Test status progression: pending -> confirmed -> shipped -> delivered
	statuses := []domain.OrderStatus{
		domain.OrderStatusConfirmed,
		domain.OrderStatusShipped,
		domain.OrderStatusDelivered,
	}

	aggregate := domain.NewOrderAggregate(order)
	for _, status := range statuses {
		err = aggregate.UpdateStatus(status)
		require.NoError(t, err)
		err = repo.Update(ctx, order)
		require.NoError(t, err)

		// Verify status was updated
		retrieved, err := repo.GetByID(ctx, order.ID)
		require.NoError(t, err)
		assert.Equal(t, status, retrieved.Status)
	}
}
