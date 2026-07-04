package main

import "fmt"

// This function demonstrates pointer usage in Go.
func pointersExample() {
    value := 10

    // Create a pointer to value.
    pointer := &value

    fmt.Println("Pointer example:")
    fmt.Println("Value:", value)
    fmt.Println("Pointer address:", pointer)
    fmt.Println("Value through pointer:", *pointer)
}
