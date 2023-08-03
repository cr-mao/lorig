package rsa_test

import (
	"fmt"
	"github.com/cr-mao/lorig/crypto/rsa"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key, err := rsa.GenerateKey(256)
	if err != nil {
		t.Fatal(err)
	}

	v, err := key.MarshalPublicKey(rsa.PKCS1)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(v))
}

func TestKey_SaveKeyPair(t *testing.T) {
	key, err := rsa.GenerateKey(1024)
	if err != nil {
		t.Fatal(err)
	}

	err = key.SaveKeyPair(rsa.PKCS1, "./pem", "key.pem")
	if err != nil {
		t.Fatal(err)
	}
}
