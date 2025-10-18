package bootstrap

import (
	orderRepository "meuApp/internal/modules/order/repository"
	productRepository "meuApp/internal/modules/product/repository"
	userRepository "meuApp/internal/modules/user/adapters/repository"

	"gorm.io/gorm"
)

func (s *Bootstrap) StartRepositories() {

	db := s.container.MustGet("database").(*gorm.DB)

	// Inicializa os repositórios e os registra no container
	s.Repositories["productRepository"] = productRepository.NewMySQLProductRepository(db)
	s.Repositories["userRepository"] = userRepository.NewMySQLUserRepository(db)
	s.Repositories["orderRepository"] = orderRepository.NewMySQLOrderRepository(db)

}
