/*
 * Copyright (c) 2019. Temple3x (temple3x@gmail.com)
 *
 * Use of this source code is governed by the MIT License
 * that can be found in the LICENSE file.
 */

package xorsimd

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const (
	kb = 1024
	mb = 1024 * 1024

	testSize = kb
)

func TestRS_Encode(t *testing.T) {
	max := testSize

	switch getCPUFeature() {
	case avx512:
		testEncode(t, max, sse2, -1)
		testEncode(t, max, avx2, sse2)
		testEncode(t, max, avx512, avx2)
	case avx2:
		testEncode(t, max, sse2, -1)
		testEncode(t, max, avx2, sse2)
	case sse2:
		testEncode(t, max, sse2, -1)
	case generic:
		testEncode(t, max, generic, -1)
	}
}

func testEncode(t *testing.T, maxSize, feat, cmpFeat int) {

	rand.Seed(time.Now().UnixNano())
	srcN := randIntn(10, 2) // Cannot be 1, see func encode(dst []byte, src [][]byte, feature int).

	fs := featToStr(feat)
	for size := 1; size <= maxSize; size++ {
		exp := make([]byte, size)
		src := make([][]byte, srcN)
		for j := 0; j < srcN; j++ {
			src[j] = make([]byte, size)
			fillRandom(src[j])
		}

		if cmpFeat < 0 {
			encodeTested(exp, src)
		} else {
			encode(exp, src, cmpFeat)
		}

		act := make([]byte, size)
		encode(act, src, feat)

		if !bytes.Equal(exp, act) {
			t.Fatalf("%s mismatched with %s, src_num: %d, size: %d",
				fs, featToStr(cmpFeat), srcN, size)
		}
	}

	t.Logf("%s pass src_num:%d, max_size: %d",
		fs, srcN, maxSize)
}

func featToStr(f int) string {
	switch f {
	case avx512:
		return "AVX512"
	case avx2:
		return "AVX2"
	case sse2:
		return "SSE2"
	case generic:
		return "Generic"
	default:
		return "Tested"
	}
}

func encodeTested(dst []byte, src [][]byte) {

	n := len(dst)
	for i := 0; i < n; i++ {
		s := src[0][i]
		for j := 1; j < len(src); j++ {
			s ^= src[j][i]
		}
		dst[i] = s
	}
}

// randIntn returns, as an int, a non-negative pseudo-random number in [min,n)
// from the default Source.
func randIntn(n, min int) int {
	m := rand.Intn(n)
	if m < min {
		m = min
	}
	return m
}

func BenchmarkEncode(b *testing.B) {
	sizes := []int{4 * kb, mb, 8 * mb}

	srcNums := []int{5, 10}

	var feats []int
	switch getCPUFeature() {
	case avx512:
		feats = append(feats, avx512)
		feats = append(feats, avx2)
		feats = append(feats, sse2)
	case avx2:
		feats = append(feats, avx2)
		feats = append(feats, sse2)
	case sse2:
		feats = append(feats, sse2)
	default:
		feats = append(feats, generic)
	}

	b.Run("", benchEncRun(benchEnc, srcNums, sizes, feats))
}

func benchEncRun(f func(*testing.B, int, int, int), srcNums, sizes, feats []int) func(*testing.B) {
	return func(b *testing.B) {
		for _, feat := range feats {
			for _, srcNum := range srcNums {
				for _, size := range sizes {
					b.Run(fmt.Sprintf("(%d+1)-%s-%s", srcNum, byteToStr(size), featToStr(feat)), func(b *testing.B) {
						f(b, srcNum, size, feat)
					})
				}
			}
		}
	}
}

func benchEnc(b *testing.B, dataNum, size, feat int) {
	dst := make([]byte, size)
	src := make([][]byte, dataNum)
	for i := 0; i < dataNum; i++ {
		src[i] = make([]byte, size)
		fillRandom(src[i])
	}

	b.SetBytes(int64(dataNum * size))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encode(dst, src, feat)
	}
}

func fillRandom(p []byte) {
	rand.Read(p)
}

func byteToStr(n int) string {
	if n >= mb {
		return fmt.Sprintf("%dMB", n/mb)
	}

	return fmt.Sprintf("%dKB", n/kb)
}
