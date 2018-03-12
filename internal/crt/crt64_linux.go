// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64

package crt

type (
	long_t    = int64
	pthread_t = uint64
	rawmem    [1<<50 - 1]byte
	size_t    = uint64
	ssize_t   = int64
	ulong_t   = uint64
)
