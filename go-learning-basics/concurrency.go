package main

import (
    "fmt"
    "time"
)

// This function demonstrates a simple goroutine.
func concurrencyExample() {
    fmt.Println("Concurrency example:")

    go func() {
        time.Sleep(100 * time.Millisecond)
        fmt.Println("Goroutine finished")
    }()

    fmt.Println("Main function continues")
    time.Sleep(200 * time.Millisecond)
}
