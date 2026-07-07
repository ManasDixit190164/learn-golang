package main // package declaration

import "fmt" // statement

// Introduce prints a greeting using the receiver.
func (p Person) Introduce() string { // function declaration
	return "Hello, I am " + p.Name // return result or error
} // statement

// This function demonstrates methods in Go.
func methodsExample() { // function declaration
	person := Person{Name: "Manas"} // declare and initialize variable

	fmt.Println("Method example:")  // statement
	fmt.Println(person.Introduce()) // statement
} // statement
