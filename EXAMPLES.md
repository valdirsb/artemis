# 📚 Exemplos de Uso do Artemis Framework

## Criando um E-commerce Simples

### 1. Criar o projeto
```bash
./bin/artemis new ecommerce-api
cd ecommerce-api
```

### 2. Gerar módulos principais
```bash
./bin/artemis make module users
./bin/artemis make module products  
./bin/artemis make module orders
./bin/artemis make module payments
```

### 3. Gerar migrations
```bash
./bin/artemis make migration create_users_table
./bin/artemis make migration create_products_table
./bin/artemis make migration create_orders_table
```

### 4. Estrutura gerada
```
ecommerce-api/
├── internal/modules/
│   ├── users/
│   │   ├── domain/         # Entidades e regras de negócio
│   │   ├── service/        # Casos de uso
│   │   ├── repository/     # Persistência  
│   │   └── handler/        # HTTP handlers
│   ├── products/
│   ├── orders/
│   └── payments/
├── migrations/             # SQL migrations
└── cmd/server/            # Entry point
```

## Comandos Úteis

### Desenvolvimento
```bash
# Iniciar servidor
./bin/artemis serve --port 8080

# Verificar help
./bin/artemis --help
./bin/artemis make --help
```

### Estrutura de um Módulo

Cada módulo gerado segue a arquitetura hexagonal:

- **domain/**: Entidades e interfaces de repositório
- **ports/**: Interfaces de casos de uso (primary e secondary ports)  
- **service/**: Implementação dos casos de uso
- **repository/**: Implementação de persistência
- **handler/**: Controllers HTTP
- **adapters/**: Adaptadores externos

### Próximos Passos

Após gerar os módulos, você pode:

1. Implementar as regras de negócio nas entidades
2. Definir os casos de uso nos services  
3. Implementar persistência nos repositories
4. Criar endpoints nos handlers
5. Configurar injeção de dependências no bootstrap
