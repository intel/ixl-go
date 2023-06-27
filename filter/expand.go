// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"runtime"
	"unsafe"

	"github.com/intel/ixl-go/errors"
)

// Expand generates an array in which the elements in the data are placed according to 1 bits in the set.
func Expand[T DataUnit](s *Context, data []T, set BitSet) (result []T, err error) {
	s.desc.Reset()
	s.cr.Reset()
	count := set.OnesCount()
	if count > len(data) {
		return nil, errors.InvalidArgument
	} else if count < len(data) {
		data = data[:count]
	}

	result = make([]T, set.Size())
	expandUints(s.desc, data, set, result, s.cr)
	status := s.ctx.Submit(uintptr(unsafe.Pointer(s.desc)), &s.cr.Header)
	runtime.KeepAlive(data)
	runtime.KeepAlive(s.desc)
	runtime.KeepAlive(s.cr)
	runtime.KeepAlive(result)
	runtime.KeepAlive(set)
	if status != 1 {
		return nil, s.cr.CheckError()
	}

	return result, nil
}
