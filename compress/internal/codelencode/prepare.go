// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package codelencode provides a function to prepare the code len code.
package codelencode

import (
	"github.com/intel/ixl-go/internal/iaa"
)

// Prepare for code len code
var Prepare func(histogram *iaa.Histogram, source []byte) (litNum uint16, distanceNum uint16)

func prepare(histogram *iaa.Histogram, source []byte) (litNum uint16, distanceNum uint16) {
	litNum = 0
	for i := 285; i >= 0; i-- {
		if histogram.LiteralCodes[i] != 0 {
			litNum = uint16(i) + 1
			break
		}
	}
	for i := 29; i >= 0; i-- {
		if histogram.DistanceCodes[i] != 0 {
			distanceNum = uint16(i) + 1
			break
		}
	}
	insertOneDistance := false
	if distanceNum == 0 {
		distanceNum = 1
		insertOneDistance = true
	}

	// prepare to combine repeat numbers
	source = source[:litNum+distanceNum+1]
	for i := uint16(0); i < litNum; i++ {
		source[i] = uint8(histogram.LiteralCodes[i] >> 15)
	}
	for i := uint16(0); i < distanceNum; i++ {
		source[litNum+i] = uint8(histogram.DistanceCodes[i] >> 15)
	}
	if insertOneDistance {
		source[litNum] = 1
	}
	return litNum, distanceNum
}
