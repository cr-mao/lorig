package rsa

import (
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/crypto/hash"
	"github.com/cr-mao/lorig/utils/xconv"
	"strings"
)

const (
	defaultEncryptorHashKey      = "config.crypto.rsa.encryptor.hash"
	defaultEncryptorPaddingKey   = "config.crypto.rsa.encryptor.padding"
	defaultEncryptorLabelKey     = "config.crypto.rsa.encryptor.label"
	defaultEncryptorBlockSizeKey = "config.crypto.rsa.encryptor.blockSize"
	defaultEncryptorPublicKeyKey = "config.crypto.rsa.encryptor.publicKey"
)

type EncryptorOption func(o *encryptorOptions)

type encryptorOptions struct {
	// hash算法。支持sha1、sha224、sha256、sha384、sha512
	// 默认为sha256
	hash hash.Hash

	// 填充规则。支持NORMAL和OAEP
	// 默认为NORMAL
	padding EncryptPadding

	// 标签。加解密时必需一致
	// 默认为空
	label []byte

	// 加密数据块大小，单位字节。由于加密数据长度限制，需要对加密数据进行分块儿加密。
	// 默认根据填充方式选择最大的长度进行切割
	blockSize int

	// 公钥。可设置文件路径或公钥串
	publicKey string
}

func defaultEncryptorOptions() *encryptorOptions {
	return &encryptorOptions{
		hash:      hash.Hash(strings.ToLower(conf.Get(defaultEncryptorHashKey))),
		padding:   EncryptPadding(strings.ToUpper(conf.Get(defaultEncryptorPaddingKey))),
		label:     []byte(conf.Get(defaultEncryptorLabelKey)),
		blockSize: conf.GetInt(defaultEncryptorBlockSizeKey),
		publicKey: conf.GetString(defaultEncryptorPublicKeyKey),
	}
}

// WithEncryptorHash 设置加密hash算法
func WithEncryptorHash(hash hash.Hash) EncryptorOption {
	return func(o *encryptorOptions) { o.hash = hash }
}

// WithEncryptorPadding 设置加密填充规则
func WithEncryptorPadding(padding EncryptPadding) EncryptorOption {
	return func(o *encryptorOptions) { o.padding = padding }
}

// WithEncryptorLabel 设置加密标签
func WithEncryptorLabel(label string) EncryptorOption {
	return func(o *encryptorOptions) { o.label = xconv.StringToBytes(label) }
}

// WithEncryptorBlockSize 设置加密数据块大小
func WithEncryptorBlockSize(blockSize int) EncryptorOption {
	return func(o *encryptorOptions) { o.blockSize = blockSize }
}

// WithEncryptorPublicKey 设置加密公钥
func WithEncryptorPublicKey(publicKey string) EncryptorOption {
	return func(o *encryptorOptions) { o.publicKey = publicKey }
}
