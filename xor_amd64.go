package xor

import "github.com/templexxx/cpu"

const nonTmpSize = 8 * 1024 // depends on CPU Cache Size

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
	default:
		if nonTmp {
			encodeSSE2NonTmp(dst, src)
		} else {
			encodeSSE2(dst, src)
		}
	}
	return
}

const (
	avx512 = iota
	avx2
	sse2
)

func getCPUFeature() int {
	if useAVX512() {
		return avx512
	} else if cpu.X86.HasAVX2 {
		return avx2
	}  else {
		return sse2	// amd64 must has sse2
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
