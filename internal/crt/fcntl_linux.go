// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// int open64(const char *pathname, int flags, ...);
func Xopen64(tls *TLS, pathname uintptr, flags int32, args ...interface{}) int32 {
	var mode uintptr
	if len(args) != 0 {
		switch x := args[0].(type) {
		case int32:
			mode = uintptr(x)
		default:
			panic(fmt.Errorf("crt.Xopen64 %T", x))
		}
	}
	r, _, err := syscall.Syscall(syscall.SYS_OPEN, pathname, uintptr(flags), mode)
	if strace {
		fmt.Fprintf(os.Stderr, "open(%q, %v, %#o) %v %v\n", GoString(pathname), modeString(flags), mode, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}
