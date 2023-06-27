// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package iaa

import (
	"fmt"
	"strings"

	"github.com/intel/ixl-go/util/mem"
)

// Descriptor represents a descriptor for IAA engine.
type Descriptor struct {
	Header              uint32      // A header field that contains the pasid, reversed and priv flags
	FlagsAndOpCode      uint32      // A field that contains the flags and opcode
	CompletionAddr      uintptr     // The completion address
	Src1Addr            uintptr     // The source 1 address
	DestAddr            uintptr     // The destination address
	Size                uint32      // The size of the descriptor
	IntHandle           uint16      // The interrupt handle
	compressionFlag     uint16      // The compression flag
	Src2Addr            uintptr     // The source 2 address
	MaxDestionationSize uint32      // The maximum size of the destination
	Src2Size            uint32      // The source 2 size
	FilterFlags         FilterFlags // Filter flags
	ElementsNumber      uint32      // The number of elements
}

// NewDescriptor returns a new descriptor.
func NewDescriptor() *Descriptor {
	return mem.Alloc64Align[Descriptor]()
}

// String returns a string representation of the descriptor.
func (d *Descriptor) String() string {
	return fmt.Sprintf("opcode:[%d] flags:[%b] comp_flag:[%b]",
		d.GetOpcode(),
		d.GetFlags(),
		d.compressionFlag,
	)
}

// SetCompressionFlag sets the compression flag.
func (d *Descriptor) SetCompressionFlag(c CompressionFlag) {
	d.compressionFlag = uint16(c)
}

// AddCompressionFlag adds the compression flag.
func (d *Descriptor) AddCompressionFlag(c CompressionFlag) {
	d.compressionFlag |= uint16(c)
}

// GetCompressionFlag returns the compression flag.
func (d *Descriptor) GetCompressionFlag() CompressionFlag {
	return CompressionFlag(d.compressionFlag)
}

// GetDecompressionFlag returns the decompression flag.
func (d *Descriptor) GetDecompressionFlag() DecompressionFlag {
	return DecompressionFlag(d.compressionFlag)
}

// SetCompleteRecord sets the completion record.
func (d *Descriptor) SetCompleteRecord(cr uintptr) {
	d.SetFlag(FlagRequestCompletionRecord | FlagCompletionRecordValid)
	d.CompletionAddr = cr
}

// SetDecompressionFlag sets the decompression flag.
func (d *Descriptor) SetDecompressionFlag(c DecompressionFlag) {
	d.compressionFlag = uint16(c)
}

// GetFlags returns the descriptor flags.
func (d *Descriptor) GetFlags() DescriptorFlag { return DescriptorFlag((d.FlagsAndOpCode << 8) >> 8) }

// SetFlag sets the descriptor flag.
func (d *Descriptor) SetFlag(value DescriptorFlag) {
	current := d.GetFlags()
	current = current | value
	d.SetFlags(current)
}

// SetFlags sets the descriptor flags.
func (d *Descriptor) SetFlags(value DescriptorFlag) {
	rv := uint32(value)

	left := (d.FlagsAndOpCode >> 24) << 24

	d.FlagsAndOpCode = left | rv
}

// GetOpcode returns the opcode.
func (d *Descriptor) GetOpcode() Opcode { return Opcode(d.FlagsAndOpCode >> 24) }

// SetOpcode sets the opcode.
func (d *Descriptor) SetOpcode(value Opcode) {
	rv := uint32(value)

	right := (d.FlagsAndOpCode << (32 - 24)) >> (32 - 24)

	d.FlagsAndOpCode = (rv << 24) | right
}

// DescriptorFlag represents a descriptor flag.
type DescriptorFlag uint32

