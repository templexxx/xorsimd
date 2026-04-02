// Copyright (c) 2019. Temple3x (temple3x@gmail.com)
//
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package xorsimd

import "github.com/templexxx/cpu"

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

func getCPUFeature() int {
	if hasAVX512() && EnableAVX512 {
		return avx512
	} else if cpu.X86.HasAVX2 {
		return avx2
	} else {
		return sse2 // amd64 must have sse2
	}
}

// Bytes8 XORs exactly 8 bytes from a and b into dst.
// Each slice must have length >= 8, otherwise it panics.
func Bytes8(dst, a, b []byte) {

	bytes8(&dst[0], &a[0], &b[0])
}

// Bytes16 XORs exactly 16 bytes from a and b into dst.
// Each slice must have length >= 16, otherwise it panics.
func Bytes16(dst, a, b []byte) {

	bytes16(&dst[0], &a[0], &b[0])
}

// Bytes8Align XORs exactly 8 bytes from a and b into dst.
// Each slice must have length >= 8, otherwise it panics.
// On amd64, explicit alignment is not required; this function exists for API
// compatibility with non-amd64 implementations.
func Bytes8Align(dst, a, b []byte) {

	bytes8(&dst[0], &a[0], &b[0])
}

// Bytes16Align XORs exactly 16 bytes from a and b into dst.
// Each slice must have length >= 16, otherwise it panics.
// On amd64, explicit alignment is not required; this function exists for API
// compatibility with non-amd64 implementations.
func Bytes16Align(dst, a, b []byte) {

	bytes16(&dst[0], &a[0], &b[0])
}

// BytesA XORs len(a) bytes from a and b into dst.
// Callers must ensure len(dst) >= len(a) and len(b) >= len(a).
//
// This helper is intended for small slices where setup overhead dominates.
// For larger slices, Bytes is usually faster.
func BytesA(dst, a, b []byte) {

	bytesN(&dst[0], &a[0], &b[0], len(a))
}

// BytesB XORs len(b) bytes from a and b into dst.
// Callers must ensure len(dst) >= len(b) and len(a) >= len(b).
//
// This helper is intended for small slices where setup overhead dominates.
// For larger slices, Bytes is usually faster.
func BytesB(dst, a, b []byte) {

	bytesN(&dst[0], &a[0], &b[0], len(b))
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
