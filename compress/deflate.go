// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"encoding/binary"
	"io"
	"runtime"
	"unsafe"

	"github.com/intel/ixl-go/compress/internal/huffman"
	"github.com/intel/ixl-go/errors"
	"github.com/intel/ixl-go/internal/device"
	"github.com/intel/ixl-go/internal/iaa"
	"github.com/intel/ixl-go/util/mem"
)

// Deflate takes data written to it and writes the deflate compressed
// form of that data to an underlying writer (see NewDeflate).
//
// Notice:
//
//  1. the history buffer used by hardware is 4KB.
//  2. the `Deflate` object should be reused as much as possible to reduce the GC overhead.
type Deflate struct {
	w   io.Writer
	ctx *device.Context

	frame  headerFrame
	output []byte

	readcache []byte

	cacheForGencode []int32
	litGen          huffman.TreeGenerator
	offsetGen       huffman.TreeGenerator
	codeGen         huffman.TreeGenerator

	dynHeader *dynamicHeader

	descriptor       iaa.Descriptor
	completionRecord *iaa.CompletionRecord
	aecs             *compressAECSPair
	toggle           uint8
	bits             uint8
	bitsNum          uint8
	crc              uint32
}

type iaaCachedObject struct {
	iaa.CompletionRecord
	compressAECSPair
}

type compressAECSPair [2]iaa.CompressAECS

// NewDeflate returns a new Deflate writing compressed data to underlying writer `w`.
func NewDeflate(w io.Writer) (*Deflate, error) {
	ctx := iaa.LoadContext()

	if ctx == nil {
		// no device found
		return nil, errors.NoHardwareDeviceDetected
	}
	deflate := &Deflate{
		ctx:             ctx,
		dynHeader:       newDynamicHeader(),
		w:               w,
		frame:           headerFrame{},
		litGen:          huffman.NewLenLimitedCode(),
		offsetGen:       huffman.NewLenLimitedCode(),
		codeGen:         huffman.NewLenLimitedCode(),
		descriptor:      iaa.Descriptor{},
		cacheForGencode: make([]int32, 16*2),
		// size(block) + storedBlockHeaderSize + lastBlockBits
		output: mem.Alloc64ByteAligned(maxBlockSize + 5 + 1),
	}
	ico := mem.Alloc64Align[iaaCachedObject]()
	deflate.completionRecord = &ico.CompletionRecord
	deflate.aecs = &ico.compressAECSPair
	return deflate, nil
}

// Reset the `Deflate` object.
func (d *Deflate) Reset(w io.Writer) {
	d.crc = 0
	d.bits = 0
	d.bitsNum = 0
	d.w = w
}

// maxBlockSize is max deflate block size.
const maxBlockSize = 32 * 1024

// ReadFrom reads all data from `r` and compresses the data and then writes compressed data into underlying writer `w`.
func (d *Deflate) ReadFrom(r io.Reader) (total int64, err error) {
	if d.readcache == nil {
		d.readcache = mem.Alloc64ByteAligned(maxBlockSize * 2)
	}

	current, prev := d.readcache[:maxBlockSize], d.readcache[maxBlockSize:]

	prevSize := 0
	for {
		size, err := r.Read(current)
		if err == nil && size == 0 {
			continue
		}
		total = int64(size) + total
		if err != nil && err != io.EOF {
			return total, err
		}
		if err == io.EOF {
			// deal with EOF
			if size == 0 {
				_, err = d.writeBlock(prev[:prevSize], true)
				if err != nil {
					return total, err
				}
				return total, io.EOF
			}
			_, err = d.writeBlock(prev[:prevSize], false)
			if err != nil {
				return total, err
			}
			_, err = d.writeBlock(current[:size], true)
			if err != nil {
				return total, err
			}
			return total, io.EOF
		}
		// check if it's the first block
		if prevSize == 0 {
			prevSize = size
			prev, current = current, prev
			continue
		}
		// flush prev block
		_, err = d.writeBlock(prev[:prevSize], false)
		if err != nil {
			return total, err
		}
		// exchange current block with prev block
		prevSize = size
		prev, current = current, prev
		continue
	}
}

