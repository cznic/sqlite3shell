// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crt provides C-runtime services. (Work In Progress)
//
// Installation
//
//     $ go get github.com/cznic/crt
//
// Documentation: http://godoc.org/github.com/cznic/crt
package crt

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/cznic/ccir/libc/errno"
	"github.com/cznic/internal/buffer"
	"github.com/cznic/mathutil"
	"github.com/cznic/memory"
)

const (
	ptrSize = mathutil.UintPtrBits / 8
)

var (
	_ io.Writer = (*memWriter)(nil)

	allocMu   sync.Mutex
	allocator memory.Allocator
	threadID  uintptr
)

type (
	T__uid_t = uint32
	T__gid_t = uint32
)

// TLS represents the C-thread local storage.
type TLS struct {
	threadID uintptr
	errno    int32
}

func (t *TLS) setErrno(err interface{}) {
	switch x := err.(type) {
	case int:
		t.errno = int32(x)
	case *os.PathError:
		t.setErrno(x.Err)
	case syscall.Errno:
		t.errno = int32(x)
	default:
		panic(fmt.Errorf("crt.setErrno %T(%#v)", x, x))
	}
}

// NewTLS returns a newly created TLS, allocated outside of the Go runtime
// heap.  To free the TLS use Free(uintptr(unsafe.Pointer(tls))).
func NewTLS() *TLS {
	tls := (*TLS)(unsafe.Pointer(MustCalloc(int(unsafe.Sizeof(TLS{})))))
	tls.threadID = atomic.AddUintptr(&threadID, 1)
	return tls
}

// void __register_stdfiles(void *, void *, void *);
func X__register_stdfiles(tls *TLS, in, out, err uintptr) {
	stdin = in
	stdout = out
	stderr = err
}

// void exit(int);
func X__builtin_exit(tls *TLS, n int32) { os.Exit(int(n)) }

// BSS allocates the the bss segment of a package/command.
func BSS(init *byte) uintptr {
	r := uintptr(unsafe.Pointer(init))
	if r%unsafe.Sizeof(uintptr(0)) != 0 {
		panic("internal error")
	}

	return r
}

// DS allocates the the data segment of a package/command.
func DS(init []byte) uintptr {
	r := (*reflect.SliceHeader)(unsafe.Pointer(&init)).Data
	if r%unsafe.Sizeof(uintptr(0)) != 0 {
		panic("internal error")
	}

	return r
}

// TS allocates the R/O text segment of a package/command.
func TS(init string) uintptr { return (*reflect.StringHeader)(unsafe.Pointer(&init)).Data }

// Free frees memory allocated by Calloc, Malloc or Realloc.
func Free(p uintptr) error {
	allocMu.Lock()
	err := allocator.UintptrFree(p)
	allocMu.Unlock()
	return err
}

// Calloc allocates zeroed memory.
func Calloc(size int) (uintptr, error) {
	allocMu.Lock()
	p, err := allocator.UintptrCalloc(size)
	allocMu.Unlock()
	return p, err
}

// Malloc allocates memory.
func Malloc(size int) (uintptr, error) {
	allocMu.Lock()
	p, err := allocator.UintptrMalloc(size)
	allocMu.Unlock()
	return p, err
}

// MustCalloc is like Calloc but panics if the allocation cannot be made.
func MustCalloc(size int) uintptr {
	p, err := Calloc(size)
	if err != nil {
		panic(fmt.Errorf("out of memory: %v", err))
	}

	return p
}

// MustMalloc is like Malloc but panics if the allocation cannot be made.
func MustMalloc(size int) uintptr {
	p, err := Malloc(size)
	if err != nil {
		panic(fmt.Errorf("out of memory: %v", err))
	}

	return p
}

// Realloc reallocates memory.
func Realloc(p uintptr, size int) (uintptr, error) {
	allocMu.Lock()
	p, err := allocator.UintptrRealloc(p, size)
	allocMu.Unlock()
	return p, err
}

type memWriter uintptr

func (m *memWriter) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	Copy(uintptr(*m), uintptr(unsafe.Pointer(&b[0])), len(b))
	*m += memWriter(len(b))
	return len(b), nil
}

func (m *memWriter) WriteByte(b byte) error {
	*(*byte)(unsafe.Pointer(*m)) = b
	*m++
	return nil
}

// Copy copies n bytes form src to dest and returns n.
func Copy(dst, src uintptr, n int) int {
	return copy((*rawmem)(unsafe.Pointer(dst))[:n], (*rawmem)(unsafe.Pointer(src))[:n])
}

// CString allocates a C string initialized from s.
func CString(s string) uintptr {
	n := len(s)
	var tls TLS
	p := malloc(&tls, n+1)
	if p == 0 {
		return 0
	}

	copy((*rawmem)(unsafe.Pointer(p))[:n], s)
	(*rawmem)(unsafe.Pointer(p))[n] = 0
	return p
}

// GoString returns a string from a C char* null terminated string s.
func GoString(s uintptr) string {
	if s == 0 {
		return ""
	}

	var b buffer.Bytes
	for {
		ch := *(*byte)(unsafe.Pointer(s))
		if ch == 0 {
			r := string(b.Bytes())
			b.Close()
			return r
		}

		b.WriteByte(ch)
		s++
	}
}

func calloc(tls *TLS, size int) uintptr {
	p, err := Calloc(size)
	if err != nil {
		tls.setErrno(errno.XENOMEM)
		return 0
	}

	return p
}

func free(tls *TLS, p uintptr) { Free(p) }

func malloc(tls *TLS, size int) uintptr {
	p, err := Malloc(size)
	if err != nil {
		tls.setErrno(errno.XENOMEM)
		return 0
	}

	return p
}

func realloc(tls *TLS, p uintptr, size int) uintptr {
	p, err := Realloc(p, size)
	if err != nil {
		tls.setErrno(errno.XENOMEM)
		return 0
	}

	return p
}
