package main // package declaration

import "fmt" // statement

// This function demonstrates generics in Go.
func printValue[T any](value T) { // function declaration
    fmt.Println("Generic value:", value) // statement
} // statement

// This function shows how to call a generic function.
func genericsExample() { // function declaration
    fmt.Println("Generics example:") // statement
    printValue("Hello") // statement
    printValue(42) // statement
} // statement
