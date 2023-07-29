package async_op

import (
	"hash/crc32"
)

// StrToWorkId 字符串转为 workID
func StrToWorkId(strVal string) int {
	v := int(crc32.ChecksumIEEE([]byte(strVal)))
	if v >= 0 {
		return v
	} else {
		return -v
	}
}
