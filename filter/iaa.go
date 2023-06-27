// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"strconv"
	"unsafe"

	"github.com/intel/ixl-go/internal/iaa"
)

func scanInt[R DataUnit](
	desc *iaa.Descriptor,
	input []R,
	output []byte,
	aesc *iaa.FilterAECS,
	cr *iaa.CompletionRecord,
) {
	desc.SetOpcode(iaa.OpScan)

	desc.SetFlag(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagCompletionRecordValid |
			iaa.FlagRequestCompletionRecord |
			iaa.FlagReadSource2Aecs,
	)

	var r R

	desc.FilterFlags.SetSource1Width(uint8(unsafe.Sizeof(r) * 8))
	desc.Src1Addr = uintptr(unsafe.Pointer(&input[0]))
	desc.Size = uint32(len(input) * int(unsafe.Sizeof(r)))

	desc.DestAddr = uintptr(unsafe.Pointer(&output[0]))
	desc.MaxDestionationSize = uint32(len(output))

	desc.Src2Addr = uintptr(unsafe.Pointer(aesc))
	desc.Src2Size = uint32(unsafe.Sizeof(iaa.FilterAECS{}))
	desc.ElementsNumber = uint32(len(input))
	desc.SetCompleteRecord(uintptr(unsafe.Pointer(cr)))
}

func scanBitPacking(
	desc *iaa.Descriptor,
	width uint8,
	size int,
	input []byte,
	output []byte,
	aesc *iaa.FilterAECS,
	cr *iaa.CompletionRecord,
) {
	desc.SetOpcode(iaa.OpScan)

	desc.SetFlag(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagCompletionRecordValid |
			iaa.FlagRequestCompletionRecord |
			iaa.FlagReadSource2Aecs,
	)

	desc.FilterFlags.SetSource1Width(width)

	desc.Src1Addr = uintptr(unsafe.Pointer(&input[0]))
	desc.Size = uint32(len(input))

	desc.DestAddr = uintptr(unsafe.Pointer(&output[0]))
	desc.MaxDestionationSize = uint32(len(output))

	desc.Src2Addr = uintptr(unsafe.Pointer(aesc))
	desc.Src2Size = uint32(unsafe.Sizeof(iaa.FilterAECS{}))
	desc.ElementsNumber = uint32(size)
	desc.SetCompleteRecord(uintptr(unsafe.Pointer(cr)))
}

func scanRLE(desc *iaa.Descriptor,
	size int,
	input []byte,
	output []byte,
	aesc *iaa.FilterAECS,
	cr *iaa.CompletionRecord,
) {
	desc.SetOpcode(iaa.OpScan)

	desc.SetFlag(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagCompletionRecordValid |
			iaa.FlagRequestCompletionRecord |
			iaa.FlagReadSource2Aecs,
	)

	desc.FilterFlags |= iaa.FilterFlagSource1ParquetRLE
	desc.Src1Addr = uintptr(unsafe.Pointer(&input[0]))
	desc.Size = uint32(len(input))

	desc.DestAddr = uintptr(unsafe.Pointer(&output[0]))
	desc.MaxDestionationSize = uint32(len(output))

	desc.Src2Addr = uintptr(unsafe.Pointer(aesc))
	desc.Src2Size = uint32(unsafe.Sizeof(iaa.FilterAECS{}))
	desc.ElementsNumber = uint32(size)
	desc.SetCompleteRecord(uintptr(unsafe.Pointer(cr)))
}

func extractRLE[R DataUnit](desc *iaa.Descriptor,
	size int,
	input []byte,
	output []R,
	aesc *iaa.FilterAECS,
	cr *iaa.CompletionRecord,
) {
	desc.SetOpcode(iaa.OpExtract)

	desc.SetFlag(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagCompletionRecordValid |
			iaa.FlagRequestCompletionRecord |
			iaa.FlagReadSource2Aecs,
	)

	desc.FilterFlags |= iaa.FilterFlagSource1ParquetRLE
	var r R
	typeSize := unsafe.Sizeof(r)
	switch typeSize {
	case 1:
		desc.FilterFlags |= iaa.FilterFlagOutputWithByte
	case 2:
		desc.FilterFlags |= iaa.FilterFlagOutputWithWord
	case 4:
		desc.FilterFlags |= iaa.FilterFlagOutputWithDword
	default:
		panic("unknown type:" + strconv.Itoa(int(typeSize)))
	}

	desc.Src1Addr = uintptr(unsafe.Pointer(&input[0]))
	desc.Size = uint32(len(input))

	desc.DestAddr = uintptr(unsafe.Pointer(&output[0]))
	desc.MaxDestionationSize = uint32(len(output) * int(typeSize))

	desc.Src2Addr = uintptr(unsafe.Pointer(aesc))
	desc.Src2Size = uint32(unsafe.Sizeof(iaa.FilterAECS{}))
	desc.ElementsNumber = uint32(size)
	desc.SetCompleteRecord(uintptr(unsafe.Pointer(cr)))
}

func selectUints[R DataUnit](desc *iaa.Descriptor, input []R, secInput []byte, output []R, cr *iaa.CompletionRecord) {
	desc.SetOpcode(iaa.OpSelect)

	desc.SetFlag(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagCompletionRecordValid |
			iaa.FlagRequestCompletionRecord |
			iaa.FlagReadSource2SecondaryInputToFilterFunction,
	)

	var r R

	desc.FilterFlags.SetSource1Width(uint8(unsafe.Sizeof(r) * 8))
	desc.Src1Addr = uintptr(unsafe.Pointer(&input[0]))
	desc.Size = uint32(len(input) * int(unsafe.Sizeof(r)))

	desc.DestAddr = uintptr(unsafe.Pointer(&output[0]))
	desc.MaxDestionationSize = uint32(len(output) * int(unsafe.Sizeof(r)))

	desc.Src2Addr = uintptr(unsafe.Pointer(&secInput[0]))
	desc.Src2Size = uint32(len(secInput))
	desc.ElementsNumber = uint32(len(input))
	desc.SetCompleteRecord(uintptr(unsafe.Pointer(cr)))
}

func expandUints[R DataUnit](desc *iaa.Descriptor, input []R, secInput []byte, output []R, cr *iaa.CompletionRecord) {
	desc.SetOpcode(iaa.OpExpand)

	desc.SetFlag(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagCompletionRecordValid |
			iaa.FlagRequestCompletionRecord |
			iaa.FlagReadSource2SecondaryInputToFilterFunction,
	)

	var r R

	desc.FilterFlags.SetSource1Width(uint8(unsafe.Sizeof(r) * 8))

	desc.Src1Addr = uintptr(unsafe.Pointer(&input[0]))
	desc.Size = uint32(len(input) * int(unsafe.Sizeof(r)))

	desc.DestAddr = uintptr(unsafe.Pointer(&output[0]))
	desc.MaxDestionationSize = uint32(len(output) * int(unsafe.Sizeof(r)))

	desc.Src2Addr = uintptr(unsafe.Pointer(&secInput[0]))
	desc.Src2Size = uint32(unsafe.Sizeof(len(secInput)))
	desc.ElementsNumber = uint32(len(output))
	desc.SetCompleteRecord(uintptr(unsafe.Pointer(cr)))
}
