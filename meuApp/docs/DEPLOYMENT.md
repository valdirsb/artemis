# 🚀 Guia de Deployment e Configuração

> **Guia completo** para configurar, executar e fazer deploy da aplicação Artemis em diferentes ambientes

## 📋 Índice

- [Requisitos](#requisitos)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Configuração Local](#configuração-local)
- [Docker e Docker Compose](#docker-e-docker-compose)
- [Build e Execução](#build-e-execução)
- [Ambientes](#ambientes)
- [Troubleshooting](#troubleshooting)

---

## 📦 Requisitos

### Desenvolvimento Local

- **Go**: 1.24 ou superior
- **MySQL**: 8.0 ou superior
- **Make**: Para comandos de build
- **Protocol Buffers Compiler**: Para gerar código gRPC
- **Git**: Para controle de versão

### Produção

- **Docker**: 20.10+ (recomendado)
- **Docker Compose**: 2.0+
- **Kubernetes**: 1.25+ (opcional, para orquestração)

---

## 🔐 Variáveis de Ambiente

### Arquivo `.env`

Crie um arquivo `.env` na raiz do projeto:

```bash
# ======================
# APPLICATION
# ======================
APP_ENV=development                    # development | staging | production
APP_NAME=meuApp
APP_VERSION=1.0.0
APP_DEBUG=true                         # Habilita logs debug

# ======================
# HTTP SERVER
# ======================
HTTP_PORT=8080
HTTP_READ_TIMEOUT=30s
HTTP_WRITE_TIMEOUT=30s
HTTP_IDLE_TIMEOUT=120s

# ======================
# gRPC SERVER
# ======================
GRPC_PORT=50051
GRPC_MAX_RECV_MSG_SIZE=4194304        # 4MB
GRPC_MAX_SEND_MSG_SIZE=4194304        # 4MB

# ======================
# DATABASE (MySQL)
# ======================
DB_HOST=localhost
DB_PORT=3306
DB_USER=app_user
DB_PASSWORD=app_password
DB_NAME=app_db
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=3600             # Segundos

# ======================
# LOGGING
# ======================
LOG_LEVEL=info                         # debug | info | warn | error
LOG_FORMAT=json                        # json | text
LOG_OUTPUT=stdout                      # stdout | file
LOG_FILE_PATH=/var/log/meuapp/app.log

# ======================
# CORS
# ======================
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Request-ID
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400                     # 24 horas

# ======================
# SECURITY
# ======================
JWT_SECRET=your-256-bit-secret-key-here-change-in-production
JWT_EXPIRATION=3600                    # 1 hora em segundos
BCRYPT_COST=10                         # 4-31 (maior = mais seguro mas mais lento)

# ======================
# RATE LIMITING
# ======================
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_SECOND=100
RATE_LIMIT_BURST=200

# ======================
# METRICS & OBSERVABILITY
# ======================
METRICS_ENABLED=true
METRICS_PORT=9090
TRACING_ENABLED=false
TRACING_ENDPOINT=http://jaeger:14268/api/traces

# ======================
# EXTERNAL SERVICES
# ======================
EMAIL_SERVICE_URL=http://email-service:8080
EMAIL_FROM=noreply@meuapp.com
PAYMENT_GATEWAY_URL=https://api.payment-provider.com
PAYMENT_API_KEY=your-payment-api-key
```

### Variáveis por Ambiente

#### `.env.development`

```bash
APP_ENV=development
APP_DEBUG=true
DB_HOST=localhost
LOG_LEVEL=debug
METRICS_ENABLED=true
```

#### `.env.production`

```bash
APP_ENV=production
APP_DEBUG=false
DB_HOST=mysql-prod.internal
DB_MAX_OPEN_CONNS=200
LOG_LEVEL=info
LOG_FORMAT=json
METRICS_ENABLED=true
TRACING_ENABLED=true
RATE_LIMIT_REQUESTS_PER_SECOND=1000
```

---

## 💻 Configuração Local

### 1. Clone o Repositório

```bash
git clone https://github.com/your-org/artemis.git
cd artemis/meuApp
```

### 2. Configure Variáveis de Ambiente

```bash
# Copie o template
cp .env.example .env

# Edite com suas configurações
vim .env  # ou seu editor preferido
```

### 3. Configure o MySQL

```bash
# Criar database
mysql -u root -p

CREATE DATABASE app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'app_user'@'localhost' IDENTIFIED BY 'app_password';
GRANT ALL PRIVILEGES ON app_db.* TO 'app_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 4. Instale Dependências

```bash
# Download de módulos Go
go mod download

# Instale ferramentas de desenvolvimento
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 5. Gere Código gRPC (se necessário)

```bash
make proto
# ou
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/*.proto
```

### 6. Execute a Aplicação

```bash
# Via Go
go run main.go

# Ou via Make
make run

# Ou compile e execute
make build
./meuApp
```

---

## 🐳 Docker e Docker Compose

### Dockerfile Multi-Stage

**Arquivo:** `Dockerfile`

```dockerfile
# ======================
# Stage 1: Builder
# ======================
FROM golang:1.24-alpine AS builder

# Instalar dependências de build
RUN apk add --no-cache git make protobuf-dev

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Download de dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o meuApp .

# ======================
# Stage 2: Runtime
# ======================
FROM alpine:latest

# Instalar CA certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Definir diretório de trabalho
WORKDIR /app

# Copiar binário do builder
COPY --from=builder /app/meuApp .

# Copiar arquivo de configuração (opcional)
COPY --from=builder /app/.env.example .env

# Mudar ownership
RUN chown -R appuser:appuser /app

# Usar usuário não-root
USER appuser

# Expor portas
EXPOSE 8080 50051 9090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando de execução
ENTRYPOINT ["./meuApp"]
```

### Docker Compose

**Arquivo:** `docker-compose.yml`

```yaml
version: '3.9'

services:
  # ======================
  # Aplicação
  # ======================
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: artemis-app
    restart: unless-stopped
    ports:
      - "8080:8080"   # HTTP
      - "50051:50051" # gRPC
      - "9090:9090"   # Metrics
    environment:
      APP_ENV: development
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: app_user
      DB_PASSWORD: app_password
      DB_NAME: app_db
    depends_on:
      mysql:
        condition: service_healthy
    networks:
      - artemis-network
    volumes:
      - ./logs:/var/log/meuapp
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s

  # ======================
  # MySQL Database
  # ======================
  mysql:
    image: mysql:8.0
    container_name: artemis-mysql
    restart: unless-stopped
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: app_db
      MYSQL_USER: app_user
      MYSQL_PASSWORD: app_password
    volumes:
      - mysql-data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - artemis-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-proot_password"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s

  # ======================
  # Redis (opcional - para cache)
  # ======================
  redis:
    image: redis:7-alpine
    container_name: artemis-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    networks:
      - artemis-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3

  # ======================
  # Prometheus (opcional - métricas)
  # ======================
  prometheus:
    image: prom/prometheus:latest
    container_name: artemis-prometheus
    restart: unless-stopped
    ports:
      - "9091:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
    networks:
      - artemis-network

  # ======================
  # Grafana (opcional - dashboards)
  # ======================
  grafana:
    image: grafana/grafana:latest
    container_name: artemis-grafana
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes:
      - grafana-data:/var/lib/grafana
      - ./monitoring/grafana/dashboards:/etc/grafana/provisioning/dashboards
    networks:
      - artemis-network
    depends_on:
      - prometheus

networks:
  artemis-network:
    driver: bridge

volumes:
  mysql-data:
  prometheus-data:
  grafana-data:
```

### Comandos Docker Compose

```bash
# Iniciar todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f app

# Parar serviços
docker-compose down

# Rebuild e restart
docker-compose up -d --build

# Executar apenas app + mysql
docker-compose up -d app mysql

# Limpar tudo (⚠️ CUIDADO: apaga volumes)
docker-compose down -v
```

---

## 🔨 Build e Execução

### Makefile

**Arquivo:** `Makefile`

```makefile
.PHONY: help build run test clean proto docker-build docker-run

# Variáveis
APP_NAME=meuApp
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

help: ## Mostra esta ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-deps: ## Instala dependências
	go mod download
	go mod tidy

proto: ## Gera código gRPC a partir dos .proto
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/*.proto

build: ## Compila a aplicação
	go build ${LDFLAGS} -o ${APP_NAME} .

build-linux: ## Compila para Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${APP_NAME}-linux .

run: ## Executa a aplicação
	go run ${LDFLAGS} main.go

dev: ## Executa com hot reload (requer air)
	air

test: ## Executa testes
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Gera relatório de cobertura
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Executa linter
	golangci-lint run ./...

fmt: ## Formata código
	go fmt ./...
	goimports -w .

clean: ## Limpa builds
	rm -f ${APP_NAME} ${APP_NAME}-linux
	rm -f coverage.out coverage.html
	go clean

docker-build: ## Build da imagem Docker
	docker build -t ${APP_NAME}:${VERSION} .
	docker tag ${APP_NAME}:${VERSION} ${APP_NAME}:latest

docker-run: ## Executa container Docker
	docker run -p 8080:8080 -p 50051:50051 --env-file .env ${APP_NAME}:latest

docker-compose-up: ## Inicia Docker Compose
	docker-compose up -d

docker-compose-down: ## Para Docker Compose
	docker-compose down

migrate-up: ## Executa migrations (up)
	@echo "Migrations são automáticas via AutoMigrate"
	go run main.go migrate

seed: ## Executa seed do banco
	go run scripts/seed.go

.DEFAULT_GOAL := help
```

### Comandos Úteis

```bash
# Ver ajuda
make help

# Build
make build

# Executar
make run

# Testes
make test
make test-coverage

# Docker
make docker-build
make docker-run

# Docker Compose
make docker-compose-up
make docker-compose-down

# Linting
make lint

# Format code
make fmt
```

---

## 🌍 Ambientes

### Development

```bash
# Usar .env.development
export APP_ENV=development
go run main.go
```

### Staging

```bash
# Usar .env.staging
export APP_ENV=staging
./meuApp
```

### Production

```bash
# Usar .env.production
export APP_ENV=production
./meuApp

# Ou via Docker
docker-compose -f docker-compose.prod.yml up -d
```

### Configuração por Ambiente

**Arquivo:** `pkg/config/config.go`

```go
func LoadConfig() (*Config, error) {
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }
    
    // Carregar .env.{environment}
    envFile := fmt.Sprintf(".env.%s", env)
    if err := godotenv.Load(envFile); err != nil {
        // Fallback para .env
        godotenv.Load()
    }
    
    // Carregar configuração
    config := &Config{
        Environment: env,
        // ... resto da config
    }
    
    return config, nil
}
```

---

## 🔍 Health Checks

### Endpoint HTTP

```bash
# Basic health check
curl http://localhost:8080/health

# Response
{
  "status": "ok",
  "timestamp": "2025-10-18T22:06:43Z",
  "version": "1.0.0",
  "database": "connected",
  "registry": {
    "http_handlers": 3,
    "grpc_services": 3,
    "repositories": 3,
    "app_services": 3
  }
}
```

### Kubernetes Probes

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 3
  failureThreshold: 3
```

---

## 🐛 Troubleshooting

### Problema: Aplicação não conecta ao MySQL

**Sintoma:**
```
Error: failed to connect to database
```

**Soluções:**
```bash
# 1. Verificar se MySQL está rodando
mysql -u app_user -p -h localhost

# 2. Verificar variáveis de ambiente
echo $DB_HOST
echo $DB_USER

# 3. Testar conexão
mysql -u app_user -papp_password -h localhost -e "SHOW DATABASES;"

# 4. Ver logs do MySQL
docker-compose logs mysql
```

### Problema: Porta já em uso

**Sintoma:**
```
Error: bind: address already in use
```

**Solução:**
```bash
# Encontrar processo usando a porta
lsof -i :8080

# Matar processo
kill -9 <PID>

# Ou mudar porta no .env
HTTP_PORT=8081
```

### Problema: gRPC não funciona

**Sintoma:**
```
Error: failed to start gRPC server
```

**Soluções:**
```bash
# 1. Verificar porta
netstat -tulpn | grep 50051

# 2. Regenerar proto
make proto

# 3. Testar com grpcurl
grpcurl -plaintext localhost:50051 list
```

### Problema: Migrations não executam

**Sintoma:**
```
Table doesn't exist
```

**Solução:**
```bash
# Migrations são automáticas
# Verificar se AutoMigrate está sendo chamado
go run main.go

# Ou manualmente
mysql -u app_user -p app_db < scripts/schema.sql
```

---

## 📚 Recursos Adicionais

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Arquitetura do sistema
- [MODULE_CREATION_GUIDE.md](./MODULE_CREATION_GUIDE.md) - Como criar módulos
- [API.md](./API.md) - Documentação da API
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Go Deployment](https://golang.org/doc/articles/wiki/)

---

**Mantido por:** Time Artemis  
**Última atualização:** 18 de Outubro de 2025
