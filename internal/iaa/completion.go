// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package iaa

import (
	"fmt"

	"github.com/intel/ixl-go/internal/device"
	"github.com/intel/ixl-go/internal/errors"
	"github.com/intel/ixl-go/util/mem"
)

// CompletionRecord is a representation of a completion record.
type CompletionRecord struct {
	Header               device.CompletionRecordHeader
	FaultAddress         uint64
	InvalidFlags         uint32
	_                    uint32
	OutputSize           uint32
	OutputBits           uint8
	_                    uint8
	XORCheckSum          uint16
	CRC                  uint32
	MinOrFirst           uint32
	MaxOrLast            uint32
	SumOrPopulationCount uint32
	_                    [2]uint64
}

// NewCompletionRecord returns a new CompletionRecord.
func NewCompletionRecord() *CompletionRecord {
	return mem.Alloc64Align[CompletionRecord]()
}

// Reset resets the CompletionRecord.
func (c *CompletionRecord) Reset() {
	*c = CompletionRecord{}
}

// String returns a string representation of the CompletionRecord.
func (c *CompletionRecord) String() string {
	return fmt.Sprintf("header:[%s] "+
		"fault_address:[%d] "+
		"invalid_flags:[%b] "+
		"output_size:[%d] "+
		"output_bits:[%d] "+
		"XORCheckSum:[%d] "+
		"CRC:[%d] "+
		"min_or_first:[%d] "+
		"max_or_last:[%d] "+
		"sum_or_poluation_count:[%d]",
		c.GetHeader().String(),
		c.FaultAddress,
		c.InvalidFlags,
		c.OutputSize,
		c.OutputBits,
		c.XORCheckSum,
		c.CRC,
		c.MinOrFirst,
		c.MaxOrLast,
		c.SumOrPopulationCount,
	)
}

// CompletionRecordHeader is the header of a CompletionRecord.
type CompletionRecordHeader struct {
	StatusCode
	StatusReadFault bool
	ErrorCode       ErrorCode
	Completed       uint32
}

// String returns a string representation of the CompletionRecordHeader.
func (crh CompletionRecordHeader) String() string {
	return fmt.Sprintf("status_code[%s] page_fault[%v] error_code[%s] completed[%d]",
		crh.StatusCode.String(),
		crh.StatusReadFault,
		crh.ErrorCode,
		crh.Completed)
}

// GetHeader returns the CompletionRecordHeader.
func (c *CompletionRecord) GetHeader() (crh CompletionRecordHeader) {
	return fromHeader(c.Header)
}

// CheckError checks if error happened.
func (c *CompletionRecord) CheckError() error {
	h := c.GetHeader()
	if h.StatusCode == Success {
		return nil
	}
	return errors.HardwareError{Status: h.StatusCode.String(), ErrorCode: h.ErrorCode.String()}
}

// StatusCode is a status code.
type StatusCode uint8

const (
	// Uncompelete indicates the operation is not completed.
	Uncompelete StatusCode = 0x00
	// Success indicates the operation is success.
	Success StatusCode = 0x01
	// AnalyticsError is an analytics error status code.
	AnalyticsError StatusCode = 0x0a
	// OutputBufferOverflow is an output buffer overflow status code.
	OutputBufferOverflow StatusCode = 0x0b
	// InvalidFlags is an invalid flags status code.
	InvalidFlags StatusCode = 0x11
	// NonZeroReservedField is a non-zero reserved field status code.
	NonZeroReservedField StatusCode = 0x12
	// InvalidSizeValue is an invalid size value status code.
	InvalidSizeValue StatusCode = 0x13
	// CompletionRecordAddressNotAligned is a completion record address not aligned status code.
	CompletionRecordAddressNotAligned StatusCode = 0x1b
	// AecsMisalignedAddress is an AECS misaligned address status code.
	AecsMisalignedAddress StatusCode = 0x1c
	// PageRequestTimeout is a page request timeout status code.
	PageRequestTimeout StatusCode = 0x23
	// WatchdogExpired is a watchdog expired status code.
	WatchdogExpired StatusCode = 0x24
	// InvalidOpFlags is an invalid operation flags status code.
	InvalidOpFlags StatusCode = 0x30
	// InvalidFilterFlags is an invalid filter flags status code.
	InvalidFilterFlags StatusCode = 0x31
	// InvalidInputSize is an invalid input size status code.
	InvalidInputSize StatusCode = 0x32
	// InvalidNumberOfElements is an invalid number of elements status code.
	InvalidNumberOfElements StatusCode = 0x33
	// InvalidSource1Width is an invalid source1 width status code.
	InvalidSource1Width StatusCode = 0x34
	// InvalidInvertOutput is an invalid invert output status code.
	InvalidInvertOutput StatusCode = 0x35
)

