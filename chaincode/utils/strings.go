package utils

import "strings"

// Function for string concatenation using builder, thus not creating a new string every time.
// Strings are immutable and it's a slow operation without builder.
func concat(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}
