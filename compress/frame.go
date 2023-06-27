// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"errors"
	"unsafe"
)

type headerFrame struct {
	err     error
	acc     int
	buffer  *[256]byte
	written int

	bits   uint64 // bits buffer for improve writeBits performance
	bitNum uint64 // bit number in the bits
}

func (b *headerFrame) start(frame *[256]byte) {
	b.buffer = frame
}

var errOutOfHeaderSize = errors.New("out of max header size")

// WriteBit write the data with bit number `n`.
// Notice: n should be less than 16.
func (b *headerFrame) WriteBit(data uint16, n int8) {
	if b.err != nil {
		return
	}
	b.acc += int(n)
	if b.acc > 2048 {
		// this will not happen
		b.err = errOutOfHeaderSize
		return
	}
	b.bits |= uint64(data) << b.bitNum
	b.bitNum += uint64(n)

	if b.bitNum > 48 {
		b.buffer[b.written] = byte(b.bits)
		b.buffer[b.written+1] = byte(b.bits >> 8)
		b.buffer[b.written+2] = byte(b.bits >> 16)
		b.buffer[b.written+3] = byte(b.bits >> 24)
		b.buffer[b.written+4] = byte(b.bits >> 32)
		b.buffer[b.written+5] = byte(b.bits >> 40)
		b.written += 6

		b.bitNum -= 48
		b.bits = b.bits >> 48
	}
}

func (b *headerFrame) flush() (acc int) {
	available := b.bitNum / 8
	hasBit := b.bitNum%8 != 0
	if hasBit {
		available++
	}

	target := b.buffer[b.written : b.written+int(available)]
	source := *(*[]byte)(unsafe.Pointer(&sliceHeader{
		ptr: &b.bits,
		len: int(available),
		cap: int(available),
	}))

	copy(target, source)
	acc = b.acc

	b.acc = 0
	b.written = 0
	b.bits = 0
	b.bitNum = 0
	b.buffer = nil

	return acc
}

type sliceHeader struct {
	ptr *uint64
	len int
	cap int
}
