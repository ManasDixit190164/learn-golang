package main // package declaration

import "fmt" // statement

// This function accepts a name and returns a greeting message.
func greet(name string) string { // function declaration
    return "Hello, " + name // return result or error
} // statement

// This function demonstrates calling another function from inside a function.
func functionsExample() { // function declaration
    message := greet("Manas") // declare and initialize variable
    fmt.Println("Function example:") // statement
    fmt.Println(message) // statement
} // statement
