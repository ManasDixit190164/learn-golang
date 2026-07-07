package utils // package declaration for the module

import ( // start import block
	"crypto/rand" // import package
	"crypto/sha256" // import package
	"encoding/base64" // import package
	"encoding/hex" // import package
) // end import block or block scope

func GenerateSecureToken(byteLength int) (string, error) { // declare function
	b := make([]byte, byteLength) // declare and initialize variable
	if _, err := rand.Read(b); err != nil { // if condition
		return "", err // return statement
	} // end block
	return base64.RawURLEncoding.EncodeToString(b), nil // return statement
} // end block

func HashToken(token string) string { // declare function
	sum := sha256.Sum256([]byte(token)) // declare and initialize variable
	return hex.EncodeToString(sum[:]) // return statement
} // end block
