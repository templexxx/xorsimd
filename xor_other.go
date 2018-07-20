// +build !amd64 noasm

package xor

func encode(dst []byte, src [][]byte) {
	encodeNoSIMD(dst, src)
}