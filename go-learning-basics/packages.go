package main // package declaration

import ( // start import block
	"fmt" // import package

	"go-learning/mathpkg" // import package
) // end import block or close block

// This function demonstrates using a custom package.
func packagesExample() { // function declaration
	fmt.Println("Package example:")                         // statement
	fmt.Println("Add result:", mathpkg.Add(2, 3))           // statement
	fmt.Println("Multiply result:", mathpkg.Multiply(2, 3)) // statement
} // statement
