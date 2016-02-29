# tta
[![Go Report Card](https://goreportcard.com/badge/github.com/zyxar/tta)](https://goreportcard.com/report/github.com/zyxar/tta)
[![GoDoc](https://godoc.org/github.com/zyxar/tta?status.svg)](https://godoc.org/github.com/zyxar/tta)
[![Build Status](https://travis-ci.org/zyxar/tta.svg?branch=master)](https://travis-ci.org/zyxar/tta)

[TTA Lossless Audio Codec](http://en.true-audio.com/TTA_Lossless_Audio_Codec_-_Realtime_Audio_Compressor) Encoder/Decoder for #golang

## `gotta` console tool

- install: `go get github.com/zyxar/tta/cmd/gotta`
- usage:

  ```
    -decode=false: decode file
    -encode=false: encode file
    -help=false: print this help
    -passwd="": specify password
  ```

## Comparison

![](https://github.com/zyxar/tta/blob/master/data/tta_comp.svg)

## TODOs

- [ ] general optimization
- [ ] SSE4 acceleration
