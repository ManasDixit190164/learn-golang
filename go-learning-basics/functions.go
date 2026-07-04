package main

import "fmt"

// This function accepts a name and returns a greeting message.
func greet(name string) string {
    return "Hello, " + name
}

// This function demonstrates calling another function from inside a function.
func functionsExample() {
    message := greet("Manas")
    fmt.Println("Function example:")
    fmt.Println(message)
}
