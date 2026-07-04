package main

import "fmt"

// This function demonstrates channels in Go.
func channelsExample() {
    fmt.Println("Channel example:")

    messages := make(chan string, 2)
    messages <- "hello"
    messages <- "world"

    fmt.Println(<-messages)
    fmt.Println(<-messages)
}