// ErrorCode is an error code.
type ErrorCode uint8

const (
	// ErrorCodeHeaderTooLarge is a header too large error code.
	ErrorCodeHeaderTooLarge ErrorCode = 1 + iota
	// ErrorCodeUndefinedClCode is an undefined CL code error code.
	ErrorCodeUndefinedClCode
	// ErrorCodeFirstCodeInLlTreeIs16 is a first code in LL tree is 16 error code.
	ErrorCodeFirstCodeInLlTreeIs16
	// ErrorCodeFirstCodeInDTreeIs16 is a first code in D tree is 16 error code.
	ErrorCodeFirstCodeInDTreeIs16
	// ErrorCodeNoValidLlCode is a no valid LL code error code.
	ErrorCodeNoValidLlCode
	// ErrorCodeWrongNumberOfLlCodes is a wrong number of LL codes error code.
	ErrorCodeWrongNumberOfLlCodes
	// ErrorCodeWrongNumberOfDistCodes is a wrong number of dist codes error code.
	ErrorCodeWrongNumberOfDistCodes
	// ErrorCodeBadClCodeLengths is a bad CL code lengths error code.
	ErrorCodeBadClCodeLengths
	// ErrorCodeBadLlCodeLengths is a bad LL code lengths error code.
	ErrorCodeBadLlCodeLengths
	// ErrorCodeBadDistCodeLengths is a bad dist code lengths error code.
	ErrorCodeBadDistCodeLengths
	// ErrorCodeBadLlCodes is a bad LL codes error code.
	ErrorCodeBadLlCodes
	// ErrorCodeBadDCode is a bad D code error code.
	ErrorCodeBadDCode
	// ErrorCodeInvalidBlockType is an invalid block type error code.
	ErrorCodeInvalidBlockType
	// ErrorCodeInvalidStoredLength is an invalid stored length error code.
	ErrorCodeInvalidStoredLength
	// ErrorCodeBadEndOfFile is a bad end of file error code.
	ErrorCodeBadEndOfFile
	// ErrorCodeBadLengthDecode is a bad length decode error code.
	ErrorCodeBadLengthDecode
	// ErrorCodeBadDistanceDecode is a bad distance decode error code.
	ErrorCodeBadDistanceDecode
	// ErrorCodeDistanceBeforeStartOfFile is a distance before start of file error code.
	ErrorCodeDistanceBeforeStartOfFile
	// ErrorCodeTimeout is a timeout error code.
	ErrorCodeTimeout
	// ErrorCodePrleFormatError is a PRLE format error code.
	ErrorCodePrleFormatError
	// ErrorCodeFilterFunctionWordOverflow is a filter function word overflow error code.
	ErrorCodeFilterFunctionWordOverflow
	// ErrorCodeAecsError is an AECS error code.
	ErrorCodeAecsError
	// ErrorCodeSource1TooSmall is a source1 too small error code.
	ErrorCodeSource1TooSmall
	// ErrorCodeSource2TooSmall is a source2 too small error code.
	ErrorCodeSource2TooSmall
	// ErrorCodeUnrecoverableOutputOverflow is an unrecoverable output overflow error code.
	ErrorCodeUnrecoverableOutputOverflow
	// ErrorCodeDistanceSpansMiniBlocks is a distance spans mini blocks error code.
	ErrorCodeDistanceSpansMiniBlocks
	// ErrorCodeLengthSpansMiniBlocks is a length spans mini blocks error code.
	ErrorCodeLengthSpansMiniBlocks
	// ErrorCodeInvalidBlockSize is an invalid block size error code.
	ErrorCodeInvalidBlockSize
	// ErrorCodeZcompressVerifyFailure is a ZCompress verify failure error code.
	ErrorCodeZcompressVerifyFailure
	// ErrorCodeInvalidHuffmanCode is an invalid Huffman code error code.
	ErrorCodeInvalidHuffmanCode
	// ErrorCodePrleBitWidthTooLarge is a PRLE bit width too large error code.
	ErrorCodePrleBitWidthTooLarge
	// ErrorCodeTooFewElementsProcessed is a too few elements processed error code.
	ErrorCodeTooFewElementsProcessed
	// ErrorCodeInvalidRleCount is an invalid RLE count error code.
	ErrorCodeInvalidRleCount
	// ErrorCodeInvalidZDecompressHeader is an invalid ZDecompress header error code.
	ErrorCodeInvalidZDecompressHeader
	// ErrorCodeTooManyLlCodes is a too many LL codes error code.
	ErrorCodeTooManyLlCodes
	// ErrorCodeTooManyDCodes is a too many D codes error code.
	ErrorCodeTooManyDCodes
	// ErrorCodeAdministrativeTimeout is an administrative timeout error code.
	ErrorCodeAdministrativeTimeout
)

