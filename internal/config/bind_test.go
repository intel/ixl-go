// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"os"
	"path"
	"testing"
)

func Test_bindSysfs(t *testing.T) {
	type testObject struct {
		Str  string `binding:""`
		Int  int    `binding:""`
		Uint uint   `binding:""`
		Hex  int    `binding:""`
		Bool bool   `binding:""`
		Skip string
	}
	tempDir := os.TempDir()
	tempDir = path.Join(tempDir, "bindSysfsTests")
	_ = os.MkdirAll(tempDir, 0o777)
	_ = os.WriteFile(path.Join(tempDir, "str"), []byte(`string`), 0o666)
	_ = os.WriteFile(path.Join(tempDir, "int"), []byte(`101`), 0o666)
	_ = os.WriteFile(path.Join(tempDir, "uint"), []byte(`101`), 0o666)
	_ = os.WriteFile(path.Join(tempDir, "hex"), []byte(`0x20`), 0o666)
	_ = os.WriteFile(path.Join(tempDir, "skip"), []byte(`skip`), 0o666)
	_ = os.WriteFile(path.Join(tempDir, "any"), []byte(`any`), 0o666)
	_ = os.WriteFile(path.Join(tempDir, "bool"), []byte(`1`), 0o666)
	defer func() { _ = os.RemoveAll(tempDir) }()
	var obj testObject
	err := bindSysfs(&obj, tempDir)
	if err != nil {
		t.Fatal(err)
	}
	if obj.Hex != 0x20 {
		t.Fatalf("expected obj.Hex == 0x20, got 0x%x", obj.Hex)
	}
	if obj.Str != "string" {
		t.Fatalf("expected obj.Str == 'string', got '%s'", obj.Str)
	}
	if obj.Int != 101 {
		t.Fatalf("expected obj.Int == 101, got %d", obj.Int)
	}
	if obj.Uint != 101 {
		t.Fatalf("expected obj.Uint == 101, got %d", obj.Uint)
	}
	if obj.Skip != "" {
		t.Fatalf("expected obj.Skip is empty, got %s", obj.Skip)
	}
}
