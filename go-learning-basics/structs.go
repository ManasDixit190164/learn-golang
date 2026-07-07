package main // package declaration

import "fmt" // statement

// Person defines a simple struct with fields.
type Person struct { // type/struct declaration
    Name string // statement
    Age  int // statement
} // statement

// This function demonstrates structs in Go.
func structsExample() { // function declaration
    // Create a Person value using a struct literal.
    p := Person{Name: "Manas", Age: 25} // declare and initialize variable

    fmt.Println("Struct values:") // statement
    fmt.Println(p) // statement
} // statement
