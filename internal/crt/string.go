// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"unsafe"
)

// char *strcat(char *dest, const char *src)
func Xstrcat(tls *TLS, dest, src uintptr) uintptr {
	ret := dest
	for *(*int8)(unsafe.Pointer(dest)) != 0 {
		dest++
	}
	for {
		c := *(*int8)(unsafe.Pointer(src))
		src++
		*(*int8)(unsafe.Pointer(dest)) = c
		dest++
		if c == 0 {
			return ret
		}
	}
}

// char *index(const char *s, int c)
func Xindex(tls *TLS, s uintptr, c int32) uintptr { return Xstrchr(tls, s, c) }

// char *strchr(const char *s, int c)
func Xstrchr(tls *TLS, s uintptr, c int32) uintptr {
	for {
		ch2 := *(*byte)(unsafe.Pointer(s))
		if ch2 == byte(c) {
			return s
		}

		if ch2 == 0 {
			return 0
		}

		s++
	}
}

// char *strchrnul(const char *s, int c);
func Xstrchrnul(tls *TLS, s uintptr, c int32) uintptr {
	for {
		ch2 := *(*byte)(unsafe.Pointer(s))
		if ch2 == 0 || ch2 == byte(c) {
			return s
		}

		s++
	}
}

// int strcmp(const char *s1, const char *s2)
func X__builtin_strcmp(tls *TLS, s1, s2 uintptr) int32 {
	for {
		ch1 := *(*byte)(unsafe.Pointer(s1))
		s1++
		ch2 := *(*byte)(unsafe.Pointer(s2))
		s2++
		if ch1 != ch2 || ch1 == 0 || ch2 == 0 {
			return int32(ch1) - int32(ch2)
		}
	}
}

// int strcmp(const char *s1, const char *s2)
func Xstrcmp(tls *TLS, s1, s2 uintptr) int32 { return X__builtin_strcmp(tls, s1, s2) }

// char *strcpy(char *dest, const char *src)
func X__builtin_strcpy(tls *TLS, dest, src uintptr) uintptr {
	r := dest
	for {
		c := *(*int8)(unsafe.Pointer(src))
		src++
		*(*int8)(unsafe.Pointer(dest)) = c
		dest++
		if c == 0 {
			return r
		}
	}
}

// char *strcpy(char *dest, const char *src)
func Xstrcpy(tls *TLS, dest, src uintptr) uintptr { return X__builtin_strcpy(tls, dest, src) }

// char *rindex(const char *s, int c)
func Xrindex(tls *TLS, s uintptr, c int32) uintptr { return Xstrrchr(tls, s, c) }

// char *strrchr(const char *s, int c)
func Xstrrchr(tls *TLS, s uintptr, c int32) uintptr {
	var ret uintptr
	for {
		ch2 := *(*byte)(unsafe.Pointer(s))
		if ch2 == 0 {
			return ret
		}

		if ch2 == byte(c) {
			ret = s
		}
		s++
	}
}

// char *strstr(const char *haystack, const char *needle);
func Xstrstr(tls *TLS, haystack, needle uintptr) int8 {
	panic("TODO strstr")
}

// char *strncpy(char *dest, const char *src, size_t n)
func Xstrncpy(tls *TLS, dest, src uintptr, n size_t) uintptr {
	ret := dest
	for c := *(*int8)(unsafe.Pointer(src)); c != 0 && n > 0; n-- {
		*(*int8)(unsafe.Pointer(dest)) = c
		dest++
		src++
		c = *(*int8)(unsafe.Pointer(src))
	}
	for ; n > 0; n-- {
		*(*int8)(unsafe.Pointer(dest)) = 0
		dest++
	}
	return ret
}

// size_t strlen(const char *s)
func X__builtin_strlen(tls *TLS, s uintptr) size_t {
	var n size_t
	for ; *(*int8)(unsafe.Pointer(s)) != 0; s++ {
		n++
	}
	return n
}

// size_t strlen(const char *s)
func Xstrlen(tls *TLS, s uintptr) size_t { return X__builtin_strlen(tls, s) }

// int strncmp(const char *s1, const char *s2, size_t n)
func Xstrncmp(tls *TLS, s1, s2 uintptr, n size_t) int32 {
	var ch1, ch2 byte
	for n != 0 {
		ch1 = *(*byte)(unsafe.Pointer(s1))
		s1++
		ch2 = *(*byte)(unsafe.Pointer(s2))
		s2++
		n--
		if ch1 != ch2 || ch1 == 0 || ch2 == 0 {
			break
		}
	}
	if n != 0 {
		return int32(ch1) - int32(ch2)
	}

	return 0
}

// void *memset(void *s, int c, size_t n)
func Xmemset(tls *TLS, s uintptr, c int32, n size_t) uintptr {
	return X__builtin_memset(tls, s, c, n)
}

// void *memset(void *s, int c, size_t n)
func X__builtin_memset(tls *TLS, s uintptr, c int32, n size_t) uintptr {
	for d := s; n > 0; n-- {
		*(*int8)(unsafe.Pointer(d)) = int8(c)
		d++
	}
	return s
}

// void *memcpy(void *dest, const void *src, size_t n)
func X__builtin_memcpy(tls *TLS, dest, src uintptr, n size_t) uintptr {
	Copy(dest, src, int(n))
	return dest
}

// void *memcpy(void *dest, const void *src, size_t n)
func Xmemcpy(tls *TLS, dest, src uintptr, n size_t) uintptr {
	return X__builtin_memcpy(tls, dest, src, n)
}

// int memcmp(const void *s1, const void *s2, size_t n)
func X__builtin_memcmp(tls *TLS, s1, s2 uintptr, n size_t) int32 {
	var ch1, ch2 byte
	for n != 0 {
		ch1 = *(*byte)(unsafe.Pointer(s1))
		s1++
		ch2 = *(*byte)(unsafe.Pointer(s2))
		s2++
		if ch1 != ch2 {
			break
		}

		n--
	}
	if n != 0 {
		return int32(ch1) - int32(ch2)
	}

	return 0
}

// int memcmp(const void *s1, const void *s2, size_t n)
func Xmemcmp(tls *TLS, s1, s2 uintptr, n size_t) int32 {
	return X__builtin_memcmp(tls, s1, s2, n)
}

// void *memmove(void *dest, const void *src, size_t n);
func Xmemmove(tls *TLS, dest, src uintptr, n size_t) uintptr {
	Copy(dest, src, int(n))
	return dest
}

// void *mempcpy(void *dest, const void *src, size_t n);
func Xmempcpy(tls *TLS, dest, src uintptr, n size_t) uintptr {
	return dest + uintptr(Copy(dest, src, int(n)))
}
