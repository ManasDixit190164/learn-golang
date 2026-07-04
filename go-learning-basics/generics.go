package main

import "fmt"

// This function demonstrates generics in Go.
func printValue[T any](value T) {
    fmt.Println("Generic value:", value)
}

// This function shows how to call a generic function.
func genericsExample() {
    fmt.Println("Generics example:")
    printValue("Hello")
    printValue(42)
}
