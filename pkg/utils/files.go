package utils

import (
    "io/ioutil"
    "os"
    "path/filepath"
)

// CreateFile creates a new file with the specified content.
func CreateFile(filePath string, content []byte) error {
    return ioutil.WriteFile(filePath, content, 0644)
}

// ReadFile reads the content of a file and returns it as a byte slice.
func ReadFile(filePath string) ([]byte, error) {
    return ioutil.ReadFile(filePath)
}

// DeleteFile removes a file from the filesystem.
func DeleteFile(filePath string) error {
    return os.Remove(filePath)
}

// EnsureDir creates a directory if it does not exist.
func EnsureDir(dirPath string) error {
    if _, err := os.Stat(dirPath); os.IsNotExist(err) {
        return os.MkdirAll(dirPath, os.ModePerm)
    }
    return nil
}

// WalkDir walks the directory tree rooted at root and calls walkFn for each file or directory in the tree.
func WalkDir(root string, walkFn filepath.WalkFunc) error {
    return filepath.Walk(root, walkFn)
}