// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// int fchmod(int fd, mode_t mode);
func Xfchmod(tls *TLS, fd int32, mode uint32) int32 {
	panic("TODO fchmod")
}

// int mkdir(const char *pathname, mode_t mode);
func Xmkdir(tls *TLS, pathname uintptr, mode uint32) int32 {
	panic("TODO mkdir")
}
