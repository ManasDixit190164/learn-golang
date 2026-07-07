package main // package declaration

import "fmt" // statement

// This function demonstrates slices in Go.
func slicesExample() { // function declaration
    // Slices are dynamic and more flexible than arrays.
    fruits := []string{"apple", "banana", "mango"} // declare and initialize variable

    fmt.Println("Slice values:") // statement
    fmt.Println(fruits) // statement
} // statement
