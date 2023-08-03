package proto

import (
	"errors"

	"github.com/gogo/protobuf/proto"
)

const Name = "proto"

var DefaultCodec = &codec{}

type codec struct{}

// Name 编解码器名称
func (codec) Name() string {
	return Name
}

// Marshal 编码
func (codec) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("can't marshal a value that not implements proto.Message interface")
	}

	return proto.Marshal(msg)
}

// Unmarshal 解码
func (codec) Unmarshal(data []byte, v interface{}) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return errors.New("can't unmarshal to a value that not implements proto.Message")
	}

	return proto.Unmarshal(data, msg)
}

// Marshal 编码
func Marshal(v interface{}) ([]byte, error) {
	return DefaultCodec.Marshal(v)
}

// Unmarshal 解码
func Unmarshal(data []byte, v interface{}) error {
	return DefaultCodec.Unmarshal(data, v)
}