const (
	_ DescriptorFlag = 1 << iota
	// FlagBlockOnFault specifies whether the device waits for page faults to be resolved and then continues the operation.
	FlagBlockOnFault
	// FlagCompletionRecordValid specifies whether the descriptor completion record is valid.
	FlagCompletionRecordValid
	// FlagRequestCompletionRecord specifies whether the descriptor requests a completion record.
	FlagRequestCompletionRecord
	// FlagRequestCompletionInterrupt specifies whether the descriptor requests a completion interrupt.
	FlagRequestCompletionInterrupt
	// FlagCompletionRecordSteeringTagSelector is for Optane.
	FlagCompletionRecordSteeringTagSelector
	_
	_
	// FlagCacheControl specifies whether to hint to direct data writes to CPU cache
	FlagCacheControl
	_ // FlagSource1TcSelector
	_ // FlagDestinationTcSelector
	_ // FlagSource2TcSelector
	_ // FlagCompletionRecordTcSelector
	// FlagStrictOrdering specifies whether to force strict ordering of all memory writes.
	FlagStrictOrdering
	_ // FlagDestinationReadback
	_ // FlagDestinationSteeringTagSelector
	// FlagReadSource2Aecs specifies whether the Source 2 is read as AECS
	FlagReadSource2Aecs
	// FlagReadSource2SecondaryInputToFilterFunction specifies
	// whether the Source 2 is read as secondary input to filter function
	FlagReadSource2SecondaryInputToFilterFunction
	// FlagWriteSource2CompletionOfOperation specifies whether the engine writes source 2 at completion of operation.
	FlagWriteSource2CompletionOfOperation
	// FlagWriteSource2OnlyIfOutputOverflow specifies whether the engine writes source 2 only if output overflow occurs.
	FlagWriteSource2OnlyIfOutputOverflow
	_ // FlagSource2SteeringTagSelector
	// FlagCRCSelectRFC3720 specifies whether the CRC job is RFC 3720 CRC.
	FlagCRCSelectRFC3720
	// FlagAecsRWToggleSelector specifies whether the engine's reads are done from (A+S) and writes are done to (A)
	FlagAecsRWToggleSelector
)

// Reset resets the descriptor
func (d *Descriptor) Reset() {
	*d = Descriptor{}
}

func (d DescriptorFlag) String() string {
	var flags []string

	if (d & FlagBlockOnFault) != 0 {
		flags = append(flags, "BLOCK_ON_FAULT")
	}

	if (d & FlagCompletionRecordValid) != 0 {
		flags = append(flags, "COMPLETION_RECORD_VALID")
	}

	if (d & FlagRequestCompletionRecord) != 0 {
		flags = append(flags, "REQUEST_COMPLETION_RECORD")
	}

	if (d & FlagRequestCompletionInterrupt) != 0 {
		flags = append(flags, "REQUEST_COMPLETION_INTERRUPT")
	}

	if (d & FlagCompletionRecordSteeringTagSelector) != 0 {
		flags = append(flags, "COMPLETION_RECORD_STEERING_TAG_SELECTOR")
	}

	if (d & FlagCacheControl) != 0 {
		flags = append(flags, "CACHE_CONTROL")
	}

	if (d & FlagStrictOrdering) != 0 {
		flags = append(flags, "STRICT_ORDERING")
	}

	if (d & FlagReadSource2Aecs) != 0 {
		flags = append(flags, "READ_SOURCE_2_AECS")
	}

	if (d & FlagReadSource2SecondaryInputToFilterFunction) != 0 {
		flags = append(flags, "READ_SOURCE_2_SECONDARY_INPUT_TO_FILTER_FUNCTION")
	}

	if (d & FlagWriteSource2CompletionOfOperation) != 0 {
		flags = append(flags, "WRITE_SOURCE_2_COMPLETION_OF_OPERATION")
	}

	if (d & FlagWriteSource2OnlyIfOutputOverflow) != 0 {
		flags = append(flags, "WRITE_SOURCE_2_ONLY_IF_OUTPUT_OVERFLOW_OCCURS")
	}

	if (d & FlagCRCSelectRFC3720) != 0 {
		flags = append(flags, "CRC_SELECT_RFC_3720")
	}

	if (d & FlagAecsRWToggleSelector) != 0 {
		flags = append(flags, "AECS_R_W_TOGGLE_SELECTOR")
	}

	return strings.Join(flags, " | ")
}
