package geetest

import (
	"testing"
	"crypto/md5"
	"fmt"
	"encoding/hex"
)

func TestGeetLib_GetVersionInfo(t *testing.T) {
	gt := NewGeetLib("abc", "123")
	fmt.Println(gt.getFailPreProcessRes())
}

func TestMd5(t *testing.T) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte("123456"))

	cipherStr := md5Ctx.Sum(nil)
	fmt.Print(hex.EncodeToString(cipherStr))
}
