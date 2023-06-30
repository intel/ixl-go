// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"bytes"
	"compress/flate"
	"crypto/rand"
	"encoding/base64"
	"io"
	"reflect"
	"strconv"
	"testing"

	"github.com/intel/ixl-go/internal/testutil"
)

func TestInflate(t *testing.T) {
	if !Ready() {
		t.Skip("IAA devices not found")
	}
	w, err := NewDeflate(io.Discard)
	if err != nil {
		t.Skip(err)
	}
	buf := bytes.NewBuffer(nil)
	r, err := NewInflate(nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 2; i <= 4096*1024; i = i * 2 {
		buf.Reset()
		text := []byte(testutil.RandomText(i))
		w.Reset(buf)
		_, err := w.ReadFrom(bytes.NewBuffer(text))
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		w.Close()
		block := make([]byte, i)
		r.Reset(buf)
		for {
			_, err = r.Read(block)
			if err != nil {
				break
			}
		}
		if err != nil && err != io.EOF {
			t.Log(base64.StdEncoding.EncodeToString(buf.Bytes()))
			t.Log("error:", i, buf.Len())
			t.Fatal(err.Error(), reflect.TypeOf(err))
		}
	}
}

func TestInflateReadZeroLengthData(t *testing.T) {
	if !Ready() {
		t.Skip("IAA devices not found")
	}
	compressed := bytes.NewBuffer(nil)
	w, err := NewDeflate(compressed)
	if err != nil {
		t.Skip(err)
	}
	_, err = w.ReadFrom(io.LimitReader(rand.Reader, 1024))
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	w.Close()

	r, err := NewInflate(compressed)
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 1024)
	_, err = r.Read(data[:0])
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Read(data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Read(nil)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
}

func FuzzInflate(f *testing.F) {
	if !Ready() {
		f.Skip("IAA devices not found")
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		buf := bytes.NewBuffer(nil)
		w, _ := NewDeflate(buf)
		_, err := w.ReadFrom(bytes.NewBuffer(data))
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		w.Close()
		rr := flate.NewReader(bytes.NewBuffer(buf.Bytes()))
		out, err := io.ReadAll(rr)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if !bytes.Equal(out, data) {
			t.Fatal("expected equals")
		}

		i, _ := NewInflate(bytes.NewBuffer(buf.Bytes()))
		out, err = io.ReadAll(i)
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if !bytes.Equal(out, data) {
			t.Fatal("expected equals")
		}
	})
}

func TestInflate_Binary(t *testing.T) {
	if !Ready() {
		t.Skip("IAA devices not found")
	}
	w, _ := NewDeflate(io.Discard)

	buf := bytes.NewBuffer(nil)
	for i := 2; i <= 4096*1024; i = i * 2 {
		buf.Reset()
		data := make([]byte, i)
		_, _ = rand.Read(data)
		w.Reset(buf)
		_, err := w.ReadFrom(bytes.NewBuffer(data))
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		w.Close()
		r, _ := NewInflate(buf)
		block := make([]byte, i)
		for {
			_, err = r.Read(block)
			if err != nil {
				break
			}
		}
		if err != nil && err != io.EOF {
			t.Log(base64.StdEncoding.EncodeToString(buf.Bytes()))
			t.Log("error:", i, buf.Len())
			t.Fatal(err.Error(), reflect.TypeOf(err))
		}
	}
}

func BenchmarkInflate(b *testing.B) {
	if !Ready() {
		b.Skip("IAA devices not found")
	}
	for i := 1; i < 1024; i = i * 2 {
		w, _ := NewDeflate(io.Discard)
		buf := bytes.NewBuffer(nil)
		buf.Reset()
		text := []byte(testutil.RandomText(i * 1024))
		w.Reset(buf)
		_, err := w.ReadFrom(bytes.NewBuffer(text))
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
		w.Close()
		data := buf.Bytes()
		b.Run("standard_inflate_"+strconv.Itoa(i)+"KB", func(b *testing.B) {
			r := flate.NewReader(nil)
			temp := make([]byte, len(data))
			for i := 0; i < b.N; i++ {
				copy(temp, data)
				_ = r.(flate.Resetter).Reset(bytes.NewReader(temp), nil)
				_, _ = io.Copy(io.Discard, r)
			}
		})

		b.Run("IAA_inflate_"+strconv.Itoa(i)+"KB", func(b *testing.B) {
			r, _ := NewInflate(buf)
			temp := make([]byte, len(data))
			for i := 0; i < b.N; i++ {
				copy(temp, data)
				r.Reset(bytes.NewReader(temp))
				_, _ = io.Copy(io.Discard, r)
			}
		})
	}
}
