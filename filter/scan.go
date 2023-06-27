// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"runtime"
	"unsafe"

	"github.com/intel/ixl-go/errors"
	"github.com/intel/ixl-go/util/mem"
)

// Scan scans the input for values within the specified range
func Scan[R DataUnit](s *Context, input []R, r Range[R]) (output BitSet, err error) {
	if r.Min > r.Max {
		return nil, errors.InvalidArgument
	}
	s.aecs.LowFilterParameter = uint32(r.Min)
	s.aecs.HighFilterParameter = uint32(r.Max)
	s.desc.Reset()
	s.cr.Reset()
	output = mem.Alloc64ByteAligned(uintptr(len(input)/8 + 1))
	scanInt(s.desc, input, output, s.aecs, s.cr)
	s.ctx.Submit(uintptr(unsafe.Pointer(s.desc)), &s.cr.Header)
	runtime.KeepAlive(s.aecs)
	runtime.KeepAlive(s.desc)
	runtime.KeepAlive(s.cr)
	runtime.KeepAlive(input)
	runtime.KeepAlive(output)
	cerr := s.cr.CheckError()
	if cerr != nil {
		return nil, cerr
	}
	return output, nil
}

// ScanBitPacking scans the input using bit packing
func (s *Context) ScanBitPacking(r Range[uint32], w int, size int, data []byte) (output BitSet, err error) {
	s.aecs.LowFilterParameter = r.Min
	s.aecs.HighFilterParameter = r.Max
	s.desc.Reset()
	s.cr.Reset()

	output = mem.Alloc64ByteAligned(uintptr(size/8 + 1))
	scanBitPacking(s.desc, uint8(w), size, data, output, s.aecs, s.cr)
	s.ctx.Submit(uintptr(unsafe.Pointer(s.desc)), &s.cr.Header)
	runtime.KeepAlive(s.aecs)
	runtime.KeepAlive(s.desc)
	runtime.KeepAlive(s.cr)
	runtime.KeepAlive(data)
	runtime.KeepAlive(output)
	cerr := s.cr.CheckError()
	if cerr != nil {
		return nil, cerr
	}
	return output, nil
}

// ScanRLE scans the input using run-length encoding
func (s *Context) ScanRLE(r Range[uint32], size int, data []byte) (output BitSet, err error) {
	s.aecs.LowFilterParameter = r.Min
	s.aecs.HighFilterParameter = r.Max
	s.desc.Reset()
	s.cr.Reset()

	output = mem.Alloc64ByteAligned(uintptr(size/8 + 1))
	scanRLE(s.desc, size, data, output, s.aecs, s.cr)
	s.ctx.Submit(uintptr(unsafe.Pointer(s.desc)), &s.cr.Header)
	runtime.KeepAlive(s.aecs)
	runtime.KeepAlive(s.desc)
	runtime.KeepAlive(s.cr)
	runtime.KeepAlive(data)
	runtime.KeepAlive(output)
	cerr := s.cr.CheckError()
	if cerr != nil {
		return nil, cerr
	}
	return output, nil
}
