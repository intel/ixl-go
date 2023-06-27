// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package mem declares functions for allocating memory that is aligned to specific byte boundaries.
package mem

import "unsafe"

// Alloc64ByteAligned returns a byte slice of the specified size that is aligned to 64 bytes.
// The returned slice may have additional unused capacity to ensure alignment.
//
//go:noinline
func Alloc64ByteAligned(size uintptr) []byte {
	data := make([]byte, size+64)
	addr := (uintptr)(unsafe.Pointer(&data[0]))
	m := addr % 64
	if m == 0 {
		return data[:size]
	}
	start := 64 - m
	return data[start:int(start+size)]
}

// Alloc64Align returns a pointer to a value of type T that is aligned to 64 bytes.
// The returned pointer may point to additional unused memory to ensure alignment.
//
//go:noinline
func Alloc64Align[T any]() *T {
	var t T
	size := unsafe.Sizeof(t)
	data := make([]byte, size+64)
	addr := (uintptr)(unsafe.Pointer(&data[0]))
	m := addr % 64
	if m == 0 {
		data = data[:size]
	} else {
		start := 64 - m
		data = data[start:int(start+size)]
	}
	return (*T)(unsafe.Pointer(&data[0]))
}

// Alloc32Align returns a pointer to a value of type T that is aligned to 32 bytes.
// The returned pointer may point to additional unused memory to ensure alignment.
//
//go:noinline
func Alloc32Align[T any]() *T {
	var t T
	size := unsafe.Sizeof(t)
	data := make([]byte, size+32)
	addr := (uintptr)(unsafe.Pointer(&data[0]))
	m := addr % 32
	if m == 0 {
		data = data[:size]
	} else {
		start := 32 - m
		data = data[start:int(start+size)]
	}
	return (*T)(unsafe.Pointer(&data[0]))
}
