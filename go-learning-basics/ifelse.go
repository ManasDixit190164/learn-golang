package main // package declaration

import "fmt" // statement

// checkAge prints a message based on the age value.
func checkAge(age int) { // function declaration
	if age >= 18 { // check condition
		fmt.Println("You are an adult") // statement
	} else { // statement
		fmt.Println("You are a minor") // statement
	} // statement
} // statement
