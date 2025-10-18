package bootstrap

// StartServices era usado para registrar os services antigos
// Agora usamos StartApplicationServices que registra os Application Services (CQRS)
// Este arquivo é mantido vazio para compatibilidade, mas pode ser removido no futuro
func (s *Bootstrap) StartServices() {
// Todos os services agora são registrados via StartApplicationServices()
// que implementa o padrão CQRS com Commands e Queries
}
