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
func FormatNodeProperty(property, title string) bool {
	if property == "ll:en" || title == "http://www.w3.org/2000/01/rdf-schema#label" {
		return true
	}

	return false
}
