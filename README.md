# xor

XOR code engine in pure Go

more than 200GB/S per core

## Introduction:

1. Use SIMD (SSE2 or AVX2 or AVX512) for speeding up
2. Go1.11

## Installation
To get the package use the standard:
```bash
go get github.com/templexxx/xorsimd
```

## Documentation

See the associated [GoDoc](http://godoc.org/github.com/templexxx/xorsimd)


## Performance

Performance depends mainly on:

1. SIMD extension
2. non-temporal size
3. hardware ( CPU RAM etc)

Example of performance on my AWS c5d.large (enable AVX512)
```
speed = ( shards * size ) / cost
```
| data_shards    | shard_size |speed (MB/S) |
|----------------|------------|-------------|
|5               |    4KB     |207751.51    |
|5               |    64KB    |98743.27     |

