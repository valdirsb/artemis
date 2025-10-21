# ✅ Validation

## 📋 Índice
- [Visão Geral](#visão-geral)
- [Validação Manual](#validação-manual)
- [Validação com go-playground/validator](#validação-com-go-playgroundvalidator)
- [Validação de Domínio](#validação-de-domínio)
- [Custom Validators](#custom-validators)
- [Validação em Commands](#validação-em-commands)
- [Error Messages](#error-messages)
- [Exemplos Práticos](#exemplos-práticos)

---

## 🎯 Visão Geral

A validação no Artemis Framework ocorre em **múltiplas camadas** para garantir a integridade dos dados:

```
┌────────────────────────────────────────────────┐
│         1. HTTP Layer (Request)                │
│   • Bind & Validate struct tags               │
│   • Basic type checking                        │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│         2. Application Layer (Command)         │
│   • Business rules validation                  │
│   • Input sanitization                         │
│   • Cross-field validation                     │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│         3. Domain Layer (Entity)               │
│   • Invariants enforcement                     │
│   • Domain-specific rules                      │
│   • Value objects validation                   │
└────────────────────────────────────────────────┘
```

### Princípios de Validação

1. **Fail Fast** - Valide o mais cedo possível
2. **Clear Messages** - Mensagens de erro claras e acionáveis
3. **Layered Validation** - Cada camada tem sua responsabilidade
4. **Type Safety** - Use tipos fortes sempre que possível

---

## 🖐️ Validação Manual

### Validação Básica em Commands

**`internal/modules/user/application/commands/create_user.go`**

```go
package commands

import (
    "context"
    "fmt"
    "regexp"
    "strings"

    "meuApp/internal/modules/user/domain"
    "meuApp/internal/modules/user/ports"
    "meuApp/pkg/errors"
)

type CreateUserCommand struct {
    Username string
    Email    string
    Password string
}

type CreateUserHandler struct {
    userRepo ports.UserRepository
    logger   contracts.Logger
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
    // 1. Validação de campos obrigatórios
    if err := h.validateRequired(cmd); err != nil {
        return nil, err
    }

    // 2. Validação de formato
    if err := h.validateFormats(cmd); err != nil {
        return nil, err
    }

    // 3. Validação de regras de negócio
    if err := h.validateBusinessRules(ctx, cmd); err != nil {
        return nil, err
    }

    // Continuar com a criação...
    user, err := domain.NewUser(cmd.Username, cmd.Email)
    // ...
}

// validateRequired valida campos obrigatórios
func (h *CreateUserHandler) validateRequired(cmd CreateUserCommand) error {
    if strings.TrimSpace(cmd.Username) == "" {
        return errors.NewValidationError("username", "username is required")
    }
    if strings.TrimSpace(cmd.Email) == "" {
        return errors.NewValidationError("email", "email is required")
    }
    if strings.TrimSpace(cmd.Password) == "" {
        return errors.NewValidationError("password", "password is required")
    }
    return nil
}

// validateFormats valida formatos de dados
func (h *CreateUserHandler) validateFormats(cmd CreateUserCommand) error {
    // Email format
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(cmd.Email) {
        return errors.NewValidationError("email", "invalid email format")
    }

    // Password strength
    if len(cmd.Password) < 8 {
        return errors.NewValidationError("password", "password must be at least 8 characters")
    }
    if !containsUppercase(cmd.Password) {
        return errors.NewValidationError("password", "password must contain at least one uppercase letter")
    }
    if !containsNumber(cmd.Password) {
        return errors.NewValidationError("password", "password must contain at least one number")
    }

    // Username length
    if len(cmd.Username) < 3 || len(cmd.Username) > 50 {
        return errors.NewValidationError("username", "username must be between 3 and 50 characters")
    }

    return nil
}

// validateBusinessRules valida regras de negócio
func (h *CreateUserHandler) validateBusinessRules(ctx context.Context, cmd CreateUserCommand) error {
    // Verificar se email já existe
    existingUser, err := h.userRepo.GetByEmail(ctx, cmd.Email)
    if err == nil && existingUser != nil {
        return errors.NewValidationError("email", "email already in use")
    }

    // Verificar se username já existe
    existingUser, err = h.userRepo.GetByUsername(ctx, cmd.Username)
    if err == nil && existingUser != nil {
        return errors.NewValidationError("username", "username already taken")
    }

    return nil
}

// Helper functions
func containsUppercase(s string) bool {
    return regexp.MustCompile(`[A-Z]`).MatchString(s)
}

func containsNumber(s string) bool {
    return regexp.MustCompile(`[0-9]`).MatchString(s)
}
```

---

## 🏷️ Validação com go-playground/validator

### Instalação

```bash
go get github.com/go-playground/validator/v10
```

### Setup do Validator

**`pkg/validation/validator.go`**

```go
package validation

import (
    "fmt"
    "reflect"
    "strings"

    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()

    // Usar nome JSON nos erros em vez do nome do campo
    validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })

    // Registrar validações customizadas
    _ = validate.RegisterValidation("username", validateUsername)
    _ = validate.RegisterValidation("strong_password", validateStrongPassword)
}

// Validate valida uma struct
func Validate(s interface{}) error {
    if err := validate.Struct(s); err != nil {
        return FormatValidationErrors(err)
    }
    return nil
}

// FormatValidationErrors formata erros de validação
func FormatValidationErrors(err error) error {
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        var errors []string
        for _, e := range validationErrors {
            errors = append(errors, formatFieldError(e))
        }
        return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
    }
    return err
}

func formatFieldError(e validator.FieldError) string {
    field := e.Field()
    
    switch e.Tag() {
    case "required":
        return fmt.Sprintf("%s is required", field)
    case "email":
        return fmt.Sprintf("%s must be a valid email", field)
    case "min":
        return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
    case "max":
        return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
    case "gte":
        return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
    case "lte":
        return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
    case "username":
        return fmt.Sprintf("%s must be 3-50 characters, alphanumeric and underscores only", field)
    case "strong_password":
        return fmt.Sprintf("%s must be at least 8 characters with uppercase, lowercase, and number", field)
    default:
        return fmt.Sprintf("%s failed validation: %s", field, e.Tag())
    }
}

// Custom validators
func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    if len(username) < 3 || len(username) > 50 {
        return false
    }
    // Apenas letras, números e underscore
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
    return matched
}

func validateStrongPassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    if len(password) < 8 {
        return false
    }
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    return hasUpper && hasLower && hasNumber
}
```

### Usando Tags de Validação

**`internal/modules/user/dto/requests.go`**

```go
package dto

type CreateUserRequest struct {
    Username string `json:"username" validate:"required,username"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,strong_password"`
}

type UpdateUserRequest struct {
    Username *string `json:"username,omitempty" validate:"omitempty,username"`
    Email    *string `json:"email,omitempty" validate:"omitempty,email"`
}

type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=100"`
    Description string  `json:"description" validate:"required,min=10,max=1000"`
    CategoryID  string  `json:"category_id" validate:"required,uuid4"`
    Price       float64 `json:"price" validate:"required,gte=0"`
    Stock       int     `json:"stock" validate:"required,gte=0"`
}

type CreateOrderRequest struct {
    UserID string      `json:"user_id" validate:"required,uuid4"`
    Items  []OrderItem `json:"items" validate:"required,min=1,dive"`
}

type OrderItem struct {
    ProductID string `json:"product_id" validate:"required,uuid4"`
    Quantity  int    `json:"quantity" validate:"required,gte=1,lte=100"`
}
```

### Validando no Handler HTTP

**`internal/modules/user/adapters/http/user_handler.go`**

```go
package http

import (
    "net/http"

    "meuApp/internal/modules/user/application/commands"
    "meuApp/internal/modules/user/dto"
    "meuApp/pkg/validation"

    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    commandBus contracts.CommandBus
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest

    // 1. Bind JSON
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
            "details": err.Error(),
        })
        return
    }

    // 2. Validate struct
    if err := validation.Validate(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Validation failed",
            "details": err.Error(),
        })
        return
    }

    // 3. Criar command
    cmd := commands.CreateUserCommand{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    // 4. Executar command
    user, err := h.commandBus.Execute(c.Request.Context(), cmd)
    if err != nil {
        handleCommandError(c, err)
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "data": user,
    })
}
```

---

## 🏛️ Validação de Domínio

### Value Objects com Validação

**`internal/modules/user/domain/email.go`**

```go
package domain

import (
    "fmt"
    "regexp"
    "strings"
)

// Email é um Value Object que representa um email válido
type Email struct {
    value string
}

// NewEmail cria um novo Email validado
func NewEmail(email string) (Email, error) {
    email = strings.TrimSpace(strings.ToLower(email))

    if email == "" {
        return Email{}, fmt.Errorf("email cannot be empty")
    }

    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return Email{}, fmt.Errorf("invalid email format: %s", email)
    }

    return Email{value: email}, nil
}

