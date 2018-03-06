// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/cznic/ccir/libc/errno"
	"github.com/cznic/ccir/libc/unistd"
	"golang.org/x/crypto/ssh/terminal"
)

// int close(int fd);
func Xclose(tls *TLS, fd int32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_CLOSE, uintptr(fd), 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "close(%v) %v %v\n", fd, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int access(const char *path, int amode);
func Xaccess(tls *TLS, path uintptr, amode int32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_ACCESS, path, uintptr(amode), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "access(%q) %v %v\n", GoString(path), r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int unlink(const char *path);
func Xunlink(tls *TLS, path uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_UNLINK, path, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "unlink(%q) %v %v\n", GoString(path), r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int rmdir(const char *pathname);
func Xrmdir(tls *TLS, pathname uintptr) int32 {
	panic("TODO rmdir")
}

// int fchown(int fd, uid_t owner, gid_t group);
func Xfchown(tls *TLS, fd int32, owner, group uint32) int32 {
	panic("TODO fchown")
}

// uid_t getuid(void);
func Xgetuid(tls *TLS) uint32 {
	r, _, _ := syscall.RawSyscall(syscall.SYS_GETUID, 0, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "getuid() %v\n", r)
	}
	return uint32(r)
}

// uid_t geteuid(void);
func Xgeteuid(tls *TLS) uint32 {
	r, _, _ := syscall.RawSyscall(syscall.SYS_GETEUID, 0, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "geteuid() %v\n", r)
	}
	return uint32(r)
}

// int fsync(int fildes);
func Xfsync(tls *TLS, fildes int32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_FSYNC, uintptr(fildes), 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "fsync(%v) %v %v\n", fildes, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int fdatasync(int fd);
func Xfdatasync(tls *TLS, fildes int32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_FDATASYNC, uintptr(fildes), 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "fdatasync(%v) %v %v\n", fildes, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// pid_t getpid(void);
func Xgetpid(tls *TLS) int32 {
	r, _, _ := syscall.RawSyscall(syscall.SYS_GETPID, 0, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "getpid() %v\n", r)
	}
	return int32(r)
}

// unsigned sleep(unsigned seconds);
func Xsleep(tls *TLS, seconds uint32) uint32 {
	time.Sleep(time.Duration(seconds) * time.Second)
	if strace {
		fmt.Fprintf(os.Stderr, "sleep(%#x)", seconds)
	}
	return 0
}

// off_t lseek64(int fildes, off_t offset, int whence);
func Xlseek64(tls *TLS, fildes int32, offset int64, whence int32) int64 {
	r, _, err := syscall.Syscall(syscall.SYS_LSEEK, uintptr(fildes), uintptr(offset), uintptr(whence))
	if strace {
		fmt.Fprintf(os.Stderr, "lseek(%v, %v, %v) %v %v\n", fildes, offset, whence, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int64(r)
}

// int ftruncate(int fildes, off_t length);
func Xftruncate64(tls *TLS, fildes int32, length int64) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_FTRUNCATE, uintptr(fildes), uintptr(length), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "ftruncate(%#x, %#x) %v, %v\n", fildes, length, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
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

// ssize_t read(int fd, void *buf, size_t count);
func Xread(tls *TLS, fd int32, buf uintptr, count size_t) ssize_t { //TODO stdin
	r, _, err := syscall.Syscall(syscall.SYS_READ, uintptr(fd), buf, uintptr(count))
	if strace {
		fmt.Fprintf(os.Stderr, "read(%v, %#x, %v) %v %v\n", fd, buf, count, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return ssize_t(r)
}

// char *getcwd(char *buf, size_t size);
func Xgetcwd(tls *TLS, buf uintptr, size size_t) uintptr {
	r, _, err := syscall.Syscall(syscall.SYS_GETCWD, buf, uintptr(size), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "getcwd(%#x, %#x) %v %v %q\n", buf, size, r, err, GoString(buf))
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return r
}

// ssize_t write(int fd, const void *buf, size_t count);
func Xwrite(tls *TLS, fd int32, buf uintptr, count size_t) ssize_t {
	switch fd {
	case unistd.XSTDOUT_FILENO:
		n, err := os.Stdout.Write((*rawmem)(unsafe.Pointer(buf))[:count])
		if err != nil {
			tls.setErrno(err)
		}
		return ssize_t(n)
	case unistd.XSTDERR_FILENO:
		n, err := os.Stderr.Write((*rawmem)(unsafe.Pointer(buf))[:count])
		if err != nil {
			tls.setErrno(err)
		}
		return ssize_t(n)
	}
	r, _, err := syscall.Syscall(syscall.SYS_WRITE, uintptr(fd), buf, uintptr(count))
	if strace {
		fmt.Fprintf(os.Stderr, "write(%v, %#x, %v) %v %v\n", fd, buf, count, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return ssize_t(r)
}

// ssize_t readlink(const char *pathname, char *buf, size_t bufsiz);
func Xreadlink(tls *TLS, pathname, buf uintptr, bufsiz size_t) ssize_t {
	panic("TODO readlink")
}

// long sysconf(int name);
func Xsysconf(tls *TLS, name int32) int64 {
	switch name {
	case unistd.X_SC_PAGESIZE:
		return int64(os.Getpagesize())
	default:
		panic(fmt.Errorf("%v(%#x)", name, name))
	}
}

// int isatty(int fd);
func Xisatty(tls *TLS, fd int32) int32 {
	if terminal.IsTerminal(int(fd)) {
		return 1
	}

	tls.setErrno(errno.XENOTTY)
	return 0
}
