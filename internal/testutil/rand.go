// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package testutil provides some useful functions for test.
package testutil

import "crypto/rand"

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
