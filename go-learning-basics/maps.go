package main

import "fmt"

// This function demonstrates maps in Go.
func mapsExample() {
    // Maps store key-value pairs.
    person := map[string]string{
        "name": "Manas",
        "city": "Delhi",
    }

    fmt.Println("Map values:")
    fmt.Println(person)
}
