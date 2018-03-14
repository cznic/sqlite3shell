// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build crt.strace

package crt

import (
	"fmt"
	"strings"

	fcntl2 "github.com/cznic/ccir/libc/fcntl"
)

const strace = true

func cmdString(cmd int32) string {
	switch cmd {
	default:
		return fmt.Sprintf("%#x", cmd)
	}
}

func modeString(flag int32) string {
	if flag == 0 {
		return "0"
	}

	var a []string
	for _, v := range []struct {
		int32
		string
	}{
		{fcntl2.XO_APPEND, "O_APPEND"},
		{fcntl2.XO_CREAT, "O_CREAT"},
		{fcntl2.XO_EXCL, "O_EXCL"},
		{fcntl2.XO_RDONLY, "O_RDONLY"},
		{fcntl2.XO_RDWR, "O_RDWR"},
		{fcntl2.XO_WRONLY, "O_WRONLY"},
	} {
		if flag&v.int32 != 0 {
			a = append(a, v.string)
		}
	}
	return strings.Join(a, "|")
}
