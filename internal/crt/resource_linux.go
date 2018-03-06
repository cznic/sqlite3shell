// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// int getrusage(int who, struct rusage *usage);
func Xgetrusage(tls *TLS, who int32, usage uintptr) int32 {
	panic("TODO getrusage")
}
