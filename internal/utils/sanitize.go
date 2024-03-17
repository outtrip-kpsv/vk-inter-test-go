package utils

import (
	"regexp"
)

func Sanitize(input string) (string, bool) {
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s-]")
	result := reg.ReplaceAllString(input, "")
	return result, result == input
}
