// Package strings provides utility functions to manipulate strings.
package strings

import "regexp"

// StripNonAlphaNumeric returns a string with all non alphanumeric characters stripped.
func StripNonAlphaNumeric(input string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(input, "")
}
