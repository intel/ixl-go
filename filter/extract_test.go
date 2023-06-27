// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"reflect"
	"testing"

	"github.com/intel/ixl-go/internal/software/filter"
)

func mustData(data []byte, _ error) []byte {
	return data
}

func TestExtractRLE(t *testing.T) {
	if !Ready() {
		t.Skip()
	}
	type args struct {
		size int
		data []byte
		r    Range[uint32]
	}
	tests := []struct {
		name       string
		args       args
		wantResult []uint16
		wantErr    bool
	}{
		{
			name: "extract rle",
			args: args{
				data: mustData(filter.GenerateRLE([]int{1, 2, 3, 4}, 3)),
				size: 4,
				r:    Range[uint32]{1, 2},
			},
			wantResult: []uint16{2, 3},
			wantErr:    false,
		},
		{
			name: "extract all",
			args: args{
				data: mustData(filter.GenerateRLE([]int{1, 2, 3, 4}, 3)),
				size: 4,
				r:    Range[uint32]{0, 3},
			},
			wantResult: []uint16{1, 2, 3, 4},
			wantErr:    false,
		},
	}
	ctx, _ := NewContext()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := ExtractRLE[uint16](ctx, tt.args.size, tt.args.data, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractRLE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("ExtractRLE() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
