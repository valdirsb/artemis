package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Module struct {
	Name string
}

func CreateModule(moduleName string) error {
	module := Module{Name: moduleName}

	// Define the module directory structure
	moduleDir := filepath.Join("internal", "modules", moduleName)
	directories := []string{
		moduleDir,
		filepath.Join(moduleDir, "domain"),
		filepath.Join(moduleDir, "ports"),
		filepath.Join(moduleDir, "service"),
		filepath.Join(moduleDir, "adapters"),
		filepath.Join(moduleDir, "repository"),
		filepath.Join(moduleDir, "handler"),
	}

	// Create directories
	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	// Generate module.go file from template
	if err := generateModuleFile("pkg/generator/templates/module.go.tmpl", filepath.Join(moduleDir, "module.go"), module); err != nil {
		return err
	}

	// Generate other files as needed (e.g., domain, service, etc.)
	// This can be expanded based on the requirements

	return nil
}

func generateModuleFile(templatePath, outputPath string, data interface{}) error {

	fmt.Println("DEBUG: generateModuleFile")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
