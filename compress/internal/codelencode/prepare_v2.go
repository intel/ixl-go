// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64.v2 && !(amd64.v3 || amd64.v4)

package codelencode

import "github.com/intel/ixl-go/compress/internal/codelencode/sse2"

func init() {
	Prepare = sse2.PrepareForCodeLenCode
}
