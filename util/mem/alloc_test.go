// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package mem

import (
	"testing"
	"unsafe"
)

func TestAlloc64ByteAligned(t *testing.T) {
	for i := 1; i <= 4096; i++ {
		data := Alloc64ByteAligned(uintptr(i))
		if uintptr(unsafe.Pointer(&data[0]))%64 != 0 {
			t.Fatal("the data addr must aligned to 64 byte")
		}
	}
}

func TestAlloc(t *testing.T) {
	type Any struct {
		data [10]byte
	}
	for i := 0; i < 1000000; i++ {
		a := Alloc32Align[Any]()
		if uintptr(unsafe.Pointer(a))%32 != 0 {
			t.Fatal("the data addr must aligned to 32 byte")
		}
		b := Alloc64Align[Any]()
		if uintptr(unsafe.Pointer(b))%64 != 0 {
			t.Fatal("the data addr must aligned to 64 byte")
		}
	}
}
