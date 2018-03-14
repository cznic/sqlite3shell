// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"github.com/cznic/ccir/libc/unistd"
)

// int access(const char *path, int amode);
func X_access(tls *TLS, path uintptr, amode int32) int32 {
	mode := 0
	if amode != unistd.XF_OK {
		panic("access mode not supported")
	}

	f := openFile(tls, GoString(path), mode)
	if f != nil {
		if err := f.Close(); err != nil {
			return -1
		}
		return 0
	}
	// TODO: potentially support more
	return -1
}

func X_setmode(tls *TLS, fd uintptr, amode int32) int32 {
	// TODO
	return 0
}

// TODO: These shouldn't be required here? ctype & stuff?
func Xfputs(tls *TLS, fd uintptr, data uintptr) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}

func Xisdigit(tls *TLS, chr int32) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}

func Xisspace(tls *TLS, chr int32) uintptr {
	if chr == int32(' ') {
		return 1
	}
	return 0
}

func Xisalpha(tls *TLS, chr int32) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}

func Xisalnum(tls *TLS, chr int32) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}
