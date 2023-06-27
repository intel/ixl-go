// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"math/rand"
	"testing"
)

func TestGenerate(t *testing.T) {
	for _, width := range []int{3, 5, 9, 17, 24, 32} {
		data := generateTestData(1024, width)
		_, err := GenerateRLE(data, width)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, size := range []int{1, 3, 5, 10, 17, 24, 32, 69, 128, 256} {
		data := generateTestData(size, 5)
		_, err := BitPack(data, 5)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func generateTestData(size int, width int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		num := rand.Intn(1 << width)
		arr[i] = num
	}
	return arr
}
