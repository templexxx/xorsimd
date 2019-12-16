// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !amd64

package xorsimd

import (
	"runtime"
	"unsafe"
)

func encode(dst []byte, src [][]byte, feature int) {
	if supportsUnaligned {
		fastXORBytes(dst, src, len(dst))
	} else {
		// TODO(hanwen): if (dst, a, b) have common alignment
		// we could still try fastXORBytes. It is not clear
		// how often this happens, and it's only worth it if
		// the block encryption itself is hardware
		// accelerated.
		safeXORBytes(dst, src, len(dst))
	}

}

// fastXORBytes xor in bulk. It only works on architectures that
// support unaligned read/writes.
func fastXORBytes(dst []byte, src [][]byte, n int) {
	w := n / wordSize
	if w > 0 {
		wordBytes := w * wordSize

		wordAlignSrc := make([][]byte, len(src))
		for i := range src {
			wordAlignSrc[i] = src[i][:wordBytes]
		}
		fastXORWords(dst[:wordBytes], wordAlignSrc)
	}
	for i := n - n%wordSize; i < n; i++ {
		s := src[0][i]
		for j := 1; j < len(src); j++ {
			s ^= src[j][i]
		}
		dst[i] = s
	}
}

const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

func safeXORBytes(dst []byte, src [][]byte, n int) {
	for i := 0; i < n; i++ {
		s := src[0][i]
		for j := 1; j < len(src); j++ {
			s ^= src[j][i]
		}
		dst[i] = s
	}
}

// fastXORWords XORs multiples of 4 or 8 bytes (depending on architecture.)
// The arguments are assumed to be of equal length.
func fastXORWords(dst []byte, src [][]byte) {
	dw := *(*[]uintptr)(unsafe.Pointer(&dst))
	sw := make([][]uintptr, len(src))
	for i := range src {
		sw[i] = *(*[]uintptr)(unsafe.Pointer(&src[i]))
	}

	n := len(dst) / wordSize
	for i := 0; i < n; i++ {
		s := sw[0][i]
		for j := 1; j < len(sw); j++ {
			s ^= sw[j][i]
		}
		dw[i] = s
	}
}
