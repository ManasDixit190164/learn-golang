package utils // package declaration for the module

import "golang.org/x/crypto/bcrypt" // execute statement

func HashPassword(password string) (string, error) { // declare function
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // declare and initialize variable
	if err != nil { // if condition
		return "", err // return statement
	} // end block
	return string(hash), nil // return statement
} // end block

func CheckPassword(password, hash string) bool { // declare function
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil // return statement
} // end block
