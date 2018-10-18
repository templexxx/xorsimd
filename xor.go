package xor

import "github.com/templexxx/cpu"

var EnableAVX512 = false

// Encode encodes elements from source slice into a
// destination slice. The source and destination may overlap.
// Encode returns the number of elements encoded, which will be the minimum of
// len(src[i]) and len(dst).
func Encode(dst []byte, src [][]byte) (n int) {
	n = len(dst)
	for i := range src {
		if len(src[i]) < n {
			n = len(src[i])
		}
	}

	if n <= 0 {
		return 0
	}

	dst = dst[:n]
	for i := range src {
		src[i] = src[i][:n]
	}
	encode(dst, src)
	return
}

const nonTmpSize = 8 * 1024 // depends on CPU Cache Size

const (
	avx512 = iota
	avx2
	sse2
)

func useAVX512() (ok bool) {
	if !cpu.X86.HasAVX512VL {
		return
	}
	if !cpu.X86.HasAVX512BW {
		return
	}
	if !cpu.X86.HasAVX512F {
		return
	}
	if !cpu.X86.HasAVX512DQ {
		return
	}
	if !EnableAVX512 {
		return
	}
	return true
}
