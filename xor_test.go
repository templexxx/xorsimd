package xor

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/templexxx/cpu"
)

const (
	kb = 1024
	mb = 1024 * 1024

	// testMaxVects is the maximum length of source vectors,
	// in test process, we will pick up a vects cnt randomly from [1, testMaxVects].
	testMaxVects = 17
	testDataCnt  = 3 // len(data)
	// testSize is a single vector's size.
	// it will cover Non-Temporal Hint & loop256b, loop128b, loop64b, loop_8b, loop_1b(in *_amd64.s)
	testSize = nonTmpSize + 256 + 16 + 8 + 7
)

func TestVerifyEncode(t *testing.T) {
	switch runtime.GOARCH {

	case `amd64`:
		verifyEncode(t, sse2)
		if cpu.X86.HasAVX2 {
			verifyEncode(t, avx2)
			if hasAVX512 {
				verifyEncode(t, avx512)
			}
		}

	default:
		verifyEncode(t, generic)
	}
}

func verifyEncode(t *testing.T, feature int) {
	rand.Seed(time.Now().UnixNano())

	for size := 1; size <= testSize; size++ {
		exp := make([]byte, size)
		vcnt := randIntn(testMaxVects, 2)
		src := make([][]byte, vcnt)
		for j := 0; j < vcnt; j++ {
			src[j] = make([]byte, size)
			fillRandom(src[j])
		}
		for i := 0; i < size; i++ {
			s := src[0][i]
			for j := 1; j < vcnt; j++ {
				s ^= src[j][i]
			}
			exp[i] = s
		}

		fs := featureToString(feature)
		act := make([]byte, size)
		encode(act, src, feature)
		if !bytes.Equal(exp, act) {
			t.Fatalf("encode mismatch; vect cnt: %d, size: %d; ext: %s", vcnt, size, fs)
		}
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

func featureToString(f int) (s string) {
	switch f {
	case avx512:
		s = "avx512"
	case avx2:
		s = "avx2"
	case sse2:
		s = "sse2"
	}
	return s
}

func TestVerifyUpdate(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for size := 1; size <= testSize; size++ {
		vcnt := randIntn(testMaxVects, 2)
		src0 := make([][]byte, vcnt)
		src1 := make([][]byte, vcnt)
		for j := 0; j < vcnt; j++ {
			src0[j] = make([]byte, size)
			fillRandom(src0[j])
			src1[j] = make([]byte, size)
			copy(src1[j], src0[j])
		}

		act := make([]byte, size)
		Encode(act, src0)

		newData := make([]byte, size)
		fillRandom(newData)
		idx := rand.Intn(vcnt)
		Update(src0[idx], newData, act)

		src1[idx] = newData
		exp := make([]byte, size)
		Encode(exp, src1)
		if !bytes.Equal(exp, act) {
			t.Fatalf("update mismatch; vect cnt: %d, size: %d", vcnt, size)
		}
	}
}

func TestVerifyReplace(t *testing.T) {
	verifyReplaceZero(t)
	verifyReplaceData(t)
}

// zero vectors -> data vectors.
func verifyReplaceZero(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for size := 1; size <= testSize; size++ {
		vcnt := randIntn(testMaxVects, 2)
		zeroV := randIntn(vcnt, 1)
		src := make([][]byte, vcnt)
		for i := range src {
			src[i] = make([]byte, size)
			if i >= zeroV {
				fillRandom(src[i])
			}
		}
		act := make([]byte, size)
		Encode(act, src)
		data := make([][]byte, zeroV)
		for i := range data {
			data[i] = make([]byte, size)
			fillRandom(data[i])
		}
		Replace(act, data)

		exp := make([]byte, size)
		for i := range data {
			src[i] = data[i]
		}
		Encode(exp, src)

		if !bytes.Equal(act, exp) {
			t.Fatalf("replace zero mismatch; vect cnt: %d, zero cnt: %d, size: %d", vcnt, zeroV, size)
		}
	}
}

// data vectors -> zero vectors.
func verifyReplaceData(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for size := 1; size <= testSize; size++ {
		vcnt := randIntn(testMaxVects, 2)
		zeroV := randIntn(vcnt, 1)
		src := make([][]byte, vcnt)
		for i := range src {
			src[i] = make([]byte, size)
			fillRandom(src[i])
		}
		act := make([]byte, size)
		Encode(act, src)
		data := make([][]byte, zeroV)
		for i := range data {
			data[i] = src[i]
		}
		Replace(act, data)

		exp := make([]byte, size)
		for i := range data {
			src[i] = make([]byte, size)
		}
		Encode(exp, src)

		if !bytes.Equal(act, exp) {
			t.Fatalf("replace data mismatch; vect cnt: %d, zero cnt: %d, size: %d", vcnt, zeroV, size)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	sizes := []int{4 * kb, 64 * kb}
	b.Run("", benchEncRun(benchEnc, 5, sizes))
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
		fillRandom(data[i])
	}
	Encode(parity, data)
	b.SetBytes(int64(dataCnt * size))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(parity, data)
	}
}

// Copy from https://github.com/klauspost/reedsolomon/blob/master/reedsolomon_test.go
// Thank him for his contribution.
func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
}
