// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"io"
	"runtime"
	"unsafe"

	"github.com/intel/ixl-go/errors"
	"github.com/intel/ixl-go/internal/device"
	"github.com/intel/ixl-go/internal/iaa"
	"github.com/intel/ixl-go/util/mem"
)

// Inflate reads data from reader r and decompresses them.
//
// Notice: the compressed data must be compressed by IAA
// or the whole stream must not larger than 4KB.
// This because the standard deflate's history buffer size is 32KB,
// but the IAA deflate's history buffer size is 4KB.
type Inflate struct {
	ctx           *device.Context
	busyPoll      bool
	buffer        []byte
	remnant       int
	desc          iaa.Descriptor
	cr            *iaa.CompletionRecord
	aecsPair      *[2]iaa.DecompressAECS
	toggle        uint8
	state         streamState
	r             io.Reader
	finished      bool
	outputRemnant []byte
}

// NewInflate creates a new Inflate with 4KB buffer size to decompress data from reader r.
func NewInflate(r io.Reader, opts ...Option) (*Inflate, error) {
	inflate, err := NewInflateWithBufferSize(r, 4096, opts...)
	if err != nil {
		return nil, err
	}
	return inflate, nil
}

const minBufferSize = 8

// NewInflateWithBufferSize creates a new Inflate with specified buffer size to decompress data from reader r.
func NewInflateWithBufferSize(r io.Reader, bufferSize int, opts ...Option) (*Inflate, error) {
	if bufferSize < minBufferSize {
		return nil, errors.BufferSizeTooSmall
	}
	opt := &option{}
	for _, f := range opts {
		f(opt)
	}
	i := &Inflate{}
	i.busyPoll = opt.busyPoll
	i.ctx = iaa.LoadContext()
	if i.ctx == nil {
		return nil, errors.NoHardwareDeviceDetected
	}
	i.cr = mem.Alloc64Align[iaa.CompletionRecord]()
	i.aecsPair = mem.Alloc64Align[[2]iaa.DecompressAECS]()
	i.r = r
	i.buffer = mem.Alloc64ByteAligned(uintptr(bufferSize))
	return i, nil
}

// Reset reset the Inflate object
func (i *Inflate) Reset(r io.Reader) {
	i.remnant = 0
	i.toggle = 0
	i.state = first
	i.r = r
	i.finished = false
	i.outputRemnant = i.outputRemnant[:0]
}

func (i *Inflate) submit() iaa.StatusCode {
	if i.busyPoll {
		return iaa.StatusCode(i.ctx.SubmitBusyPoll(uintptr(unsafe.Pointer(&i.desc)), &i.cr.Header))
	} else {
		return iaa.StatusCode(i.ctx.Submit(uintptr(unsafe.Pointer(&i.desc)), &i.cr.Header))
	}
}

// DecompressAll decompress all compressed data and write result into raw.
// The caller should make sure that `raw` has enough space.
func (i *Inflate) DecompressAll(compressed []byte, raw []byte) (int, error) {
	if len(compressed) > int(i.ctx.MaxTransferSize()) || len(raw) > int(i.ctx.MaxTransferSize()) {
		return 0, errors.DataSizeTooLarge
	}
	i.decompressJob(compressed, raw, &i.aecsPair[0])
	status := i.submit()
	if status != iaa.Success {
		return 0, i.cr.CheckError()
	}
	runtime.KeepAlive(compressed)
	runtime.KeepAlive(i.buffer)
	return int(i.cr.OutputSize), nil
}

