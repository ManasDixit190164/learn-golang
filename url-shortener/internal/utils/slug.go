package utils // package declaration for the module

import ( // start import block
	"crypto/rand" // import package
	"math/big" // import package
) // end import block or block scope

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // declare constant

func GenerateShortCode(length int) (string, error) { // declare function
	if length <= 0 { // if condition
		length = 7 // assign value
	} // end block

	result := make([]byte, length) // declare and initialize variable
	max := big.NewInt(int64(len(base62Chars))) // declare and initialize variable

	for i := 0; i < length; i++ { // for loop
		n, err := rand.Int(rand.Reader, max) // declare and initialize variable
		if err != nil { // if condition
			return "", err // return statement
		} // end block
		result[i] = base62Chars[n.Int64()] // assign value
	} // end block

	return string(result), nil // return statement
} // end block
