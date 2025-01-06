// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"

	"github.com/intel/ixl-go/compress"
)

var f = flag.String("f", "", "file to compress")

func main() {
	flag.Parse()
	if !compress.Ready() {
		log.Fatalln("IAA workqueue not found")
		return
	}
	if *f == "" {
		log.Fatalln("must give a file to compress")
	}

	gzipExample(*f)
	deflateExample(*f)
}

// gzip example:
func gzipExample(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("open file failed:", err)
	}
	defer file.Close()

	output, err := os.Create(filename + ".gz")
	if err != nil {
		log.Fatalln("create file failed:", err)
	}
	defer output.Close()

	w := compress.NewGzip(output)
	w.ReadFrom(file)
	w.Close()
}

// deflate example:
func deflateExample(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("open file failed:", err)
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	d, err := compress.NewDeflate(buf)
	if err != nil {
		log.Fatalln("NewDeflate failed:", err)
	}
	d.ReadFrom(file)
	d.Close()

	inflateExample(buf)
}

// inflate example:
func inflateExample(r io.Reader) {
	i, err := compress.NewInflate(r)
	if err != nil {
		log.Fatalln("NewInflate failed:", err)
	}
	block := make([]byte, 1024)
	for {
		_, err := i.Read(block)
		if err != nil {
			if err != io.EOF {
				log.Fatalln("error:", err)
			}
			break
		}
	}
}