// Value retorna o valor do email
func (e Email) Value() string {
    return e.value
}

// String implementa fmt.Stringer
func (e Email) String() string {
    return e.value
}

// Equals compara dois emails
func (e Email) Equals(other Email) bool {
    return e.value == other.value
}
```

**`internal/modules/user/domain/password.go`**

```go
package domain

import (
    "fmt"
    "regexp"
)

// Password é um Value Object que representa uma senha forte
type Password struct {
    value string
}

// NewPassword cria uma nova senha validada
func NewPassword(password string) (Password, error) {
    if len(password) < 8 {
        return Password{}, fmt.Errorf("password must be at least 8 characters")
    }

    if len(password) > 128 {
        return Password{}, fmt.Errorf("password too long (max 128 characters)")
    }

    if !hasUppercase(password) {
        return Password{}, fmt.Errorf("password must contain at least one uppercase letter")
    }

    if !hasLowercase(password) {
        return Password{}, fmt.Errorf("password must contain at least one lowercase letter")
    }

    if !hasNumber(password) {
        return Password{}, fmt.Errorf("password must contain at least one number")
    }

    return Password{value: password}, nil
}

// Value retorna o valor da senha
func (p Password) Value() string {
    return p.value
}

func hasUppercase(s string) bool {
    return regexp.MustCompile(`[A-Z]`).MatchString(s)
}

