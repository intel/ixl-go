// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package iaa

// Opcode specifies the operation to be executed
type Opcode uint8

const (
	Noop         Opcode = 0x00 // Noop specifies No operation.
	OpDrain      Opcode = 0x02 // OpDrain specifies Drain operation.
	OpDecompress Opcode = 0x42 // OpDecompress specifies Decompress operation
	OpCompress   Opcode = 0x43 // OpCompress specifies Compress operation
	OpCRC64      Opcode = 0x44 // OpCRC64 specifies Calculate CRC64 checksum operation
	OpScan       Opcode = 0x50 // OpScan specifies Scan operation
	OpExtract    Opcode = 0x52 // OpExtract specifies Extract operation
	OpSelect     Opcode = 0x53 // OpSelect specifies Select operation
	OpExpand     Opcode = 0x56 // OpExpand specifies Expand operation
)
