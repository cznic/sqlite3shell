// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
)

func X__builtin_assert_fail(tls *TLS, file uintptr, line int32, fn, msg uintptr) {
	panic(fmt.Errorf("%s.%s:%d: assertion failure: %s", GoString(file), GoString(fn), line, GoString(msg)))
}
