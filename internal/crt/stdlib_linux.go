// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"os/exec"
	"syscall"
	"unsafe"
)

// int system(const char *command);
func Xsystem(tls *TLS, command uintptr) int32 {
	if command == 0 || *(*int8)(unsafe.Pointer(command)) == 0 {
		return 1
	}

	cmd := exec.Command("sh", "-c", GoString(command))
	if err := cmd.Run(); err != nil {
		return int32(cmd.ProcessState.Sys().(syscall.WaitStatus))
	}

	return 0
}
