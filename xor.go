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
	base
)

func encode(dst []byte, src [][]byte) {

	var nonTmp bool
	if len(dst) > nonTmpSize {
		nonTmp = true
	}

	feat := getCPUFeature()
	switch feat {
	case avx512:
		if nonTmp {
			encodeAVX512NonTmp(dst, src)
		} else {
			encodeAVX512(dst, src)
		}
	case avx2:
		if nonTmp {
			encodeAVX2NonTmp(dst, src)
		} else {
			encodeAVX2(dst, src)
		}
	case sse2:
		if nonTmp {
			encodeSSE2NonTmp(dst, src)
		} else {
			encodeSSE2(dst, src)
		}
	case base:
		encodeNoSIMD(dst, src)
	}
	return
}

func getCPUFeature() int {
	if useAVX512() {
		return avx512
	} else if cpu.X86.HasAVX2 {
		return avx2
	} else if cpu.X86.HasSSE2 {
		return sse2
	} else {
		return base
	}
}

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
