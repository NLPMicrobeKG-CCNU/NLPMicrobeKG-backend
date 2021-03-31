package util

import (
	"encoding/base64"
	"fmt"
	"strings"
	"unicode"
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

// FormatRequestQuery format request query.
func FormatRequestQuery(query string) string {
	// provide panic
	if len(query) < 1 {
		return ""
	}

	query = strings.TrimLeft(query, " ")
	query = strings.TrimRight(query, " ")

	lowerSlice := []rune(strings.ToLower(query))
	res := fmt.Sprintf("%c%s", unicode.ToUpper(lowerSlice[0]), string(lowerSlice[1:]))

	return res
}

// FormatRequestQueryCaseInsensitivity format request query with case insensitive.
func FormatRequestQueryCaseInsensitivity(query string) string {
	// provide panic
	if len(query) < 1 {
		return ""
	}

	query = strings.TrimLeft(query, " ")
	query = strings.TrimRight(query, " ")

	return query
}

// ParseBase64 parses a base64 encoded string to usual string.
func ParseBase64(base64Str string) (string, error) {
	if base64Str == "" {
		return "", nil
	}

	res, err := base64.StdEncoding.DecodeString(base64Str)

	return string(res), err
}
