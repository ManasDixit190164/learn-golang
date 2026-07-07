package main // package declaration

import "fmt" // statement

// This function demonstrates pointer usage in Go.
func pointersExample() { // function declaration
    value := 10 // declare and initialize variable

    // Create a pointer to value.
    pointer := &value // declare and initialize variable

    fmt.Println("Pointer example:") // statement
    fmt.Println("Value:", value) // statement
    fmt.Println("Pointer address:", pointer) // statement
    fmt.Println("Value through pointer:", *pointer) // statement
} // statement
