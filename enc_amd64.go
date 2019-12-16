/*
 * Copyright (c) 2019. Temple3x (temple3x@gmail.com)
 *
 * Use of this source code is governed by the MIT License
 * that can be found in the LICENSE file.
 */

package xorsimd

func encode(dst []byte, src [][]byte, feature int) {

	switch feature {
	case avx512:
		encodeAVX512(dst, src)
	case avx2:
		encodeAVX2(dst, src)
	default:
		encodeSSE2(dst, src)
	}
	return
}

//go:noescape
func encodeAVX512(dst []byte, src [][]byte)

//go:noescape
func encodeAVX2(dst []byte, src [][]byte)

//go:noescape
func encodeSSE2(dst []byte, src [][]byte)
