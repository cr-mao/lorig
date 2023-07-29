package xfile_test

import (
	"testing"

	"github.com/cr-mao/lorig/utils/xfile"
)

func TestWriteFile(t *testing.T) {
	err := xfile.WriteFile("./run/test.txt", []byte("hello world"))
	if err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}
