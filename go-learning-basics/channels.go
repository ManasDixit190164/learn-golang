package main // package declaration

import "fmt" // statement

// This function demonstrates channels in Go.
func channelsExample() { // function declaration
    fmt.Println("Channel example:") // statement

    messages := make(chan string, 2) // declare and initialize variable
    messages <- "hello" // statement
    messages <- "world" // statement

    fmt.Println(<-messages) // statement
    fmt.Println(<-messages) // statement
} // statement