// writeBlock write one block into compression stream,
// and the compressed data will be written to the underlying `w`.
//
// Notice:
//  1. The block first byte address must be aligned to a multiple of 64 bytes.
//     You can use `mem.Alloc64ByteAligned` function to alloc a 64 bytes aligned bytes.
//  2. The `last` argument must be true if the block is the last block in the stream.
//  3. For most scenarios, you should use the `ReadFrom` method.
func (d *Deflate) writeBlock(block []byte, last bool) (n int, err error) {
	if len(block) == 0 {
		err = d.writeStoredBlock(block, last)
		return 0, err
	}

	aecs := &d.aecs[d.toggle]
	aecs.Reset()
	histogram := &d.aecs[d.toggle].Histogram
	d.descriptor.Reset()
	d.completionRecord.Reset()

	// statistic the block
	err = d.statisticBlock(block, histogram)
	if err != nil {
		return 0, err
	}
	// generate the block header into accumulator data
	headerBits := d.generateHeader(histogram, &aecs.OutputAccumulatorData, last)

	// encode the block using huffman code
	err = d.encodeBlock(aecs, block, last, headerBits)
	if err != nil {
		return 0, err
	}

	d.toggle ^= 1
	return len(block), err
}

func (d *Deflate) statisticBlock(block []byte, histogram *iaa.Histogram) error {
	d.descriptor.Reset()
	d.completionRecord.Reset()
	d.statsJob(block, histogram)
	ptr := (unsafe.Pointer(&d.descriptor))
	status := iaa.StatusCode(d.ctx.Submit(uintptr(ptr), &d.completionRecord.Header))
	runtime.KeepAlive(histogram)
	runtime.KeepAlive(d.completionRecord)
	runtime.KeepAlive(d.aecs)
	runtime.KeepAlive(&d.descriptor)
	if status != iaa.Success {
		return d.completionRecord.CheckError()
	}
	return nil
}

func (d *Deflate) statsJob(block []byte, histogram *iaa.Histogram) {
	desc := &d.descriptor

	desc.SetOpcode(iaa.OpCompress)
	desc.SetFlags(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagRequestCompletionRecord |
			iaa.FlagCompletionRecordValid)

	desc.SetCompressionFlag(iaa.CompressionFlagStatsMode | iaa.CompressionFlagEndAppendEOB)
	desc.Src1Addr = uintptr(unsafe.Pointer(&block[0]))
	desc.Size = uint32(len(block))
	desc.DestAddr = uintptr(unsafe.Pointer(histogram))
	desc.MaxDestionationSize = uint32(unsafe.Sizeof(iaa.Histogram{}))
	desc.SetCompleteRecord(uintptr(unsafe.Pointer(d.completionRecord)))
}

func (d *Deflate) encodeJob(block []byte, output []byte, aesc *iaa.CompressAECS) {
	desc := &d.descriptor
	desc.SetOpcode(iaa.OpCompress)
	if d.toggle == 1 {
		desc.SetFlag(iaa.FlagAecsRWToggleSelector)
	}
	desc.SetFlag(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagCompletionRecordValid |
			iaa.FlagReadSource2Aecs |
			iaa.FlagWriteSource2CompletionOfOperation |
			iaa.FlagRequestCompletionRecord,
	)

	desc.SetCompressionFlag(
		iaa.CompressionFlagEndAppendEOB | iaa.CompressionFlagFlushOutput,
	)

	if len(block) == 0 {
		block = block[:1]
		desc.Src1Addr = uintptr(unsafe.Pointer(&block[0]))
		desc.Size = 0
	} else {
		desc.Src1Addr = uintptr(unsafe.Pointer(&block[0]))
		desc.Size = uint32(len(block))
	}

	desc.DestAddr = uintptr(unsafe.Pointer(&output[0]))
	desc.MaxDestionationSize = uint32(len(output))

	desc.Src2Addr = uintptr(unsafe.Pointer(aesc))
	desc.Src2Size = uint32(unsafe.Sizeof(iaa.CompressAECS{}))

	desc.SetCompleteRecord(uintptr(unsafe.Pointer(d.completionRecord)))
}

