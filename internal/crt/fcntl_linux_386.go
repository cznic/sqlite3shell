// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// int fcntl(int fildes, int cmd, ...);
func Xfcntl(tls *TLS, fildes, cmd int32, args ...interface{}) int32 {
	var arg uintptr
	if len(args) != 0 {
		switch x := args[0].(type) {
		case int32:
			arg = uintptr(x)
		case uintptr:
			arg = x
		default:
			panic(fmt.Errorf("crt.Xfcntl %T", x))
		}
	}
	r, _, err := syscall.Syscall(syscall.SYS_FCNTL64, uintptr(fildes), uintptr(cmd), arg)
	if strace {
		fmt.Fprintf(os.Stderr, "fcntl(%v, %v, %#x) %v %v\n", fildes, cmdString(cmd), arg, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}
