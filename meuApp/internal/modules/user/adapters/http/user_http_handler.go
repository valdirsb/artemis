package http

import (
	"net/http"

	"meuApp/internal/modules/user/dto"
	"meuApp/internal/modules/user/ports"
	"meuApp/pkg/adapters/http/middleware"

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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with username, email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 409 {object} map[string]string "User already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users [post]
func (h *UserHTTPHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.userService.CreateUser(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	// Converter domain para DTO response
	response := dto.ToUserResponse(createdUser)
	middleware.RespondWithJSON(c.Writer, http.StatusCreated, response)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users/{id} [get]
func (h *UserHTTPHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	// Converter domain para DTO response
	response := dto.ToUserResponse(user)
	middleware.RespondWithJSON(c.Writer, http.StatusOK, response)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information (username and email)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "User data"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users/{id} [put]
func (h *UserHTTPHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), id, req.Username, req.Email)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	// Converter domain para DTO response
	response := dto.ToUserResponse(updatedUser)
	middleware.RespondWithJSON(c.Writer, http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users/{id} [delete]
func (h *UserHTTPHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ValidateUser godoc
// @Summary Validate user credentials
// @Description Validate user email and password for authentication
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users/validate [post]
func (h *UserHTTPHandler) ValidateUser(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.ValidateCredentials(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		middleware.RespondWithAppError(c.Writer, err)
		return
	}

	// Converter domain para DTO response
	response := dto.ToUserResponse(user)
	middleware.RespondWithJSON(c.Writer, http.StatusOK, response)
}
