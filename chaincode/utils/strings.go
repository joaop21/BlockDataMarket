package utils

import "strings"

// Function for string concatenation using builder, thus not creating a new string every time.
// Strings are immutable and it's a slow operation without builder.
func Concat(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// Remove repetitions from an array of strings
func RemoveRepetitions(stringArray []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _ , entry := range stringArray {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// check if an array of strings contains a specific string
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}