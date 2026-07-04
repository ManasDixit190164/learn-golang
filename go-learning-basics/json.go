package main

import (
    "encoding/json"
    "fmt"
)

// PersonInfo represents a simple JSON structure.
type PersonInfo struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

// This function demonstrates JSON encoding and decoding.
func jsonExample() {
    fmt.Println("JSON example:")

    person := PersonInfo{Name: "Manas", Age: 25}

    data, err := json.Marshal(person)
    if err != nil {
        fmt.Println("Marshal error:", err)
        return
    }

    fmt.Println("JSON output:", string(data))

    var decoded PersonInfo
    err = json.Unmarshal(data, &decoded)
    if err != nil {
        fmt.Println("Unmarshal error:", err)
        return
    }

    fmt.Println("Decoded name:", decoded.Name)
    fmt.Println("Decoded age:", decoded.Age)
}
