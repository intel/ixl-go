// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package huffman

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestQuickSort(t *testing.T) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for j := 0; j < 100; j++ {
		var source []uint32
		size := rand.Intn(32 * 1024)
		for i := 0; i < size; i++ {
			source = append(source, rand.Uint32())
		}

		arr := make([]uint32, size)
		copy(arr, source)
		quickSort(arr, 0, len(arr)-1)
		if !sort.SliceIsSorted(arr, func(i, j int) bool { return arr[i] < arr[j] }) {
			t.Fatal("unsorted data found")
		}
		copy(arr, source)
		quickSortReverse(arr, 0, len(arr)-1)
		if !sort.SliceIsSorted(arr, func(i, j int) bool { return arr[i] > arr[j] }) {
			t.Fatal("unsorted data found")
		}
	}
	// check  the empty array case
	arr := make([]uint32, 0)
	quickSort(arr, 0, -1)

	quickSort(nil, 0, -1)
}

func TestSortDecLitCounts(t *testing.T) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for j := 0; j < 1000; j++ {
		var testData decLitCounts
		// generate random test data
		size := rand.Intn(1024) + 1
		for i := 0; i < size; i++ {
			testData = append(testData, litCount{count: uint16(rand.Uint32())})
		}
		copied := make(decLitCounts, len(testData))
		copy(copied, testData)
		sortDecLitCounts(copied)

		sort.Sort(testData)
		if !reflect.DeepEqual(testData, copied) {
			t.Log(testData, copied)
			t.Fatal("the sortDecLitCounts result does not match the sort.Sort result")
		}
	}
}
