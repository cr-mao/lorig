package rsa

import (
	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/crypto/hash"
	"strings"
)

const (
	defaultVerifierHashKey      = "crypto.rsa.verifier.hash"
	defaultVerifierPaddingKey   = "crypto.rsa.verifier.padding"
	defaultVerifierPublicKeyKey = "crypto.rsa.verifier.publicKey"
)

type VerifierOption func(o *verifierOption)

type verifierOption struct {
	// hash算法。支持sha1、sha224、sha256、sha384、sha512
	// 默认为sha256
	hash hash.Hash

	// 填充规则。支持NORMAL和OAEP
	// 默认为NORMAL
	padding SignPadding

	// 公钥。可设置文件路径或公钥串
	publicKey string
}

func defaultVerifierOptions() *verifierOption {
	return &verifierOption{
		hash:      hash.Hash(strings.ToLower(conf.Get(defaultVerifierHashKey, ""))),
		padding:   SignPadding(strings.ToUpper(conf.Get(defaultVerifierPaddingKey, ""))),
		publicKey: conf.Get(defaultVerifierPublicKeyKey, ""),
	}
}

// WithVerifierHash 设置加密hash算法
func WithVerifierHash(hash hash.Hash) VerifierOption {
	return func(o *verifierOption) { o.hash = hash }
}

// WithVerifierPadding 设置加密填充规则
func WithVerifierPadding(padding SignPadding) VerifierOption {
	return func(o *verifierOption) { o.padding = padding }
}

// WithVerifierPublicKey 设置验签公钥
func WithVerifierPublicKey(publicKey string) VerifierOption {
	return func(o *verifierOption) { o.publicKey = publicKey }
}
