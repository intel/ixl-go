// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64.v4

package codelencode

import (
	"github.com/intel/ixl-go/compress/internal/codelencode/avx512"
)

func init() {
	Prepare = avx512.PrepareForCodeLenCode
}
