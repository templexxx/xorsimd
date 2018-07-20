package xor

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	cpu "github.com/templexxx/cpufeat"
)

const (
	kb          = 1024
	mb          = 1024 * 1024
	testDataCnt = 3 // len(data)
	// it will cover Non-Temporal Hint & loop128b, loop64b, loop_8b, loop_1b(in xor_amd64.s)
	verifySize = nonTmpSize + 256 + 16 + 8 + 7
)

func TestVerifyEncode(t *testing.T) {
	if cpu.X86.HasAVX512 {
		verifyEncode(t, testDataCnt, avx512)
		verifyEncode(t, testDataCnt, avx2)
		verifyEncode(t, testDataCnt, sse2)
		verifyEncode(t, testDataCnt, base)
	}
	if cpu.X86.HasAVX2 {
		verifyEncode(t, testDataCnt, avx2)
		verifyEncode(t, testDataCnt, sse2)
		verifyEncode(t, testDataCnt, base)
	} else if cpu.X86.HasSSE2 {
		verifyEncode(t, testDataCnt, sse2)
		verifyEncode(t, testDataCnt, base)
	} else {
		verifyEncode(t, testDataCnt, base)
	}
}

func verifyEncode(t *testing.T, dataCnt int, cpuFeature int) {
	for size := 1; size <= verifySize; size++ {
		expect := make([]byte, size)
		result := make([]byte, size)
		src := make([][]byte, dataCnt)
		for j := 0; j < dataCnt; j++ {
			src[j] = make([]byte, size)
			rand.Seed(int64(j))
			fillRandom(src[j])
		}
		for j := 0; j < size; j++ {
			expect[j] = src[0][j] ^ src[1][j]
		}
		for j := 2; j < dataCnt; j++ {
			for k := 0; k < size; k++ {
				expect[k] ^= src[j][k]
			}
		}
		var nonTmp bool
		if size > nonTmpSize {
			nonTmp = true
		}

		switch cpuFeature {
		case avx512:
			if nonTmp {
				encodeAVX512NonTmp(result, src)
			} else {
				encodeAVX512(result, src)
			}
		case avx2:
			if nonTmp {
				encodeAVX2NonTmp(result, src)
			} else {
				encodeAVX2(result, src)
			}
		case sse2:
			if nonTmp {
				encodeSSE2NonTmp(result, src)
			} else {
				encodeSSE2(result, src)
			}
		case base:
			encodeNoSIMD(result, src)
		}

		if !bytes.Equal(expect, result) {
			t.Fatalf("encode mismatch; size: %d; ext: %s", size, cpuFeature)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	sizes := []int{4 * kb, 64 * kb, mb}
	b.Run("", benchEncRun(benchEnc, testDataCnt, sizes))
}

func benchEncRun(f func(*testing.B, int, int), dataCnt int, sizes []int) func(*testing.B) {
	return func(b *testing.B) {
		for _, s := range sizes {
			b.Run(fmt.Sprintf("%d+1_%dKB", dataCnt, s/kb), func(b *testing.B) {
				f(b, dataCnt, s)
			})
		}
	}
}

func benchEnc(b *testing.B, dataCnt, size int) {
	parity := make([]byte, size)
	data := make([][]byte, dataCnt)
	for i := 0; i < dataCnt; i++ {
		data[i] = make([]byte, size)
		rand.Seed(int64(i))
		fillRandom(data[i])
	}
	Encode(parity, data)
	b.SetBytes(int64(dataCnt * size))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(parity, data)
	}
}

func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
}
