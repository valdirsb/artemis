# ⚙️ Configuration & Bootstrap

## 📋 Índice
- [Visão Geral](#visão-geral)
- [Framework.yaml](#frameworkyaml)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Processo de Bootstrap](#processo-de-bootstrap)
- [Configuração de Módulos](#configuração-de-módulos)
- [Ambientes (Dev, Staging, Prod)](#ambientes)
- [Secrets e Segurança](#secrets-e-segurança)
- [Exemplos Práticos](#exemplos-práticos)

---

## 🎯 Visão Geral

O Artemis Framework usa uma abordagem de **configuração declarativa** através do arquivo `framework.yaml` combinado com **variáveis de ambiente** para configurações sensíveis.

### Fluxo de Inicialização

```
┌────────────────────────────────────────────────┐
│              main.go                           │
│   Ponto de entrada da aplicação               │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│         1. Load Configuration                  │
│   • framework.yaml                             │
│   • .env file                                  │
│   • Environment variables                      │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│         2. Bootstrap Framework                 │
│   • Initialize database                        │
│   • Create event bus                           │
│   • Create command/query buses                 │
│   • Setup logger                               │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│         3. Register Modules                    │
│   • Auto-discover modules                      │
│   • Initialize each module                     │
│   • Register handlers                          │
│   • Register routes                            │
└──────────────────┬─────────────────────────────┘
                   │
                   ▼
┌────────────────────────────────────────────────┐
│         4. Start Server                        │
│   • HTTP server (Gin)                          │
│   • gRPC server (optional)                     │
│   • Health checks                              │
└────────────────────────────────────────────────┘
```

---

## 📄 Framework.yaml

### Estrutura Completa

**`framework.yaml`**

```yaml
# Informações do Framework
framework:
  name: "meuApp"
  version: "1.0.0"
  description: "Application built with Artemis Framework"

# Features do Core
core:
  http_server: true
  dependency_injection: true
  event_system: true
  logging: true
  database_migrations: true

# Configuração de Banco de Dados
database:
  mysql: true
  postgresql: false
  mongodb: false
  sqlite: false

# Cache
cache:
  redis: true
  memcached: false
  in_memory: true  # Fallback

# Protocolos
protocols:
  http: true
  grpc: true
  websockets: false
  graphql: false

# Integrações
integrations:
  jwt: true
  oauth2: false
  firebase_auth: false
  stripe: false
  sendgrid: false
  aws_s3: false

# Segurança
security:
  cors: true
  rate_limiting: true
  csrf_protection: false

# Desenvolvimento
development:
  swagger: true
  debug_mode: true
  hot_reload: false

# Módulos da Aplicação
modules:
  user: true
  product: true
  order: true
  category: true
  payment: false
  notification: false

# Observabilidade
observability:
  prometheus: false
  jaeger: false
  sentry: false
```

### Seções Detalhadas

#### 1. Framework Info

```yaml
framework:
  name: "meuApp"              # Nome da aplicação
  version: "1.0.0"            # Versão semântica
  description: "My awesome app"
  environment: "development"   # development, staging, production
```

#### 2. Core Features

```yaml
core:
  http_server: true            # Habilita servidor HTTP (Gin)
  dependency_injection: true   # Habilita DI Container
  event_system: true          # Habilita Event Bus
  logging: true               # Habilita sistema de logs
  database_migrations: true   # Auto-run migrations
```

#### 3. Database

```yaml
database:
  mysql: true                 # Usar MySQL
  postgresql: false           # Não usar PostgreSQL
  mongodb: false              # Não usar MongoDB
  sqlite: false               # Não usar SQLite
  
  # Configurações (em .env)
  # DB_HOST=localhost
  # DB_PORT=3306
  # DB_NAME=artemis
  # DB_USER=root
  # DB_PASSWORD=secret
```

#### 4. Modules

```yaml
modules:
  user: true        # Módulo de usuários habilitado
  product: true     # Módulo de produtos habilitado
  order: true       # Módulo de pedidos habilitado
  category: false   # Módulo de categorias desabilitado
```

---

## 🔐 Variáveis de Ambiente

### Arquivo .env

**`.env`** (nunca commitar!)

```bash
# Application
APP_NAME=meuApp
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# Database
DB_HOST=localhost
DB_PORT=3306
DB_NAME=artemis_dev
DB_USER=root
DB_PASSWORD=secret123
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=720h

# Email (SendGrid)
SENDGRID_API_KEY=SG.xxxxx
EMAIL_FROM=noreply@meuapp.com

# AWS S3
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=AKIAXXXXX
AWS_SECRET_ACCESS_KEY=xxxxx
AWS_S3_BUCKET=meuapp-uploads

# Stripe
STRIPE_PUBLIC_KEY=pk_test_xxxxx
STRIPE_SECRET_KEY=sk_test_xxxxx

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
LOG_OUTPUT=stdout
```

### Carregamento de .env

**`pkg/config/config.go`**

```go
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    App      AppConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    Email    EmailConfig
    AWS      AWSConfig
    CORS     CORSConfig
}

type AppConfig struct {
    Name        string
    Environment string
    Port        int
    Debug       bool
}

type DatabaseConfig struct {
    Host            string
    Port            int
    Name            string
    User            string
    Password        string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

type JWTConfig struct {
    Secret            string
    Expiration        time.Duration
    RefreshExpiration time.Duration
}

type EmailConfig struct {
    SendGridAPIKey string
    FromAddress    string
}

type AWSConfig struct {
    Region          string
    AccessKeyID     string
    SecretAccessKey string
    S3Bucket        string
}

type CORSConfig struct {
    AllowedOrigins []string
    AllowedMethods []string
    AllowedHeaders []string
}

// LoadConfig carrega configurações do .env
func LoadConfig() (*Config, error) {
    // Carregar .env (ignora erro se não existir)
    _ = godotenv.Load()

    config := &Config{
        App: AppConfig{
            Name:        getEnv("APP_NAME", "artemis"),
            Environment: getEnv("APP_ENV", "development"),
            Port:        getEnvAsInt("APP_PORT", 8080),
            Debug:       getEnvAsBool("APP_DEBUG", true),
        },
        Database: DatabaseConfig{
            Host:            getEnv("DB_HOST", "localhost"),
            Port:            getEnvAsInt("DB_PORT", 3306),
            Name:            getEnv("DB_NAME", "artemis"),
            User:            getEnv("DB_USER", "root"),
            Password:        getEnv("DB_PASSWORD", ""),
            MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
            MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
            ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvAsInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvAsInt("REDIS_DB", 0),
        },
        JWT: JWTConfig{
            Secret:            getEnv("JWT_SECRET", "change-this-secret"),
            Expiration:        getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),
            RefreshExpiration: getEnvAsDuration("JWT_REFRESH_EXPIRATION", 720*time.Hour),
        },
        Email: EmailConfig{
            SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
            FromAddress:    getEnv("EMAIL_FROM", "noreply@example.com"),
        },
        AWS: AWSConfig{
            Region:          getEnv("AWS_REGION", "us-east-1"),
            AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
            SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
            S3Bucket:        getEnv("AWS_S3_BUCKET", ""),
        },
        CORS: CORSConfig{
            AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
            AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE"}),
            AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
        },
    }

    // Validar configurações críticas
    if err := config.Validate(); err != nil {
        return nil, err
    }

    return config, nil
}

// Validate valida configurações obrigatórias
func (c *Config) Validate() error {
    if c.Database.Host == "" {
        return fmt.Errorf("DB_HOST is required")
    }
    if c.Database.Name == "" {
        return fmt.Errorf("DB_NAME is required")
    }
    if c.JWT.Secret == "change-this-secret" && c.App.Environment == "production" {
        return fmt.Errorf("JWT_SECRET must be changed in production")
    }
    return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
    valueStr := getEnv(key, "")
    if value, err := strconv.ParseBool(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    valueStr := getEnv(key, "")
    if value, err := time.ParseDuration(valueStr); err == nil {
        return value
    }
    return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    return strings.Split(valueStr, ",")
}
```

---

## 🚀 Processo de Bootstrap

### main.go

**`main.go`**

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "meuApp/internal/bootstrap"
    "meuApp/pkg/config"

    "github.com/gin-gonic/gin"
)

func main() {
    log.Println("🚀 Starting Artemis Application...")

    // 1. Load configuration
    appConfig, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    log.Printf("📝 Config loaded: %s (%s)", appConfig.App.Name, appConfig.App.Environment)

    // 2. Bootstrap framework
    _, registry, framework, err := bootstrap.FrameworkBootstrapWithRegistry("framework.yaml")
    if err != nil {
        log.Fatalf("Failed to bootstrap: %v", err)
    }

    log.Printf("✅ Framework initialized with %d modules", registry.ModuleCount())

    // 3. Setup HTTP router
    router := setupRouter(appConfig, registry, framework)

    // 4. Start server
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", appConfig.App.Port),
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // 5. Graceful shutdown
    go func() {
        log.Printf("🌐 Server listening on :%d", appConfig.App.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("🛑 Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Server forced to shutdown: %v", err)
    }

    log.Println("✅ Server stopped gracefully")
}

func setupRouter(
    appConfig *config.Config,
    registry *container.ModuleRegistry,
    framework *framework.Framework,
) *gin.Engine {
    // Gin mode
    if !appConfig.App.Debug {
        gin.SetMode(gin.ReleaseMode)
    }

    router := gin.New()
    router.Use(gin.Logger(), gin.Recovery())

    // CORS
    if len(appConfig.CORS.AllowedOrigins) > 0 {
        router.Use(corsMiddleware(appConfig.CORS))
    }

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "app":    appConfig.App.Name,
            "env":    appConfig.App.Environment,
        })
    })

    // API routes
    api := router.Group("/api/v1")
    registry.RegisterHTTPRoutes(api)

    return router
}
```

### Bootstrap Package

**`internal/bootstrap/bootstrap_registry.go`**

```go
package bootstrap

import (
    "fmt"
    "log"

    "meuApp/internal/modules/user"
    "meuApp/internal/modules/product"
    "meuApp/internal/modules/order"
    "meuApp/pkg/container"
    "meuApp/pkg/framework"
    "meuApp/pkg/config"

    "gorm.io/gorm"
)

func FrameworkBootstrapWithRegistry(
    frameworkConfigPath string,
) (*gorm.DB, *container.ModuleRegistry, *framework.Framework, error) {
    // 1. Load framework config
    frameworkConfig, err := framework.LoadConfig(frameworkConfigPath)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to load framework config: %w", err)
    }

    // 2. Load app config
    appConfig, err := config.LoadConfig()
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to load app config: %w", err)
    }

    // 3. Initialize framework components
    fw := framework.New(frameworkConfig)

    // 4. Setup database
    db, err := setupDatabase(appConfig)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to setup database: %w", err)
    }

    // 5. Create module registry
    registry := container.NewModuleRegistry(
        db,
        fw.GetCommandBus(),
        fw.GetQueryBus(),
        fw.GetEventBus(),
        fw.GetLogger(),
    )

    // 6. Register modules
    if err := registerModules(registry, frameworkConfig); err != nil {
        return nil, nil, nil, fmt.Errorf("failed to register modules: %w", err)
    }

    // 7. Initialize all modules
    if err := registry.InitializeAll(); err != nil {
        return nil, nil, nil, fmt.Errorf("failed to initialize modules: %w", err)
    }

    log.Printf("✅ Bootstrap complete: %d modules registered", registry.ModuleCount())

    return db, registry, fw, nil
}

