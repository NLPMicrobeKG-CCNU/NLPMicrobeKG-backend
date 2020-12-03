package util

import "net/url"

// EncodeURL 编码URL
func EncodeURL(str string) string {
	encodedURL := url.QueryEscape(str)
	return encodedURL
}
