package main

import "fmt"

// Person defines a simple struct with fields.
type Person struct {
    Name string
    Age  int
}

// This function demonstrates structs in Go.
func structsExample() {
    // Create a Person value using a struct literal.
    p := Person{Name: "Manas", Age: 25}

    fmt.Println("Struct values:")
    fmt.Println(p)
}
