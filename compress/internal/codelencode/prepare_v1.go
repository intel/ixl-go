// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64.v1 && !(amd64.v2 || amd64.v3 || amd64.v4)

package codelencode

func init() {
	Prepare = prepare
}
