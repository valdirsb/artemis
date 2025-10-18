package generator

import (
	"os"
	"path/filepath"
	"text/template"
)

// Project represents the structure of a new Artemis project.
type Project struct {
	Name string
}

// CreateProject generates the project structure and files based on the provided project name.
func CreateProject(projectName string) error {
	project := Project{Name: projectName}

	// Define the project directories
	directories := []string{
		filepath.Join(projectName, "cmd", "artemis"),
		filepath.Join(projectName, "pkg", "cli"),
		filepath.Join(projectName, "pkg", "generator"),
		filepath.Join(projectName, "internal", "config"),
	}

	// Create the directories
	for _, dir := range directories {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	// Generate main.go file
	if err := generateProjectFile(filepath.Join(projectName, "cmd", "artemis", "main.go"), "main.go.tmpl", project); err != nil {
		return err
	}

	// Additional project setup can be added here

	return nil
}

// generateFile creates a file from a template.
func generateProjectFile(filePath, templateName string, project Project) error {
	tmpl, err := template.ParseFiles(filepath.Join("pkg", "generator", "templates", templateName))
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, project)
}
