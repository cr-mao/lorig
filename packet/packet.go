package packet

var globalPacker Packer

func init() {
	globalPacker = NewPacker()
}

// SetPacker 设置打包器
func SetPacker(packer Packer) {
	globalPacker = packer
}

// GetPacker 获取打包器
func GetPacker() Packer {
	return globalPacker
}

// Pack 打包消息
func Pack(message *Message) ([]byte, error) {
	return globalPacker.Pack(message)
}

// Unpack 解包消息
func Unpack(data []byte) (*Message, error) {
	return globalPacker.Unpack(data)
}
