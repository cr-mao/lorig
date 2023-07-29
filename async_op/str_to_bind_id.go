package async_op

import (
	"hash/crc32"
)

// StrToBindId 字符串转为 bindId
func StrToBindId(strVal string) int {
	v := int(crc32.ChecksumIEEE([]byte(strVal)))

	if v >= 0 {
		return v
	} else {
		return -v
	}
}
