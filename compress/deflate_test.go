// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/intel/ixl-go/internal/testutil"
)

func TestDeflate(t *testing.T) {
	if !Ready() {
		t.Skip("IAA devices not found")
	}
	w, _ := NewDeflate(io.Discard)
	buf := bytes.NewBuffer(nil)
	for i := 2; i <= 4096*1024; i = i * 2 {
		buf.Reset()
		text := []byte(testutil.RandomText(i))
		w.Reset(buf)
		_, err := w.ReadFrom(bytes.NewBuffer(text))
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		w.Close()
		r := flate.NewReader(buf)
		_, err = r.Read(make([]byte, i*1024))
		r.Close()
		if err != nil && err != io.EOF {
			t.Log(base64.StdEncoding.EncodeToString(buf.Bytes()))
			t.Log("error:", i, buf.Len())
			t.Fatal(err.Error(), reflect.TypeOf(err))
		}
	}
}

func FuzzDeflate(f *testing.F) {
	if !Ready() {
		f.Skip("IAA devices not found")
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		w, _ := NewDeflate(io.Discard)
		_, err := w.ReadFrom(bytes.NewBuffer(data))
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		w.Close()
	})
}

func BenchmarkDeflate4k(b *testing.B) {
	if !Ready() {
		b.Skip("IAA devices not found")
	}
	text := []byte(testutil.RandomText(4 * 1024))
	w, _ := NewDeflate(io.Discard)
	for j := 0; j < b.N; j++ {
		w.Reset(io.Discard)
		_, _ = w.ReadFrom(bytes.NewBuffer(text))
		w.Close()
	}
}

func BenchmarkPDeflate4k(b *testing.B) {
	if !Ready() {
		b.Skip("IAA devices not found")
	}
	text := []byte(testutil.RandomText(4 * 1024))
	b.RunParallel(func(p *testing.PB) {
		w, _ := NewDeflate(io.Discard)
		for p.Next() {
			w.Reset(io.Discard)
			_, _ = w.ReadFrom(bytes.NewBuffer(text))
			w.Close()
		}
	})
}

func BenchmarkReusedRandomTextCompress(b *testing.B) {
	if !Ready() {
		b.Skip("IAA devices not found")
	}
	for i := 4; i <= 4096; i = i * 2 {
		text := []byte(testutil.RandomText(i * 1024))
		fw, _ := flate.NewWriter(io.Discard, flate.DefaultCompression)
		b.Run(fmt.Sprintf("standard deflate[%dk]", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				fw.Reset(io.Discard)
				_, _ = fw.Write(text)
				fw.Close()
			}
		})
		w, _ := NewDeflateWriter(io.Discard)
		b.Run(fmt.Sprintf("IAA deflate[%dk]", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				w.Reset(io.Discard)
				_, _ = w.Write(text)
				w.Close()
			}
		})
	}
}

func BenchmarkRandomTextCompress(b *testing.B) {
	if !Ready() {
		b.Skip("IAA devices not found")
	}
	for i := 4; i <= 4096; i = i * 2 {
		text := []byte(testutil.RandomText(i * 1024))
		b.Run(fmt.Sprintf("standard deflate[%dk]", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				fw, _ := flate.NewWriter(io.Discard, flate.DefaultCompression)
				_, _ = fw.Write(text)
				fw.Close()
			}
		})
		b.Run(fmt.Sprintf("IAA deflate[%dk]", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				w, _ := NewDeflateWriter(io.Discard)
				_, _ = w.Write(text)
				w.Close()
			}
		})
	}
}
