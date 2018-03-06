// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// void *mmap(void *addr, size_t len, int prot, int flags, int fildes, off_t off);
func Xmmap64(tls *TLS, addr uintptr, len size_t, prot, flags, fildes int32, off int64) uintptr {
	r, _, err := syscall.Syscall6(syscall.SYS_MMAP, addr, uintptr(len), uintptr(prot), uintptr(flags), uintptr(fildes), uintptr(off))
	if strace {
		fmt.Fprintf(os.Stderr, "mmap(%#x, %#x, %#x, %#x, %#x, %#x) (%#x, %v)\n", addr, len, prot, flags, fildes, off, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}

	return r
}

// int munmap(void *addr, size_t len);
func Xmunmap(tls *TLS, addr uintptr, len size_t) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_MUNMAP, addr, uintptr(len), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "munmap(%#x, %#x) (%#x, %v)\n", addr, len, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}
