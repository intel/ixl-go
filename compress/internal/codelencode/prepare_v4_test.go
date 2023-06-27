// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64.v4

package codelencode

import (
	"bytes"
	"testing"

	"github.com/intel/ixl-go/compress/internal/codelencode/avx2"
	"github.com/intel/ixl-go/compress/internal/codelencode/avx512"
	"github.com/intel/ixl-go/compress/internal/codelencode/sse2"
)

func TestPrepareCodeLenCodeConsistentAVX512(t *testing.T) {
	for i := 0; i < 100; i++ {
		h := randHistogram()
		data := make([]uint8, 288+32)
		l, d := prepare(&h, data)

		data2 := make([]uint8, 288+32)
		l2, d2 := avx512.PrepareForCodeLenCode(&h, data2)
		if !bytes.Equal(data2, data) || l2 != l || d2 != d {
			t.Fatal("avx512 function is inconsistent with go function")
		}
	}
}

func BenchmarkPrepareCodeLenCode(b *testing.B) {
	h := randHistogram()
	data := make([]uint8, 288+32)
	b.Run("sse2_asm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sse2.PrepareForCodeLenCode(&h, data)
		}
	})
	b.Run("avx2_asm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			avx2.PrepareForCodeLenCode(&h, data)
		}
	})
	b.Run("avx512_asm", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			avx512.PrepareForCodeLenCode(&h, data)
		}
	})
	b.Run("go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			prepare(&h, data)
		}
	})
}
