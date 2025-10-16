package templates

import "fmt"

func MainTemplate(projectName string) string {
return fmt.Sprintf("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello from %s!\")\n}\n", projectName)
}

func BootstrapTemplate(projectName string) string {
return "package bootstrap\n\n// TODO: Implement bootstrap\n"
}

func ConfigTemplate(projectName string) string {
return "package config\n\n// TODO: Implement config\n"
}

func ContractsTemplate(projectName string) string {
return "package contracts\n\n// TODO: Implement contracts\n"
}

func ContainerTemplate(projectName string) string {
return "package container\n\n// TODO: Implement container\n"
}

func GoModTemplate(projectName string) string {
return fmt.Sprintf("module %s\n\ngo 1.21\n", projectName)
}

func ReadmeTemplate(projectName string) string {
return fmt.Sprintf("# %s\n\nProjeto Artemis\n", projectName)
}

func GitignoreTemplate() string {
return "*.exe\n*.log\n.env\n"
}

func DockerfileTemplate(projectName string) string {
return "FROM golang:1.21-alpine\nWORKDIR /app\nCOPY . .\nCMD [\"./main\"]\n"
}

func DockerComposeTemplate(projectName string) string {
return fmt.Sprintf("version: '3.8'\nservices:\n  %s:\n    build: .\n", projectName)
}
