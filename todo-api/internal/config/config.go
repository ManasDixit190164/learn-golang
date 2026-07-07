package config // package declaration

import "os" // statement

type Config struct { // type/struct declaration
	Port string // statement
} // statement

func Load() Config { // function declaration
	port := os.Getenv("PORT") // declare and initialize variable

	if port == "" { // check condition
		port = "8080" // assign value
	} // statement

	return Config{ // return result or error
		Port: port, // statement
	} // statement
} // statement
