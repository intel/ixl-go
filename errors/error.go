// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package errors provides detailed error types for ixl-go.
package errors

import (
	"github.com/intel/ixl-go/internal/errors"
)

// An Error represents a ixl-go error.
type Error = errors.Error

var (
	// DataSizeTooLarge represents that data size is large than device's max_transfer_size
	DataSizeTooLarge error = errors.SimpleError("data size is large than device's max_transfer_size")
	// InvalidArgument represents invalid arguments error.
	InvalidArgument error = errors.SimpleError("invalid argument")
	// NoHardwareDeviceDetected represents no device found.
	NoHardwareDeviceDetected error = errors.SimpleError("no hardware device detected")
	// BufferSizeTooSmall represents that buffer size is too small.
	BufferSizeTooSmall error = errors.SimpleError("buffer size too small")
)

var (

	// ErrNonLatin1Header means the header string should be Latin-1 encoded.
	ErrNonLatin1Header = errors.SimpleError("gzip: non-Latin-1 header string")
	// ErrZeroByte means the header string must not contains any zero byte.
	ErrZeroByte = errors.SimpleError("gzip: header string contains zero byte")
)