func setupDatabase(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Name,
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

    return db, nil
}

func registerModules(
    registry *container.ModuleRegistry,
    frameworkConfig *framework.Config,
) error {
    // Register modules based on framework.yaml
    if frameworkConfig.Modules["user"] {
        userModule := user.NewUserModule(
            registry.GetDB(),
            registry.GetCommandBus(),
            registry.GetQueryBus(),
            registry.GetEventBus(),
        )
        registry.Register(userModule)
        log.Println("✅ User module registered")
    }

    if frameworkConfig.Modules["product"] {
        productModule := product.NewProductModule(
            registry.GetDB(),
            registry.GetCommandBus(),
            registry.GetQueryBus(),
            registry.GetEventBus(),
        )
        registry.Register(productModule)
        log.Println("✅ Product module registered")
    }

    if frameworkConfig.Modules["order"] {
        orderModule := order.NewOrderModule(
            registry.GetDB(),
            registry.GetCommandBus(),
            registry.GetQueryBus(),
            registry.GetEventBus(),
        )
        registry.Register(orderModule)
        log.Println("✅ Order module registered")
    }

    return nil
}
```

---

## 🏭 Ambientes

### Configurações por Ambiente

**`.env.development`**
```bash
APP_ENV=development
APP_DEBUG=true
DB_HOST=localhost
DB_NAME=artemis_dev
LOG_LEVEL=debug
```

**`.env.staging`**
```bash
APP_ENV=staging
APP_DEBUG=false
DB_HOST=staging-db.example.com
DB_NAME=artemis_staging
LOG_LEVEL=info
```

**`.env.production`**
```bash
APP_ENV=production
APP_DEBUG=false
DB_HOST=prod-db.example.com
DB_NAME=artemis_prod
LOG_LEVEL=warn
```

### Carregar por Ambiente

```go
func LoadConfigForEnvironment(env string) (*Config, error) {
    // Tentar carregar .env específico do ambiente
    envFile := fmt.Sprintf(".env.%s", env)
    if _, err := os.Stat(envFile); err == nil {
        _ = godotenv.Load(envFile)
    } else {
        // Fallback para .env padrão
        _ = godotenv.Load()
    }

    return LoadConfig()
}
```

---

## 🔒 Secrets e Segurança

### Boas Práticas

1. **Nunca commitar secrets**

```bash
# .gitignore
.env
.env.*
!.env.example
secrets/
*.key
*.pem
```

2. **Usar .env.example como template**

```bash
# .env.example
APP_NAME=meuApp
APP_ENV=development
DB_HOST=localhost
DB_PASSWORD=change-this
JWT_SECRET=change-this-secret
```

3. **Validar secrets em produção**

```go
func (c *Config) ValidateProduction() error {
    if c.App.Environment != "production" {
        return nil
    }

    if c.JWT.Secret == "change-this-secret" {
        return errors.New("JWT_SECRET must be set in production")
    }

    if c.Database.Password == "" {
        return errors.New("DB_PASSWORD must be set in production")
    }

    return nil
}
```

4. **Usar secrets managers em produção**

```go
// AWS Secrets Manager
func loadFromAWSSecretsManager(secretName string) (string, error) {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess)

    input := &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    }

    result, err := svc.GetSecretValue(input)
    if err != nil {
        return "", err
    }

    return *result.SecretString, nil
}

