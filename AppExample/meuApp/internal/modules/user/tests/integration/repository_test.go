package integration

import (
	"context"
	"testing"
	"time"

	"meuApp/internal/modules/user/domain"
	"meuApp/internal/modules/user/repository"

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
	err = db.AutoMigrate(&repository.UserModel{})
	require.NoError(t, err)

	return db
}

// createTestUser cria um usuário de teste
func createTestUser(username, email string) *domain.User {
	user, err := domain.NewUser(generateID(), username, email)
	if err != nil {
		panic("Failed to create test user: " + err.Error())
	}
	user.Password = "hashedPassword123" // Set password after creation
	return user
}

// generateID gera um ID simples para testes
func generateID() string {
	return "test-" + time.Now().Format("20060102150405.999999999")
}

func TestUserRepository_Create(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	user := createTestUser("johndoe", "john@example.com")

	// Act
	err := repo.Create(ctx, user)

	// Assert
	require.NoError(t, err)

	// Verify user was created in database
	var count int64
	db.Model(&repository.UserModel{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	user1 := createTestUser("johndoe", "john@example.com")
	user2 := createTestUser("janedoe", "john@example.com") // Same email

	// Act
	err1 := repo.Create(ctx, user1)
	err2 := repo.Create(ctx, user2)

	// Assert
	require.NoError(t, err1)
	assert.Error(t, err2) // Should fail due to unique constraint
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	user := createTestUser("johndoe", "john@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Act
	retrievedUser, err := repo.GetByID(ctx, user.ID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	// Act
	retrievedUser, err := repo.GetByID(ctx, "non-existent-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
	assert.Contains(t, err.Error(), "not found")
}

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	user := createTestUser("johndoe", "john@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Act
	retrievedUser, err := repo.GetByEmail(ctx, "john@example.com")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Email, retrievedUser.Email)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	// Act
	retrievedUser, err := repo.GetByEmail(ctx, "nonexistent@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedUser)
	assert.Contains(t, err.Error(), "not found")
}

func TestUserRepository_Update_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	user := createTestUser("johndoe", "john@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Act - Update username
	user.Username = "JohnUpdated"
	err = repo.Update(ctx, user)

	// Assert
	require.NoError(t, err)

	// Verify update
	updatedUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "JohnUpdated", updatedUser.Username)
}

func TestUserRepository_Delete_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	user := createTestUser("johndoe", "john@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// Act
	err = repo.Delete(ctx, user.ID)

	// Assert
	require.NoError(t, err)

	// Verify deletion
	deletedUser, err := repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserRepository_List_Success(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	// Create multiple users
	users := []*domain.User{
		createTestUser("user1", "user1@example.com"),
		createTestUser("user2", "user2@example.com"),
		createTestUser("user3", "user3@example.com"),
		createTestUser("user4", "user4@example.com"),
		createTestUser("user5", "user5@example.com"),
	}

	for _, user := range users {
		err := repo.Create(ctx, user)
		require.NoError(t, err)
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	}

	// Act
	retrievedUsers, err := repo.List(ctx, 0, 10)

	// Assert
	require.NoError(t, err)
	assert.Len(t, retrievedUsers, 5)
}

func TestUserRepository_List_WithPagination(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	// Create 5 users
	for i := 1; i <= 5; i++ {
		username := "user" + string(rune('0'+i))
		email := username + "@example.com"
		user := createTestUser(username, email)
		err := repo.Create(ctx, user)
		require.NoError(t, err)
		time.Sleep(1 * time.Millisecond)
	}

	// Act - Get first page (2 items)
	page1, err := repo.List(ctx, 0, 2)
	require.NoError(t, err)

	// Act - Get second page (2 items)
	page2, err := repo.List(ctx, 2, 2)
	require.NoError(t, err)

	// Assert
	assert.Len(t, page1, 2)
	assert.Len(t, page2, 2)

	// Verify pages don't overlap
	assert.NotEqual(t, page1[0].ID, page2[0].ID)
}

func TestUserRepository_List_EmptyDatabase(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	// Act
	users, err := repo.List(ctx, 0, 10)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepository_CRUD_FullFlow(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := repository.NewMySQLUserRepository(db)
	ctx := context.Background()

	// 1. Create
	user := createTestUser("johndoe", "john@example.com")
	err := repo.Create(ctx, user)
	require.NoError(t, err)

	// 2. Read
	retrievedUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Username, retrievedUser.Username)

	// 3. Update
	retrievedUser.Username = "JohnUpdated"
	err = repo.Update(ctx, retrievedUser)
	require.NoError(t, err)

	// 4. Verify Update
	updatedUser, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, "JohnUpdated", updatedUser.Username)

	// 5. Delete
	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	// 6. Verify Deletion
	deletedUser, err := repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
}

func TestUserRepository_ConcurrentCreates(t *testing.T) {
	t.Skip("Skipping concurrent test with SQLite in-memory - requires real database")

	// This test would work with a real MySQL/PostgreSQL database
	// SQLite in-memory has limitations with concurrent writes
}
