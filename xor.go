package xor

import "github.com/templexxx/cpu"

// AVX512 may slow down CPU Clock (maybe not).
// TODO need more research:
// https://lemire.me/blog/2018/04/19/by-how-much-does-avx-512-slow-down-your-cpu-a-first-experiment/
var EnableAVX512 = true

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

	if len(src) == 1 {
		copy(dst, src[0])
		return
	}

	f := getCPUFeature()
	encode(dst, src, f)
	return
}

const (
	avx512 = iota
	avx2
	sse2
	generic
)

// TODO: add arm feature...
func getCPUFeature() int {
	if hasAVX512 && EnableAVX512 {
		return avx512
	} else if cpu.X86.HasAVX2 {
		return avx2
	} else {
		return sse2 // amd64 must has sse2
	}
}

var hasAVX512 = cpu.X86.HasAVX512VL &&
	cpu.X86.HasAVX512BW &&
	cpu.X86.HasAVX512F &&
	cpu.X86.HasAVX512DQ

// if size > nonTmpSize, it will use Non-Temporal Hint store
const nonTmpSize = 8 * 1024 // depends on CPU Cache Size

// Update updates parity when there is one data vector being changed.
// Returns the minimum length of oldData and newData and parity.
// e.g.
// oldData ⊕ A = parity
// Pass oldData, newData, parity as args,
// then the parity will be:
// parity = newData ⊕ A
func Update(oldData, newData, parity []byte) int {
	src := make([][]byte, 3)
	src[0], src[1], src[2] = oldData, newData, parity
	return Encode(parity, src)
}

// Replace replaces oldData vectors with 0 or replaces 0 with newData vectors.
// Returns the minimum length of data[i] and parity.
//
// It's used in two situations:
// 1. We didn't have enough data for filling in a stripe, but still did xor encode,
// we need replace several zero vectors with new vectors which have data after we get enough data finally.
// 2. After compact, we may have several useless vectors in a stripe,
// we need replaces these useless vectors with zero vectors for free space.
func Replace(parity []byte, data [][]byte) int {
	vects := make([][]byte, len(data)+1)
	vects[0] = parity
	for i := range data {
		vects[i+1] = data[i]
	}
	return Encode(parity, vects)
}
