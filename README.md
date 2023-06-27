# Intel Accelerators User Library for Go
--- 
[![Coverage Status](./coverage.svg)]()

The project is aimed to develop Go library of IAA/DSA to empower cloud native industry. IAA/DSA is introduced by 4th Gen Intel(R) Xeon(R) Scalable Processor which was launched on Jan'23.


## Supported Hardware Accelerator Features

- Compression/Decompression
  - Deflate
  - Gzip
- CRC calculation
- Data Filter (Bitpack / RLE format / Int Array): 
  - Expand
  - Select
  - Scan
  - Extract
- Data Move

## Advantages

The library is designed to be a lightweight pure Go language implementation with zero dependencies.


- Lower barriers to use for users, no need to install dependency libraries in advance.
- No CGO performance overhead.
- Easy to learn and start.
# How to use

## Installation

Use the following commands to add ixl-go as a dependency to your existing Go project.

```bash
go get github.com/intel/ixl-go
```
## Quick Start

To use the accelerated functions you need IAA or DSA devices on your machine.

You can use `Ready` function to check if the devices are OK for ixl-go. 

Notice: 

> we don't support non-svm environment for now. 
> 
> We only support workqueues which enabled **block_on_fault** and **SVM**.

If you don't know about how to config IAA/DSA workqueues, you can follow [this simple recipe](./enable-iaa.md).

### Example1: Compress/Decompress

```go
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

```



### Example2: CRC Calculation

```go
c, err := crc.NewCalculator()
if err != nil{
	return err
}
result, err := c.CheckSum64(data, crc64.ISO)
```

### Example3: Data Copy
```go
datamove.Copy(output, input)

```

# Documentation

| module   | doc                           | accelerator | 
| -------- | ----------------------------- | ----------- |
| compress | [compress](./compress/doc.md) | IAA         |
| crc      | [crc](./crc/doc.md)           | IAA         |
| filter   | [filter](./filter/doc.md)     | IAA         |
| datamove | [datamove](./datamove/doc.md) | DSA         |
| util/mem | [util/mem](./util/mem/doc.md) | _           |