func hasLowercase(s string) bool {
    return regexp.MustCompile(`[a-z]`).MatchString(s)
}

func hasNumber(s string) bool {
    return regexp.MustCompile(`[0-9]`).MatchString(s)
}
```

### Entidades com Invariantes

**`internal/modules/product/domain/product.go`**

```go
package domain

import (
    "fmt"
    "time"
)

type Product struct {
    id          string
    name        string
    description string
    categoryID  string
    price       float64
    stock       int
    createdAt   time.Time
    updatedAt   time.Time
}

// NewProduct cria um novo produto validado
func NewProduct(
    id, name, description, categoryID string,
    price float64,
    stock int,
) (*Product, error) {
    p := &Product{
        id:         id,
        createdAt:  time.Now(),
        updatedAt:  time.Now(),
    }

    // Validar e setar campos
    if err := p.SetName(name); err != nil {
        return nil, err
    }
    if err := p.SetDescription(description); err != nil {
        return nil, err
    }
    if err := p.SetCategoryID(categoryID); err != nil {
        return nil, err
    }
    if err := p.SetPrice(price); err != nil {
        return nil, err
    }
    if err := p.SetStock(stock); err != nil {
        return nil, err
    }

    return p, nil
}

// SetName valida e define o nome
func (p *Product) SetName(name string) error {
    name = strings.TrimSpace(name)
    if name == "" {
        return fmt.Errorf("product name cannot be empty")
    }
    if len(name) < 3 {
        return fmt.Errorf("product name must be at least 3 characters")
    }
    if len(name) > 100 {
        return fmt.Errorf("product name too long (max 100 characters)")
    }
    p.name = name
    p.updatedAt = time.Now()
    return nil
}

