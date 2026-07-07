package utils // package declaration for the module

import ( // start import block
	"net/mail" // import package
	"net/url" // import package
	"regexp" // import package
	"strings" // import package
) // end import block or block scope

var aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,64}$`) // declare variable

func IsValidEmail(email string) bool { // declare function
	_, err := mail.ParseAddress(email) // declare and initialize variable
	return err == nil // return statement
} // end block

func IsValidHTTPURL(rawURL string) bool { // declare function
	parsed, err := url.ParseRequestURI(rawURL) // declare and initialize variable
	if err != nil { // if condition
		return false // return statement
	} // end block

	scheme := strings.ToLower(parsed.Scheme) // lowercase text
	return (scheme == "http" || scheme == "https") && parsed.Host != "" // return statement
} // end block

func IsValidAlias(alias string) bool { // declare function
	return aliasRegex.MatchString(alias) // return statement
} // end block
