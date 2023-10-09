// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package device provides an interface for interacting with DSA/IAA devices.
package device

import (
	"os"
	"runtime"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/intel/ixl-go/internal/config"
	"github.com/intel/ixl-go/internal/log"
)

// Context represents the context of the device.
type Context struct {
	typ             config.DeviceType   // Device type.
	wqs             []*config.WorkQueue // Work queues.
	wqFiles         []int               // Work queue files.
	registers       [][]byte            // Registers.
	queue           uint64              // Queue.
	processors      []submitter         // Processors.
	maxTransferSize uint32
}
type submitter interface {
	Submit(desc uintptr, comp *CompletionRecordHeader) (status uint8)
	SubmitBusyPoll(desc uintptr, comp *CompletionRecordHeader) (status uint8)
}

// CreateContext creates a new context instance given the device type.
func CreateContext(typ config.DeviceType) *Context {
	c := &Context{typ: typ}
	c.init()
	if len(c.wqs) == 0 {
		log.Debug("empty workqueues")
		return nil
	}
	return c
}

// Ready returns true if the context is ready.
func (c *Context) Ready() bool {
	return c != nil
}

// init initializes the context.
func (c *Context) init() {
	var selector string
	switch c.typ {
	case config.DSA:
		selector = os.Getenv("DSA_WQ_SELECTOR")
	case config.IAA:
		selector = os.Getenv("IAA_WQ_SELECTOR")
	}
	var m matcher
	if selector == "" || selector == "*" {
		m = matchAll{}
	} else {
		matcher, err := getMatcher(selector)
		if err != nil {
			log.Debug("[%s] format error: %v \n", selector, err)
			log.Debug("fallback to use all workqueues\n")
			m = matchAll{}
		} else {
			m = matcher
		}
	}
	ctx := (&config.Context{})
	ctx.Init()

	wqs := ctx.WorkQueues(c.typ)

	for _, wq := range wqs {
		// currently we only support block on fault devices
		if wq.BlockOnFault != 1 {
			continue
		}
		if !m.match(wq) {
			continue
		}
		fd, err := syscall.Open(wq.DevicePath(), syscall.O_RDWR, 0)
		if err != nil {
			log.Debug("open %s failed: %v\n", wq.DevicePath(), err)
			continue
		}
		register, err := initWQRegister(fd)
		if err != nil {
			log.Debug("init wq register failed: %v\n", err)
			continue
		}
		if c.maxTransferSize == 0 {
			c.maxTransferSize = uint32(wq.MaxTransferSize)
		} else if wq.MaxTransferSize < uint64(c.maxTransferSize) {
			c.maxTransferSize = uint32(wq.MaxTransferSize)
		}
		c.wqs = append(c.wqs, wq)
		c.wqFiles = append(c.wqFiles, fd)
		c.registers = append(c.registers, register)
		var s submitter
		if wq.Mode == config.ModeDedicated {
			s = newDWQSubmitter(int32(wq.Size), register)
		} else {
			s = newSWQSubmitter(register)
		}
		c.processors = append(c.processors, s)
	}
}

// MaxTransferSize is the max transfer size supported by device.
func (c *Context) MaxTransferSize() uint32 {
	return c.maxTransferSize
}

// SetMaxTransferSize set the maxTransferSize, the method only used for test.
func (c *Context) SetMaxTransferSize(s uint32) {
	c.maxTransferSize = s
}

// Submit submits a new request with the given descriptor and completion record header.
func (c *Context) Submit(desc uintptr, comp *CompletionRecordHeader) uint8 {
	idx := int(atomic.AddUint64(&c.queue, 1) % uint64(len(c.processors)))
	return c.processors[idx].Submit(desc, comp)
}

// SubmitBusyPoll submits a new quick request with the given descriptor and completion record header.
// This method may cause higher CPU cost.
func (c *Context) SubmitBusyPoll(desc uintptr, comp *CompletionRecordHeader) uint8 {
	idx := int(atomic.AddUint64(&c.queue, 1) % uint64(len(c.processors)))
	return c.processors[idx].SubmitBusyPoll(desc, comp)
}

// initWQRegister initializes a new work queue register.
func initWQRegister(fd int) ([]byte, error) {
	return syscall.Mmap(fd, 0, 0x1000, syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_POPULATE)
}

