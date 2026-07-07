package main // package declaration

import "fmt" // statement

// Speaker describes anything that can speak.
type Speaker interface { // type declaration
	Speak() string // statement
} // statement

// Dog implements the Speaker interface.
type Dog struct{} // type/struct declaration

func (d Dog) Speak() string { // function declaration
	return "Woof" // return result or error
} // statement

// Cat implements the Speaker interface.
type Cat struct{} // type/struct declaration

func (c Cat) Speak() string { // function declaration
	return "Meow" // return result or error
} // statement

// This function demonstrates interfaces in Go.
func interfacesExample() { // function declaration
	animals := []Speaker{Dog{}, Cat{}} // declare and initialize variable

	fmt.Println("Interface example:") // statement
	for _, animal := range animals {  // declare and initialize variable
		fmt.Println(animal.Speak()) // statement
	} // statement
} // statement
