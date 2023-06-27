// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package dsa

// Flag is the flag of descriptor.
type Flag uint32

const (
	// OpFlagFence indicates the device waits for previous descriptors in the same batch to complete before beginning
	// work on this descriptor.
	OpFlagFence Flag = 0x0001
	// OpFlagBlockOnFault indicates the device waits for page faults to be resolved and then continues the operation.
	OpFlagBlockOnFault Flag = 0x0002
	// OpFlagCRAddrValid indicates the completion record address is valid
	OpFlagCRAddrValid Flag = 0x0004
	// OpFlagReqCR indicates the job request a completion record.
	OpFlagReqCR        Flag = 0x0008
	opFlagReqCompIntr  Flag = 0x0010
	opFlagCrsts        Flag = 0x0020
	opFlagCr           Flag = 0x0080
	opFlagCacheControl Flag = 0x0100
	opFlagAddr1Tcs     Flag = 0x0200
	opFlagAddr2Tcs     Flag = 0x0400
	opFlagAddr3Tcs     Flag = 0x0800
	opFlagCrTcs        Flag = 0x1000
	opFlagStord        Flag = 0x2000
	opFlagDrdbk        Flag = 0x4000
	opFlagDsts         Flag = 0x8000
)

// Opcode specifies the operation to be executed
type Opcode uint8

// Opcodes which is not used have been commentted out.
const (
	// OpcodeNoOp indicates no operation
	OpcodeNoOp Opcode = 0x00
	// OpcodeBatch indicates the operation is batch.
	OpcodeBatch Opcode = 0x01
	// Drain the queue
	OpcodeDrain Opcode = 0x02
	// OpcodeMemmove indicates that the operation is moving the memory.
	OpcodeMemmove           Opcode = 0x03
	OpcodeMemfill           Opcode = 0x04
	OpcodeCompare           Opcode = 0x05
	OpcodeComparePattern    Opcode = 0x06
	OpcodeCreateDeltaRecord Opcode = 0x07
	OpcodeApplyDeltaRecord  Opcode = 0x08
	OpcodeCopyWithDualcast  Opcode = 0x09
)

// const (
// OpcodeCrcgen     Opcode = 0x10
// OpcodeCopyCrc    Opcode = 0x11
// OpcodeDIFCheck   Opcode = 0x12
// OpcodeDIFInsert  Opcode = 0x13
// OpcodeDIFStrip   Opcode = 0x14
// OpcodeDIFUpdate  Opcode = 0x15
// OpcodeCacheflush Opcode = 0x20
// )

// StatusCode is operation status.
type StatusCode uint8

