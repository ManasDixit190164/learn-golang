package main // package declaration

import ( // start import block
	"errors" // import package
	"fmt"    // import package
) // end import block or close block

// This function returns an error for demonstration purposes.
func divide(a, b int) (int, error) { // function declaration
	if b == 0 { // check condition
		return 0, errors.New("cannot divide by zero") // return result or error
	} // statement
	return a / b, nil // return result or error
} // statement

// This function demonstrates basic error handling in Go.
func errorsExample() { // function declaration
	fmt.Println("Error handling example:") // statement

	result, err := divide(10, 2) // declare and initialize variable
	if err != nil {              // check condition
		fmt.Println("Error:", err) // statement
	} else { // statement
		fmt.Println("Result:", result) // statement
	} // statement

	result, err = divide(10, 0) // assign value
	if err != nil {             // check condition
		fmt.Println("Error:", err) // statement
	} else { // statement
		fmt.Println("Result:", result) // statement
	} // statement
} // statement
