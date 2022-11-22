package utils

import "unicode"

// Contains check if slice contain specific element. Could be improved with generics for other types.
func Contains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

// LcFirst converts the first character of a string to lowercase.
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}
