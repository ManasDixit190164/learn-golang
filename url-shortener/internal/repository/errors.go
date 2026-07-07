package repository // package declaration for the module

import "errors" // execute statement

var ( // declare variable
	ErrNotFound = errors.New("record not found") // handle missing database record
	ErrConflict = errors.New("record already exists") // handle duplicate database entry
) // end import block or block scope
