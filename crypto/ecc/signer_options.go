package ecc

import (
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/crypto/hash"
	"strings"
)

const (
	defaultSignerHashKey       = "crypto.rsa.signer.hash"
	defaultSignerDelimiterKey  = "crypto.rsa.signer.delimiter"
	defaultSignerPrivateKeyKey = "crypto.rsa.signer.privateKey"
)

type SignerOption func(o *signerOptions)

type signerOptions struct {
	// hash算法。支持sha1、sha224、sha256、sha384、sha512
	// 默认为sha256
	hash hash.Hash

	// 签名分隔符。
	delimiter string

	// 私钥。可设置文件路径或私钥串
	privateKey string
}

func defaultSignerOptions() *signerOptions {
	return &signerOptions{
		hash:       hash.Hash(strings.ToLower(conf.Get(defaultSignerHashKey, ""))),
		delimiter:  conf.Get(defaultSignerDelimiterKey, " "),
		privateKey: conf.Get(defaultSignerPrivateKeyKey),
	}
}

// WithSignerHash 设置加密hash算法
func WithSignerHash(hash hash.Hash) SignerOption {
	return func(o *signerOptions) { o.hash = hash }
}

// WithSignerDelimiter 设置签名分割符
func WithSignerDelimiter(delimiter string) SignerOption {
	return func(o *signerOptions) { o.delimiter = delimiter }
}

// WithSignerPrivateKey 设置解密私钥
func WithSignerPrivateKey(privateKey string) SignerOption {
	return func(o *signerOptions) { o.privateKey = privateKey }
}
