// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/cznic/ccir/libc/errno"
	"github.com/cznic/ccir/libc/pthread"
)

type mu struct {
	*sync.Cond
	attr  int32
	count int
	owner uintptr
	sync.Mutex
}

type mutexMap struct {
	m map[uintptr]*mu
	sync.Mutex
}

func (m *mutexMap) mu(p uintptr) *mu {
	m.Lock()
	r := m.m[p]
	if r == nil {
		r = &mu{}
		r.Cond = sync.NewCond(&r.Mutex)
		m.m[p] = r
	}
	m.Unlock()
	return r
}

type threadState struct {
	c        chan struct{}
	detached bool
	retval   uintptr
}

type threadMap struct {
	m map[uintptr]*threadState
	sync.Mutex
}

var (
	mutexes = &mutexMap{m: map[uintptr]*mu{}}
	threads = &threadMap{m: map[uintptr]*threadState{}}
)

// extern int pthread_mutexattr_init(pthread_mutexattr_t * __attr);
func Xpthread_mutexattr_init(tls *TLS, attr uintptr) int32 {
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_init(%#x) %v\n", attr, r)
	}
	return r
}

// extern int pthread_mutexattr_settype(pthread_mutexattr_t * __attr, int __kind);
func Xpthread_mutexattr_settype(tls *TLS, attr uintptr, kind int32) int32 {
	*(*int32)(unsafe.Pointer(attr)) = kind
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_settype(%#x, %v) %v\n", attr, kind, r)
	}
	return r
}

// extern int pthread_mutex_init(pthread_mutex_t * __mutex, pthread_mutexattr_t * __mutexattr);
func Xpthread_mutex_init(tls *TLS, mutex, mutexattr uintptr) int32 {
	attr := int32(pthread.XPTHREAD_MUTEX_NORMAL)
	if mutexattr != 0 {
		attr = *(*int32)(unsafe.Pointer(mutexattr))
	}
	mutexes.mu(mutex).attr = attr
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_init(%#x, %#x) %v\n", mutex, mutexattr, r)
	}
	return r
}

// extern int pthread_mutexattr_destroy(pthread_mutexattr_t * __attr);
func Xpthread_mutexattr_destroy(tls *TLS, attr uintptr) int32 {
	*(*int32)(unsafe.Pointer(attr)) = -1
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_destroy(%#x) %v\n", attr, r)
	}
	return r
}

// extern int pthread_mutex_destroy(pthread_mutex_t * __mutex);
func Xpthread_mutex_destroy(tls *TLS, mutex uintptr) int32 {
	mutexes.Lock()
	delete(mutexes.m, mutex)
	mutexes.Unlock()
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_destroy(%#x) %v\n", mutex, r)
	}
	return r
}

// extern int pthread_mutex_lock(pthread_mutex_t * __mutex);
func Xpthread_mutex_lock(tls *TLS, mutex uintptr) int32 {
	threadID := tls.threadID
	mu := mutexes.mu(mutex)
	var r int32
	mu.Lock()
	switch mu.attr {
	case pthread.XPTHREAD_MUTEX_NORMAL:
		if mu.count == 0 {
			mu.owner = threadID
			mu.count = 1
			break
		}

		for mu.count != 0 {
			mu.Cond.Wait()
		}
		mu.owner = threadID
		mu.count = 1
	case pthread.XPTHREAD_MUTEX_RECURSIVE:
		if mu.count == 0 {
			mu.owner = threadID
			mu.count = 1
			break
		}

		if mu.owner == threadID {
			mu.count++
			break
		}

		panic("TODO")
	default:
		panic(fmt.Errorf("attr %#x", mu.attr))
	}
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_lock(%#x: %+v [thread id %v]) %v\n", mutex, mu, threadID, r)
	}
	mu.Unlock()
	return r
}

