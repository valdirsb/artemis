package integration

import (
	"context"
	"fmt"
	"testing"

	"meuApp/internal/modules/product/ports"
	"meuApp/internal/modules/product/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductRepository_ListPaginated(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	// Criar 25 produtos de teste
	for i := 1; i <= 25; i++ {
		product := createTestProduct(
			fmt.Sprintf("Product %d", i),
			"cat-test",
			float64(i*10),
			i,
		)
		err := repo.Create(ctx, product)
		require.NoError(t, err)
	}

	t.Run("Primeira página com 10 itens", func(t *testing.T) {
		filters := ports.ProductFilters{
			Page:     1,
			PageSize: 10,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 10, len(result.Items))
		assert.Equal(t, int64(25), result.TotalItems)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 3, result.TotalPages)
	})

	t.Run("Segunda página com 10 itens", func(t *testing.T) {
		filters := ports.ProductFilters{
			Page:     2,
			PageSize: 10,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 10, len(result.Items))
		assert.Equal(t, int64(25), result.TotalItems)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 3, result.TotalPages)
	})

	t.Run("Última página com 5 itens", func(t *testing.T) {
		filters := ports.ProductFilters{
			Page:     3,
			PageSize: 10,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, len(result.Items))
		assert.Equal(t, int64(25), result.TotalItems)
		assert.Equal(t, 3, result.Page)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 3, result.TotalPages)
	})

	t.Run("Página inexistente retorna vazio", func(t *testing.T) {
		filters := ports.ProductFilters{
			Page:     10,
			PageSize: 10,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
		assert.Equal(t, int64(25), result.TotalItems)
		assert.Equal(t, 10, result.Page)
		assert.Equal(t, 3, result.TotalPages)
	})

	t.Run("Page size diferente", func(t *testing.T) {
		filters := ports.ProductFilters{
			Page:     1,
			PageSize: 5,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, len(result.Items))
		assert.Equal(t, int64(25), result.TotalItems)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 5, result.PageSize)
		assert.Equal(t, 5, result.TotalPages)
	})

	t.Run("Valores padrão quando page=0", func(t *testing.T) {
		filters := ports.ProductFilters{
			Page:     0,
			PageSize: 10,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Page) // Deve usar page=1 como padrão
		assert.Equal(t, 10, len(result.Items))
	})

	t.Run("Valores padrão quando page_size=0", func(t *testing.T) {
		filters := ports.ProductFilters{
			Page:     1,
			PageSize: 0,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 10, result.PageSize) // Deve usar 10 como padrão
		assert.Equal(t, 10, len(result.Items))
	})

	t.Run("Paginação com filtro de categoria", func(t *testing.T) {
		// Criar produtos de categoria diferente
		for i := 1; i <= 5; i++ {
			product := createTestProduct(
				fmt.Sprintf("Product Cat2 %d", i),
				"cat-test-2",
				float64(i*10),
				i,
			)
			err := repo.Create(ctx, product)
			require.NoError(t, err)
		}

		categoryID := "cat-test-2"
		filters := ports.ProductFilters{
			CategoryID: &categoryID,
			Page:       1,
			PageSize:   10,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, len(result.Items))
		assert.Equal(t, int64(5), result.TotalItems)
		assert.Equal(t, 1, result.TotalPages)

		// Verificar que todos são da categoria correta
		for _, item := range result.Items {
			assert.Equal(t, "cat-test-2", item.CategoryID)
		}
	})

	t.Run("Paginação com filtro de preço", func(t *testing.T) {
		minPrice := 100.0
		maxPrice := 200.0
		filters := ports.ProductFilters{
			MinPrice: &minPrice,
			MaxPrice: &maxPrice,
			Page:     1,
			PageSize: 10,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Verificar que todos os produtos estão na faixa de preço
		for _, item := range result.Items {
			assert.GreaterOrEqual(t, item.Price, minPrice)
			assert.LessOrEqual(t, item.Price, maxPrice)
		}
	})

	t.Run("Paginação com filtro de estoque", func(t *testing.T) {
		// Criar produtos sem estoque
		for i := 1; i <= 3; i++ {
			product := createTestProduct(
				fmt.Sprintf("Product No Stock %d", i),
				"cat-test",
				float64(i*10),
				0, // sem estoque
			)
			err := repo.Create(ctx, product)
			require.NoError(t, err)
		}

		inStock := true
		filters := ports.ProductFilters{
			InStock:  &inStock,
			Page:     1,
			PageSize: 100,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Verificar que todos os produtos têm estoque
		for _, item := range result.Items {
			assert.Greater(t, item.Stock, 0)
		}
	})

	t.Run("Paginação com múltiplos filtros", func(t *testing.T) {
		categoryID := "cat-test"
		minPrice := 50.0
		maxPrice := 150.0
		inStock := true

		filters := ports.ProductFilters{
			CategoryID: &categoryID,
			MinPrice:   &minPrice,
			MaxPrice:   &maxPrice,
			InStock:    &inStock,
			Page:       1,
			PageSize:   5,
		}

		result, err := repo.ListPaginated(ctx, filters)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Verificar que todos os produtos atendem aos filtros
		for _, item := range result.Items {
			assert.Equal(t, categoryID, item.CategoryID)
			assert.GreaterOrEqual(t, item.Price, minPrice)
			assert.LessOrEqual(t, item.Price, maxPrice)
			assert.Greater(t, item.Stock, 0)
		}
	})
}

func TestProductRepository_ListPaginated_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMySQLProductRepository(db)
	ctx := context.Background()

	filters := ports.ProductFilters{
		Page:     1,
		PageSize: 10,
	}

	result, err := repo.ListPaginated(ctx, filters)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Items))
	assert.Equal(t, int64(0), result.TotalItems)
	assert.Equal(t, 0, result.TotalPages)
}

func TestProductRepository_ListPaginated_CalculatesCorrectTotalPages(t *testing.T) {
	testCases := []struct {
		name               string
		totalProducts      int
		pageSize           int
		expectedTotalPages int
	}{
		{"10 produtos, 10 por página", 10, 10, 1},
		{"11 produtos, 10 por página", 11, 10, 2},
		{"20 produtos, 10 por página", 20, 10, 2},
		{"25 produtos, 10 por página", 25, 10, 3},
		{"100 produtos, 20 por página", 100, 20, 5},
		{"7 produtos, 5 por página", 7, 5, 2},
		{"1 produto, 10 por página", 1, 10, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Criar novo banco para cada teste
			db := setupTestDB(t)
			repo := repository.NewMySQLProductRepository(db)
			ctx := context.Background()

			// Criar produtos
			for i := 1; i <= tc.totalProducts; i++ {
				product := createTestProduct(
					fmt.Sprintf("Product %d", i),
					"cat-test",
					float64(i*10),
					i,
				)
				err := repo.Create(ctx, product)
				require.NoError(t, err)
			}

			filters := ports.ProductFilters{
				Page:     1,
				PageSize: tc.pageSize,
			}

			result, err := repo.ListPaginated(ctx, filters)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedTotalPages, result.TotalPages,
				"TotalPages incorreto para %s", tc.name)
		})
	}
}
