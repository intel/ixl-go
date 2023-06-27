// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package iaa

import (
	"strings"
)

// DecompressionFlag is used to specify options for decompression.
type DecompressionFlag uint16

// String returns a string representation of the decompression flags.
func (d DecompressionFlag) String() string {
	var flags []string

	if (d & DecompressionFlagEnableDecompression) != 0 {
		flags = append(flags, "ENABLE_DECOMPRESSION")
	}
	if (d & DecompressionFlagFlushOutput) != 0 {
		flags = append(flags, "ENABLE_DECOMPRESSION")
	}
	if (d & DecompressionFlagStopOnEOB) != 0 {
		flags = append(flags, "STOP_ON_EOB")
	}
	if (d & DecompressionFlagCheckForEOB) != 0 {
		flags = append(flags, "CHECK_FOR_EOB")
	}
	if (d & DecompressionFlagSelectBFinalEOB) != 0 {
		flags = append(flags, "SELECT_B_FINAL_EOB")
	}
	if (d & DecompressionFlagBigEndian) != 0 {
		flags = append(flags, "DECOMPRESS_BIT_ORDER")
	}
	if (d & DecompressionFlagIgnoreEndBits) != 0 {
		flags = append(flags, "IGNORE_END_BITS")
	}
	if (d & DecompressionFlagSupressOutput) != 0 {
		flags = append(flags, "SUPPRESS_OUTPUT")
	}
	return strings.Join(flags, " | ")
}

// Indexing sets the indexing type.
func (d DecompressionFlag) Indexing(t IndexingType) DecompressionFlag {
	return d | DecompressionFlag(uint16(t)<<10)
}

// Decompression flags.
const (
	// DecompressionFlagEnableDecompression is used to Enable Decompression.
	// If Operation is Decompress, this flag must be set.
	DecompressionFlagEnableDecompression DecompressionFlag = 1 << iota
	// A partial output word is written to the output stream. If it would overflow the output buffer, it
	// is saved in the AECS, so that the job can be completed by a subsequent descriptor. This value
	// should be used for a Decompress descriptor that is the last (or only) descriptor in a job. For filter
	// operations, output flushing is automatic and this flag is ignored.
	DecompressionFlagFlushOutput
	// Stop decompression when an EOB code is encountered.
	DecompressionFlagStopOnEOB
	// Check for an EOB code and stop decompression if it is encountered.
	DecompressionFlagCheckForEOB
	// Selects the B-final EOB code.
	DecompressionFlagSelectBFinalEOB
	// Decompress big endian.
	DecompressionFlagBigEndian
	// Ignore end bits.
	DecompressionFlagIgnoreEndBits
	// Suppress output.
	DecompressionFlagSupressOutput
)

// CompressionFlag is used to specify options for compression.
type CompressionFlag uint16

// Compression flags.
const (
	// StatsMode.
	CompressionFlagStatsMode CompressionFlag = 1
	// A partial output word is written to the output stream. If it would overflow the output buffer, it
	// is saved in the AECS, so that the job can be completed by a subsequent descriptor. This value
	// should be used for a Decompress descriptor that is the last (or only) descriptor in a job. For filter
	// operations, output flushing is automatic and this flag is ignored.
	CompressionFlagFlushOutput = 1 << 1
	// End append EOB.
	CompressionFlagEndAppendEOB = 1 << 2
	// End append EOB non-B-final.
	CompressionFlagEndAppendEOBNonBFinal = 2 << 2
	// End append EOB and B-final.
	CompressionFlagEndAppendEOBAndBFinal = 3 << 2
	// Generate all literals.
	CompressionFlagGenerateAllLiterals = 1 << 4
	// Set compress using big endian
	CompressionFlagCompressBigEndian = 1 << 5
)

// Indexing sets the indexing type.
func (c CompressionFlag) Indexing(t IndexingType) CompressionFlag {
	return c | CompressionFlag(uint16(t)<<6)
}

// IndexingType is used to specify the indexing type.
type IndexingType uint8

// Indexing types.
const (
	// Disable indexing.
	DisableIndexing IndexingType = iota
	// Enable indexing (512 bytes).
	EnableIndexing512
	// Enable indexing (1 Kb).
	EnableIndexing1Kb
	// Enable indexing (2 Kb).
	EnableIndexing2Kb
	// Enable indexing (4 Kb).
	EnableIndexing4Kb
	// Enable indexing (8 Kb).
	EnableIndexing8Kb
	// Enable indexing (16 Kb).
	EnableIndexing16Kb
	// Enable indexing (32 Kb).
	EnableIndexing32Kb
)
