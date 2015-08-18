# tta
[TTA Lossless Audio Codec](http://en.true-audio.com/TTA_Lossless_Audio_Codec_-_Realtime_Audio_Compressor) Encoder/Decoder for #golang

## [API doc](https://godoc.org/github.com/zyxar/tta)

## *gotta* console tool

- install: `go get github.com/zyxar/tta/cmd/gotta`
- usage:

  ```
    -decode=false: decode file
    -encode=false: encode file
    -help=false: print this help
    -passwd="": specify password
  ```

## TODOs

- [ ] general optimization
- [ ] SSE4 acceleration
