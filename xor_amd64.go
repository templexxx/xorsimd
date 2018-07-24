package xor

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
