# xorsimd

[![GoDoc][1]][2] [![MIT licensed][3]][4] [![Build Status][5]][6] [![Go Report Card][7]][8] [![Sourcegraph][9]][10]

[1]: https://godoc.org/github.com/templexxx/xorsimd?status.svg
[2]: https://godoc.org/github.com/templexxx/xorsimd
[3]: https://img.shields.io/badge/license-MIT-blue.svg
[4]: LICENSE
[5]: https://github.com/templexxx/xorsimd/workflows/unit-test/badge.svg
[6]: https://github.com/templexxx/xorsimd
[7]: https://goreportcard.com/badge/github.com/templexxx/xorsimd
[8]: https://goreportcard.com/report/github.com/templexxx/xorsimd
[9]: https://sourcegraph.com/github.com/templexxx/xorsimd/-/badge.svg
[10]: https://sourcegraph.com/github.com/templexxx/xorsimd?badge

`xorsimd` provides high-throughput XOR operations for byte slices.
On amd64 it selects AVX-512/AVX2/SSE2 implementations, and falls back to
portable Go code on other architectures.

## Install

```bash
go get github.com/templexxx/xorsimd
```

## Quick Start

```go
package main

import "github.com/templexxx/xorsimd"

func main() {
	a := []byte{1, 2, 3, 4}
	b := []byte{9, 8, 7, 6}
	dst := make([]byte, 4)

	// dst = a ^ b
	n := xorsimd.Bytes(dst, a, b)
	_ = n

	// dst = src[0] ^ src[1] ^ src[2] ...
	src := [][]byte{a, b, []byte{1, 1, 1, 1}}
	n = xorsimd.Encode(dst, src)
	_ = n
}
```

## API Notes

- `Encode(dst, src)` XORs all rows in `src` into `dst` and returns the number
  of bytes processed.
- `Encode` requires `len(src) >= 1`.
- `Bytes(dst, a, b)` is a convenience wrapper around `Encode`.
- `Bytes8`, `Bytes16`, `Bytes8Align`, and `Bytes16Align` are fixed-size helpers.
- `BytesA` processes `len(a)` bytes; `BytesB` processes `len(b)` bytes.

All APIs may reuse overlapping slices. For correctness, callers must ensure the
input slices are long enough for the selected function variant.

## Performance

Performance mainly depends on:

- CPU SIMD instruction support.
- Number of input source rows.
- Slice size.

Benchmark formula:

`I/O = (src_num + 1) * vector_size / cost`

Test environment:

- AWS c5d.xlarge
- Intel Xeon Platinum 8124M @ 3.00GHz
- Single physical core

| Src Num | Vector Size | AVX512 I/O (MB/s) | AVX2 I/O (MB/s) | SSE2 I/O (MB/s) |
| --- | --- | --- | --- | --- |
| 5 | 4KB | 270403.73 | 142825.25 | 74443.91 |
| 5 | 1MB | 26948.34 | 26887.37 | 26950.65 |
| 5 | 8MB | 17881.32 | 17212.56 | 16402.97 |
| 10 | 4KB | 190445.30 | 102953.59 | 53244.04 |
| 10 | 1MB | 26424.44 | 26618.65 | 26094.39 |
| 10 | 8MB | 15471.31 | 14866.72 | 13565.80 |

## Run Tests and Benchmarks

```bash
go test ./...
go test -bench=. -benchmem
```
