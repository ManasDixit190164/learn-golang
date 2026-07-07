package main // package declaration

import ( // start import block
	"fmt" // import package
	"os"  // import package
) // end import block or close block

// This function demonstrates basic file writing and reading.
func fileIOExample() { // function declaration
	fmt.Println("File I/O example:") // statement

	filename := "sample.txt"            // declare and initialize variable
	content := "Hello from Go file I/O" // declare and initialize variable

	err := os.WriteFile(filename, []byte(content), 0o644) // declare and initialize variable
	if err != nil {                                       // check condition
		fmt.Println("Write error:", err) // statement
		return                           // return
	} // statement

	data, err := os.ReadFile(filename) // declare and initialize variable
	if err != nil {                    // check condition
		fmt.Println("Read error:", err) // statement
		return                          // return
	} // statement

	fmt.Println("File content:", string(data)) // statement
} // statement
