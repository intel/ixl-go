// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"reflect"
	"testing"
)

func TestParseGenCap(t *testing.T) {
	gc := &GenCap{
		BlockOnFault:                   true,
		Overlapping:                    true,
		CacheControlMemory:             true,
		CacheControlCacheFlush:         false,
		CommandCapabilities:            false,
		DestinationReadback:            true,
		DrainDescriptorReadbackAddress: false,
		ConfigurationSupport:           true,
		MaxTransferSize:                0x80000000,
		MaxBatchSize:                   0x400,
		InterruptMessageStorageSize:    0x2f,
	}
	g := ParseGenCap(0x40915f010f)
	if !reflect.DeepEqual(g, gc) {
		t.Fatal("expected parse gen cap correctly")
	}
}

func TestCheckAllFunctionsNoPanic(t *testing.T) {
	ctx := NewContext()
	ctx.Engines(IAA)
	ctx.Engines(DSA)

	for _, d := range ctx.Devices {
		_, _ = d.Clients()
		_ = d.ReadOPCap()
	}
	checkWQs(ctx.DedicatedWQs(IAA))
	checkWQs(ctx.SharedWQs(IAA))
	checkWQs(ctx.DedicatedWQs(DSA))
	checkWQs(ctx.SharedWQs(DSA))
}

func checkWQs(wqs []*WorkQueue) {
	for _, wq := range wqs {
		_, _ = wq.Clients()
		_ = wq.DevicePath()
	}
}
