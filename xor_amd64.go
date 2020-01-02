// Copyright (c) 2019. Temple3x (temple3x@gmail.com)
//
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package xorsimd

func encode(dst []byte, src [][]byte) {

	switch cpuFeature {
	case avx512:
		encodeAVX512(dst, src)
	case avx2:
		encodeAVX2(dst, src)
	default:
		encodeSSE2(dst, src)
	}
	return
}

// Bytes8 XORs of word (8 Bytes).
// The slice arguments a, b, dst's lengths are assumed to be at least 8,
// if not, Bytes8 will panic.
func Bytes8(dst, a, b []byte) {

	bytes8(&dst[0], &a[0], &b[0])
}

// Bytes16 XORs of packed doubleword (16 Bytes).
// The slice arguments a, b, dst's lengths are assumed to be at least 16,
// if not, Bytes16 will panic.
func Bytes16(dst, a, b []byte) {

	bytes16(&dst[0], &a[0], &b[0])
}

// BytesU XORs the bytes in a and b into a
// destination slice. The source and destination may overlap.
//
// BytesU returns the number of bytes encoded, which will be the minimum of
// len(dst), len(a), len(b).
//
// It's used for encoding small bytes slices (< dozens bytes),
// and the slices may not be aligned to 8 bytes or 16 bytes.
// If the length is big, use 'func Bytes(dst, a, b []byte)' instead.
func BytesU(dst, a, b []byte) int {

	n := checkLen(dst, [][]byte{a, b})

	if n == 0 {
		return 0
	}

	bytesN(&dst[0], &a[0], &b[0], n)
	return n
}

//go:noescape
func encodeAVX512(dst []byte, src [][]byte)

//go:noescape
func encodeAVX2(dst []byte, src [][]byte)

//go:noescape
func encodeSSE2(dst []byte, src [][]byte)

//go:noescape
func bytesN(dst, a, b *byte, n int)

//go:noescape
func bytes8(dst, a, b *byte)

//go:noescape
func bytes16(dst, a, b *byte)
