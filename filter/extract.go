// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"runtime"
	"unsafe"

	"github.com/intel/ixl-go/errors"
)

// ExtractRLE extract specified range of data.
func ExtractRLE[T DataUnit](s *Context, size int, data []byte, r Range[uint32]) (result []T, err error) {
	if int(r.Min) >= size || r.Min > r.Max {
		return nil, errors.InvalidArgument
	}
	result = make([]T, r.Max-r.Min+1)
	s.aecs.Reset()
	s.aecs.LowFilterParameter = r.Min
	s.aecs.HighFilterParameter = r.Max
	s.desc.Reset()
	s.cr.Reset()
	extractRLE(s.desc, size, data, result, s.aecs, s.cr)
	status := s.ctx.Submit(uintptr(unsafe.Pointer(s.desc)), &s.cr.Header)
	runtime.KeepAlive(s.aecs)
	runtime.KeepAlive(s.desc)
	runtime.KeepAlive(s.cr)
	runtime.KeepAlive(data)
	runtime.KeepAlive(result)
	if status != 1 {
		return nil, s.cr.CheckError()
	}

	return result, nil
}
