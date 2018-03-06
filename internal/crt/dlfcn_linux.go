// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// void *dlopen(const char *filename, int flags);
func Xdlopen(tls *TLS, filename uintptr, flags int32) uintptr {
	panic("TODO dlopen")
}

// char *dlerror(void);
func Xdlerror(tls *TLS) uintptr {
	panic("TODO dlerror")
}

// int dlclose(void *handle);
func Xdlclose(tls *TLS, handle uintptr) int32 {
	panic("TODO dlclose")
}

// void *dlsym(void *handle, const char *symbol);
func Xdlsym(tls *TLS, handle, symbol uintptr) uintptr {
	panic("TODO dlsym")
}
