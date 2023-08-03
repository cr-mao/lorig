package ecc

import (
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/utils/xconv"
)

const (
	defaultEncryptorShareInfo1Key = "crypto.ecc.encryptor.s1"
	defaultEncryptorShareInfo2Key = "crypto.ecc.encryptor.s2"
	defaultEncryptorPublicKeyKey  = "crypto.ecc.encryptor.publicKey"
)

type EncryptorOption func(o *encryptorOptions)

type encryptorOptions struct {
	// 共享信息。加解密时必需一致
	// 默认为空
	s1 []byte

	// 共享信息。加解密时必需一致
	// 默认为空
	s2 []byte

	// 公钥。可设置文件路径或公钥串
	publicKey string
}

func defaultEncryptorOptions() *encryptorOptions {
	return &encryptorOptions{
		s1:        []byte(conf.Get(defaultEncryptorShareInfo1Key, "")),
		s2:        []byte(conf.Get(defaultEncryptorShareInfo2Key, "")),
		publicKey: conf.Get(defaultEncryptorPublicKeyKey, ""),
	}
}

// WithEncryptorShareInfo 设置共享信息
func WithEncryptorShareInfo(s1, s2 string) EncryptorOption {
	return func(o *encryptorOptions) { o.s1, o.s2 = xconv.StringToBytes(s1), xconv.StringToBytes(s2) }
}

// WithEncryptorPublicKey 设置加密公钥
func WithEncryptorPublicKey(publicKey string) EncryptorOption {
	return func(o *encryptorOptions) { o.publicKey = publicKey }
}
