package main

import "fmt"

// Introduce prints a greeting using the receiver.
func (p Person) Introduce() string {
    return "Hello, I am " + p.Name
}

// This function demonstrates methods in Go.
func methodsExample() {
    person := Person{Name: "Manas"}

    fmt.Println("Method example:")
    fmt.Println(person.Introduce())
}
