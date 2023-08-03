package ecc

import (
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

type Decryptor struct {
	err        error
	opts       *decryptorOptions
	privateKey *ecies.PrivateKey
}

var DefaultDecryptor = NewDecryptor()

func NewDecryptor(opts ...DecryptorOption) *Decryptor {
	o := defaultDecryptorOptions()
	for _, opt := range opts {
		opt(o)
	}

	d := &Decryptor{opts: o}
	d.privateKey, d.err = parseECIESPrivateKey(d.opts.privateKey)

	return d
}

// Name 名称
func (d *Decryptor) Name() string {
	return Name
}

// Decrypt 解密
func (d *Decryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	if d.err != nil {
		return nil, d.err
	}

	return d.privateKey.Decrypt(ciphertext, d.opts.s1, d.opts.s2)
}

// Decrypt 解密
func Decrypt(ciphertext []byte) ([]byte, error) {
	return DefaultDecryptor.Decrypt(ciphertext)
}
