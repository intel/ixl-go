// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package iaa

// FilterFlags represents a collection of filter operation flags.
type FilterFlags uint32

const (
	// FilterFlagSource1BigEndian indicates whether source 1 is big-endian.
	FilterFlagSource1BigEndian FilterFlags = 1
	// FilterFlagSource1ParquetRLE indicates whether source 1 is Parquet RLE encoded.
	FilterFlagSource1ParquetRLE FilterFlags = 1 << 1
	// FilterFlagSource2BigEndian indicates whether source 2 is big-endian.
	FilterFlagSource2BigEndian FilterFlags = 1 << 12
	// FilterFlagOutputWithByte indicates output with byte.
	FilterFlagOutputWithByte FilterFlags = 1 << 13
	// FilterFlagOutputWithWord indicates output with word.
	FilterFlagOutputWithWord FilterFlags = 2 << 13
	// FilterFlagOutputWithDword indicates output with dword.
	FilterFlagOutputWithDword FilterFlags = 3 << 13
	// FilterFlagOutputBigEndian indicates output is big-endian.
	FilterFlagOutputBigEndian FilterFlags = 1 << 15
	// FilterFlagInvertOutput indicates whether to invert output.
	FilterFlagInvertOutput FilterFlags = 1 << 16
)

// SetSource1Width sets the width of source 1 in bits.
// If the width is greater than 32, it cannot be handled.
func (f *FilterFlags) SetSource1Width(width uint8) {
	if width > 32 {
		// todo: cannot handle
		return
	}
	mask := ^(uint32(31) << 2)
	*f &= FilterFlags(mask)
	*f |= FilterFlags(uint32(width-1) << 2)
}

// SetSource2Width sets the width of source 2 in bits.
// If the width is greater than 32, it cannot be handled.
func (f *FilterFlags) SetSource2Width(width uint8) {
	if width > 32 {
		// todo: cannot handle
		return
	}

	mask := ^(uint32(31) << 7)
	*f &= FilterFlags(mask)
	*f |= FilterFlags(uint32(width-1) << 7)
}

// SetDropLowBits sets the number of low bits to be dropped from the output.
// If the number of bits is greater than 32, it cannot be handled.
func (f *FilterFlags) SetDropLowBits(width uint8) {
	if width > 32 {
		// todo: cannot handle
		return
	}

	mask := ^(uint32(31) << 17)
	*f &= FilterFlags(mask)
	*f |= FilterFlags(uint32(width-1) << 17)
}

// SetDropHighBits sets the number of high bits to be dropped from the output.
// If the number of bits is greater than 32, it cannot be handled.
func (f *FilterFlags) SetDropHighBits(width uint8) {
	if width > 32 {
		// todo: cannot handle
		return
	}

	mask := ^(uint32(31) << 22)
	*f &= FilterFlags(mask)
	*f |= FilterFlags(uint32(width-1) << 22)
}
