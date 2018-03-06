// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"unsafe"

	"github.com/cznic/ccir/libc/errno"
	"github.com/cznic/ccir/libc/stdio"
	"github.com/cznic/internal/buffer"
	"github.com/cznic/mathutil"
)

var (
	stdin, stdout, stderr uintptr

	files = &fmap{
		m: map[uintptr]*os.File{},
	}
	nullReader = bytes.NewBuffer(nil)
)

type fmap struct {
	m  map[uintptr]*os.File
	mu sync.Mutex
}

func (m *fmap) add(f *os.File, u uintptr) {
	m.mu.Lock()
	m.m[u] = f
	m.mu.Unlock()
}

func (m *fmap) reader(u uintptr) io.Reader {
	switch u {
	case stdin:
		return os.Stdin
	case stdout, stderr:
		return nullReader
	}

	m.mu.Lock()
	f := m.m[u]
	m.mu.Unlock()
	return f
}

func (m *fmap) file(u uintptr) *os.File {
	switch u {
	case stdin:
		return os.Stdin
	case stdout:
		return os.Stdout
	case stderr:
		return os.Stderr
	}

	m.mu.Lock()
	f := m.m[u]
	m.mu.Unlock()
	return f
}

func (m *fmap) writer(u uintptr) io.Writer {
	switch u {
	case stdin:
		return ioutil.Discard
	case stdout:
		return os.Stdout
	case stderr:
		return os.Stderr
	}

	m.mu.Lock()
	f := m.m[u]
	m.mu.Unlock()
	return f
}

func (m *fmap) extract(u uintptr) *os.File {
	m.mu.Lock()
	f := m.m[u]
	delete(m.m, u)
	m.mu.Unlock()
	return f
}

// int printf(const char *format, ...);
func Xprintf(tls *TLS, format uintptr /* *int8 */, args ...interface{}) int32 {
	return X__builtin_printf(tls, format, args...)
}

// int printf(const char *format, ...);
func X__builtin_printf(tls *TLS, format uintptr /* *int8 */, args ...interface{}) int32 {
	return goFprintf(os.Stdout, format, args...)
}

func goFprintf(w io.Writer, format uintptr /* *int8 */, ap ...interface{}) int32 {
	var b buffer.Bytes
	written := 0
	for {
		c := *(*int8)(unsafe.Pointer(format))
		format++
		switch c {
		case 0:
			_, err := b.WriteTo(w)
			b.Close()
			if err != nil {
				return -1
			}

			return int32(written)
		case '%':
			modifiers := ""
			long := 0
			short := 0
			hash := false
			var w []interface{}
		more:
			c := *(*int8)(unsafe.Pointer(format))
			format++
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-':
				modifiers += string(c)
				goto more
			case '#':
				hash = true
				goto more
			case '*':
				w = append(w, VAint32(&ap))
				modifiers += string(c)
				goto more
			case 'c':
				arg := VAint32(&ap)
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sc", modifiers), append(w, arg)...)
				written += n
			case 'd', 'i':
				var arg interface{}
				switch long {
				case 0:
					arg = VAint32(&ap)
				case 1:
					arg = vaLong(&ap)
				default:
					arg = VAint64(&ap)
				}
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sd", modifiers), append(w, arg)...)
				written += n
			case 'l':
				long++
				goto more
			case 'f':
				if hash {
					panic(fmt.Errorf("TODO #f"))
				}

				arg := VAfloat64(&ap)
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sf", modifiers), append(w, arg)...)
				written += n
			case 'h':
				short++
				goto more
			case 'o':
				nz := false
				var arg interface{}
				switch long {
				case 0:
					v := VAuint32(&ap)
					nz = v != 0
					arg = v
				case 1:
					v := vaULong(&ap)
					nz = v != 0
					arg = v
				default:
					v := VAuint64(&ap)
					nz = v != 0
					arg = v
				}
				if hash && nz {
					b.WriteByte('0')
					written++
				}

				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%so", modifiers), append(w, arg)...)
				written += n
			case 'p':
				arg := VAuintptr(&ap)
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sp", modifiers), append(w, arg)...)
				written += n
			case 'g':
				if hash {
					panic(fmt.Errorf("TODO #g"))
				}

				arg := VAfloat64(&ap)
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sg", modifiers), append(w, arg)...)
				written += n
			case 's':
				arg := VAuintptr(&ap)
				if arg == 0 {
					break
				}

				var b2 buffer.Bytes
				for {
					c := *(*int8)(unsafe.Pointer(arg))
					arg++
					if c == 0 {
						n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%ss", modifiers), append(w, b2.Bytes())...)
						b2.Close()
						written += n
						break
					}

					b2.WriteByte(byte(c))
				}
			case 'u':
				var arg interface{}
				switch long {
				case 0:
					arg = VAuint32(&ap)
				case 1:
					arg = vaULong(&ap)
				default:
					arg = VAuint64(&ap)
				}
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sd", modifiers), append(w, arg)...)
				written += n
			case 'x':
				nz := false
				var arg interface{}
				switch long {
				case 0:
					v := VAuint32(&ap)
					nz = v != 0
					arg = v
				case 1:
					v := vaULong(&ap)
					nz = v != 0
					arg = v
				default:
					v := VAuint64(&ap)
					nz = v != 0
					arg = v
				}
				if hash && nz {
					b.WriteString("0x")
					written += 2
				}

				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sx", modifiers), append(w, arg)...)
				written += n
			case 'X':
				nz := false
				var arg interface{}
				switch long {
				case 0:
					v := VAuint32(&ap)
					nz = v != 0
					arg = v
				case 1:
					v := vaULong(&ap)
					nz = v != 0
					arg = v
				default:
					v := VAuint64(&ap)
					nz = v != 0
					arg = v
				}
				if hash && nz {
					b.WriteString("0X")
					written += 2
				}

				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sX", modifiers), append(w, arg)...)
				written += n
			default:
				panic(fmt.Errorf("TODO %q", "%"+string(c)))
			}
		default:
			b.WriteByte(byte(c))
			written++
			if c == '\n' {
				if _, err := b.WriteTo(w); err != nil {
					b.Close()
					return -1
				}
				b.Reset()
			}
		}
	}
}

