package scrypt

import (
	"testing"
	"fmt"
	"log"
)

func TestCompareHashAndPassword(t *testing.T) {
	// e.g. r.PostFormValue("password")
	passwordFromForm := "ymex@foxmail.com"

	// Generates a derived key of the form "N$r$p$salt$dk" where N, r and p are defined as per
	// Colin Percival's scrypt paper: http://www.tarsnap.com/scrypt/scrypt.pdf
	// scrypt.Defaults (N=16384, r=8, p=1) makes it easy to provide these parameters, and
	// (should you wish) provide your own values via the scrypt.Params type.
	hash, err := GenerateFromPassword([]byte(passwordFromForm), DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	// Print the derived key with its parameters prepended.
	fmt.Printf("%s\n", hash)
	fmt.Print("加密后长度：",len(string(hash)),"\n")
	// Uses the parameters from the existing derived key. Return an error if they don't match.
	_err := CompareHashAndPassword(hash, []byte(passwordFromForm))
	if _err != nil {
		log.Fatal(_err)
	}else {
		log.Println("密码相同..")
	}

}
