// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"testing"
)

func TestBitSet_OnesCount(t *testing.T) {
	tests := []struct {
		name string
		b    BitSet
		want int
	}{
		{
			name: "one_byte",
			b:    BitSet{0b001},
			want: 1,
		},
		{
			name: "two_byte",
			b:    BitSet{0b001, 0b111},
			want: 4,
		},
		{
			name: "three_byte",
			b:    BitSet{0b010, 0b00, 0b01},
			want: 2,
		},
		{
			name: "18_byte",
			b: BitSet{
				0b010, 0b00, 0b01, 0b010, 0b00,
				0b01, 0b010, 0b00, 0b01, 0b010,
				0b00, 0b01, 0b010, 0b00, 0b01,
				0b010, 0b00, 0b01,
			},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.OnesCount(); got != tt.want {
				t.Errorf("BitSet.OnesCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitSet_Size(t *testing.T) {
	tests := []struct {
		name string
		b    BitSet
		want int
	}{
		{
			name: "one byte",
			b:    BitSet{0b010},
			want: 2,
		},
		{
			name: "zero bytes",
			b:    BitSet{0b010, 0, 0},
			want: 2,
		},
		{
			name: "more than one byte",
			b:    BitSet{0b110, 0, 0b010, 0},
			want: 18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Size(); got != tt.want {
				t.Errorf("BitSet.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
