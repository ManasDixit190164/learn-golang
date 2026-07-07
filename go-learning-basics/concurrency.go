package main // package declaration

import ( // start import block
    "fmt" // import package
    "time" // import package
) // end import block or close block

// This function demonstrates a simple goroutine.
func concurrencyExample() { // function declaration
    fmt.Println("Concurrency example:") // statement

    go func() { // start asynchronous goroutine
        time.Sleep(100 * time.Millisecond) // statement
        fmt.Println("Goroutine finished") // statement
    }() // statement

    fmt.Println("Main function continues") // statement
    time.Sleep(200 * time.Millisecond) // statement
} // statement
