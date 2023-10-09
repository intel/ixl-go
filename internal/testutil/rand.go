// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package testutil provides some useful functions for test.
package testutil

import (
	"bytes"
	"crypto/rand"
	"math"
	"math/big"
)

var chars []rune

func init() {
	for i := 'a'; i <= 'z'; i++ {
		chars = append(chars, i)
	}
	for i := 'A'; i <= 'Z'; i++ {
		chars = append(chars, i)
	}
	for i := '0'; i <= '9'; i++ {
		chars = append(chars, i)
	}
	chars = append(chars, ',', '.', '[', ']', '(', ')', '=', '-', '+', '_', '\\', '/')
}

// RandomText generates random text with specified size.
func RandomText(size int) string {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	result := make([]rune, size)
	for i := range data {
		result[i] = chars[int(data[i])%len(chars)]
	}
	return string(result)
}

// RandomByRatio generates random byte with specified size and ratio.
func RandomByRatio(size int, ratio float64) []byte {
	if ratio < 1 {
		ratio = 1
	}
	base := int(float64(size) / ratio)
	data := make([]byte, base)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	var arr [][]byte
	total := 0
	repeat := math.Floor(ratio)
	rsize := int(float64(base) * (ratio - repeat))

	for total < base {
		n, _ := rand.Int(rand.Reader, big.NewInt(128))
		num := int(n.Int64() + 3)
		if total+num > base {
			num = base - total
		}
		p := data[total : total+num]
		for i := 0; i < int(repeat); i++ {
			arr = append(arr, p)
		}
		if rsize > 0 {
			arr = append(arr, p)
			rsize -= num
		}
		total += num
	}
	for i := range arr {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(i)+1))
		j := n.Int64()
		arr[i], arr[j] = arr[j], arr[i]
	}
	result := bytes.Join(arr, nil)
	if len(result) < size {
		data := make([]byte, size-len(result))
		_, _ = rand.Read(data)
		result = append(result, data...)
	} else {
		result = result[:size]
	}
	return result
}
