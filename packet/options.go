package packet

import (
	"encoding/binary"
	"strings"

	"github.com/cr-mao/lorig/conf"
)

const (
	littleEndian = "little"
	bigEndian    = "big"
)

const (
	defaultSeqBytesLen    = 2
	defaultRouteBytesLen  = 2
	defaultBufferBytesLen = 5000
)

const (
	defaultEndianKey         = "packet.endian"
	defaultSeqBytesLenKey    = "packet.seqBytesLen"
	defaultRouteBytesLenKey  = "packet.routeBytesLen"
	defaultBufferBytesLenKey = "packet.bufferBytesLen"
)

type options struct {
	// 字节序
	// 默认为binary.LittleEndian
	byteOrder binary.ByteOrder

	// 序列号字节长度（字节），长度为0时不开启序列号编码
	// 默认为2字节，最大值为65535
	seqBytesLen int

	// 路由字节长度（字节）
	// 默认为2字节，最大值为65535
	routeBytesLen int

	// 消息字节长度（字节）
	// 默认为5000字节
	bufferBytesLen int
}

type Option func(o *options)

func defaultOptions() *options {
	opts := &options{
		byteOrder:      binary.BigEndian,
		seqBytesLen:    conf.GetInt(defaultSeqBytesLenKey, defaultSeqBytesLen),
		routeBytesLen:  conf.GetInt(defaultRouteBytesLenKey, defaultRouteBytesLen),
		bufferBytesLen: conf.GetInt(defaultBufferBytesLenKey, defaultBufferBytesLen),
	}
	endian := conf.GetString(defaultEndianKey, bigEndian)
	switch strings.ToLower(endian) {
	case littleEndian:
		opts.byteOrder = binary.LittleEndian
	default:
		opts.byteOrder = binary.BigEndian
	}
	return opts
}

// WithByteOrder 设置字节序
func WithByteOrder(byteOrder binary.ByteOrder) Option {
	return func(o *options) { o.byteOrder = byteOrder }
}

// WithSeqBytesLen 设置序列号字节长度
func WithSeqBytesLen(seqBytesLen int) Option {
	return func(o *options) { o.seqBytesLen = seqBytesLen }
}

// WithRouteBytesLen 设置路由字节长度
func WithRouteBytesLen(routeBytesLen int) Option {
	return func(o *options) { o.routeBytesLen = routeBytesLen }
}

// WithBufferBytesLen 设置消息字节长度
func WithBufferBytesLen(bufferBytesLen int) Option {
	return func(o *options) { o.bufferBytesLen = bufferBytesLen }
}
