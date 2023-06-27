// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import "testing"

var parseSnakeToCamlTests = []struct {
	input  string
	output string
}{
	{"happy_x", "HappyX"},
	{"happy_1", "Happy1"},
	{"happy_", "Happy"},
	{"happy_world_Good_boy", "HappyWorldGoodBoy"},
}

func TestParseToCamel(t *testing.T) {
	for _, test := range parseSnakeToCamlTests {
		output := parseSnakeToCamel(test.input)
		if output != test.output {
			t.Fatalf("expected %s got %s", test.output, output)
		}
	}
}
