package util

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

func SHA1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return hex.EncodeToString(t.Sum(nil))
}