// Uso
func LoadConfig() (*Config, error) {
    if os.Getenv("APP_ENV") == "production" {
        dbPassword, _ := loadFromAWSSecretsManager("prod/db/password")
        os.Setenv("DB_PASSWORD", dbPassword)
    }

    return loadConfigFromEnv()
}
```

---

## 📚 Exemplos Práticos

### Exemplo 1: Multi-Database

```yaml
# framework.yaml
database:
  mysql: true
  postgresql: true  # Usar ambos
  redis: true
```

```go
// bootstrap
func setupDatabases(cfg *config.Config) (*Databases, error) {
    dbs := &Databases{}

    // MySQL (principal)
    if cfg.MySQL.Enabled {
        mysql, err := setupMySQL(cfg.MySQL)
        if err != nil {
            return nil, err
        }
        dbs.MySQL = mysql
    }

    // PostgreSQL (analytics)
    if cfg.PostgreSQL.Enabled {
        postgres, err := setupPostgreSQL(cfg.PostgreSQL)
        if err != nil {
            return nil, err
        }
        dbs.PostgreSQL = postgres
    }

    // Redis (cache)
    if cfg.Redis.Enabled {
        redis, err := setupRedis(cfg.Redis)
        if err != nil {
            return nil, err
        }
        dbs.Redis = redis
    }

    return dbs, nil
}
```

### Exemplo 2: Feature Flags

```yaml
# framework.yaml
features:
  new_checkout: false
  ai_recommendations: true
  beta_features: false
```

```go
type FeatureFlags struct {
    NewCheckout       bool
    AIRecommendations bool
    BetaFeatures      bool
}

func LoadFeatureFlags(frameworkConfig *framework.Config) *FeatureFlags {
    return &FeatureFlags{
        NewCheckout:       frameworkConfig.Features["new_checkout"],
        AIRecommendations: frameworkConfig.Features["ai_recommendations"],
        BetaFeatures:      frameworkConfig.Features["beta_features"],
    }
}

// Uso no handler
if featureFlags.NewCheckout {
    return h.newCheckoutFlow(ctx, cmd)
}
return h.legacyCheckoutFlow(ctx, cmd)
```

---

## 📚 Próximos Passos

- **[Módulos](06-modules-system.md)** - Como criar e registrar módulos
- **[Testes](21-testing-strategy.md)** - Testar configurações
- **[Deployment](DEPLOYMENT.md)** - Deploy em produção

---

**[⬅️ Eventos](15-working-with-events.md)** | **[Índice](README.md)** | **[Paginação ➡️](17-pagination.md)**
