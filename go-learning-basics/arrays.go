package main // package declaration

import "fmt" // statement

// This function demonstrates arrays in Go.
func arraysExample() { // function declaration
	// Arrays have a fixed size.
	numbers := [3]int{10, 20, 30} // declare and initialize variable

	fmt.Println("Array values:") // statement
	fmt.Println(numbers)         // statement
} // statement
