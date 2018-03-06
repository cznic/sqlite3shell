// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// int ioctl(int fd, unsigned long request, ...);
func Xioctl(tls *TLS, fd int32, request ulong_t, va ...interface{}) int32 {
	panic("TODO ioctl")
}
