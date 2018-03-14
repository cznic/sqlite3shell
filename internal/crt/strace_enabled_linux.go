package crt

import (
	"fmt"

	fcntl2 "github.com/cznic/ccir/libc/fcntl"
)

// +build crt.strace

func cmdString(cmd int32) string {
	switch cmd {
	case fcntl2.XF_DUPFD:
		return "F_DUPFD"
	case fcntl2.XF_GETFD:
		return "F_GETFD"
	case fcntl2.XF_GETFL:
		return "F_GETFL"
	case fcntl2.XF_GETLK:
		return "F_GETLK"
	case fcntl2.XF_GETOWN:
		return "F_GETOWN"
	case fcntl2.XF_SETFD:
		return "F_SETFD"
	case fcntl2.XF_SETFL:
		return "F_SETFL"
	case fcntl2.XF_SETLK:
		return "F_SETLK"
	case fcntl2.XF_SETLKW:
		return "F_SETLKW"
	case fcntl2.XF_SETOWN:
		return "F_SETOWN"
	default:
		return fmt.Sprintf("%#x", cmd)
	}
}
