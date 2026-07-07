package main // package declaration

import ( // start import block
    "encoding/json" // import package
    "fmt" // import package
) // end import block or close block

// PersonInfo represents a simple JSON structure.
type PersonInfo struct { // type/struct declaration
    Name string `json:"name"` // statement
    Age  int    `json:"age"` // statement
} // statement

// This function demonstrates JSON encoding and decoding.
func jsonExample() { // function declaration
    fmt.Println("JSON example:") // statement

    person := PersonInfo{Name: "Manas", Age: 25} // declare and initialize variable

    data, err := json.Marshal(person) // declare and initialize variable
    if err != nil { // check condition
        fmt.Println("Marshal error:", err) // statement
        return // return
    } // statement

    fmt.Println("JSON output:", string(data)) // statement

    var decoded PersonInfo // variable declaration
    err = json.Unmarshal(data, &decoded) // assign value
    if err != nil { // check condition
        fmt.Println("Unmarshal error:", err) // statement
        return // return
    } // statement

    fmt.Println("Decoded name:", decoded.Name) // statement
    fmt.Println("Decoded age:", decoded.Age) // statement
} // statement
