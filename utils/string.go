package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

// StringEncMd5 ...
func StringEncMd5(s string) string {
	m := md5.New()
	io.WriteString(m, s)
	return hex.EncodeToString(m.Sum(nil))
}

// Fstring ...
func Fstring(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}
