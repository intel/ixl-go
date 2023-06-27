// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package errors provides detailed error types for ixl-go, should be used by internal packages.
package errors

import (
	"fmt"
)

// An Error represents a ixl-go error.
type Error interface {
	error
	isIXLGoError()
}

var (
	_ Error = SimpleError("")
	_ Error = HardwareError{}
)

// SimpleError is a simple implementation of Error.
type SimpleError string

func (e SimpleError) Error() string {
	return string(e)
}

func (e SimpleError) isIXLGoError() {}

// HardwareError is a hardware related implementation of Error.
type HardwareError struct {
	Status    string // Status represents job descriptor status
	ErrorCode string // ErrorCode represents job error code.
}

func (h HardwareError) Error() string {
	if h.ErrorCode == "" {
		return fmt.Sprintf("status: %s", h.Status)
	}
	return fmt.Sprintf("status: %s error code: %s", h.Status, h.ErrorCode)
}

func (h HardwareError) isIXLGoError() {}
