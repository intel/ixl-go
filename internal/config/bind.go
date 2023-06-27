// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/intel/ixl-go/internal/log"
)

// bindSysfs binding dir content to object
// object must be a pointer to struct.
// currently we only support int/uint/string types fileds.
func bindSysfs(obj interface{}, path string) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr && val.Elem().Kind() != reflect.Struct {
		panic("unsupported")
	}
	val = val.Elem()

	dir, _ := os.ReadDir(path)
	for _, de := range dir {
		if de.IsDir() {
			continue
		}
		fName := parseSnakeToCamel(de.Name())
		ft, ok := val.Type().FieldByName(fName)
		if !ok {
			continue
		}

		if _, ok := ft.Tag.Lookup("binding"); !ok {
			continue
		}

		f := val.FieldByName(fName)

		str, err := readTrimString(filepath.Join(path, de.Name()))
		if err != nil {
			continue
		}

		switch f.Kind() {
		case reflect.Bool:
			switch str {
			case "1", "enabled", "true":
				f.SetBool(true)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			base := 10
			if len(str) > 2 && str[:2] == "0x" {
				base = 16
				str = str[2:]
			}
			num, err := strconv.ParseInt(str, base, 64)
			if err != nil {
				return fmt.Errorf("%s :%w", fName, err)
			}
			f.SetInt(num)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			base := 10
			if len(str) > 2 && str[:2] == "0x" {
				base = 16
				str = str[2:]
			}
			num, err := strconv.ParseUint(str, base, 64)

			f.SetUint(num)
			if err != nil {
				return fmt.Errorf("%s :%w", fName, err)
			}
		case reflect.String:
			f.SetString(str)
		default:
			log.Debug("unsupported type:", f.Kind())
		}
	}
	return nil
}