// String returns a string representation of the StatusCode.
func (c StatusCode) String() string {
	switch c {
	case Uncompelete:
		return "UNCOMPELETE"
	case Success:
		return "SUCCESS"
	case AnalyticsError:
		return "ANALYTICS_ERROR"
	case OutputBufferOverflow:
		return "OUTPUT_BUFFER_OVERFLOW"
	case InvalidFlags:
		return "INVALID_FLAGS"
	case NonZeroReservedField:
		return "NON_ZERO_RESERVED_FIELD"
	case InvalidSizeValue:
		return "INVALID_SIZE_VALUE"
	case CompletionRecordAddressNotAligned:
		return "COMPLETION_RECORD_ADDRESS_NOT_ALIGNED"
	case AecsMisalignedAddress:
		return "AECS_MISALIGNED_ADDRESS"
	case PageRequestTimeout:
		return "PAGE_REQUEST_TIMEOUT"
	case WatchdogExpired:
		return "WATCHDOG_EXPIRED"
	case InvalidOpFlags:
		return "INVALID_OP_FLAGS"
	case InvalidFilterFlags:
		return "INVALID_FILTER_FLAGS"
	case InvalidInputSize:
		return "INVALID_INPUT_SIZE"
	case InvalidNumberOfElements:
		return "INVALID_NUMBER_OF_ELEMENTS"
	case InvalidSource1Width:
		return "INVALID_SOURCE_1_WIDTH"
	case InvalidInvertOutput:
		return "INVALID_INVERT_OUTPUT"
	}
	return ""
}

