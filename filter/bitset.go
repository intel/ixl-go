// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"bytes"
	"fmt"
	"math/bits"
	"unsafe"
)

// DataUnit represents a type that can be used in filter operations
type DataUnit interface {
	uint8 | uint16 | uint32
}

// Range represents a range of values
type Range[R DataUnit] struct {
	Min, Max R
}

// BitSet is a fixed-size collection of bits that can be manipulated individually.
// It is a data structure that is used to represent a set of elements,
// where each element is represented by a single bit.
type BitSet []byte

// OnesCount returns the number of one bits ("population count") in b.
func (b BitSet) OnesCount() int {
	idx := 0
	num := len(b) / 8
	c := 0
	if num > 0 {
		nums := unsafe.Slice((*uint64)(unsafe.Pointer(&b[idx])), num)
		for _, v := range nums {
			c += bits.OnesCount64(v)
		}
		idx = num * 8
	}
	for _, v := range b[idx:] {
		c += bits.OnesCount8(v)
	}
	return c
}

// Size return the min bitset size
func (b BitSet) Size() int {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] == 0 {
			continue
		}
		prev := (i) * 8
		last := (8 - bits.LeadingZeros8(b[i]))
		return prev + last
	}
	return 0
}

// String returns a string representation of the bit set
func (b BitSet) String() string {
	buf := bytes.NewBuffer(nil)
	for _, v := range b {
		buf.WriteString(fmt.Sprintf("%08b", v))
		buf.WriteByte('\n')
	}
	return buf.String()
}