/* Completion record status */
const (
	// StatusNone is the default StatusCode.
	StatusNone StatusCode = iota
	// StatusSuccess means the operation is success.
	StatusSuccess StatusCode = 0x01
	// Success with false predicate
	StatusSuccessPred StatusCode = 0x02
	// Partial completion due to page fault, when the Block on Fault flag in the descriptor is 0.
	StatusPageFaultNoBOF StatusCode = 0x03
	// Partial completion due to an Invalid Request response to a Page Request.
	StatusPageFaultIr StatusCode = 0x04
	// One or more operations in the batch completed with Status not equal to Success. This value is
	// used only for a Batch descriptor
	StatusBatchFail StatusCode = 0x05
	// Partial completion of batch due to page fault while translating the Descriptor List Address in a
	// Batch descriptor and either:
	// - Page Request Services are disabled; or
	// - An Invalid Request response was received for the Page Request for the Descriptor List
	// Address.
	// This value is used only for a Batch descriptor
	StatusBatchPageFault StatusCode = 0x06
	// Offsets in the delta record were not in increasing order. This value is used only for an Apply
	// Delta Record operation.
	StatusDrOffsetNoinc StatusCode = 0x07
	// An offset in the delta record was greater than or equal to the Transfer Size of the descriptor.
	// This value is used only for an Apply Delta Record operation
	StatusDrOffsetErange StatusCode = 0x08
	// DIF error. This value is used for the DIF Check, DIF Strip, and DIF Update operations.
	StatusDIFErr StatusCode = 0x09
	// Unsupported operation code
	StatusBadOpcode StatusCode = 0x10
	// Invalid flags. One or more flags in the descriptor Flags field contain an unsupported or reserved
	// value.
	StatusInvalidFlags StatusCode = 0x11
	// Non-zero reserved field (other than a flag in the Flags field).
	StatusNoZeroReserve StatusCode = 0x12
	// Invalid Transfer Size.
	StatusInvalidTransferSize StatusCode = 0x13
	// Descriptor Count out of range (less than 2 or greater than the maximum batch size for the WQ).
	StatusDescCountOutOfRange StatusCode = 0x14
	// Maximum Delta Record Size or Delta Record Size out of range
	StatusDROutOfRange StatusCode = 0x15
	// Overlapping buffers.
	StatusOverlapBuffers StatusCode = 0x16
	// Bits 11:0 of the two destination buffers differ in Memory Copy with Dualcast
	StatusDcastErr StatusCode = 0x17
	// Misaligned Descriptor List Address
	StatusDescListAlign StatusCode = 0x18
	// Invalid Completion Interrupt Handle.
	// - If the Request Interrupt Handle command is not supported:
	// 		o The handle is out of range of the MSI-X or IMS table.
	// - If the Request Interrupt Handle command is supported:
	// 		o The interrupt handle was not returned by the Request Interrupt Handle command.
	// 		o The interrupt handle has been revoked. See section 3.7.
	// - The PASID Enable and PASID fields in the selected interrupt table entry don’t match
	// those of the descriptor.
	StatusInvalidIntHandle StatusCode = 0x19
	// A page fault occurred while translating a Completion Record Address
	StatusCraXlat StatusCode = 0x1a
	// Completion Record Address is not 32-byte aligned
	StatusCraAlign StatusCode = 0x1b
	// Misaligned address:
	// - In a Create Delta Record or Apply Delta Record operation: Source1 Address, Source2
	// Address, Destination Address, or Transfer Size is not 8-byte aligned.
	// - In a CRC Generation or Copy with CRC Generation operation: CRC Seed Address is not
	// 4-byte aligned.
	StatusAddrAlign StatusCode = 0x1c
	// In a descriptor submitted to an SWQ, Priv is 1 and the Privileged Mode Enable field of the PCI
	// Express PASID capability is 0
	StatusPrivBad StatusCode = 0x1d
	// Incorrect Traffic Class configuration:
	StatusInvalidTrafficClassConf StatusCode = 0x1e
	// A page fault occurred while translating a Readback Address in a Drain descriptor
	StatusPageFaultRADD StatusCode = 0x1f
	// The operation failed due to a hardware error
	StatusHwErr1 StatusCode = 0x20
	// Hardware error (completion timeout or unsuccessful completion status) on a destination
	// readback operation
	StatusHWErrDRB StatusCode = 0x21
	// An error occurred during address translation
	StatusTranslationFail StatusCode = 0x22
)