// enqcmd enqueues a descriptor.
func enqcmd(ctx *byte, desc uintptr) bool

// endcmdWithRetry retries enqueuing a new command.
func endcmdWithRetry(ctx *byte, desc uintptr) bool

// waitForComplete waits for the request to complete.
func waitForComplete(comp *uint64) (status uint8)

// movdir64b enqueues a descriptor
func movdir64b(ctx *byte, desc uintptr)

func newDWQSubmitter(max int32, register []byte) (p *dwqSubmitter) {
	p = &dwqSubmitter{}
	p.max = max
	p.register = register
	return p
}

type dwqSubmitter struct {
	max      int32
	sem      atomic.Int32
	register []byte
}

// Submit submits a new request with the given descriptor and completion record header and wait the result.
func (p *dwqSubmitter) Submit(desc uintptr, comp *CompletionRecordHeader) (status uint8) {
	// clear status
	comp.ComplexStatus = 0
	uip := (*uint64)(unsafe.Pointer(comp))
	for {
		s := p.sem.Load()
		if s >= p.max {
			runtime.Gosched()
			continue
		}
		if p.sem.CompareAndSwap(s, s+1) {
			break
		}
	}
	movdir64b(&p.register[0], desc)
	for {
		runtime.Gosched()
		hdr := atomic.LoadUint64(uip)
		h := (*CompletionRecordHeader)(unsafe.Pointer(&hdr))
		if h.ComplexStatus == 0 {
			continue
		}
		p.sem.Add(-1)
		status := h.ComplexStatus & 0b00011111
		return status
	}
}

// SubmitBusyPoll submits a new request with the given descriptor and completion record header
// and wait the result by busy-polling.
// This method may cause higher CPU cost.
func (p *dwqSubmitter) SubmitBusyPoll(desc uintptr, comp *CompletionRecordHeader) (status uint8) {
	// clear status
	comp.ComplexStatus = 0
	uip := (*uint64)(unsafe.Pointer(comp))
	for {
		s := p.sem.Load()
		if s >= p.max {
			runtime.Gosched()
			continue
		}
		if p.sem.CompareAndSwap(s, s+1) {
			break
		}
	}
	movdir64b(&p.register[0], desc)
	status = waitForComplete(uip)
	p.sem.Add(-1)
	return status
}

// swqSubmitter represents a swqSubmitter of the context.
type swqSubmitter struct {
	register []byte // Register.
}

// newSWQSubmitter creates a new processor instance.
func newSWQSubmitter(register []byte) (p *swqSubmitter) {
	p = &swqSubmitter{}
	p.register = register
	return p
}

// SubmitBusyPoll submits a new request with the given descriptor and completion record header
// and wait the result by busy-polling.
// This method may cause higher CPU cost.
func (p *swqSubmitter) SubmitBusyPoll(desc uintptr, comp *CompletionRecordHeader) (status uint8) {
	// clear status
	comp.ComplexStatus = 0

	uip := (*uint64)(unsafe.Pointer(comp))
	ret := endcmdWithRetry(&p.register[0], desc)
	if ret {
		panic("unexpected ENQCMD return value")
	}
	status = waitForComplete(uip)
	return status
}

// Submit submits a new request with the given descriptor and completion record header and wait the result.
func (p *swqSubmitter) Submit(desc uintptr, comp *CompletionRecordHeader) (status uint8) {
	// clear status
	comp.ComplexStatus = 0

	uip := (*uint64)(unsafe.Pointer(comp))
	for enqcmd(&p.register[0], desc) {
		// go cannot setup goroutine's priority
		runtime.Gosched()
	}
	for {
		runtime.Gosched()
		hdr := atomic.LoadUint64(uip)
		h := (*CompletionRecordHeader)(unsafe.Pointer(&hdr))
		if h.ComplexStatus == 0 {
			continue
		}
		status := h.ComplexStatus & 0b00011111
		return status
	}
}

// CompletionRecordHeader represents the completion record header.
type CompletionRecordHeader struct {
	ComplexStatus  uint8  // Status.
	ErrorCode      uint8  // ErrorCode.
	Rsvd           uint16 // Reserved.
	BytesCompleted uint32 // Bytes completed.
}

// Status return header's status
func (c *CompletionRecordHeader) Status() uint8 {
	return c.ComplexStatus & 0b00011111
}
