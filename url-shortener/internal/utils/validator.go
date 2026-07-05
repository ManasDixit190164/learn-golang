package utils

import (
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

var aliasRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,64}$`)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidHTTPURL(rawURL string) bool {
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	scheme := strings.ToLower(parsed.Scheme)
	return (scheme == "http" || scheme == "https") && parsed.Host != ""
}

func IsValidAlias(alias string) bool {
	return aliasRegex.MatchString(alias)
}
