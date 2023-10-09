// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"io"

	"github.com/intel/ixl-go/util/mem"
)

var _ io.WriteCloser = &BufWriter{}

// BufWriter is a buffer writer for wrapping Deflate/Gzip as a io.Writer.
type BufWriter struct {
	buffer []byte
	offset int
	bw     blockWriter
}

type blockWriter interface {
	writeBlock(block []byte, last bool) (n int, err error)
	Reset(w io.Writer)
	Close() error
}

// NewWriter create a new BufWriter.
// The argument should be Gzip or Deflate.
func NewWriter(bw blockWriter) *BufWriter {
	return &BufWriter{
		buffer: mem.Alloc64ByteAligned(maxBlockSize),
		bw:     bw,
	}
}

// Reset writer.
func (w *BufWriter) Reset(writer io.Writer) {
	w.offset = 0
	w.bw.Reset(writer)
}

// Write data to underlying block writer.
func (w *BufWriter) Write(data []byte) (n int, err error) {
	size := len(data)
CONSUME:
	if w.offset+len(data) < maxBlockSize {
		copy(w.buffer[w.offset:], data)
		w.offset += len(data)
		return size, nil
	}

	copy(w.buffer[w.offset:], data)
	if w.offset+len(data) >= maxBlockSize {
		copiedSize := maxBlockSize - w.offset
		data = data[copiedSize:]
		w.offset = 0
		_, err := w.bw.writeBlock(w.buffer, false)
		if err != nil {
			return 0, err
		}
		goto CONSUME
	}
	return size, nil
}

// Flush immediately write all buffered data to underlying block writer.
func (w *BufWriter) Flush() error {
	_, err := w.bw.writeBlock(w.buffer[:w.offset], false)
	w.offset = 0
	return err
}

// Close flush all buffered data to underlying block writer and close it.
func (w *BufWriter) Close() error {
	_, err := w.bw.writeBlock(w.buffer[:w.offset], true)
	w.offset = 0
	if err != nil {
		uerr := w.bw.Close()
		if uerr != nil {
			return uerr
		}
		return err
	}
	return w.bw.Close()
}

// NewDeflateWriter create a deflate writer
func NewDeflateWriter(w io.Writer, opts ...Option) (*BufWriter, error) {
	d, err := NewDeflate(w, opts...)
	if err != nil {
		return nil, err
	}
	return NewWriter(d), nil
}

// NewGzipWriter create a gzip writer
func NewGzipWriter(w io.Writer, opts ...Option) *BufWriter {
	return NewWriter(NewGzip(w, opts...))
}
