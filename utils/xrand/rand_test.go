package xrand_test

import (
	"testing"

	"github.com/cr-mao/lorig/utils/xrand"
)

func Test_Str(t *testing.T) {
	t.Log(xrand.Str("您好中国AJCKEKD", 5))
}

func Test_Symbols(t *testing.T) {
	t.Log(xrand.Symbols(5))
}

func Test_Int(t *testing.T) {
	t.Log(xrand.Int(1, 2))
}

func Test_Float32(t *testing.T) {
	t.Log(xrand.Float32(-50, 5))
}
