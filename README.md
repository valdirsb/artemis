# 🏹 Artemis Go Framework

![Logo do Projeto](logo.png)

**Build monoliths, modularly.**

Artemis é um framework modular para Go que traz uma estrutura inspirada no Laravel,
com foco em produtividade e organização de código.

## ✨ Recursos
- Estrutura modular de alto desempenho.
- CLI poderoso (`artemis make:module users`).
- Suporte a múltiplas tecnologias (REST, gRPC, GraphQL, etc).
- Facilidade de configuração e extensibilidade.

## 🚀 Instalação
```bash
go install github.com/seunome/artemis@latest
```

## 🧩 Criando um novo projeto
```bash
artemis new myapp
```

## ⚙️ Gerando módulos
```bash
cd myapp
artemis make:module users
```

## 📦 Estrutura
```
app/
 ├── modules/
 │   ├── users/
 │   │   ├── controller.go
 │   │   ├── repository.go
 │   │   └── service.go
 ├── core/
 │   ├── http/
 │   ├── config/
 │   └── database/
```

## 🛠️ Em desenvolvimento
- [ ] Gerador de migrations
- [ ] Integração com Docker
- [ ] Middleware e Providers
- [ ] CLI para testes e seeds

## 📜 Licença
MIT © 2025 - Valdir Barbosa
