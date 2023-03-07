package utils

import (
	"unicode"
)

func IsValid(s string) bool {
	var (
		hasMinLen = false
		hasNumber = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}
	return hasMinLen && hasNumber
}
