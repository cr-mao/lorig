package xstring

import (
	"hash/crc32"
	"unicode"
	"unicode/utf8"
)

// FirstLetterIsUpper 首字母是否大写
func FirstLetterIsUpper(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return r != utf8.RuneError && unicode.IsUpper(r)
}

// FirstLetterIsLower 首字母是否小写
func FirstLetterIsLower(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return r != utf8.RuneError && unicode.IsLower(r)
}

// StrToId 字符串转为 int
func StrToId(strVal string) int {
	v := int(crc32.ChecksumIEEE([]byte(strVal)))
	if v >= 0 {
		return v
	} else {
		return -v
	}
}
