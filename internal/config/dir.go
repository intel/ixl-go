// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/intel/ixl-go/internal/log"
)

var devDirName = regexp.MustCompile(`^\S+?(\d+)(\.\d+)?$`)

func getIdsFromName(name string) (numbers []int) {
	match := devDirName.FindStringSubmatch(name)
	if match == nil {
		return numbers
	}
	var n0, n1 int
	n0, _ = strconv.Atoi(match[1])

	if match[2] != "" {
		n1, _ = strconv.Atoi(match[2][1:])
	}

	return []int{n0, n1}
}

func devDirs(path string, prefix string, fn func(ids []int, path string)) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, de := range dir {
		if !strings.HasPrefix(de.Name(), prefix) {
			continue
		}
		path := filepath.Join(path, de.Name())

		isdir, err := isDir(path, de)
		if err != nil {
			log.Debug("warn:", err)
			continue
		}
		if !isdir {
			log.Debug(de.Name(), "is not a dir")
			continue
		}
		if strings.Contains(de.Name(), "!") {
			continue
		}

		ids := getIdsFromName(de.Name())
		if len(ids) == 0 {
			log.Debug(de.Name())
			continue
		}
		fn(ids, path)
	}
	return nil
}

func isDir(path string, e os.DirEntry) (bool, error) {
	if e.IsDir() {
		return true, nil
	}
	mod := e.Type()
	if mod&os.ModeSymlink == os.ModeSymlink {
		rPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return false, err
		}
		fi, err := os.Stat(rPath)
		if err != nil {
			return false, err
		}
		return fi.IsDir(), nil
	}
	return false, nil
}
