// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"reflect"
	"testing"
)

func TestExpand(t *testing.T) {
	if !Ready() {
		t.Skip()
	}
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
			name: "expand",
			args: args{
				data: []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				set:  BitSet{0b0110101},
			},
			wantResult: []uint16{1, 0, 2, 0, 3, 4},
			wantErr:    false,
		},
	}
	ctx, _ := NewContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := Expand(ctx, tt.args.data, tt.args.set)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Expand() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
