package helper

import "strings"

// StringSplitArg splits CSV string into interface values that can be passed to SQLBoiler
func StringSplitArg(stringToSplit, separator string) []interface{} {
	split := strings.Split(stringToSplit, separator)
	splitInterface := make([]interface{}, len(split))
	for i, s := range split {
		splitInterface[i] = s
	}
	return splitInterface
}

// ReverseString reverses a string
func ReverseString(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
