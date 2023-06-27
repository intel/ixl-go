// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package crc

import (
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"testing"

	"github.com/intel/ixl-go/internal/testutil"
)

func TestCRC64(t *testing.T) {
	if !Ready() {
		t.Skip()
	}
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i)
	}
	cases := []struct {
		name     string
		data     []byte
		poly     uint64
		expected uint64
	}{
		{
			name: "ISO",
			data: data,
			poly: crc64.ISO,
		},

		{
			name: "ECMA",
			data: data,
			poly: crc64.ECMA,
		},
	}
	calc, _ := NewCalculator()
	for _, tc := range cases {
		crc, err := calc.CheckSum64(tc.data, tc.poly)
		if err != nil {
			t.Fatalf("test case %s Failed: %s", tc.name, err.Error())
		}
		expected := crc64.Checksum(tc.data, crc64.MakeTable(tc.poly))
		if crc != expected {
			t.Fatalf("test case %s Failed: [expected %x] [actual %x]", tc.name, expected, crc)
		}
	}
}

func TestCRC32(t *testing.T) {
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i)
	}
	cases := []struct {
		name     string
		data     []byte
		poly     uint32
		expected uint32
	}{
		{
			name: "IEEE",
			data: data,
			poly: crc32.IEEE,
		},

		{
			name: "Koopman",
			data: data,
			poly: crc32.Koopman,
		},

		{
			name: "Castagnoli",
			data: data,
			poly: crc32.Castagnoli,
		},
	}
	calc, _ := NewCalculator()
	for _, tc := range cases {
		crc, err := calc.CheckSum32(tc.data, tc.poly)
		if err != nil {
			t.Fatalf("test case %s Failed: %s", tc.name, err.Error())
		}
		expected := crc32.Checksum(tc.data, crc32.MakeTable(tc.poly))
		if crc != expected {
			t.Fatalf("test case %s Failed: [expected %x] [actual %x]", tc.name, expected, crc)
		}
	}
}

func TestCRC16(t *testing.T) {
	data := make([]byte, 1000)
	for i := range data {
		data[i] = byte(i)
	}
	cases := []struct {
		name     string
		data     []byte
		poly     uint16
		expected uint16
	}{
		{
			name:     "T10DIF",
			data:     data,
			poly:     0x8BB7,
			expected: 0x8439,
		},

		{
			name:     "CCITT",
			data:     data,
			poly:     0x8408,
			expected: 0x245f,
		},
	}
	calc, _ := NewCalculator()
	for _, tc := range cases {
		crc, err := calc.CheckSum16(tc.data, tc.poly)
		if err != nil {
			t.Fatalf("test case %s Failed: %s", tc.name, err.Error())
		}
		if crc != tc.expected {
			t.Fatalf("test case %s Failed: [expected %x] [actual %x]", tc.name, tc.expected, crc)
		}
	}
}

func BenchmarkCRC64(b *testing.B) {
	for i := 4; i <= 4096; i = i * 2 {
		text := []byte(testutil.RandomText(i * 1024))
		table := crc64.MakeTable(crc64.ISO)
		b.Run(fmt.Sprintf("standard crc64 ISO[%dk]", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				crc64.Checksum(text, table)
			}
		})
		calc, _ := NewCalculator()
		b.Run(fmt.Sprintf("IAA crc64 ISO[%dk]", i), func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				_, _ = calc.CheckSum64(text, crc64.ISO)
			}
		})
	}
}

func FuzzCRC64(f *testing.F) {
	if !Ready() {
		f.Skip()
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		c, _ := NewCalculator()
		num, err := c.CheckSum64(data, crc64.ISO)
		if err != nil {
			t.Fatal(err)
		}
		if target := crc64.Checksum(data, crc64.MakeTable(crc64.ISO)); target != num {
			t.Errorf("expected %d got %d", target, num)
		}
	})
}

func FuzzCRC32(f *testing.F) {
	if !Ready() {
		f.Skip()
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		c, _ := NewCalculator()
		num, err := c.CheckSum32(data, IEEE)
		if err != nil {
			t.Fatal(err)
		}
		if target := crc32.Checksum(data, crc32.IEEETable); target != num {
			t.Errorf("expected %d got %d", target, num)
		}
	})
}