// SetPrice valida e define o preço
func (p *Product) SetPrice(price float64) error {
    if price < 0 {
        return fmt.Errorf("price cannot be negative")
    }
    if price > 1000000 {
        return fmt.Errorf("price too high (max 1,000,000)")
    }
    p.price = price
    p.updatedAt = time.Now()
    return nil
}

// SetStock valida e define o estoque
func (p *Product) SetStock(stock int) error {
    if stock < 0 {
        return fmt.Errorf("stock cannot be negative")
    }
    p.stock = stock
    p.updatedAt = time.Now()
    return nil
}

// DecreaseStock diminui o estoque com validação
func (p *Product) DecreaseStock(quantity int) error {
    if quantity <= 0 {
        return fmt.Errorf("quantity must be positive")
    }
    if p.stock < quantity {
        return fmt.Errorf("insufficient stock: have %d, need %d", p.stock, quantity)
    }
    p.stock -= quantity
    p.updatedAt = time.Now()
    return nil
}

// IncreaseStock aumenta o estoque com validação
func (p *Product) IncreaseStock(quantity int) error {
    if quantity <= 0 {
        return fmt.Errorf("quantity must be positive")
    }
    p.stock += quantity
    p.updatedAt = time.Now()
    return nil
}
```

---

## 🎨 Custom Validators

### Validator Complexo: CPF

**`pkg/validation/cpf_validator.go`**

```go
package validation

import (
    "regexp"
    "strconv"

    "github.com/go-playground/validator/v10"
)

func init() {
    validate.RegisterValidation("cpf", validateCPF)
}

func validateCPF(fl validator.FieldLevel) bool {
    cpf := fl.Field().String()
    return isValidCPF(cpf)
}

func isValidCPF(cpf string) bool {
    // Remove caracteres não numéricos
    cpf = regexp.MustCompile(`[^0-9]`).ReplaceAllString(cpf, "")

    // CPF deve ter 11 dígitos
    if len(cpf) != 11 {
        return false
    }

    // Verificar sequências inválidas (111.111.111-11, etc)
    allSame := true
    for i := 1; i < len(cpf); i++ {
        if cpf[i] != cpf[0] {
            allSame = false
            break
        }
    }
    if allSame {
        return false
    }

    // Validar dígitos verificadores
    return validateCPFCheckDigits(cpf)
}

func validateCPFCheckDigits(cpf string) bool {
    // Cálculo do primeiro dígito verificador
    sum := 0
    for i := 0; i < 9; i++ {
        digit, _ := strconv.Atoi(string(cpf[i]))
        sum += digit * (10 - i)
    }
    firstDigit := (sum * 10) % 11
    if firstDigit == 10 {
        firstDigit = 0
    }
    if firstDigit != int(cpf[9]-'0') {
        return false
    }

    // Cálculo do segundo dígito verificador
    sum = 0
    for i := 0; i < 10; i++ {
        digit, _ := strconv.Atoi(string(cpf[i]))
        sum += digit * (11 - i)
    }
    secondDigit := (sum * 10) % 11
    if secondDigit == 10 {
        secondDigit = 0
    }
    return secondDigit == int(cpf[10]-'0')
}
```

### Validator para Datas

**`pkg/validation/date_validator.go`**

```go
package validation

import (
    "time"

    "github.com/go-playground/validator/v10"
)

func init() {
    validate.RegisterValidation("future_date", validateFutureDate)
    validate.RegisterValidation("past_date", validatePastDate)
    validate.RegisterValidation("min_age", validateMinAge)
}

// validateFutureDate valida se a data é futura
func validateFutureDate(fl validator.FieldLevel) bool {
    dateStr := fl.Field().String()
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return false
    }
    return date.After(time.Now())
}

// validatePastDate valida se a data é passada
func validatePastDate(fl validator.FieldLevel) bool {
    dateStr := fl.Field().String()
    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return false
    }
    return date.Before(time.Now())
}

