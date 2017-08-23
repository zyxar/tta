package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	"time"

	"github.com/zyxar/tta"
)

var (
	help, decode, encode bool
	passwd               string
	cpuprofile           string
)

func init() {
	flag.BoolVar(&encode, "encode", false, "encode file")
	flag.BoolVar(&decode, "decode", false, "decode file")
	flag.BoolVar(&help, "help", false, "print this help")
	flag.StringVar(&passwd, "passwd", "", "specify password (optional)")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
}

func main() {
	fmt.Fprintf(os.Stderr, "\r\nTTA1 lossless audio encoder/decoder, version %s\n\n", tta.Version)
	flag.Parse()
	if help || flag.NArg() < 1 || (!decode && !encode) {
		fmt.Fprintf(os.Stderr, "\rUsage of gotta: [encode|decode] [passwd PASSWORD] INPUT_FILE OUTPUT_FILE\n\n")
		flag.PrintDefaults()
		return
	}
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			panic(err.Error())
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	infile := flag.Arg(0)
	outfile := flag.Arg(1)
	input, err := os.Open(infile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer input.Close()
	if outfile == "" {
		outfile = path.Base(infile)
		outfile = outfile[:len(outfile)-len(path.Ext(outfile))]
		ext := path.Ext(outfile)
		if decode {
			if ext != ".wav" {
				outfile += ".wav"
			}
		} else {
			if ext != ".tta" {
				outfile += ".tta"
			}
		}
	}
	if _, err = os.Stat(outfile); err == nil || !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, outfile, "exists")
		return
	}
	output, err := os.Create(outfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer output.Close()
	callback := func(rate, fnum, frames uint32) {
		pcnt := uint32(float32(fnum) * 100. / float32(frames))
		if (pcnt % 10) == 0 {
			fmt.Fprintf(os.Stderr, "\rProgress: %02d%%", pcnt)
		}
	}
	if decode {
		fmt.Fprintf(os.Stderr, "Decoding: \"%v\" to \"%v\"\n", infile, outfile)
		beginTime := time.Now()
		if err = tta.Decompress(input, output, passwd, callback); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Printf("\rTime: %.3f sec.\n", float64(time.Now().UnixNano()-beginTime.UnixNano())/1000000000)
		return
	}
	if encode {
		fmt.Fprintf(os.Stderr, "Encoding: \"%v\" to \"%v\"\n", infile, outfile)
		beginTime := time.Now()
		if err = tta.Compress(input, output, passwd, callback); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Printf("\rTime: %.3f sec.\n", float64(time.Now().UnixNano()-beginTime.UnixNano())/1000000000)
	}
}
