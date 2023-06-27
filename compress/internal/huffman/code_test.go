// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package huffman

import (
	"reflect"
	"testing"
)

func TestCode(t *testing.T) {
	for _, x := range codeTests {
		rcodes := make([]uint16, len(x.lengths))
		GenerateCode(make([]int32, (15+1)*2), 15, x.lengths, rcodes)
		if !reflect.DeepEqual(rcodes, x.rcodes) {
			t.Fatal(rcodes, x.rcodes)
		}
	}
}

func TestCodeForIAA(t *testing.T) {
	for _, x := range codeTests {
		codes := make([]int32, len(x.lengths))
		copy(codes, x.lengths)
		GenerateCodeForIAA(codes, make([]int32, (15+1)*2))
		for i, c := range codes {
			if c>>15 != x.lengths[i] {
				t.Fatalf("%d %d", c>>15, x.lengths[i])
			}
			c = (c << 17) >> 17
			if uint16(c) != x.codes[i] {
				t.Fatalf("%b %b", uint16(c), x.codes[i])
			}
		}
	}
}

var codeTests = []struct {
	lengths []int32
	codes   []uint16
	rcodes  []uint16
}{
	{
		[]int32{2, 1, 3, 3},
		[]uint16{0b10, 0b0, 0b110, 0b111},
		[]uint16{0b01, 0b0, 0b011, 0b111},
	},
	{
		[]int32{3, 3, 3, 3, 3, 2, 4, 4},
		[]uint16{0b010, 0b011, 0b100, 0b101, 0b110, 0b00, 0b1110, 0b1111},
		[]uint16{0b010, 0b110, 0b001, 0b101, 0b011, 0b00, 0b0111, 0b1111},
	},
}
