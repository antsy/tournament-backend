package utils

import "strings"

/**
 * Check whether string is all whitespace
 */
func IsEmptyOrWhitespace(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