// validateMinAge valida idade mínima
func validateMinAge(fl validator.FieldLevel) bool {
    dateStr := fl.Field().String()
    minAge, err := strconv.Atoi(fl.Param())
    if err != nil {
        return false
    }

    birthDate, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return false
    }

    age := time.Now().Year() - birthDate.Year()
    return age >= minAge
}
```

**Uso:**

```go
type CreateUserRequest struct {
    Username  string `json:"username" validate:"required,username"`
    Email     string `json:"email" validate:"required,email"`
    CPF       string `json:"cpf" validate:"required,cpf"`
    BirthDate string `json:"birth_date" validate:"required,past_date,min_age=18"`
}
```

---

## ⚠️ Error Messages

### Estrutura de Erro de Validação

**`pkg/errors/validation_error.go`**

```go
package errors

import (
    "fmt"
    "strings"
)

// ValidationError representa um erro de validação
type ValidationError struct {
    Field   string
    Message string
}

// Error implementa a interface error
func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// NewValidationError cria um novo erro de validação
func NewValidationError(field, message string) error {
    return &ValidationError{
        Field:   field,
        Message: message,
    }
}

// ValidationErrors representa múltiplos erros de validação
type ValidationErrors struct {
    Errors []ValidationError
}

// Error implementa a interface error
func (e ValidationErrors) Error() string {
    var messages []string
    for _, err := range e.Errors {
        messages = append(messages, err.Error())
    }
    return strings.Join(messages, "; ")
}

// Add adiciona um erro de validação
func (e *ValidationErrors) Add(field, message string) {
    e.Errors = append(e.Errors, ValidationError{
        Field:   field,
        Message: message,
    })
}

// HasErrors verifica se há erros
func (e *ValidationErrors) HasErrors() bool {
    return len(e.Errors) > 0
}

// ToMap converte para map para JSON
func (e *ValidationErrors) ToMap() map[string]string {
    result := make(map[string]string)
    for _, err := range e.Errors {
        result[err.Field] = err.Message
    }
    return result
}
```

### Response de Erro Padronizado

**`pkg/adapters/http/error_response.go`**

```go
package http

import (
    "net/http"

    "meuApp/pkg/errors"

    "github.com/gin-gonic/gin"
)

type ErrorResponse struct {
    Error   string            `json:"error"`
    Message string            `json:"message,omitempty"`
    Fields  map[string]string `json:"fields,omitempty"`
}

func HandleError(c *gin.Context, err error) {
    switch e := err.(type) {
    case *errors.ValidationError:
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:  "validation_error",
            Fields: map[string]string{e.Field: e.Message},
        })

    case *errors.ValidationErrors:
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Error:  "validation_error",
            Fields: e.ToMap(),
        })

    case *errors.NotFoundError:
        c.JSON(http.StatusNotFound, ErrorResponse{
            Error:   "not_found",
            Message: e.Error(),
        })

    case *errors.ConflictError:
        c.JSON(http.StatusConflict, ErrorResponse{
            Error:   "conflict",
            Message: e.Error(),
        })

    default:
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Error:   "internal_error",
            Message: "An unexpected error occurred",
        })
    }
}
```

---

## 📚 Exemplos Práticos

### Exemplo 1: Validação Completa de Command

**`internal/modules/order/application/commands/create_order.go`**

```go
package commands

import (
    "context"
    "fmt"

    "meuApp/internal/modules/order/domain"
    "meuApp/internal/modules/order/ports"
    "meuApp/pkg/contracts"
    "meuApp/pkg/errors"
)

type CreateOrderCommand struct {
    UserID string
    Items  []OrderItemCommand
}

type OrderItemCommand struct {
    ProductID string
    Quantity  int
}

type CreateOrderHandler struct {
    orderRepo    ports.OrderRepository
    productRepo  ports.ProductRepository
    userRepo     ports.UserRepository
    logger       contracts.Logger
}

