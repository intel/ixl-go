// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package device

import (
	"strconv"
	"strings"
	"testing"

	"github.com/intel/ixl-go/internal/config"
)

var selectorsTable = []struct {
	selector string
	match    string
	result   bool
}{
	{"3.1", "3.1", true},
	{"3.1", "3.2", false},
	{"3.*", "3.1", true},
	{"3", "3.1", true},
	{"4.1 | (3.* & !(3.1~3.4))", "3.7", true},
	{"4.1 | (3.* & !(3.1~3.4))", "3.2", false},
	{"4.1 | (3.* & !(3.1~3.4))", "4.1", true},
	{"4.1 | (3.* & !(3.1~3.4))", "4.1", true},
	{"(4.2 | 4.3) | (3.* & !(3.1~3.4))", "4.1", false},
	{"4.2,4.3,(3.* & !(3.1~3.4))", "4.1", false},
	{"4.2,4.3,(3.* & !(3.1~3.4))", "4.3", true},
	{"4.2,4.3,(3.* & !(3.1~3.4))", "3.5", true},
	{"!*", "3.5", false},
	{"*", "3.5", true},
}

func TestSelector(t *testing.T) {
	for _, i := range selectorsTable {
		dw := strings.Split(i.match, ".")
		did, _ := strconv.Atoi(dw[0])
		wid, _ := strconv.Atoi(dw[1])

		m, err := getMatcher(i.selector)
		if err != nil {
			t.Fatal(err)
		}
		w := &config.WorkQueue{
			ID:     wid,
			Device: &config.Device{ID: uint64(did)},
		}
		if m.match(w) != i.result {
			t.Fatal(i)
		}
	}
}

func FuzzSelector(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string, wq int, deviceid int) {
		m, err := getMatcher(s)
		if err == nil {
			m.match(&config.WorkQueue{
				ID: wq,
				Device: &config.Device{
					ID: uint64(deviceid),
				},
			})
		}
	})
}
