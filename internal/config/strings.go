// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"bytes"
	"os"
	"strconv"
)

func parseSnakeToCamel(filename string) string {
	parts := bytes.Split([]byte(filename), []byte("_"))
	buf := bytes.NewBuffer(make([]byte, 0, len(parts)))
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		if bytes.Equal(p, []byte("id")) {
			_, err := buf.Write([]byte("ID"))
			if err != nil {
				// it can not happen, just for coverity check
				return filename
			}
			continue
		}
		if p[0] >= 'a' && p[0] <= 'z' {
			p[0] = p[0] - byte('a'-'A')
		}
		_, err := buf.Write(p)
		if err != nil {
			// it can not happen, just for coverity check
			return filename
		}
	}
	return buf.String()
}

func readTrimString(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	data = bytes.TrimSpace(data)
	return string(data), nil
}

func readInt(filename string) (num int, err error) {
	str, err := readTrimString(filename)
	if err != nil {
		return
	}
	return strconv.Atoi(str)
}
