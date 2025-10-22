package utils

// Func CreateDir cria um diretório se ele não existir
import (
	"os"
)

// CreateDir cria um diretório se ele não existir
func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