func (s StatusCode) String() string {
	switch s {
	case StatusSuccessPred:
		return "SuccessPred"
		// Partial completion due to page fault, when the Block on Fault flag in the descriptor is 0.
	case StatusPageFaultNoBOF:
		return "PageFaultNoBOF"
		// Partial completion due to an Invalid Request response to a Page Request.
	case StatusPageFaultIr:
		return "PageFaultIr"
		// One or more operations in the batch completed with Status not equal to Success. This value is
		// used only for a Batch descriptor
	case StatusBatchFail:
		return "BatchFail"
		// Partial completion of batch due to page fault while translating the Descriptor List Address in a
		// Batch descriptor and either:
		// - Page Request Services are disabled; or
		// - An Invalid Request response was received for the Page Request for the Descriptor List
		// Address.
		// This value is used only for a Batch descriptor
	case StatusBatchPageFault:
		return "BatchPageFault"
		// Offsets in the delta record were not in increasing order. This value is used only for an Apply
		// Delta Record operation.
	case StatusDrOffsetNoinc:
		return "DrOffsetNoinc"
		// An offset in the delta record was greater than or equal to the Transfer Size of the descriptor.
		// This value is used only for an Apply Delta Record operation
	case StatusDrOffsetErange:
		return "DrOffsetErange"
		// DIF error. This value is used for the DIF Check, DIF Strip, and DIF Update operations.
	case StatusDIFErr:
		return "DIFErr"
		// Unsupported operation code
	case StatusBadOpcode:
		return "BadOpcode"
		// Invalid flags. One or more flags in the descriptor Flags field contain an unsupported or reserved
		// value.
	case StatusInvalidFlags:
		return "InvalidFlags"
		// Non-zero reserved field (other than a flag in the Flags field).
	case StatusNoZeroReserve:
		return "NoZeroReserve"
		// Invalid Transfer Size.
	case StatusInvalidTransferSize:
		return "InvalidTransferSize"
		// Descriptor Count out of range (less than 2 or greater than the maximum batch size for the WQ).
	case StatusDescCountOutOfRange:
		return "DescCountOutOfRange"
		// Maximum Delta Record Size or Delta Record Size out of range
	case StatusDROutOfRange:
		return "DROutOfRange"
		// Overlapping buffers.
	case StatusOverlapBuffers:
		return "OverlapBuffers"
		// Bits 11:0 of the two destination buffers differ in Memory Copy with Dualcast
	case StatusDcastErr:
		return "DcastErr"
		// Misaligned Descriptor List Address
	case StatusDescListAlign:
		return "DescListAlign"
		// Invalid Completion Interrupt Handle.
		// - If the Request Interrupt Handle command is not supported:
		// 		o The handle is out of range of the MSI-X or IMS table.
		// - If the Request Interrupt Handle command is supported:
		// 		o The interrupt handle was not returned by the Request Interrupt Handle command.
		// 		o The interrupt handle has been revoked. See section 3.7.
		// - The PASID Enable and PASID fields in the selected interrupt table entry don’t match
		// those of the descriptor.
	case StatusInvalidIntHandle:
		return "InvalidIntHandle"
		// A page fault occurred while translating a Completion Record Address
	case StatusCraXlat:
		return "CraXlat"
		// Completion Record Address is not 32-byte aligned
	case StatusCraAlign:
		return "CraAlign"
		// Misaligned address:
		// - In a Create Delta Record or Apply Delta Record operation: Source1 Address, Source2
		// Address, Destination Address, or Transfer Size is not 8-byte aligned.
		// - In a CRC Generation or Copy with CRC Generation operation: CRC Seed Address is not
		// 4-byte aligned.
	case StatusAddrAlign:
		return "AddrAlign"
		// In a descriptor submitted to an SWQ, Priv is 1 and the Privileged Mode Enable field of the PCI
		// Express PASID capability is 0
	case StatusPrivBad:
		return "PrivBad"
		// Incorrect Traffic Class configuration:
	case StatusInvalidTrafficClassConf:
		return "InvalidTrafficClassConf"
		// A page fault occurred while translating a Readback Address in a Drain descriptor
	case StatusPageFaultRADD:
		return "PageFaultRADD"
		// The operation failed due to a hardware error
	case StatusHwErr1:
		return "HwErr1"
		// Hardware error (completion timeout or unsuccessful completion status) on a destination
		// readback operation
	case StatusHWErrDRB:
		return "HWErrDRB"
		// An error occurred during address translation
	case StatusTranslationFail:
		return "TranslationFail"
	default:
		return "unknown"
	}
}
