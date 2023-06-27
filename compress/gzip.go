// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"compress/flate"
	"compress/gzip"
	"encoding/binary"
	"io"

	"github.com/intel/ixl-go/errors"
)

// Header is same as gzip.Header
type Header = gzip.Header

// Gzip is an object to hold the state for compress data using gzip format.
type Gzip struct {
	Header
	UTF8        bool // can be used with gzip command
	wroteHeader bool
	level       int
	w           io.Writer
	buf         []byte
	sum         int64
	compressor  *Deflate
}

var fixedGzipHeader = [3]byte{0x1f, 0x8b, 8}

// NewGzip create a new Gzip.
func NewGzip(w io.Writer) *Gzip {
	g := &Gzip{Header: Header{OS: 255}, w: w}
	g.w = w
	return g
}

// Gzip format: https://www.rfc-editor.org/rfc/rfc1952#page-4
func (g *Gzip) writeHeader() (err error) {
	flag := gzipFlag(0)
	bCap := 10 // |ID1|ID2|CM|FLG|MTIME(4 byte)|XFL|OS|
	if len(g.Extra) != 0 {
		bCap += 2 + len(g.Extra) // xlen + bytes
		flag |= gzipFileExtra
	}
	if len(g.Name) != 0 {
		bCap += len(g.Name) + 1 // zero-terminated
		flag |= gzipFileName
	}
	if len(g.Comment) != 0 {
		bCap += len(g.Comment) + 1 // zero-terminated
		flag |= gzipFileComment
	}
	if cap(g.buf) > bCap {
		g.buf = g.buf[:bCap]
	} else {
		g.buf = make([]byte, bCap)
	}
	// ID1|ID2|CM
	copy(g.buf[:3], fixedGzipHeader[:])
	g.buf[3] = uint8(flag)
	// MTIME
	sec := g.ModTime.Unix()
	if sec < 0 {
		sec = 0
	}
	binary.LittleEndian.PutUint32(g.buf[4:8], uint32(sec))

	level := 0
	if g.level == flate.BestCompression {
		level = 2
	} else if g.level == flate.BestSpeed {
		level = 4
	}
	g.buf[8] = byte(level)
	g.buf[9] = g.OS
	idx := 10
	if len(g.Extra) != 0 {
		binary.LittleEndian.PutUint16(g.buf[idx:idx+2], uint16(len(g.Extra)))
		idx += 2
		copy(g.buf[idx:], g.Extra)
		idx += len(g.Extra)
	}
	if len(g.Name) != 0 {
		idx, err = g.writeHeaderStr(idx, g.Name)
		if err != nil {
			return err
		}
	}
	if len(g.Comment) != 0 {
		idx, err = g.writeHeaderStr(idx, g.Comment)
		if err != nil {
			return err
		}
	}
	_, err = g.w.Write(g.buf[:idx])
	return err
}

func (g *Gzip) writeHeaderStr(idx int, str string) (nidx int, err error) {
	if g.UTF8 {
		for _, char := range str {
			if char == 0 {
				return idx, errors.ErrZeroByte
			}
		}
		copy(g.buf[idx:], []byte(g.Name))
		idx += len(g.Name)
		g.buf[idx] = 0
		idx++
		return idx, nil
	}
	safe := true
	for _, char := range str {
		if char == 0 || char > 0xff {
			return idx, errors.ErrNonLatin1Header
		}
		if char > 0x7f {
			safe = false
			break
		}
	}
	if safe {
		copy(g.buf[idx:], []byte(g.Name))
		idx += len(g.Name)
		g.buf[idx] = 0
		idx++
		return idx, nil
	}
	for _, char := range str {
		g.buf[idx] = byte(char)
		idx++
	}
	g.buf[idx] = 0
	idx++
	return idx, nil
}

// writeBlock compresses the block and writes it to underlying writer.
//
// Notice:
//  1. The block first byte address must be aligned to a multiple of 64 bytes.
//     You can use `mem.Alloc64ByteAligned` function to alloc a 64 bytes aligned bytes.
//  2. The `last` argument must be true if the block is the last block in the stream.
//  3. For most scenarios, you should use the `ReadFrom` method.
func (g *Gzip) writeBlock(block []byte, last bool) (n int, err error) {
	if !g.wroteHeader {
		err = g.writeHeader()
		if err != nil {
			return 0, err
		}
		g.wroteHeader = true
	}
	if g.compressor == nil {
		g.compressor, err = NewDeflate(g.w)
		if err != nil {
			return 0, err
		}
	}
	g.sum += int64(len(block))
	n, err = g.compressor.writeBlock(block, last)
	if err != nil {
		return n, err
	}
	if last {
		err = g.writeTailer(g.sum)
		if err != nil {
			return 0, err
		}
	}
	return
}

// ReadFrom reads all data from `r` and compresses the data and then writes compressed data into underlying writer `w`.
func (g *Gzip) ReadFrom(reader io.Reader) (n int64, err error) {
	if !g.wroteHeader {
		err = g.writeHeader()
		if err != nil {
			return 0, err
		}
	}
	if g.compressor == nil {
		g.compressor, err = NewDeflate(g.w)
		if err != nil {
			return 0, err
		}
	}
	n, err = g.compressor.ReadFrom(reader)
	if err != nil && err != io.EOF {
		return n, err
	}
	err = g.writeTailer(n)
	return n, err
}

func (g *Gzip) writeTailer(n int64) error {
	crc := g.compressor.crc
	buf := g.buf[:8]
	binary.LittleEndian.PutUint32(buf[:4], crc)
	binary.LittleEndian.PutUint32(buf[4:], uint32(n))
	_, err := g.w.Write(buf)
	return err
}

// Reset the internal states for reusing the object.
func (g *Gzip) Reset(w io.Writer) {
	g.w = w
	g.sum = 0
	g.wroteHeader = false
	g.Header = Header{OS: 255}
	if g.compressor != nil {
		g.compressor.Reset(g.w)
	}
}

// Close the writer.
func (g *Gzip) Close() error {
	closer, ok := g.w.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}

type gzipFlag uint8

const (
	gzipFileText gzipFlag = 1 << iota
	gzipFileHCRC
	gzipFileExtra
	gzipFileName
	gzipFileComment
)
