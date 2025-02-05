package dilithium

import "crypto/rand"

func RandomBytes(out []uint8, outlen uint64) {
	var i uint64
	token := make([]byte, outlen)
	rand.Read(token)
	for i = 0; i < outlen; i++ {
		out[i] = token[i]
	}
}
