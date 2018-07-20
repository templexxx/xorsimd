package xor

import cpu "github.com/templexxx/cpufeat"

const nonTmpSize = 8 * 1024

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

const (
	avx512 = iota
	avx2
	sse2
	base
)

func getCPUFeature() int {
	if cpu.X86.HasAVX512 {
		return avx512
	} else if cpu.X86.HasAVX2 {
		return avx2
	} else if cpu.X86.HasSSE2 {
		return sse2
	} else {
		return base
	}
}

//go:noescape
func encodeAVX512NonTmp(dst []byte, src [][]byte)

//go:noescape
func encodeAVX512(dst []byte, src [][]byte)

//go:noescape
func encodeAVX2NonTmp(dst []byte, src [][]byte)

//go:noescape
func encodeAVX2(dst []byte, src [][]byte)

//go:noescape
func encodeSSE2NonTmp(dst []byte, src [][]byte)

//go:noescape
func encodeSSE2(dst []byte, src [][]byte)

