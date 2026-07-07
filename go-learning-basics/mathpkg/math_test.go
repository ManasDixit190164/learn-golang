package mathpkg // package declaration

import "testing" // statement

func TestAdd(t *testing.T) { // function declaration
	got := Add(2, 3) // declare and initialize variable
	want := 5        // declare and initialize variable

	if got != want { // check condition
		t.Fatalf("Add(2, 3) = %d; want %d", got, want) // assign value
	} // statement
} // statement

func TestMultiply(t *testing.T) { // function declaration
	got := Multiply(2, 3) // declare and initialize variable
	want := 6             // declare and initialize variable

	if got != want { // check condition
		t.Fatalf("Multiply(2, 3) = %d; want %d", got, want) // assign value
	} // statement
} // statement
