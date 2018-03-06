// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"unsafe"
)

type Spasswd struct {
	Xpw_name   uintptr // *int8
	Xpw_passwd uintptr // *int8
	Xpw_uid    T__uid_t
	Xpw_gid    T__gid_t
	Xpw_gecos  uintptr // *int8
	Xpw_dir    uintptr // *int8
	Xpw_shell  uintptr // *int8
}

// struct passwd *getpwuid(uid_t uid);
func Xgetpwuid(tls *TLS, uid uint32) uintptr {
	u, err := user.LookupId(fmt.Sprint(uid))
	if err != nil {
		tls.setErrno(err)
		return 0
	}

	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		tls.setErrno(err) //TODO Exxx
		return 0
	}

	p, err := Malloc(int(unsafe.Sizeof(Spasswd{})))
	if err != nil {
		tls.setErrno(err) //TODO Exxx
		return 0
	}

	*(*Spasswd)(unsafe.Pointer(p)) = Spasswd{
		Xpw_name:   CString(u.Username),
		Xpw_passwd: CString("x"),
		Xpw_uid:    uid,
		Xpw_gid:    T__gid_t(gid),
		Xpw_gecos:  CString(u.Name),
		Xpw_dir:    CString(u.HomeDir),
		Xpw_shell:  CString(os.Getenv("SHELL")),
	}
	return p
}