// Read decompressed data from the underlying compressed reader.
func (i *Inflate) Read(data []byte) (n int, err error) {
	if len(i.outputRemnant) != 0 {
		n = copy(data, i.outputRemnant)
		if n == len(i.outputRemnant) {
			i.outputRemnant = i.outputRemnant[:0]
		} else {
			copy(i.outputRemnant, i.outputRemnant[n:])
			i.outputRemnant = i.outputRemnant[:len(i.outputRemnant)-n]
		}
		return n, nil
	}

	if i.finished {
		return 0, io.EOF
	}
	if i.remnant == 0 && i.state != last {
		i.remnant, err = i.r.Read(i.buffer)
		if err == io.EOF {
			i.state = last
		} else if err != nil {
			return 0, err
		}
	}

	if i.remnant == 0 && i.state == last && len(data) == 0 {
		i.finished = true
		return 0, io.EOF
	}

	if len(data) == 0 {
		return 0, nil
	}

	input := i.buffer[:i.remnant]
	i.cr.Reset()
	i.desc.Reset()
	if len(data) > int(i.ctx.MaxTransferSize()) {
		data = data[:i.ctx.MaxTransferSize()]
	}
	i.decompressJob(input, data, &i.aecsPair[0])
	i.submit()
	status := i.cr.GetHeader().StatusCode
RETRY:
	switch status {
	case iaa.Success:
		i.remnant = 0
		if i.state == last {
			i.finished = true
		}
	case iaa.OutputBufferOverflow:
		i.remnant = len(input) - int(i.cr.Header.BytesCompleted)
		temp := append(make([]byte, 0, i.remnant), input[i.cr.Header.BytesCompleted:]...)
		copy(i.buffer, temp)
		if i.cr.OutputSize == 0 {
			// TODO: we should fix this problem!
			if i.outputRemnant == nil {
				i.outputRemnant = make([]byte, 258)
			} else {
				i.outputRemnant = i.outputRemnant[:258]
			}
			i.decompressJob(i.buffer[:i.remnant], i.outputRemnant, &i.aecsPair[0])
			i.submit()
			size := copy(data, i.outputRemnant[:i.cr.OutputSize])
			copy(i.outputRemnant, i.outputRemnant[size:])
			i.outputRemnant = i.outputRemnant[:int(i.cr.OutputSize)-size]
			i.cr.OutputSize = uint32(size)
			status = i.cr.GetHeader().StatusCode
			goto RETRY
		}
	default:
		return 0, i.cr.CheckError()
	}
	runtime.KeepAlive(input)
	runtime.KeepAlive(data)
	runtime.KeepAlive(i.cr)
	runtime.KeepAlive(i.aecsPair)
	outsize := i.cr.OutputSize
	i.toggle ^= 1
	if i.state == first {
		i.state = middle
	}

	return int(outsize), nil
}

type streamState uint8

const (
	first  streamState = 0
	middle streamState = 1
	last   streamState = 2
)

var emptyBlock = make([]byte, 1)

func (i *Inflate) decompressJob(block []byte, output []byte, aesc *iaa.DecompressAECS) {
	d := &i.desc
	d.SetOpcode(iaa.OpDecompress)
	if i.toggle == 1 {
		d.SetFlag(iaa.FlagAecsRWToggleSelector)
	}
	d.SetFlag(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagCompletionRecordValid |
			iaa.FlagRequestCompletionRecord,
	)
	switch i.state {
	case middle:
		d.SetFlag(iaa.FlagReadSource2Aecs)
		d.SetFlag(iaa.FlagWriteSource2CompletionOfOperation)
	case first:
		d.SetFlag(iaa.FlagWriteSource2CompletionOfOperation)
	case last:
		d.SetFlag(iaa.FlagReadSource2Aecs)
		d.SetFlag(iaa.FlagWriteSource2OnlyIfOutputOverflow)
	}

	d.SetDecompressionFlag(
		iaa.DecompressionFlagEnableDecompression |
			iaa.DecompressionFlagStopOnEOB |
			iaa.DecompressionFlagFlushOutput |
			iaa.DecompressionFlagSelectBFinalEOB,
	)

	if len(block) == 0 {
		d.Src1Addr = uintptr(unsafe.Pointer(&emptyBlock[0]))
		d.Size = 0
	} else {
		d.Src1Addr = uintptr(unsafe.Pointer(&block[0]))
		d.Size = uint32(len(block))
	}
	if len(output) == 0 {
		d.DestAddr = uintptr(unsafe.Pointer(&emptyBlock[0]))
	} else {
		d.DestAddr = uintptr(unsafe.Pointer(&output[0]))
	}
	d.MaxDestionationSize = uint32(len(output))

	d.Src2Addr = uintptr(unsafe.Pointer(aesc))
	d.Src2Size = uint32(unsafe.Sizeof(iaa.DecompressAECS{}))

	d.SetCompleteRecord(uintptr(unsafe.Pointer(i.cr)))
}
