package http

import (
	"net/http"

	"meuApp/internal/modules/user/dto"
	"meuApp/internal/modules/user/ports"

	"github.com/gin-gonic/gin"
)

// UserHTTPHandler é o adapter HTTP para o módulo de usuário
// Converte requests HTTP para comandos/queries da camada de aplicação
type UserHTTPHandler struct {
	userService ports.UserService
}

// NewUserHTTPHandler cria uma nova instância do handler HTTP
func NewUserHTTPHandler(userService ports.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{userService: userService}
}

// CreateUser é o endpoint HTTP para criação de usuário
func (h *UserHTTPHandler) CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Converter domain para DTO response
	response := dto.ToUserResponse(createdUser)
	c.JSON(http.StatusCreated, response)
}

// GetUser é o endpoint HTTP para buscar usuário por ID
func (h *UserHTTPHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Converter domain para DTO response
	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

// UpdateUser é o endpoint HTTP para atualizar usuário
func (h *UserHTTPHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), id, req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Converter domain para DTO response
	response := dto.ToUserResponse(updatedUser)
	c.JSON(http.StatusOK, response)
}

// DeleteUser é o endpoint HTTP para deletar usuário
func (h *UserHTTPHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ValidateUser é o endpoint HTTP para validar credenciais
func (h *UserHTTPHandler) ValidateUser(c *gin.Context) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.ValidateCredentials(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Converter domain para DTO response
	response := dto.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}