func (h *CreateOrderHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (*domain.Order, error) {
    // 1. Validação estrutural
    if err := h.validateStructure(cmd); err != nil {
        return nil, err
    }

    // 2. Validação de existência de recursos
    if err := h.validateResources(ctx, cmd); err != nil {
        return nil, err
    }

    // 3. Validação de regras de negócio
    if err := h.validateBusinessRules(ctx, cmd); err != nil {
        return nil, err
    }

    // Continuar com criação do pedido...
}

func (h *CreateOrderHandler) validateStructure(cmd CreateOrderCommand) error {
    validationErrs := &errors.ValidationErrors{}

    if cmd.UserID == "" {
        validationErrs.Add("user_id", "user_id is required")
    }

    if len(cmd.Items) == 0 {
        validationErrs.Add("items", "at least one item is required")
    }

    for i, item := range cmd.Items {
        if item.ProductID == "" {
            validationErrs.Add(fmt.Sprintf("items[%d].product_id", i), "product_id is required")
        }
        if item.Quantity <= 0 {
            validationErrs.Add(fmt.Sprintf("items[%d].quantity", i), "quantity must be positive")
        }
        if item.Quantity > 100 {
            validationErrs.Add(fmt.Sprintf("items[%d].quantity", i), "quantity cannot exceed 100")
        }
    }

    if validationErrs.HasErrors() {
        return validationErrs
    }

    return nil
}

func (h *CreateOrderHandler) validateResources(ctx context.Context, cmd CreateOrderCommand) error {
    // Verificar se usuário existe
    user, err := h.userRepo.GetByID(ctx, cmd.UserID)
    if err != nil || user == nil {
        return errors.NewValidationError("user_id", "user not found")
    }

    // Verificar se produtos existem
    for i, item := range cmd.Items {
        product, err := h.productRepo.GetByID(ctx, item.ProductID)
        if err != nil || product == nil {
            return errors.NewValidationError(
                fmt.Sprintf("items[%d].product_id", i),
                "product not found",
            )
        }
    }

    return nil
}

func (h *CreateOrderHandler) validateBusinessRules(ctx context.Context, cmd CreateOrderCommand) error {
    validationErrs := &errors.ValidationErrors{}

    // Verificar estoque de cada produto
    for i, item := range cmd.Items {
        product, _ := h.productRepo.GetByID(ctx, item.ProductID)
        
        if product.Stock() < item.Quantity {
            validationErrs.Add(
                fmt.Sprintf("items[%d].quantity", i),
                fmt.Sprintf("insufficient stock: available %d, requested %d", product.Stock(), item.Quantity),
            )
        }
    }

    // Verificar duplicatas (mesmo produto múltiplas vezes)
    productIDs := make(map[string]bool)
    for i, item := range cmd.Items {
        if productIDs[item.ProductID] {
            validationErrs.Add(
                fmt.Sprintf("items[%d].product_id", i),
                "duplicate product in order",
            )
        }
        productIDs[item.ProductID] = true
    }

    if validationErrs.HasErrors() {
        return validationErrs
    }

    return nil
}
```

### Exemplo 2: Cross-Field Validation

```go
type UpdatePasswordRequest struct {
    OldPassword        string `json:"old_password" validate:"required"`
    NewPassword        string `json:"new_password" validate:"required,strong_password"`
    ConfirmNewPassword string `json:"confirm_new_password" validate:"required"`
}

func (r UpdatePasswordRequest) Validate() error {
    if r.NewPassword != r.ConfirmNewPassword {
        return errors.NewValidationError("confirm_new_password", "passwords do not match")
    }
    if r.OldPassword == r.NewPassword {
        return errors.NewValidationError("new_password", "new password must be different from old password")
    }
    return nil
}
```

---

## 📚 Próximos Passos

- **[Error Handling](19-error-handling.md)** - Tratamento de erros
- **[Testing](21-testing-strategy.md)** - Testar validações
- **[Best Practices](23-best-practices.md)** - Boas práticas

---

**[⬅️ Configuration](16-configuration-bootstrap.md)** | **[Índice](README.md)** | **[Error Handling ➡️](19-error-handling.md)**
