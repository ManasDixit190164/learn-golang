package main

import "fmt"

// checkAge prints a message based on the age value.
func checkAge(age int) {
    if age >= 18 {
        fmt.Println("You are an adult")
    } else {
        fmt.Println("You are a minor")
    }
}
