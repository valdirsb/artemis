package templates

import "fmt"

func EntityTemplate(moduleNameTitle, moduleNameLower string) string {
return fmt.Sprintf("package domain\n\ntype %s struct {\n\tID   string\n\tName string\n}\n", moduleNameTitle)
}

func DomainRepositoryTemplate(moduleNameTitle, moduleNameLower string) string {
return fmt.Sprintf("package domain\n\nimport \"context\"\n\ntype %sRepository interface {\n\tCreate(ctx context.Context, entity *%s) error\n}\n", moduleNameTitle, moduleNameTitle)
}

func PortsTemplate(moduleNameTitle, moduleNameLower string) string {
return fmt.Sprintf("package ports\n\n// TODO: Implement %s ports\n", moduleNameLower)
}

func ServiceTemplate(moduleNameTitle, moduleNameLower string) string {
return fmt.Sprintf("package service\n\n// TODO: Implement %s service\n", moduleNameLower)
}

func RepositoryTemplate(moduleNameTitle, moduleNameLower string) string {
return fmt.Sprintf("package repository\n\n// TODO: Implement %s repository\n", moduleNameLower)
}

func HandlerTemplate(moduleNameTitle, moduleNameLower string) string {
return fmt.Sprintf("package handler\n\n// TODO: Implement %s handler\n", moduleNameLower)
}

func MigrationTemplate(migrationName, timestamp string) string {
return fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n-- TODO: Add SQL statements\n", migrationName, timestamp)
}
