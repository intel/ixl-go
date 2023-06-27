// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"reflect"
	"testing"
)

func TestGetIDS(t *testing.T) {
	if nums := getIdsFromName("eg1.1"); !reflect.DeepEqual(nums, []int{1, 1}) {
		t.Fatal(nums)
	}
	if nums := getIdsFromName("eg1"); !reflect.DeepEqual(nums, []int{1, 0}) {
		t.FailNow()
	}
}
