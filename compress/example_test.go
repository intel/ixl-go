// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"log"

	"github.com/intel/ixl-go/compress"
)

func Example() {
	input := make([]byte, 1024)
	_, _ = rand.Read(input)
	var w io.WriteCloser
	d, err := compress.NewDeflate(bytes.NewBuffer(nil))
	if err != nil {
		// IAA devices not found
		log.Fatalln(err)
	}
	w = compress.NewWriter(d)
	_, _ = w.Write(input)
	_ = w.Close()
}
