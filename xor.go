// Copyright (c) 2019. Temple3x (temple3x@gmail.com)
//
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package xorsimd

import (
	"unsafe"

	"github.com/templexxx/cpu"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))

// EnableAVX512 controls whether AVX-512 is considered during feature detection.
// Note: instruction-path selection is decided at package initialization, so
// changing this value afterwards does not refresh the selected backend.
var EnableAVX512 = true

// cpuFeature indicates which instruction set will be used.
var cpuFeature = getCPUFeature()

const (
	avx512 = iota
	avx2
	sse2
	generic
)

// TODO: Add ARM feature...

func hasAVX512() (ok bool) {

	return cpu.X86.HasAVX512VL &&
		cpu.X86.HasAVX512BW &&
		cpu.X86.HasAVX512F &&
		cpu.X86.HasAVX512DQ
}

// Encode XORs all source rows into dst.
// The source and destination may overlap.
//
// Encode returns the number of bytes processed, which is the minimum length of
// dst and every source row.
//
// Callers must provide at least one source row: len(src) >= 1.
func Encode(dst []byte, src [][]byte) (n int) {
	n = checkLen(dst, src)
	if n == 0 {
		return
	}

	dst = dst[:n]
	for i := range src {
		src[i] = src[i][:n]
	}

	if len(src) == 1 {
		copy(dst, src[0])
		return
	}

	encode(dst, src)
	return
}

func checkLen(dst []byte, src [][]byte) int {
	n := len(dst)
	for i := range src {
		if len(src[i]) < n {
			n = len(src[i])
		}
	}

	if n <= 0 {
		return 0
	}
	return n
}

// Bytes XORs a and b into dst.
// The source and destination may overlap.
//
// Bytes returns the number of bytes processed, which is the minimum of
// len(dst), len(a), and len(b).
func Bytes(dst, a, b []byte) int {
	return Encode(dst, [][]byte{a, b})
}
