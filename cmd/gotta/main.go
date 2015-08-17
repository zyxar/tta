package main

import (
	"../../../tta"
	"flag"
	"fmt"
	"os"
	"time"
)

var help, decode, encode bool
var passwd string

func init() {
	flag.BoolVar(&encode, "encode", false, "encode file")
	flag.BoolVar(&decode, "decode", false, "decode file")
	flag.BoolVar(&help, "help", false, "print this help")
	flag.StringVar(&passwd, "passwd", "", "specify password")
}

func main() {
	flag.Parse()
	if help || flag.NArg() < 2 || (!decode && !encode) {
		flag.Usage()
		return
	}
	infile, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	defer infile.Close()
	outfile, err := os.Create(flag.Arg(1))
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	callback := func(rate, fnum, frames uint32) {
		pcnt := uint32(float32(fnum) * 100. / float32(frames))
		if (pcnt % 10) == 0 {
			fmt.Printf("\rProgress: %02d%%", pcnt)
		}
	}
	if decode {
		println("Decoding:", flag.Arg(0), "to", flag.Arg(1))
		beginTime := time.Now()
		if err = tta.Decompress(infile, outfile, passwd, callback); err != nil {
			panic(err)
		}
		fmt.Printf("\rTime: %.3f sec.\n", float64(time.Now().UnixNano()-beginTime.UnixNano())/1000000000)
		return
	}
	if encode {
		println("Encoding:", flag.Arg(0), "to", flag.Arg(1))
		beginTime := time.Now()
		if err = tta.Compress(infile, outfile, passwd, callback); err != nil {
			panic(err)
		}
		fmt.Printf("\rTime: %.3f sec.\n", float64(time.Now().UnixNano()-beginTime.UnixNano())/1000000000)
	}
}