// int pthread_mutex_trylock(pthread_mutex_t *mutex);
func Xpthread_mutex_trylock(tls *TLS, mutex uintptr) int32 {
	threadID := tls.threadID
	mu := mutexes.mu(mutex)
	var r int32
	mu.Lock()
	switch mu.attr {
	case pthread.XPTHREAD_MUTEX_NORMAL:
		if mu.count == 0 {
			mu.count = 1
			mu.owner = threadID
			break
		}

		r = errno.XEBUSY
	default:
		panic(fmt.Errorf("attr %#x", mu.attr))
	}
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_trylock(%#x: %+v [thread id %v]) %v\n", mutex, mu, threadID, r)
	}
	mu.Unlock()
	return r
}

// extern int pthread_mutex_unlock(pthread_mutex_t * __mutex);
func Xpthread_mutex_unlock(tls *TLS, mutex uintptr) int32 {
	threadID := tls.threadID
	mu := mutexes.mu(mutex)
	var r int32
	mu.Lock()
	switch mu.attr {
	case pthread.XPTHREAD_MUTEX_NORMAL:
		if mu.count == 0 {
			panic("TODO")
		}

		mu.owner = 0
		mu.count = 0
		mu.Cond.Broadcast()
	case pthread.XPTHREAD_MUTEX_RECURSIVE:
		if mu.count == 0 {
			panic("TODO")
		}

		if mu.owner == threadID {
			mu.count--
			if mu.count != 0 {
				break
			}

			mu.owner = 0
			mu.Cond.Broadcast()
			break
		}

		panic("TODO")
	default:
		panic(fmt.Errorf("TODO %#x", mu.attr))
	}
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_unlock(%#x: %+v [thread id %v]) %v\n", mutex, mu, threadID, r)
	}
	mu.Unlock()
	return r
}

// pthread_t pthread_self(void);
func Xpthread_self(tls *TLS) pthread_t {
	threadID := tls.threadID
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_self() %v\n", threadID)
	}
	return pthread_t(threadID)
}

// extern int pthread_equal(pthread_t __thread1, pthread_t __thread2);
func Xpthread_equal(tls *TLS, thread1, thread2 pthread_t) int32 {
	if thread1 == thread2 {
		return 1
	}

	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_equal(%v, %v) %v\n", thread1, thread2, r)
	}
	return r
}

// int pthread_join(pthread_t thread, void **value_ptr);
func Xpthread_join(tls *TLS, thread pthread_t, value_ptr uintptr) int32 {
	panic("TODO pthread_join")
	threads.Lock()
	t := threads.m[uintptr(thread)]
	threads.Unlock()
	if t != nil {
		<-t.c
		if value_ptr != 0 {
			*(*uintptr)(unsafe.Pointer(value_ptr)) = t.retval
		}
		threads.Lock()
		delete(threads.m, uintptr(thread))
		threads.Unlock()
	}
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_join(%v, %#x) %v\n", thread, value_ptr, r)
	}
	return r
}

// int pthread_create(pthread_t *restrict thread, const pthread_attr_t *restrict attr, void *(*start_routine)(void*), void *restrict arg);
func Xpthread_create(tls *TLS, thread, attr, start_routine, arg uintptr) int32 {
	panic("TODO pthread_create")
	if attr != 0 {
		panic("TODO")
	}

	new := NewTLS()
	*(*uint64)(unsafe.Pointer(thread)) = uint64(new.threadID)
	threads.Lock()
	t := &threadState{c: make(chan struct{})}
	threads.m[new.threadID] = t
	threads.Unlock()
	ch := make(chan struct{})
	go func() {
		close(ch)
		t.retval = (*(*func(*TLS, uintptr) uintptr)(unsafe.Pointer(&start_routine)))(new, arg)
		if ptrace {
			fmt.Fprintf(os.Stderr, "thread #%#x finished: %#x\n", new.threadID, t.retval)
		}
		close(t.c)
		if t.detached {
			threads.Lock()
			delete(threads.m, new.threadID)
			threads.Unlock()
			if ptrace {
				fmt.Fprintf(os.Stderr, "thread #%#x was detached", new.threadID)
			}
		}
	}()
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_create(%#x, %#x, fn, %#x) #%#x %v\n", thread, attr, arg, new.threadID, r)
	}
	<-ch
	return r
}
