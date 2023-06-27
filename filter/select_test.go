// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"reflect"
	"testing"
)

func TestSelect(t *testing.T) {
	type args struct {
		data []uint16
		set  BitSet
	}
	tests := []struct {
		name       string
		args       args
		wantResult []uint16
		wantErr    bool
	}{
		{
			name: "select",
			args: args{
				data: []uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
				set:  BitSet{0b00100},
			},
			wantResult: []uint16{2},
			wantErr:    false,
		},
		{
			name: "select",
			args: args{
				data: []uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
				set:  BitSet{0b00100},
			},
			wantResult: []uint16{2},
			wantErr:    false,
		},
	}
	ctx, _ := NewContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := Select(ctx, tt.args.data, tt.args.set)
			if (err != nil) != tt.wantErr {
				t.Errorf("Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Select() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
