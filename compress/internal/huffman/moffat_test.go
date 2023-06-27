// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package huffman

import (
	"math/rand"
	"testing"
)

func TestMoffat(t *testing.T) {
	for i := 0; i < 10000; i++ {
		mfc := MoffatHuffmanCode{}
		hist := []int32{}
		for i := 0; i < 220; i++ {
			hist = append(hist, rand.Int31n(256))
		}
		mfc.Generate(hist)
		if !validCodes(hist) {
			t.Fail()
		}
	}
}

type number interface {
	uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | int | uint | float32 | float64
}

func max[N number](arr []N) N {
	var max N
	for _, n := range arr {
		if max < n {
			max = n
		}
	}
	return max
}

// sum(2 ^ (maxLength - length[0]) , ... 2 ^ (maxLength - length[n])) == 2 ^ maxLength
func validCodes(lengths []int32) bool {
	maxLength := 0
	for _, v := range lengths {
		if v > int32(maxLength) {
			maxLength = int(v)
		}
	}
	if pow2(maxLength) < 0 {
		panic(maxLength)
	}
	sum := 0
	for _, v := range lengths {
		if v == 0 {
			continue
		}
		sum += pow2(maxLength - int(v))
	}
	return sum == pow2(maxLength)
}

func pow2(num int) int {
	if num == 0 {
		return 1
	}
	return 2 << (num - 1)
}