// String returns a string representation of the ErrorCode.
func (e ErrorCode) String() string {
	switch e {
	case ErrorCodeHeaderTooLarge:
		return "ERROR_CODE_HEADER_TOO_LARGE"
	case ErrorCodeUndefinedClCode:
		return "ERROR_CODE_UNDEFINED_CL_CODE"
	case ErrorCodeFirstCodeInLlTreeIs16:
		return "ERROR_CODE_FIRST_CODE_IN_LL_TREE_IS_16"
	case ErrorCodeFirstCodeInDTreeIs16:
		return "ERROR_CODE_FIRST_CODE_IN_D_TREE_IS_16"
	case ErrorCodeNoValidLlCode:
		return "ERROR_CODE_NO_VALID_LL_CODE"
	case ErrorCodeWrongNumberOfLlCodes:
		return "ERROR_CODE_WRONG_NUMBER_OF_LL_CODES"
	case ErrorCodeWrongNumberOfDistCodes:
		return "ERROR_CODE_WRONG_NUMBER_OF_DIST_CODES"
	case ErrorCodeBadClCodeLengths:
		return "ERROR_CODE_BAD_CL_CODE_LENGTHS"
	case ErrorCodeBadLlCodeLengths:
		return "ERROR_CODE_BAD_LL_CODE_LENGTHS"
	case ErrorCodeBadDistCodeLengths:
		return "ERROR_CODE_BAD_DIST_CODE_LENGTHS"
	case ErrorCodeBadLlCodes:
		return "ERROR_CODE_BAD_LL_CODES"
	case ErrorCodeBadDCode:
		return "ERROR_CODE_BAD_D_CODE"
	case ErrorCodeInvalidBlockType:
		return "ERROR_CODE_INVALID_BLOCK_TYPE"
	case ErrorCodeInvalidStoredLength:
		return "ERROR_CODE_INVALID_STORED_LENGTH"
	case ErrorCodeBadEndOfFile:
		return "ERROR_CODE_BAD_END_OF_FILE"
	case ErrorCodeBadLengthDecode:
		return "ERROR_CODE_BAD_LENGTH_DECODE"
	case ErrorCodeBadDistanceDecode:
		return "ERROR_CODE_BAD_DISTANCE_DECODE"
	case ErrorCodeDistanceBeforeStartOfFile:
		return "ERROR_CODE_DISTANCE_BEFORE_START_OF_FILE"
	case ErrorCodeTimeout:
		return "ERROR_CODE_TIMEOUT"
	case ErrorCodePrleFormatError:
		return "ERROR_CODE_PRLE_FORMAT_ERROR"
	case ErrorCodeFilterFunctionWordOverflow:
		return "ERROR_CODE_FILTER_FUNCTION_WORD_OVERFLOW"
	case ErrorCodeAecsError:
		return "ERROR_CODE_AECS_ERROR"
	case ErrorCodeSource1TooSmall:
		return "ERROR_CODE_SOURCE_1_TOO_SMALL"
	case ErrorCodeSource2TooSmall:
		return "ERROR_CODE_SOURCE_2_TOO_SMALL"
	case ErrorCodeUnrecoverableOutputOverflow:
		return "ERROR_CODE_UNRECOVERABLE_OUTPUT_OVERFLOW"
	case ErrorCodeDistanceSpansMiniBlocks:
		return "ERROR_CODE_DISTANCE_SPANS_MINI_BLOCKS"
	case ErrorCodeLengthSpansMiniBlocks:
		return "ERROR_CODE_LENGTH_SPANS_MINI_BLOCKS"
	case ErrorCodeInvalidBlockSize:
		return "ERROR_CODE_INVALID_BLOCK_SIZE"
	case ErrorCodeZcompressVerifyFailure:
		return "ERROR_CODE_ZCOMPRESS_VERIFY_FAILURE"
	case ErrorCodeInvalidHuffmanCode:
		return "ERROR_CODE_INVALID_HUFFMAN_CODE"
	case ErrorCodePrleBitWidthTooLarge:
		return "ERROR_CODE_PRLE_BIT_WIDTH_TOO_LARGE"
	case ErrorCodeTooFewElementsProcessed:
		return "ERROR_CODE_TOO_FEW_ELEMENTS_PROCESSED"
	case ErrorCodeInvalidRleCount:
		return "ERROR_CODE_INVALID_RLE_COUNT"
	case ErrorCodeInvalidZDecompressHeader:
		return "ERROR_CODE_INVALID_Z_DECOMPRESS_HEADER"
	case ErrorCodeTooManyLlCodes:
		return "ERROR_CODE_TOO_MANY_LL_CODES"
	case ErrorCodeTooManyDCodes:
		return "ERROR_CODE_TOO_MANY_D_CODES"
	case ErrorCodeAdministrativeTimeout:
		return "ERROR_CODE_ADMINISTRATIVE_TIMEOUT"
	}
	return ""
}
