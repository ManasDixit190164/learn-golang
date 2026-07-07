package main // package declaration

import "fmt" // statement

// This function demonstrates maps in Go.
func mapsExample() { // function declaration
    // Maps store key-value pairs.
    person := map[string]string{ // declare and initialize variable
        "name": "Manas", // statement
        "city": "Delhi", // statement
    } // statement

    fmt.Println("Map values:") // statement
    fmt.Println(person) // statement
} // statement
