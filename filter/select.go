// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"runtime"
	"unsafe"
)

// Select selects the elements in the data whose indices correspond to 1-bits in the set
func Select[T DataUnit](s *Context, data []T, set BitSet) (result []T, err error) {
	s.desc.Reset()
	s.cr.Reset()

	result = make([]T, set.OnesCount())
	bitsetSize := len(data)/8 + 1
	if len(data)%8 == 0 {
		bitsetSize++
	}
	for len(set) < bitsetSize {
		set = append(set, 0)
	}
	selectUints(s.desc, data, set, result, s.cr)
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
