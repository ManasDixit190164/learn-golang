package main

import "fmt"

// Speaker describes anything that can speak.
type Speaker interface {
    Speak() string
}

// Dog implements the Speaker interface.
type Dog struct{}

func (d Dog) Speak() string {
    return "Woof"
}

// Cat implements the Speaker interface.
type Cat struct{}

func (c Cat) Speak() string {
    return "Meow"
}

// This function demonstrates interfaces in Go.
func interfacesExample() {
    animals := []Speaker{Dog{}, Cat{}}

    fmt.Println("Interface example:")
    for _, animal := range animals {
        fmt.Println(animal.Speak())
    }
}