// int sprintf(char *str, const char *format, ...);
func X__builtin_sprintf(tls *TLS, str, format uintptr, args ...interface{}) int32 {
	w := memWriter(str)
	n := goFprintf(&w, format, args...)
	w.WriteByte(0)
	return n
}

// int sprintf(char *str, const char *format, ...);
func Xsprintf(tls *TLS, str, format uintptr, args ...interface{}) int32 {
	return X__builtin_sprintf(tls, str, format, args...)
}

// int fputc(int c, FILE *stream);
func Xfputc(tls *TLS, c int32, stream uintptr) int32 {
	w := files.writer(stream)
	if _, err := w.Write([]byte{byte(c)}); err != nil {
		return stdio.XEOF
	}

	return int32(byte(c))
}

// int putc(int c, FILE *stream);
func Xputc(tls *TLS, c int32, stream uintptr) int32 {
	panic("TODO putc")
}

// int putc(int c, FILE *stream);
func X_IO_putc(tls *TLS, c int32, stream uintptr) int32 { return Xputc(tls, c, stream) }

// FILE *fopen64(const char *path, const char *mode);
func Xfopen64(tls *TLS, path, mode uintptr) uintptr {
	p := GoString(path)
	var u uintptr
	switch p {
	case os.Stderr.Name():
		u = stderr
	case os.Stdin.Name():
		u = stdin
	case os.Stdout.Name():
		u = stdout
	default:
		var f *os.File
		var err error
		switch mode := GoString(mode); mode {
		case "a":
			if f, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
				switch {
				case os.IsPermission(err):
					tls.setErrno(errno.XEPERM)
				default:
					tls.setErrno(errno.XEACCES)
				}
			}
		case "r", "rb":
			if f, err = os.OpenFile(p, os.O_RDONLY, 0666); err != nil {
				switch {
				case os.IsNotExist(err):
					tls.setErrno(errno.XENOENT)
				case os.IsPermission(err):
					tls.setErrno(errno.XEPERM)
				default:
					tls.setErrno(errno.XEACCES)
				}
			}
		case "w":
			if f, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
				switch {
				case os.IsPermission(err):
					tls.setErrno(errno.XEPERM)
				default:
					tls.setErrno(errno.XEACCES)
				}
			}
		default:
			panic(mode)
		}
		if f != nil {
			u = Xmalloc(tls, ptrSize)
			files.add(f, u)
		}
	}
	return u
}

// size_t fwrite(const void *ptr, size_t size, size_t nmemb, FILE *stream);
func fwrite(tls *TLS, ptr uintptr, size, nmemb size_t, stream uintptr) size_t {
	hi, lo := mathutil.MulUint128_64(uint64(size), uint64(nmemb))
	if hi != 0 || lo > uint64(len(rawmem{})) {
		tls.setErrno(errno.XE2BIG)
		return 0
	}

	n, err := files.writer(stream).Write((*rawmem)(unsafe.Pointer(ptr))[:lo])
	if err != nil {
		tls.setErrno(errno.XEIO)
	}
	return size_t(n) / size
}

// int fclose(FILE *stream);
func Xfclose(tls *TLS, stream uintptr) int32 {
	switch stream {
	case stdin, stdout, stderr:
		tls.setErrno(errno.XEIO)
		return stdio.XEOF
	}

	f := files.extract(stream)
	if f == nil {
		tls.setErrno(errno.XEBADF)
		return stdio.XEOF
	}

	Xfree(tls, stream)
	if err := f.Close(); err != nil {
		tls.setErrno(errno.XEIO)
		return stdio.XEOF
	}

	return 0
}

