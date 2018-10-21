package xor

func encode(dst []byte, src [][]byte, feature int) {

	var nonTmp bool
	if len(dst) > nonTmpSize {
		nonTmp = true
	}

	switch feature {
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