func (d *Deflate) generateHeader(histogram *iaa.Histogram, data *[256]byte, last bool) (headerBits int) {
	// generate huffman tree
	litCodes := histogram.LiteralCodes[:]
	offsetCodes := histogram.DistanceCodes[:]

	d.litGen.Generate(15, litCodes, litCodes)
	//
	huffman.GenerateCodeForIAA(litCodes, d.cacheForGencode)
	d.offsetGen.Generate(15, offsetCodes, offsetCodes)
	huffman.GenerateCodeForIAA(offsetCodes, d.cacheForGencode)

	d.frame.start(data)

	if d.bitsNum != 0 {
		d.frame.WriteBit(uint16(d.bits), int8(d.bitsNum))
	}

	d.dynHeader.writeTo(histogram, last, &d.frame)
	return d.frame.flush()
}

func (d *Deflate) encodeBlock(aecs *iaa.CompressAECS, block []byte, last bool, headerBits int) (err error) {
	d.descriptor.Reset()
	d.completionRecord.Reset()
	aecs.NumAccBitsValid = uint32(headerBits)
	// set prev crc result
	aecs.CRC = d.crc
	d.encodeJob(block, d.output, &d.aecs[0])
	status := iaa.StatusCode(d.ctx.Submit(uintptr(unsafe.Pointer(&d.descriptor)), &d.completionRecord.Header))
	if status != iaa.Success {
		if status == iaa.OutputBufferOverflow {
			return d.writeStoredBlock(block, last)
		}
		if status == iaa.AnalyticsError &&
			d.completionRecord.GetHeader().ErrorCode == iaa.ErrorCodeUnrecoverableOutputOverflow {
			return d.writeStoredBlock(block, last)
		}
		return d.completionRecord.CheckError()
	}
	d.crc = d.completionRecord.CRC
	if d.completionRecord.OutputSize == 0 {
		return
	}
	// check for best compression
	if d.completionRecord.OutputSize > uint32(len(block)+5) {
		return d.writeStoredBlock(block, last)
	}
	// copy the final bits to next output
	// note: the final block of this stream must keep the last byte
	if d.completionRecord.OutputBits != 0 && !last {
		d.completionRecord.OutputSize--
		if d.completionRecord.OutputBits >= 8 {
			panic("completion record outputBits must not greater that 7")
		}
		d.bits = d.output[d.completionRecord.OutputSize]
		d.bitsNum = d.completionRecord.OutputBits
	} else {
		// clear bits
		d.bits = 0
		d.bitsNum = 0
	}
	_, err = d.w.Write(d.output[:d.completionRecord.OutputSize])
	return
}

// Close the underlying writer.
func (d *Deflate) Close() error {
	closer, ok := d.w.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}

func (d *Deflate) writeStoredBlock(block []byte, last bool) error {
	blockHdr := 0b000
	if last {
		blockHdr = 0b001
	}
	hdr := uint16(d.bits)
	hdr |= uint16(blockHdr) << d.bitsNum
	offset := 1
	if d.bitsNum+3 <= 8 {
		d.output[0] = uint8(hdr)
	} else {
		d.output[0] = uint8(hdr)
		d.output[1] = uint8(hdr >> 8)
		offset = 2
	}

	binary.LittleEndian.PutUint16(d.output[offset:offset+2], uint16(len(block)))
	offset += 2

	binary.LittleEndian.PutUint16(d.output[offset:offset+2], ^uint16(len(block)))
	offset += 2
	copy(d.output[offset:], block)

	// clear bits
	d.bits = 0
	d.bitsNum = 0

	_, err := d.w.Write(d.output[:len(block)+offset])
	return err
}

var hclenOrder = []uint32{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}
