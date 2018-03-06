// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"syscall"
	"unsafe"

	"github.com/cznic/mathutil"
)

// void exit(int);
func Xexit(tls *TLS, n int32) { X__builtin_exit(tls, n) }

// // void exit(int);
// func X__builtin_exit(tls *TLS, n int32) {
// 	os.Exit(int(n))
// }

// void free(void *ptr);
func Xfree(tls *TLS, ptr uintptr) {
	free(tls, ptr)
	if strace {
		fmt.Fprintf(os.Stderr, "free(%#x)\n", ptr)
	}
}

// void abort();
func Xabort(tls *TLS) { X__builtin_abort(tls) }

// void __builtin_trap();
func X__builtin_trap(tls *TLS) { os.Exit(1) }

// void abort();
func X__builtin_abort(tls *TLS) { X__builtin_trap(tls) }

// char *getenv(const char *name);
func Xgetenv(tls *TLS, name uintptr) uintptr {
	nm := GoString(name)
	v := os.Getenv(nm)
	var p uintptr
	if v != "" {
		p = CString(v)
	}
	return p
}

// int abs(int j);
func X__builtin_abs(tls *TLS, j int32) int32 {
	if j < 0 {
		return -j
	}

	return j
}

// int abs(int j);
func Xabs(tls *TLS, j int32) int32 { return X__builtin_abs(tls, j) }

type sorter struct {
	base   uintptr
	compar func(tls *TLS, a, b uintptr) int32
	nmemb  int
	size   uintptr
	tls    *TLS
	buf    []byte
}

func (s *sorter) Len() int { return s.nmemb }

func (s *sorter) Less(i, j int) bool {
	return s.compar(s.tls, s.base+uintptr(i)*s.size, s.base+uintptr(j)*s.size) < 0
}

func (s *sorter) Swap(i, j int) {
	p := s.base + uintptr(i)*s.size
	q := s.base + uintptr(j)*s.size
	switch s.size {
	case 1:
		*(*int8)(unsafe.Pointer(p)), *(*int8)(unsafe.Pointer(q)) = *(*int8)(unsafe.Pointer(q)), *(*int8)(unsafe.Pointer(p))
	case 2:
		*(*int16)(unsafe.Pointer(p)), *(*int16)(unsafe.Pointer(q)) = *(*int16)(unsafe.Pointer(q)), *(*int16)(unsafe.Pointer(p))
	case 4:
		*(*int32)(unsafe.Pointer(p)), *(*int32)(unsafe.Pointer(q)) = *(*int32)(unsafe.Pointer(q)), *(*int32)(unsafe.Pointer(p))
	case 8:
		*(*int64)(unsafe.Pointer(p)), *(*int64)(unsafe.Pointer(q)) = *(*int64)(unsafe.Pointer(q)), *(*int64)(unsafe.Pointer(p))
	default:
		size := int(s.size)
		if s.buf == nil {
			s.buf = make([]byte, size)
		}
		Copy(uintptr(unsafe.Pointer(&s.buf[0])), p, size)
		Copy(p, q, size)
		Copy(q, uintptr(unsafe.Pointer(&s.buf[0])), size)
	}
}

// void qsort(void *base, size_t nmemb, size_t size, int (*compar)(const void *, const void *));
func qsort(tls *TLS, base uintptr, nmemb, size size_t, compar uintptr) {
	if size > mathutil.MaxInt {
		panic("size overflow")
	}

	if nmemb > mathutil.MaxInt {
		panic("nmemb overflow")
	}

	s := sorter{base, *(*func(*TLS, uintptr, uintptr) int32)(unsafe.Pointer(&compar)), int(nmemb), uintptr(size), tls, nil}
	sort.Sort(&s)
}

// int system(const char *command);
func Xsystem(tls *TLS, command uintptr) int32 {
	if command == 0 || *(*int8)(unsafe.Pointer(command)) == 0 {
		return 1
	}

	cmd := exec.Command("sh", "-c", GoString(command))
	if err := cmd.Run(); err != nil {
		return int32(cmd.ProcessState.Sys().(syscall.WaitStatus))
	}

	return 0
}

// void *calloc(size_t nmemb, size_t size);
func Xcalloc(tls *TLS, nmemb, size size_t) (p uintptr) {
	hi, lo := mathutil.MulUint128_64(uint64(nmemb), uint64(size))
	if hi == 0 && lo <= mathutil.MaxInt {
		p = calloc(tls, int(lo))
	}
	if strace {
		fmt.Fprintf(os.Stderr, "calloc(%#x) %#x\n", size, p)
	}
	return p
}

// void *malloc(size_t size);
func X__builtin_malloc(tls *TLS, size size_t) (p uintptr) {
	if size < mathutil.MaxInt {
		p = malloc(tls, int(size))
	}
	if strace {
		fmt.Fprintf(os.Stderr, "malloc(%#x) %#x\n", size, p)
	}
	return p
}

// void *malloc(size_t size);
func Xmalloc(tls *TLS, size size_t) uintptr { return X__builtin_malloc(tls, size) }

// void *realloc(void *ptr, size_t size);
func Xrealloc(tls *TLS, ptr uintptr, size size_t) (p uintptr) {
	if size < mathutil.MaxInt {
		p = realloc(tls, ptr, int(size))
	}
	if strace {
		fmt.Fprintf(os.Stderr, "realloc(%#x, %#x) %#x\n", ptr, size, p)
	}
	return p
}

// size_t malloc_usable_size (void *ptr);
// func Xmalloc_usable_size(tls *TLS, ptr uintptr) size_t { return size_t(memory.UintptrUsableSize(ptr)) }

// void qsort(void *base, size_t nmemb, size_t size, int (*compar)(const void *, const void *));
func Xqsort(tls *TLS, base uintptr, nmemb, size size_t, compar uintptr) {
	qsort(tls, base, nmemb, size, compar)
}

// long int strtol(const char *nptr, char **endptr, int base);
func Xstrtol(tls *TLS, nptr, endptr uintptr, base int32) int64 {
	panic("TODO strtol")
}
