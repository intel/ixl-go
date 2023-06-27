// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package codelencode

import (
	"math/rand"

	"github.com/intel/ixl-go/internal/iaa"
)

func randHistogram() iaa.Histogram {
	h := iaa.Histogram{}
	for i := 0; i < 230; i++ {
		h.LiteralCodes[i] = (rand.Int31n(3) + 1) << 15
	}
	for i := 0; i < 30; i++ {
		if i < 10 {
			h.DistanceCodes[i] = (rand.Int31n(15) + 1) << 15
		} else {
			h.DistanceCodes[i] = 0
		}
	}
	return h
}
