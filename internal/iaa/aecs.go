// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package iaa provides utils for using IAA abilities.
package iaa

// FilterAECS is AECS for filter operation
type FilterAECS struct {
	// On input this field contains the CRC seed. On output it is the CRC value.
	CRC uint32
	// Initial (on input) or final (on output) XOR Checksum value
	XORCheckSum         uint16
	_                   uint16 // Padding
	LowFilterParameter  uint32
	HighFilterParameter uint32

	// Base index associated with first output bit. When the output is a bit-vector
	// that is being modified, this value offsets the indices written to the output
	// and the values aggregated
	OutputModifierIndex uint32
	// The number of initial bytes in the decompressed output that should be
	// dropped before starting the filter operation.
	DropInitialDecompressOutBytes uint16
	_                             uint16 // Padding
	_                             uint64 // Padding
}

// Reset resets the values of FilterAECS to their zero values.
func (s *FilterAECS) Reset() {
	*s = FilterAECS{}
}

// DecompressAECS is AECS for decompression operation
type DecompressAECS struct {
	FilterAECS
	_ [17]uint64 // reserved
	// The IAA output accumulator is 8 bytes in size
	OutputAccumulatorData [8]byte
	// Number of valid data bits in output accumulator max 63
	OutputBitsValid uint8
	_               [3]uint8 // reserved
	// Total number of consumed bits on input.
	BitOffsetForIndexing uint32
	// The IAA input accumulator is 256 bytes in size.
	InputAccumulatorData [256]byte
	// The number of bytes valid in the corresponding Quadword in the Input
	// Accumulator. Valid values are 0 to 64
	SizeQWs [32]uint8
	_       [4880]byte // see decompressionInternalState
	_       [3]uint64  // padding
}

// type decompressionInternalState struct {
// 	eobCamEntry               uint32
// 	_                         uint32 // reserved
// 	aluFirstTableIndex        [5]uint32
// 	aluNumCodes               [5]uint32
// 	aluFirstCode              [5]uint32
// 	aluFirstLenCode           [5]uint32
// 	llCamEntries              [21]uint32
// 	_                         uint32 // Padding
// 	llCamTotalLengths         [4]uint32
// 	distanceCamEntries        [30]uint32
// 	distanceCamTotalLengths   [5]uint32
// 	minLengthCodeLength       uint32
// 	llMappingTable            [268]byte
// 	decompressState           uint32
// 	storedBlockBytesRemaining uint32
// 	_                         [168]byte // Padding
// 	historyBufferWritePointer uint32
// 	historyBuffer             [4096]byte
// }

// CompressAECS is AECS for compression operation
type CompressAECS struct {
	CRC                   uint32   // CRC
	XORChecksum           uint16   // XOR checksum
	_                     [22]byte // Reserved
	NumAccBitsValid       uint32
	OutputAccumulatorData [256]byte // output accumulator data
	Histogram             Histogram // Histogram
}

// Reset resets the values of CompressAECS to their zero values.
func (c *CompressAECS) Reset() {
	*c = CompressAECS{}
}

// ResetKeepHistogram resets the values of CompressAECS to their zero values but keep the histogram.
func (c *CompressAECS) ResetKeepHistogram() {
	c.CRC = 0
	c.XORChecksum = 0
	c.NumAccBitsValid = 0
	c.OutputAccumulatorData = [256]byte{}
}

// Histogram is a struct that contains the Huffman tables for compression.
type Histogram struct {
	LiteralCodes  [286]int32 // Literal Codes
	_             [2]uint32
	DistanceCodes [30]int32 // Distance Codes
	_             [2]uint32
}
