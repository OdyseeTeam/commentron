package util

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
