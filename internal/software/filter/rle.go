// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package filter provides some useful function for test filter.
package filter

import (
	"encoding/binary"
	"errors"
)

type data struct {
	times uint32
	value uint32
}

// GenerateRLE genrerates rle data format from input `arr“ and using specified width.
func GenerateRLE(arr []int, width int) ([]byte, error) {
	if width > 32 {
		return nil, errors.New("invalid width")
	}
	datas := []data{}
	var last int
	var count int
	for i, v := range arr {
		if i == 0 {
			last = v
			count = 1
			continue
		}
		if v == last {
			count++
			continue
		}
		datas = append(datas, data{
			times: uint32(count),
			value: uint32(last),
		})
		last = v
		count = 1
	}
	datas = append(datas, data{
		times: uint32(count),
		value: uint32(last),
	})
	buf := make([]byte, 0, 1024)
	buf = append(buf, byte(width))
	for _, d := range datas {
		// header
		buf = binary.AppendUvarint(buf, uint64(d.times<<1))
		switch {
		case width <= 8:
			buf = append(buf, byte(d.value))
		case width <= 16:
			buf = binary.LittleEndian.AppendUint16(buf, uint16(d.value))
		case width <= 32:
			buf = binary.LittleEndian.AppendUint32(buf, d.value)
		}
	}
	return buf, nil
}

// BitPack pack data as bitpack format using specified width
func BitPack(data []int, width int) ([]byte, error) {
	if width > 32 {
		return nil, errors.New("invalid width")
	}
	var buf []byte
	var curr uint64
	var num int
	mask := uint64(1<<width) - 1
	for _, v := range data {
		curr |= (uint64(v) & mask) << num
		num += width
		if num > 32 {
			// flush
			data := uint32(curr)
			buf = binary.LittleEndian.AppendUint32(buf, data)
			curr >>= 32
			num -= 32
		}
	}
	switch {
	case num <= 8:
		buf = append(buf, byte(curr))
	case num <= 16:
		buf = binary.LittleEndian.AppendUint16(buf, uint16(curr))
	case num <= 24:
		buf = append(buf, byte(curr))
		curr >>= 8
		buf = binary.LittleEndian.AppendUint16(buf, uint16(curr))
	case num <= 32:
		buf = binary.LittleEndian.AppendUint32(buf, uint32(curr))
	}
	return buf, nil
}
