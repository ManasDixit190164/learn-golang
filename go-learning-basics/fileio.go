package main

import (
    "fmt"
    "os"
)

// This function demonstrates basic file writing and reading.
func fileIOExample() {
    fmt.Println("File I/O example:")

    filename := "sample.txt"
    content := "Hello from Go file I/O"

    err := os.WriteFile(filename, []byte(content), 0o644)
    if err != nil {
        fmt.Println("Write error:", err)
        return
    }

    data, err := os.ReadFile(filename)
    if err != nil {
        fmt.Println("Read error:", err)
        return
    }

    fmt.Println("File content:", string(data))
}
