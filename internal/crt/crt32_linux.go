// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build 386

//TODO drop the _linux suffix when the Windows port is ready again.

package crt

type (
	long_t    = int32
	pthread_t = uint32
	rawmem    [1<<31 - 1]byte
	size_t    = uint32
	ssize_t   = int32
	ulong_t   = uint32
)
