// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package datamove

import (
	"bytes"
	"crypto/rand"
	"sync"
	"testing"
)

var longRand = func() []byte {
	data := make([]byte, 1024*1024)
	_, _ = rand.Read(data)
	return data
}()

func TestCopyTransferSize(t *testing.T) {
	if !Ready() {
		t.Skip()
	}
	size := globalCtx.MaxTransferSize()
	globalCtx.SetMaxTransferSize(1024)
	defer globalCtx.SetMaxTransferSize(size)
	dest := make([]byte, 1024*1024)
	src := longRand
	if !Copy(dest, src) {
		t.Fatal("Copy failed")
	}
	if !bytes.Equal(dest, src) {
		t.Fatal("Copy not works")
	}
}

func TestContext_Copy(t *testing.T) {
	if !Ready() {
		t.Skip()
	}
	type args struct {
		dest []byte
		src  []byte
	}
	ctx := NewContext()
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty copy",
			args: args{
				dest: make([]byte, 0),
				src:  make([]byte, 0),
			},
			want: true,
		},
		{
			name: "normal copy",
			args: args{
				dest: []byte{1, 2, 3, 4},
				src:  []byte{2, 3, 4, 5, 6},
			},
			want: true,
		},
		{
			name: "1MB copy",
			args: args{
				dest: make([]byte, 1024*1024),
				src:  longRand,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := append([]byte{}, tt.args.dest...)
			if got := ctx.Copy(output, tt.args.src); got != tt.want {
				t.Errorf("Context.Copy() = %v, want %v", got, tt.want)
			}
			size := len(tt.args.src)
			if len(tt.args.dest) < size {
				size = len(tt.args.dest)
			}
			if !bytes.Equal(output[:size], tt.args.src[:size]) {
				t.Errorf("Context.Copy() doesn't work")
			}
		})
	}

	for _, tt := range tests {
		t.Run("COPY_"+tt.name, func(t *testing.T) {
			output := append([]byte{}, tt.args.dest...)
			if got := Copy(output, tt.args.src); got != tt.want {
				t.Errorf("Context.Copy() = %v, want %v", got, tt.want)
			}
			size := len(tt.args.src)
			if len(tt.args.dest) < size {
				size = len(tt.args.dest)
			}
			if !bytes.Equal(output[:size], tt.args.src[:size]) {
				t.Errorf("Context.Copy() doesn't work")
			}
		})
	}
}

func TestCopy(t *testing.T) {
	if !Ready() {
		t.Skip()
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 1024; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			input := make([]byte, 1024*1024)
			_, _ = rand.Read(input)
			output := make([]byte, 1024*1024)
			Copy(output, input)
		}()
	}
	wg.Wait()
}

func BenchmarkCopy(b *testing.B) {
	if !Ready() {
		b.Skip()
	}
	b.Run("DSA", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			ctx := NewContext()
			input := make([]byte, 1024*1024)
			_, _ = rand.Read(input)
			output := make([]byte, 1024*1024)
			for p.Next() {
				ctx.Copy(output, input)
			}
		})
	})

	b.Run("CPU", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			input := make([]byte, 1024*1024)
			_, _ = rand.Read(input)
			output := make([]byte, 1024*1024)
			for p.Next() {
				copy(output, input)
			}
		})
	})
}

func FuzzCopy(f *testing.F) {
	if !Ready() {
		f.Skip()
	}
	ctx := NewContext()

	f.Fuzz(func(t *testing.T, data []byte) {
		output := make([]byte, len(data))
		err := ctx.CopyCheckError(output, data)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(output, data) {
			t.Fatal()
		}
	})
}
