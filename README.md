# XOR SIMD

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

## Introduction:

>- XOR code engine in pure Go.
>
>- [High Performance](https://github.com/templexxx/xorsimd#performance): 
More than 200GB/s per physics core. 

## Performance

Performance depends mainly on:

>- CPU instruction extension.
>
>- Number of source row vectors.

**Platform:** 

*AWS c5d.xlarge (Intel(R) Xeon(R) Platinum 8124M CPU @ 3.00GHz)*

**All test run on a single Core.**

`I/O = (src_num + 1) * vector_size / cost`

| Src Num  | Vector size | AVX512 I/O (MB/S) |  AVX2 I/O (MB/S) |SSE2 I/O (MB/S) |Generic I/O (MB/S) |
|-------|-------------|-------------|---------------|---------------|---------------|
|5|4KB|         |         |        |       |
|5|1MB|         |    	      |          |        |
|5|8MB|          |          |        |        |
|10|4KB|         |         |          |        |
|10|1MB|        |        |        |       |
|10|8MB|         |           |          |       |
