package xor

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"io"
	"testing"

	"github.com/templexxx/cpu"
)

const (
	kb          = 1024
	mb          = 1024 * 1024
	testDataCnt = 3 // len(data)
	// it will cover Non-Temporal Hint & loop128b, loop64b, loop_8b, loop_1b(in xor_amd64.s)
	verifySize = nonTmpSize + 256 + 16 + 8 + 7
)

func TestVerifyEncode(t *testing.T) {
	if useAVX512() {
		verifyEncode(t, testDataCnt, avx512)
		verifyEncode(t, testDataCnt, avx2)
		verifyEncode(t, testDataCnt, sse2)
	}
	if cpu.X86.HasAVX2 {
		verifyEncode(t, testDataCnt, avx2)
		verifyEncode(t, testDataCnt, sse2)
	} else {
		verifyEncode(t, testDataCnt, sse2)
	}
}

func verifyEncode(t *testing.T, dataCnt int, cpuFeature int) {
	for size := 1; size <= verifySize; size++ {
		expect := make([]byte, size)
		result := make([]byte, size)
		src := make([][]byte, dataCnt)
		for j := 0; j < dataCnt; j++ {
			src[j] = make([]byte, size)
			err := fillRandom(src[j])
			if err != nil {
				t.Fatal(err)
			}
		}
		for j := 0; j < size; j++ {
			expect[j] = src[0][j] ^ src[1][j]
		}
		for j := 2; j < dataCnt; j++ {
			for k := 0; k < size; k++ {
				expect[k] ^= src[j][k]
			}
		}

		var cpuStr string
		switch cpuFeature {
		case avx512:
			cpuStr = "avx512"
		case avx2:
			cpuStr = "avx2"
		case sse2:
			cpuStr = "sse2"
		}
		encode(result, src, cpuFeature)
		if !bytes.Equal(expect, result) {
			t.Fatalf("encode mismatch; size: %d; ext: %s", size, cpuStr)
		}
	}
}

func TestVerifyUpdate(t *testing.T) {
	verifyUpdate(t, testDataCnt)
}

func verifyUpdate(t *testing.T, dataCnt int) {
	for size := 1; size <= verifySize; size++ {
		expect := make([]byte, size)
		result := make([]byte, size)
		src0 := make([][]byte, dataCnt)
		src1 := make([][]byte, dataCnt)
		for j := 0; j < dataCnt; j++ {
			src0[j] = make([]byte, size)
			err := fillRandom(src0[j])
			if err != nil {
				t.Fatal(err)
			}
			src1[j] = make([]byte, size)
			copy(src1[j], src0[j])
		}
		Encode(result, src0)
		newData := make([]byte, size)
		err := fillRandom(newData)
		if err != nil {
			t.Fatal(err)
		}
		Update(src0[0], newData, result)
		src1[0] = newData
		Encode(expect, src1)
		if !bytes.Equal(expect, result) {
			t.Fatal("update mismatch")
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
		err := fillRandom(data[i])
		if err != nil {
			b.Fatal(err)
		}
	}
	Encode(parity, data)
	b.SetBytes(int64(dataCnt * size))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(parity, data)
	}
}

func fillRandom(p []byte) (err error) {
	_, err = io.ReadFull(crand.Reader, p)
	return
}
