package templates

// Embedded contains the embedded templates or resources used throughout the application.
var Embedded = map[string]string{
    "main.go.tmpl": `package main

import "fmt"

func main() {
    fmt.Println("Hello, Artemis!")
}
`,
    "module.go.tmpl": `package {{.ModuleName}}

import "fmt"

func Init() {
    fmt.Println("Initializing module: {{.ModuleName}}")
}
`,
    "service.go.tmpl": `package {{.ModuleName}}

type {{.ServiceName}} struct {}

func New{{.ServiceName}}() *{{.ServiceName}} {
    return &{{.ServiceName}}{}
}
`,
    "repository.go.tmpl": `package {{.ModuleName}}

type {{.RepositoryName}} struct {}

func New{{.RepositoryName}}() *{{.RepositoryName}} {
    return &{{.RepositoryName}}{}
}
`,
    "handler.go.tmpl": `package {{.ModuleName}}

import "net/http"

func Handle{{.HandlerName}}(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Handling {{.HandlerName}}"))
}
`,
    "domain.go.tmpl": `package {{.ModuleName}}

type {{.DomainName}} struct {
    ID string
}
`,
}