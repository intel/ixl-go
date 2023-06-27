// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"bytes"
	"compress/flate"
	"io"
	"testing"
)

type writerTestsOutput struct {
	data int
	last bool
}
type writerTestsInputAction int

const (
	testWriteAction writerTestsInputAction = 0
	testFlushAction writerTestsInputAction = 1
)

type writerTestsInput struct {
	Action writerTestsInputAction
	Value  int
}

var writerTests = []struct {
	inputs  []writerTestsInput
	outputs []writerTestsOutput
}{
	{
		inputs: []writerTestsInput{
			{Value: 32 * 1024},
		},
		outputs: []writerTestsOutput{
			{32 * 1024, false},
			{0, true},
		},
	},
	{
		inputs: []writerTestsInput{{Value: 16 * 1024}},
		outputs: []writerTestsOutput{
			{16 * 1024, true},
		},
	},
	{
		inputs: []writerTestsInput{{Value: 16 * 1024}, {Value: 2 * 1024}},
		outputs: []writerTestsOutput{
			{18 * 1024, true},
		},
	},
	{
		inputs: []writerTestsInput{
			{Value: 16 * 1024},
			{Value: 16 * 1024},
		},
		outputs: []writerTestsOutput{
			{32 * 1024, false},
			{0, true},
		},
	},
	{
		inputs: []writerTestsInput{
			{Value: 16 * 1024},
			{Action: testFlushAction},
			{Value: 16 * 1024},
		},
		outputs: []writerTestsOutput{
			{16 * 1024, false},
			{16 * 1024, true},
		},
	},
	{
		inputs: []writerTestsInput{
			{Value: 64 * 1024},
		},
		outputs: []writerTestsOutput{
			{32 * 1024, false},
			{32 * 1024, false},
			{0, true},
		},
	},
	{
		inputs: []writerTestsInput{
			{Value: 64*1024 + 2},
		},
		outputs: []writerTestsOutput{
			{32 * 1024, false},
			{32 * 1024, false},
			{2, true},
		},
	},
}

func TestWriter_Write(t *testing.T) {
	mbw := &mockBlockWriter{}
	w := NewWriter(mbw)
	for _, test := range writerTests {
		w.Reset(nil)
		for _, v := range test.inputs {
			switch v.Action {
			case testWriteAction:
				data := make([]byte, v.Value)
				_, _ = w.Write(data)
			case testFlushAction:
				_ = w.Flush()
			}
		}
		w.Close()
		if len(test.outputs) != len(mbw.records) {
			t.Fatal(test.outputs, mbw.toWriterTestsOutput())
		}

		for i, mbwr := range mbw.records {
			if test.outputs[i].data != len(mbwr.data) {
				t.Fatal(test.outputs[i].data, len(mbwr.data))
			}
		}
	}
}

func FuzzWriterWrite(f *testing.F) {
	if !Ready() {
		f.Skip("no IAA device detected")
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		buf := bytes.NewBuffer(nil)
		d, _ := NewDeflate(buf)
		w := NewWriter(d)
		_, err := w.Write(data)
		if err != nil {
			t.Fatal(err)
		}
		err = w.Close()
		if err != nil {
			t.Fatal(err)
		}
		r := flate.NewReader(buf)

		out, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(out, data) {
			t.Fatal("expected equal")
		}
	})
}

type mockBlockWriterRecord struct {
	data []byte
	last bool
}

type mockBlockWriter struct {
	records []mockBlockWriterRecord
}

func (b *mockBlockWriter) toWriterTestsOutput() (outputs []writerTestsOutput) {
	for _, mbwr := range b.records {
		outputs = append(outputs, writerTestsOutput{
			data: len(mbwr.data),
			last: mbwr.last,
		})
	}
	return outputs
}

func (b *mockBlockWriter) writeBlock(block []byte, last bool) (n int, err error) {
	b.records = append(b.records, mockBlockWriterRecord{data: block, last: last})
	return len(block), nil
}

func (b *mockBlockWriter) Reset(w io.Writer) {
	b.records = nil
}

func (b *mockBlockWriter) Close() error {
	return nil
}
