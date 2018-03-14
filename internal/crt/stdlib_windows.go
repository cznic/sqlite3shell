package crt

import (
	"os/exec"
	"syscall"
	"unsafe"
)

func Xsystem(tls *TLS, command uintptr) int32 {
	if command == 0 || *(*int8)(unsafe.Pointer(command)) == 0 {
		return 1
	}

	cmd := exec.Command("sh", "-c", GoString(command))
	if err := cmd.Run(); err != nil {
		return int32(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitCode)
	}

	return 0
}