// size_t fread(void *ptr, size_t size, size_t nmemb, FILE *stream);
func fread(tls *TLS, ptr uintptr, size, nmemb size_t, stream uintptr) size_t {
	hi, lo := mathutil.MulUint128_64(uint64(size), uint64(nmemb))
	if hi != 0 || lo > uint64(len(rawmem{})) {
		tls.setErrno(errno.XE2BIG)
		return 0
	}

	n, err := files.reader(stream).Read((*rawmem)(unsafe.Pointer(ptr))[:lo])
	if err != nil {
		tls.setErrno(errno.XEIO)
	}
	return size_t(n) / size
}

func fseek(tls *TLS, stream uintptr, offset long_t, whence int32) int32 {
	f := files.file(stream)
	if f == nil {
		tls.setErrno(errno.XEBADF)
		return -1
	}

	if _, err := f.Seek(int64(offset), int(whence)); err != nil {
		tls.setErrno(errno.XEINVAL)
		return -1
	}

	return 0
}

func ftell(tls *TLS, stream uintptr) long_t {
	f := files.file(stream)
	if f == nil {
		tls.setErrno(errno.XEBADF)
		return -1
	}

	n, err := f.Seek(0, os.SEEK_CUR)
	if err != nil {
		tls.setErrno(errno.XEBADF)
		return -1
	}

	return long_t(n)
}

// int fgetc(FILE *stream);
func Xfgetc(tls *TLS, stream uintptr) int32 {
	p := buffer.Get(1)
	if _, err := files.reader(stream).Read(*p); err != nil {
		buffer.Put(p)
		return stdio.XEOF
	}

	r := int32((*p)[0])
	buffer.Put(p)
	return r
}

// char *fgets(char *s, int size, FILE *stream);
func Xfgets(tls *TLS, s uintptr, size int32, stream uintptr) uintptr {
	f := files.reader(stream)
	p := buffer.Get(1)
	b := *p
	w := memWriter(s)
	ok := false
	for i := int(size) - 1; i > 0; i-- {
		_, err := f.Read(b)
		if err != nil {
			if !ok {
				buffer.Put(p)
				return 0
			}

			break
		}

		ok = true
		w.WriteByte(b[0])
		if b[0] == '\n' {
			break
		}
	}
	w.WriteByte(0)
	buffer.Put(p)
	return s

}

// int __builtin_fprintf(void* stream, const char *format, ...);
func X__builtin_fprintf(tls *TLS, stream, format uintptr, args ...interface{}) int32 {
	return goFprintf(files.writer(stream), format, args...)
}

// int fprintf(FILE * stream, const char *format, ...);
func Xfprintf(tls *TLS, stream, format uintptr, args ...interface{}) int32 {
	return X__builtin_fprintf(tls, stream, format, args...)
}

// int fflush(FILE *stream);
func Xfflush(tls *TLS, stream uintptr) int32 {
	f := files.file(stream)
	if f == nil {
		tls.setErrno(stdio.XEOF)
		return -1
	}

	if err := f.Sync(); err != nil {
		tls.setErrno(err)
		return -1
	}

	return 0
}

// int vprintf(const char *format, va_list ap);
func Xvprintf(tls *TLS, format uintptr, ap []interface{}) int32 {
	return goFprintf(os.Stdout, format, ap...)
}

// int vfprintf(FILE *stream, const char *format, va_list ap);
func Xvfprintf(tls *TLS, stream, format uintptr, ap []interface{}) int32 {
	return goFprintf(files.writer(stream), format, ap...)
}

// void rewind(FILE *stream);
func Xrewind(tls *TLS, stream uintptr) { fseek(tls, stream, 0, int32(os.SEEK_SET)) }

// FILE *popen(const char *command, const char *type);
func Xpopen(tls *TLS, command, typ uintptr) uintptr {
	panic("TODO popen")
}

// int pclose(FILE *stream);
func Xpclose(tls *TLS, stream uintptr) int32 {
	panic("TODO pclose")
}

// size_t fwrite(const void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfwrite(tls *TLS, ptr uintptr, size, nmemb size_t, stream uintptr) size_t {
	return fwrite(tls, ptr, size, nmemb, stream)
}

// size_t fread(void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfread(tls *TLS, ptr uintptr, size, nmemb size_t, stream uintptr) size_t {
	return fread(tls, ptr, size, nmemb, stream)
}

// int fseek(FILE *stream, long offset, int whence);
func Xfseek(tls *TLS, stream uintptr, offset long_t, whence int32) int32 {
	return fseek(tls, stream, offset, whence)
}

// long ftell(FILE *stream);
func Xftell(tls *TLS, stream uintptr) long_t { return ftell(tls, stream) }

// int setvbuf(FILE *stream, char *buf, int mode, size_t size);
func Xsetvbuf(tls *TLS, stream, buf uintptr, mode int32, size size_t) int32 {
	return 0 //TODO setvbuf
}
