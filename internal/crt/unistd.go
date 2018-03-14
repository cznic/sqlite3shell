// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"time"

	"github.com/cznic/ccir/libc/errno"
	"golang.org/x/crypto/ssh/terminal"
)

// int rmdir(const char *pathname);
func Xrmdir(tls *TLS, pathname uintptr) int32 {
	panic("TODO rmdir")
}

// int fchown(int fd, uid_t owner, gid_t group);
func Xfchown(tls *TLS, fd int32, owner, group uint32) int32 {
	panic("TODO fchown")
}

// unsigned sleep(unsigned seconds);
func Xsleep(tls *TLS, seconds uint32) uint32 {
	time.Sleep(time.Duration(seconds) * time.Second)
	if strace {
		fmt.Fprintf(os.Stderr, "sleep(%#x)", seconds)
	}
	return 0
}

// int usleep(useconds_t usec);
func Xusleep(tls *TLS, usec uint32) int32 {
	time.Sleep(time.Duration(usec) * time.Microsecond)
	if strace {
		fmt.Fprintf(os.Stderr, "usleep(%#x)", usec)
	}
	return 0
}

// int chdir(const char *path);
func Xchdir(tls *TLS, path uintptr) int32 {
	panic("TODO chdir")
}

// ssize_t readlink(const char *pathname, char *buf, size_t bufsiz);
func Xreadlink(tls *TLS, pathname, buf uintptr, bufsiz size_t) ssize_t {
	panic("TODO readlink")
}

// int isatty(int fd);
func Xisatty(tls *TLS, fd int32) int32 {
	if terminal.IsTerminal(int(fd)) {
		return 1
	}

	tls.setErrno(errno.XENOTTY)
	return 0
}

func X_isatty(tls *TLS, fd int32) int32 { return Xisatty(tls, fd) }
