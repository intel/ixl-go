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
	// gzip example:
	file, err := os.Open(*f)
	if err != nil {
		log.Fatalln("open file failed:", err)
	}
	output, err := os.Create(*f + ".gz")
	if err != nil {
		log.Fatalln("create file failed:", err)
	}
	w := compress.NewGzip(output)
	w.ReadFrom(file)
	w.Close()
	file.Close()
	output.Close()

	// deflate example:
	file, err = os.Open(*f)
	if err != nil {
		log.Fatalln("open file failed:", err)
	}
	buf := bytes.NewBuffer(nil)
	d, err := compress.NewDeflate(buf)
	if err != nil {
		log.Fatalln("NewDeflate failed:", err)
	}
	d.ReadFrom(file)
	d.Close()

	// inflate example
	i, err := compress.NewInflate(buf)
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
