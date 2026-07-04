package main

import (
    "fmt"

    "go-learning/mathpkg"
)

// This function demonstrates using a custom package.
func packagesExample() {
    fmt.Println("Package example:")
    fmt.Println("Add result:", mathpkg.Add(2, 3))
    fmt.Println("Multiply result:", mathpkg.Multiply(2, 3))
}
