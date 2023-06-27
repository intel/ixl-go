// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/intel/ixl-go/internal/testutil"
)

// TestGzip_ReadFrom compatibility testing
func TestGzip_ReadFrom(t *testing.T) {
	if !Ready() {
		t.Skip("IAA devices not found")
	}
	g := NewGzip(io.Discard)
	for i := 2; i < 1024*1024; i = i * 2 {
		size := i
		input := testutil.RandomText(size)
		output := bytes.NewBuffer(nil)
		g.Reset(output)
		_, err := g.ReadFrom(bytes.NewBuffer([]byte(input)))
		g.Close()
		if err != io.EOF && err != nil {
			t.Fatal("error happened while gzip:", err)
		}

		sg, _ := gzip.NewReader(output)
		soutput := bytes.NewBuffer(nil)
		_, err = io.Copy(soutput, sg)
		if err != io.EOF && err != nil {
			t.Fatal("error happened while gunzip", size, err)
		}
		if soutput.String() != input {
			t.Fatal("decompressed data is not consistent with input")
		}
	}
}

// TestGzip_ReadFromWithName compatibility testing
func TestGzip_ReadFromWithName(t *testing.T) {
	if !Ready() {
		t.Skip("IAA devices not found")
	}
	g := NewGzip(io.Discard)
	for i := 2; i < 1024*1024; i = i * 2 {
		size := i
		input := testutil.RandomText(size)
		output := bytes.NewBuffer(nil)
		g.Reset(output)
		g.Name = "hallo.txt"
		_, err := g.ReadFrom(bytes.NewBuffer([]byte(input)))
		g.Close()
		if err != io.EOF && err != nil {
			t.Fatal("error happened while gzip:", err)
		}

		sg, _ := gzip.NewReader(output)
		soutput := bytes.NewBuffer(nil)
		_, err = io.Copy(soutput, sg)
		if err != io.EOF && err != nil {
			t.Fatal("error happened while gunzip", size, err)
		}
		if soutput.String() != input {
			t.Fatal("decompressed data is not consistent with input")
		}
	}
}

// TestGzip_WriteBlock compatibility testing
func TestGzip_WriteBlock(t *testing.T) {
	if !Ready() {
		t.Skip("IAA devices not found")
	}
	g := NewGzip(io.Discard)
	w := NewWriter(g)

	for i := 2; i < 1024*1024; i = i * 2 {
		size := i
		input := testutil.RandomText(size)
		output := bytes.NewBuffer(nil)
		w.Reset(output)
		_, err := w.Write([]byte(input))
		w.Close()
		if err != io.EOF && err != nil {
			t.Fatal("error happened while gzip:", err)
		}

		sg, _ := gzip.NewReader(output)
		soutput := bytes.NewBuffer(nil)
		_, err = io.Copy(soutput, sg)
		if err != io.EOF && err != nil {
			t.Fatal("error happened while gunzip", size, err)
		}
		if soutput.String() != input {
			t.Fatal("decompressed data is not consistent with input")
		}
	}
}
