/**
随机串
 */
package random

import (
	"time"
	"encoding/hex"
	"crypto/rand"
	mr "math/rand"
)

var (
	dash byte = '-'
	defalpha = []byte(`ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz`)
)

//uuid based on random numbers (RFC 4122)
func RandomUUID() (ruuid string) {
	dest := []byte(RandomChars(16))
	setVersion(dest,4)
	SetVariant(dest)
	ruuid = uuid(dest)
	return
}

//generate random []byte by custom chars.
func RandomChars(n int, alphabets ...byte) string {
	if len(alphabets) == 0 {
		alphabets = defalpha[0:]
	}
	dest, err := safeRand(n)
	if err != nil {
		mr.Seed(time.Now().UnixNano())
		for i := range dest {
			dest[i] = alphabets[mr.Intn(len(alphabets))]
		}
	} else {
		for i, d := range dest {
			dest[i] = alphabets[d % byte(len(alphabets))]
		}
	}
	return string(dest)
}


// SetVersion sets version bits.
func setVersion(bs []byte, v byte) {
	bs[6] = (bs[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits as described in RFC 4122.
func SetVariant(bs []byte) {
	bs[8] = (bs[8] & 0xbf) | 0x80
}

func uuid(u []byte) string {
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], u[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u[10:])
	return string(buf)
}

//use package "math/rand" for  safe rand.read
func safeRand(n int) ([]byte, error) {
	var dest = make([]byte, n)
	if _, err := rand.Read(dest); err != nil {
		//rand.Read(dest) It always returns len(p) and a nil error.
		if m, err := mr.Read(dest); n != m || err != nil {

			return dest, err
		}
	}
	return dest, nil
}

