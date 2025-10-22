# Changelog

## [0.1.0] - 2025-10-16

### ✨ Added
- CLI base com Cobra framework
- Comando `artemis new <projeto>` para criação de projetos
- Comando `artemis make module <nome>` para geração de módulos
- Comando `artemis make migration <nome>` para geração de migrations
- Comando `artemis serve` para servidor de desenvolvimento
- Sistema de templates flexível e extensível
- Arquitetura hexagonal nos módulos gerados
- Estrutura modular inspirada no Laravel
- Suporte a Go 1.21+

### 🏗️ Infrastructure
- Estrutura de projeto organizada
- Sistema de geração de código automático  
- Templates para diferentes componentes
- Build system configurado

### 📦 Dependencies
- github.com/spf13/cobra v1.10.1
- github.com/spf13/viper v1.21.0
- github.com/gorilla/mux v1.8.0

### 🎯 Features Implementadas
- [x] CLI completo e funcional
- [x] Geração de projetos com estrutura padrão
- [x] Geração de módulos com arquitetura hexagonal
- [x] Sistema de templates personalizáveis
- [x] Servidor de desenvolvimento integrado
- [x] Suporte a migrations SQL

### 🔄 Em Desenvolvimento
- [ ] Integração com banco de dados
- [ ] Sistema de middleware avançado
- [ ] CLI para testes automatizados
- [ ] Documentação interativa
- [ ] Plugin system
