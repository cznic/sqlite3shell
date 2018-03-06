// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"math/bits"
)

// uint64_t __builtin_bswap64 (uint64_t x)
func X__builtin_bswap64(tls *TLS, x uint64) uint64 { return bits.ReverseBytes64(x) }
