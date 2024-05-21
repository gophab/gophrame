package util

import (
	"encoding/base64"
)

// 先base64，然后MD5
func Base64(params string) string {
	return base64.StdEncoding.EncodeToString([]byte(params))
}

// 先base64，然后MD5
func Base64Bytes(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}
