// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64.v3

package codelencode

import (
	"bytes"
	"testing"

	"github.com/intel/ixl-go/compress/internal/codelencode/avx2"
)

func TestPrepareCodeLenCodeConsistentAVX2(t *testing.T) {
	for i := 0; i < 100; i++ {
		h := randHistogram()
		data := make([]uint8, 288+32)
		l, d := prepare(&h, data)

		data2 := make([]uint8, 288+32)
		l2, d2 := avx2.PrepareForCodeLenCode(&h, data2)
		if !bytes.Equal(data2, data) || l2 != l || d2 != d {
			t.Fatal("avx2 function is inconsistent with go function")
		}
	}
}
