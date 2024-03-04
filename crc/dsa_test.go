// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package crc

import (
	"crypto/rand"
	_ "embed"
	"hash/crc32"
	"strconv"
	"testing"
)

//go:embed dsa_test.go
var testcode []byte

func TestDSACrc32C(t *testing.T) {
	hash, err := NewCRC32C()
	if err != nil {
		t.Skip(err)
	}
	hash.Write(testcode)

	hash2 := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	hash2.Write(testcode)
	if hash2.Sum32() != hash.Sum32() {
		t.Fatalf("expected crc32c sum32 result equals to hash/crc32 sum32 result,but got %d != %d", hash2.Sum32(), hash.Sum32())
	}
}

func BenchmarkCRC32C(b *testing.B) {
	dsaHasher, err := NewCRC32C()
	if err != nil {
		b.Skip(err)
	}
	hasher := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	for i := 16; i <= 1024; i *= 2 {
		data := make([]byte, i*1024)
		rand.Read(data)
		b.Run("std_"+strconv.Itoa(i)+"k", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hasher.Write(data)
				hasher.Sum32()
				hasher.Reset()
			}
		})
		b.Run("dsa_"+strconv.Itoa(i)+"k", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := dsaHasher.Write(data)
				if err != nil {
					b.Log(err)
				}
				dsaHasher.Sum32()
				dsaHasher.Reset()
			}
		})
	}
}
