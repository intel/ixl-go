// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package huffman

import "math/bits"

// GenerateCode generate code in reversed format.
// This method can be used for deflate header written
func GenerateCode(reuseCache []int32, maxLen int, lens []int32, rcodes []uint16) {
	blCount := reuseCache[:maxLen+1]
	for i := range blCount {
		blCount[i] = 0
	}
	maxBits := 0
	for _, v := range lens {
		blCount[v]++
		if v > int32(maxBits) {
			maxBits = int(v)
		}
	}

	blCount[0] = 0

	nextCodes := reuseCache[maxLen+1 : maxLen+1+maxBits+1]
	code := int32(0)
	for bits := 1; bits <= maxBits; bits++ {
		code = (code + blCount[bits-1]) << 1
		nextCodes[bits] = code
	}
	for i, l := range lens {
		if l != 0 {
			code := nextCodes[l]
			rcodes[i] = bits.Reverse16(uint16(code)) >> (16 - lens[i])
			nextCodes[l]++
		}
	}
}

// GenerateCodeForIAA generate code in non-reversed format.
// Result will be inplace.
func GenerateCodeForIAA(lens []int32, reuseCache []int32) {
	blCount := reuseCache[:15+1]
	for i := range blCount {
		blCount[i] = 0
	}
	maxBits := 0
	for _, v := range lens {
		blCount[v]++
		if v > int32(maxBits) {
			maxBits = int(v)
		}
	}

	blCount[0] = 0
	nextCodes := reuseCache[15+1 : 15+1+maxBits+1]
	for i := range nextCodes {
		nextCodes[i] = 0
	}

	code := int32(0)
	for bits := 1; bits <= maxBits; bits++ {
		code = (code + blCount[bits-1]) << 1
		nextCodes[bits] = code
	}
	for i, l := range lens {
		if l != 0 {
			code := nextCodes[l]
			lens[i] = (l << 15) | code
			nextCodes[l]++
		}
	}
}
