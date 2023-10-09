// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import "github.com/intel/ixl-go/internal/iaa"

var fixedHistogram = generateFixedHistogram()

// genertae fixed histogram follow: https://www.rfc-editor.org/rfc/rfc1951#section-3.2.6
func generateFixedHistogram() (his iaa.Histogram) {
	for i := range his.LiteralCodes {
		i := uint16(i)
		var bits uint16
		var size uint16
		switch {
		case i < 144:
			bits = i + 48
			size = 8
		case i < 256:
			bits = i + 400 - 144
			size = 9
		case i < 280:
			bits = i - 256
			size = 7
		default:
			bits = i + 192 - 280
			size = 8
		}
		his.LiteralCodes[i] = int32(bits) | int32(size)<<15
	}
	for i := range his.DistanceCodes {
		his.DistanceCodes[i] = int32(i) | 5<<15
	}
	return
}
