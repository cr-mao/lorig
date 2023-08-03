package ecc

import (
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/crypto/hash"
	"strings"
)

const (
	defaultVerifierHashKey      = "crypto.rsa.verifier.hash"
	defaultVerifierDelimiterKey = "crypto.rsa.verifier.delimiter"
	defaultVerifierPublicKeyKey = "crypto.rsa.verifier.publicKey"
)

type VerifierOption func(o *verifierOption)

type verifierOption struct {
	// hash算法。支持sha1、sha224、sha256、sha384、sha512
	// 默认为sha256
	hash hash.Hash

	// 签名分隔符。
	delimiter string

	// 公钥。可设置文件路径或公钥串
	publicKey string
}

func defaultVerifierOptions() *verifierOption {
	return &verifierOption{
		hash:      hash.Hash(strings.ToLower(conf.Get(defaultVerifierHashKey))),
		delimiter: conf.Get(defaultVerifierDelimiterKey, " "),
		publicKey: conf.Get(defaultVerifierPublicKeyKey),
	}
}

// WithVerifierHash 设置加密hash算法
func WithVerifierHash(hash hash.Hash) VerifierOption {
	return func(o *verifierOption) { o.hash = hash }
}

// WithVerifierDelimiter 设置签名分割符
func WithVerifierDelimiter(delimiter string) VerifierOption {
	return func(o *verifierOption) { o.delimiter = delimiter }
}

// WithVerifierPublicKey 设置验签公钥
func WithVerifierPublicKey(publicKey string) VerifierOption {
	return func(o *verifierOption) { o.publicKey = publicKey }
}
