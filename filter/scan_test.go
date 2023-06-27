// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/intel/ixl-go/errors"
	"github.com/intel/ixl-go/internal/software/filter"
)

type scanArgs struct {
	input []int
	r     Range[uint32]
}

var scanTests = []struct {
	name       string
	args       scanArgs
	wantOutput BitSet
	wantErr    bool
}{
	{
		name: "in_range",
		args: scanArgs{
			input: []int{1, 2, 3, 4},
			r:     Range[uint32]{1, 3},
		},
		wantOutput: BitSet{0b111},
		wantErr:    false,
	},
	{
		name: "out_range",
		args: scanArgs{
			input: []int{1, 2, 3, 4},
			r:     Range[uint32]{5, 8},
		},
		wantOutput: BitSet{0},
		wantErr:    false,
	},
	{
		name: "in_range_large",
		args: scanArgs{
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			r:     Range[uint32]{3, 12},
		},
		wantOutput: BitSet{0b11111100, 0b1111},
		wantErr:    false,
	},
}

func TestScan(t *testing.T) {
	if !Ready() {
		t.Skip("no IAA device found")
	}

	ctx, err := NewContext()
	if err != nil {
		t.Skip(err)
	}
	for _, tt := range scanTests {
		t.Run(tt.name+"_uint8", func(t *testing.T) {
			gotOutput, err := Scan(ctx, convertArr[uint8](tt.args.input), Range[uint8]{
				Min: uint8(tt.args.r.Min),
				Max: uint8(tt.args.r.Max),
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Scan() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
		t.Run(tt.name+"_uint16", func(t *testing.T) {
			gotOutput, err := Scan(ctx, convertArr[uint16](tt.args.input), Range[uint16]{
				Min: uint16(tt.args.r.Min),
				Max: uint16(tt.args.r.Max),
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Scan() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestContext_ScanBitPacking(t *testing.T) {
	if !Ready() {
		t.Skip("no IAA device found")
	}
	s, _ := NewContext()
	for _, tt := range scanTests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := filter.BitPack(tt.args.input, 5)
			gotOutput, err := s.ScanBitPacking(tt.args.r, 5, len(tt.args.input), data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Context.ScanBitPacking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Context.ScanBitPacking() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestContext_ScanRLE(t *testing.T) {
	if !Ready() {
		t.Skip("no IAA device found")
	}
	s, _ := NewContext()
	for _, tt := range scanTests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := filter.GenerateRLE(tt.args.input, 4)
			gotOutput, err := s.ScanRLE(tt.args.r, len(tt.args.input), data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Context.ScanRLE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Context.ScanRLE() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func convertArr[T DataUnit](arr []int) (list []T) {
	for _, v := range arr {
		list = append(list, T(v))
	}
	return list
}

func FuzzScan(f *testing.F) {
	_, err := NewContext()
	if err != nil {
		f.Skip(err)
	}
	f.Fuzz(func(t *testing.T, data []byte, from, to uint32) {
		ctx, _ := NewContext()
		if len(data) < 4 {
			return
		}
		nSize := len(data) / 4
		arr := unsafe.Slice((*uint32)(unsafe.Pointer(&data[0])), nSize)
		_, err := Scan(ctx, arr, Range[uint32]{from, to})
		if from > to {
			if err == errors.InvalidArgument {
				return
			} else {
				t.Fatalf("excepted invalid argument error, got %v", err)
			}
		} else if err != nil {
			t.Fatalf("error occursed: %v", err)
		}
	})
}
