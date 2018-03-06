// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

// sighandler_t sysv_signal(int signum, sighandler_t handler);
func X__sysv_signal(tls *TLS, signum int32, handler uintptr) {
	ch := make(chan os.Signal)
	go func() {
		<-ch
		(*(*func(*TLS, int32))(unsafe.Pointer(handler)))(tls, signum)
	}()
	signal.Notify(ch, syscall.Signal(signum))
}
