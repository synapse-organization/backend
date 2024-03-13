package utils

import (
	"strings"
)

func SplitMethodPrefix(methodName string) (string, string) {
	for i, char := range methodName {
		if i > 0 && strings.ToUpper(string(char)) == string(char) {
			return methodName[:i], methodName[i:]
		}
	}
	return "", methodName
}
