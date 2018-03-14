// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

func X_wunlink(tls *TLS, ptr uintptr) int32 {
	panic("TODO wunlink")
}

func X_fileno(tls *TLS, fd uintptr) uintptr {
	return fd
}
