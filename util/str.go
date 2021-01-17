package util

import (
	"fmt"
	"strings"
)

// FormatNodeTypeStr format node's type.
func FormatNodeTypeStr(strs []string) []string {
	var res []string
	for _, str := range strs {
		index := strings.IndexFunc(str, func(r rune) bool {
			if r == '#' {
				return true
			}
			return false
		})

		res = append(res, fmt.Sprintf("entity: %s", str[index+1:]))
	}

	return res
}

// FormatNodeProperty format node property.
func FormatNodeProperty(property string) bool {
	if property == "lt:http://www.w3.org/2001/XMLSchema#string" || property == "ll:en" {
		return true
	}

	return false
}
