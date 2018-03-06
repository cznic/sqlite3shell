// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// int gettimeofday(struct timeval *restrict tp, void *restrict tzp);
func Xgettimeofday(tls *TLS, tp, tzp uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_GETTIMEOFDAY, tp, tzp, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "gettimeofday(%#x, %#x) %v %v\n", tp, tzp, r, err)
	}
	return int32(r)
}

// int utimes(const char *filename, const struct timeval times[2]);
func Xutimes(tls *TLS, filename, times uintptr) int32 {
	panic("TODO utimes")
}
