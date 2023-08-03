package ecc

import (
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/utils/xconv"
)

const (
	defaultDecryptorShareInfo1Key = "crypto.ecc.decryptor.s1"
	defaultDecryptorShareInfo2Key = "crypto.ecc.decryptor.s2"
	defaultDecryptorPrivateKeyKey = "crypto.ecc.decryptor.privateKey"
)

type DecryptorOption func(o *decryptorOptions)

type decryptorOptions struct {
	// 共享信息。加解密时必需一致
	// 默认为空
	s1 []byte

	// 共享信息。加解密时必需一致
	// 默认为空
	s2 []byte

	// 私钥。可设置文件路径或私钥串
	privateKey string
}

func defaultDecryptorOptions() *decryptorOptions {
	return &decryptorOptions{
		s1:         []byte(conf.Get(defaultDecryptorShareInfo1Key, "")),
		s2:         []byte(conf.Get(defaultDecryptorShareInfo2Key, "")),
		privateKey: conf.Get(defaultDecryptorPrivateKeyKey, ""),
	}
}

// WithDecryptorShareInfo 设置共享信息
func WithDecryptorShareInfo(s1, s2 string) DecryptorOption {
	return func(o *decryptorOptions) { o.s1, o.s2 = xconv.StringToBytes(s1), xconv.StringToBytes(s2) }
}

// WithDecryptorPrivateKey 设置解密私钥
func WithDecryptorPrivateKey(privateKey string) DecryptorOption {
	return func(o *decryptorOptions) { o.privateKey = privateKey }
}
